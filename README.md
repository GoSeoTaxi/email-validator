# Email Validator

This application provides a gRPC-based email validation service using Golang. It checks the validity of email addresses by verifying their format and checking the associated domain's MX records. Additionally, it includes a Gateway for handling HTTP requests and forwarding them to the gRPC server.

## Features

- Validates email format.
- Checks MX records of the email domain.
- Caches results using Redis.
- Supports both gRPC and HTTP interfaces.
- Implements a gRPC Gateway for RESTful API access.

## Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.22 or later)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/) (optional, as it can be included in Docker installation)

## Setup

### Generate Protocol Buffers

The application uses Protocol Buffers to define the service and message types. Run the following command to generate the Go code from the `.proto` files:

```
make generate
```

### Build the Docker Images

You can build the Docker images for both the email validator and the gateway using:

```
make build
```

### Running the Application

To run the application with Docker Compose, use the following command:

```
make docker-run
```

This will start the email validator, gateway, and Redis in separate containers.

### Stopping the Application

To stop and remove the containers, run:

```
make docker-stop
```

## Testing

To run the tests for the application, use:

```
make test
```

## API Endpoints

### gRPC Endpoint

The gRPC server listens for requests on the configured port. You can use a gRPC client to send requests to validate emails.

### HTTP Gateway

The application also exposes an HTTP Gateway. You can send requests to the Gateway, which will forward them to the gRPC server. 

- **POST /v1/validate**

**Request Body:**
```
{
  "email": "example@example.com"
}
```

**Response:**
- **200 OK**: If the email is valid.
- **503 Service Unavailable**: If the email server cannot be reached.
- **400 Bad Request**: If the email format is invalid.

## Environment Variables

You can set the following environment variables to configure the application:

- `GRPC_SERVER_ENDPOINT`: The endpoint for the gRPC server (default: `email-validator:50051`).
- `HTTP_PORT`: The port for the HTTP server (default: `8080`).
- `REDIS_HOST`: The hostname for the Redis server (default: `redis`).
- `REDIS_PORT`: The port for the Redis server (default: `6379`).
- `REDIS_DB`: The Redis database number (default: `0`).
- `REDIS_MAXMEMORY`: The maximum memory limit for Redis (default: `100mb`).
- `DNS_HOSTS`: The DNS servers to use for MX record lookups (default: `1.1.1.1,1.0.0.1`).

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for suggestions and improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
