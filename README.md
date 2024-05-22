# Go Multiple Query Template Project

## Configuration

The service uses environment variables for configuration. The following variables are used:

| Name             | Description                                         | Default Value   | Required |
| ---------------- | --------------------------------------------------- | --------------- | -------- |
| `HOST`           | The host on which the service is running.           | localhost       | false    |
| `PORT`           | The port on which the service is running.           | 8080            | false    |
| `PROXY_HEADER`   | The header to use for proxying requests.            | X-Forwarded-For | false    |
| `IS_DEVELOPMENT` | Whether the service is running in development mode. | true            | false    |
| `MONGODB_URI`    | The URI of the MongoDB instance to connect to.      |                 | true     |

## Getting Started

To run the service, you can use the following command:

```bash
go run ./cmd/app/main.go
```

Note: postman collection in the root directory of the project.
