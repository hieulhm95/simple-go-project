# Install and run on your computer

Run service: `go run main.go service`

Run worker with Kafka: `go run main.go worker_kafka`

# Testing instructions

### Local run

- Spin up dependent services. This will spin up all containers which are necessary services such as Database, Queue,… to run test locally.
  `make dev-up`

- Export necessary environment variables
  `export $(xargs < ./setup/.local.env)`

- Run test locally
  `make test`

- Terminate all containers.
  `make dev-down`
