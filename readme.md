## Golang Clean Architecture Project

This project exemplifies a robust implementation of Golang Clean Architecture, a software design paradigm that distinctly separates application concerns into various layers. This architecture underscores the significance of maintaining the business logic's independence and decoupling it from external frameworks and technologies. Employing this approach ensures the project's high level of testability, scalability, and maintainability. The application integrates cutting-edge technologies such as Docker for containerization, Fiber as the web framework, JWT for authentication, and OpenTelemetry for distributed tracing and monitoring. The project's structure is meticulously organized into distinct folders, aligning with its architectural responsibilities.

### How to Run:

1. Clone this repository.
2. Execute `docker-compose up -d` to launch the Golang application within a Docker container.
3. Run `make migrate-up` to apply database migrations.
4. Access the application at `http://localhost:8080`.

### Folder Structure:

- `api/`: Contains API specifications or contract APIs.
- `cmd/`: Stores the `main.go` file executed by Docker.
- `db/`: Houses migration files or scripts for running migrations.
- `internal/`: Contains code restricted from external package access.
  - `config/`: Manages application configuration code.
  - `delivery/`: Handles code related to the delivery layer.
  - `entity/`: Contains entity data-related code.
  - `model/`: Manages model data-related code.
  - `gateway/`: Interfaces with external APIs.
  - `repository/`: Manages repository data-related code.
  - `usecase/`: Handles use case-related code.

### Important Files:

- `main.go`: Entry point executed by Docker.
- `docker-compose.yml`: Configures Docker container execution for the application.
- `Dockerfile`: Builds the Golang application within a Docker container.
- `Makefile`: Manages application scripts.
- `go.mod`: Manages Golang dependencies.
- `go.sum`: Verifies Golang dependency checksums.

### References:

- [Clean Architecture](https://github.com/khannedy/golang-clean-architecture)
- [Golang Docker](https://docs.docker.com/language/golang/)
- [Golang Fiber](https://docs.gofiber.io/)
- [JWT Token](https://pkg.go.dev/github.com/golang-jwt/jwt/v4)
- [Validator](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [Jaeger](https://github.com/jaegertracing/jaeger-client-go)
- [Prometheus](https://github.com/prometheus/client_golang)
- [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-go)
- [Air](https://github.com/go-air/air)
