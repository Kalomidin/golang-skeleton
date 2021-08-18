package gateway

type Config struct {
	// list any key value pairs to configure the system here
}

type Components struct {
	// list any repos, and other stuff the endpoints might need
}

type server struct {
	pb.UnimplementedExampleServer
	Config
	Components
}

func NewServer(cfg Config, components Components) *server {
	s := server{
		Config:     cfg,
		Components: components,
	}
	return &s
}
