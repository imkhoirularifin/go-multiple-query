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

## Documentation

base path: /api/v1

Query Params

| Name        | Required | Default | Description     |
| ----------- | -------- | ------- | --------------- |
| `page`      | false    | 1       |                 |
| `limit`     | false    | 10      |                 |
| `orderBy`   | false    | -       |                 |
| `sortOrder` | false    | asc     | `asc` or `desc` |
| `$key`      | false    | -       | filter key      |
| `$value`    | false    | -       | filter value    |

Note:

`$key` should be in camelCase

## Filter Criteria

- `equal`
- `notEqual`
- `greaterThanOrEqual`
- `lessThanOrEqual`
- `greaterThan`
- `lessThan`
- `in` (Support multiple values separated by comma)

## Example

```bash
curl -X GET "http://localhost:8080/api/v1/vouchers?page=1&limit=10&orderBy=skuName&sortOrder=asc&sku.equal=IDMR100"
```
