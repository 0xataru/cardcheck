# Cardcheck API

[![Go Report Card](https://goreportcard.com/badge/github.com/markraiter/cardcheck)](https://goreportcard.com/report/github.com/markraiter/cardcheck)


This is an API for validating credit cards.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go version 1.22.0
- Docker (optional)

### Installing

A step by step series of examples that tell you how to get a development environment running.

1. Clone the repository
2. Install the dependencies with `go mod download`
3. Create `.env` file and copy values from `.env_example`
4. Follow the instructions to install [Taskfile](https://taskfile.dev/installation/) utility
5. Run the server with `task run`

## Running the tests

1. Run the tests with `task test`
2. Also you can proceed with the [OpenAPI](https://swagger.io/) docs by link `localhost:3000/swagger`

## API Examples

### Successful Request
```bash
curl -X POST http://localhost:3000/check \
  -H "Content-Type: application/json" \
  -d '{
    "card_number": "5167803252097675",
    "expiration_month": "12",
    "expiration_year": "2029"
  }'
```

Response:
```json
{
  "valid": true,
  "error": null
}
```

### Failed Request (Invalid Card Number)
```bash
curl -X POST http://localhost:3000/check \
  -H "Content-Type: application/json" \
  -d '{
    "card_number": "5167803252097676",
    "expiration_month": "12",
    "expiration_year": "2029"
  }'
```

Response:
```json
{
  "valid": false,
  "error": {
    "code": "001",
    "message": "invalid card number"
  }
}
```

## Deployment

You can also run the service in Docker container with `task run-container`

## Built With

- [Go](https://golang.org/) - The programming language used.
- [Docker](https://www.docker.com/) - Used for containerization.
