package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ridebeam/golang-skeleton/pkg/server/grpc/gateway"

	"github.com/ridebeam/golang-skeleton/pkg/example"
	"github.com/ridebeam/golang-skeleton/pkg/serde"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/ridebeam/genproto/gengo/ridebeam/vehicle"
	commonGRPC "github.com/ridebeam/go-common/grpc"
	"google.golang.org/grpc"

	"github.com/ridebeam/golang-skeleton/pkg/monitoring"

	commonConfig "github.com/ridebeam/go-common/config"
	"github.com/ridebeam/go-common/kafka"
	"github.com/ridebeam/go-common/kv"
	"github.com/ridebeam/go-common/monitor"
	"github.com/ridebeam/go-common/otel"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
)

func main() {
	// --- we start by initializing monitoring ---
	rawLogger := monitor.InitLogger()
	defer ignoreError(rawLogger.Sync())

	logger := rawLogger.Sugar()
	ctx := monitor.WithLogger(context.Background(), logger)

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				logger.Errorw(fmt.Sprintf("main panicked because: %s", err.Error()), zap.Error(err))
			} else {
				logger.Errorw(fmt.Sprintf("main panicked because: %v", r), "panic", r)
			}
		}
	}()

	if err := run(ctx, logger); err != nil {
		logger.Fatalw(fmt.Sprintf("main panicked because: %s", err.Error()), zap.Error(err))
	}
}

func run(ctx context.Context, logger *zap.SugaredLogger) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	scope, monitorCloser, err := initMonitoring(cfg, logger)
	if err != nil {
		return err
	}
	defer monitorCloser()

	ctx = monitor.WithScope(ctx, scope)
	monitoring.RegisterServiceViews(ctx)

	// --- next we initialize any storage, like postgres, redis and kafka producers ---
	// --- essentially all the things that our components might depend on ---
	kafkaCfg := cfg.Kafka.Connection
	kafkaCfg.Service = cfg.Service

	kafkaProducerCfg := kafkaCfg
	// add extra config for producer here

	// we create an sync producer for each of the main topics
	_, syncProducer, err := kafka.Producer(ctx, kafkaProducerCfg)
	if err != nil {
		return err
	}
	exampleTopicProducer := kafka.MessageProducer(syncProducer).WithTopic(cfg.Kafka.TopicExample)

	exampleEmitter := func(ctx context.Context, output serde.ExampleOutput) error {
		data, err := output.MarshalJSON()
		if err != nil {
			scope.Increment("output-serialize-error")
			logger.DPanicw("could not serialize example message", kafka.TagKeyTargetTopic, cfg.Kafka.TopicExample)
			return err
		}
		return exampleTopicProducer(ctx, nil, data)
	}

	repo, repoCloser, err := kv.InitVersionedKV(ctx, cfg.Service, kv.Config{
		EnabledInDevelopment: cfg.DB.EnabledInDevelopment,
		Postgres:             cfg.DB.Postgres,
	}, cfg.DB.TableExample)
	if err != nil {
		return err
	}
	// TODO: fix it so we don't return a nil closer in go-common/kv
	defer func() {
		if repoCloser != nil {
			repoCloser()
		}
	}()

	// --- around here, we initialize our main service components, using above dependencies ---
	ex := example.NewExample(repo, exampleEmitter)

	// -- finally we are setting up all the sources of incoming requests like kafka consumers and grpc server ---
	consumeExample, err := kafka.CloseableConsumer(
		ctx,
		kafka.ConfigConsumer{
			Config: kafkaCfg,
			Topics: []string{cfg.Kafka.TopicExample},
		},
		func(ctx context.Context, data []byte, _ kafka.Metadata) {
			var m serde.ExampleModel
			if err := m.UnmarshalJSON(data); err != nil {
				monitor.Measurement(ctx).Increment("kafka-deserialize-error")
				monitor.Logger(ctx).DPanicw("could no deserialize device action message", zap.Error(err))
				return
			}
			ex.Handle(ctx, m)
		},
	)
	if err != nil {
		return err
	}
	// stop listening to kafka (forces all partitions to move to other instances)
	defer consumeExample.Close()

	// initializing listeners for shutdown requests. Actual graceful shutdown begins once quit channel gets closed
	quit := initGracefulShutdown(ctx)

	// we initialize http servers last, as this will enable health checks, and signals that we are ready to receive
	go startGRPC(ctx, cfg, gateway.Components{})

	// waiting for graceful shutdown to begin
	<-quit

	return nil
}

// TODO move into common repo
func initMonitoring(cfg Config, logger *zap.SugaredLogger) (monitor.Scope, monitor.OnShutdown, error) {
	if commonConfig.Development() {
		return nil, func() {}, fmt.Errorf("not going to try connecting to GC from a non-prod environment")
	}

	scope, closerOC, err := monitor.NewOpenCensusScope(cfg.Service, 0)
	if err != nil {
		logger.Warnw("could not initialize GC Metrics, trying StatsD next", zap.Error(err))

		scope, err = initStatsD(cfg)
		if err != nil {
			logger.Warnw("could not initialize StatsD, fallback to NoOpScope", zap.Error(err))
			scope = monitor.NewNoOpScope()
		}

		return scope, func() {}, nil
	}

	closerOtel, err := otel.InitGCPExporter(context.Background(), otel.Config{
		SampleRate: cfg.Monitor.TracingSampleRate,
	})

	closer := func() {
		if closerOC != nil {
			closerOC()
		}
		if closerOtel != nil {
			closerOtel()
		}
	}

	return scope, closer, err
}

// TODO move to common repo (or consider dropping statsd)
func initStatsD(cfg Config) (monitor.Scope, error) {
	options := []statsd.Option{
		statsd.Address(cfg.Monitor.StatsD.Address),
		statsd.Tags("env", cfg.Service.Env),
		statsd.Tags("service", cfg.Service.Name),
	}

	switch cfg.Monitor.StatsD.TagsFormat {
	case "datadog", "cloudwatch":
		options = append(options, statsd.TagsFormat(statsd.Datadog))
	case "influxdb":
		options = append(options, statsd.TagsFormat(statsd.InfluxDB))
	}

	return monitor.NewStatsDScope(options...)
}

func initGracefulShutdown(ctx context.Context) chan struct{} {
	quit := make(chan struct{})
	go func() {
		// We initiate a graceful shutdowns when quit via
		//  SIGINT  (Ctrl+C)
		//  SIGTERM (kill -15)
		cSig := make(chan os.Signal, 1)
		signal.Notify(cSig, os.Interrupt)
		signal.Notify(cSig, syscall.SIGTERM)

		// we block until we get a signal
		<-cSig

		monitor.Logger(ctx).Infof("received shutdown signal, will initiate graceful shutdown")
		// signal quit
		close(quit)
	}()
	return quit
}

func startGRPC(ctx context.Context, cfg Config, components gateway.Components) {
	err := commonGRPC.InitializeAndStart(
		ctx,
		cfg.Server,
		[]commonGRPC.RegisterService{
			func(srv grpc.ServiceRegistrar) {
				pb.RegisterGatewayServer(
					srv,
					gateway.NewServer(gateway.Config{}, components),
				)
			},
		},
		[]commonGRPC.RegisterJsonProxy{
			func(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
				return pb.RegisterGatewayHandlerFromEndpoint(ctx, mux, address, options)
			},
		},
	)
	if err != nil {
		monitor.Logger(ctx).Fatalf("grpc server listen and serve error: %v\n", err)
	}
}

// helper function for ignoring errors and satisfying the linter
func ignoreError(err error) {}
