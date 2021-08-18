search replace `change-me-service-name` with your service name
search for 'change-me' and adjust or remove code, then remove this line.

CI/CD setup using GCP

first time setup:
```
   go env -w GOPRIVATE=github.com/ridebeam
   git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"
```

to build and run locally:

```
   # create certs and local dir, only necessay on clean repo
   mkdir -p bin
   scripts/certs-run.sh

   go build -o bin ./... && bin/change-me-service-name > log 2>&1
   
   # for log output in json set
   LOGGER_ENCODING=json
```

or to build docker image
```
    docker build --build-arg "SSH_PRIVATE_KEY=$(cat ~/.ssh/id_rsa)" -t change-me-service-name .
```

development default overrides go into configs/defaults.yaml.
adjust settings in configs/overrides.yaml based on template

json encoding is done using code generated marshaller for improved speed. When changing
any of the JSON schema definitions make sure to execute generators:
```
    go generate ./...
```

## Connect to Cloud SQL database
To connect to the remote database you'll need to use Cloud SQL proxy.

- Download the proxy at https://cloud.google.com/sql/docs/postgres/sql-proxy#macos-64-bit or install using `gcloud components install cloud_sql_proxy`
- Then in a terminal run `./cloud_sql_proxy -instances=ridebeam-staging:asia-southeast1:pg-sg-02=tcp:5432` to run the proxy for dev.
- You can now connect to the database on `localhost:5432`
Username `postgres`.

### Get credentials for the database
- Go to Kubernetes configuration and look for `cloudsql-pg-sg-02`
- Get the YAML
- You should find `PG_PASSWORD`. It's Base64 encoded so you'll need to decode it.
- With kubectl installed you can also run `kubectl get secret/cloudsql-pg-sg-02 -n dev -o json | jq '.data.PG_PASSWORD' -r | base64 -d`

## gRPC and protobuf

Please refer to https://beammobility.atlassian.net/wiki/spaces/BE/pages/2122907651/gRPC+and+protobuf

### debugging production
It is possible to reach services directly, bypassing any external facing load balancers.
This allows us to have debugging endpoints, otherwise not reachable. We can use `kubectl port-forward` for this purpose.
To create a tunnel execute:

```
kubectl port-forward service/<service name> -n <env> <local-port>:80
```

### k8s

Pre-install certs by running 
```
### dev 
scripts/certs-run.sh
kubectl create secret generic https-credentials --namespace=dev \
  --from-file=srv-cert.pem=configs/certs/srv-cert.pem \
  --from-file=srv-key.pem=configs/certs/srv-key.pem

### staging 
scripts/certs-run.sh
kubectl create secret generic https-credentials --namespace=staging \
  --from-file=srv-cert.pem=configs/certs/srv-cert.pem \
  --from-file=srv-key.pem=configs/certs/srv-key.pem

### production 
scripts/certs-run.sh
kubectl create secret generic https-credentials --namespace=production \
  --from-file=srv-cert.pem=configs/certs/srv-cert.pem \
  --from-file=srv-key.pem=configs/certs/srv-key.pem
```
