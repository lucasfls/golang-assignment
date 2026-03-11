# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

## Project Structure
These are the current implemented layer of the project:

- domain: core business entities
- application: use cases
- ports: HTTP handlers
- infrastructure: persistence/repository

This separation helps keep business logic independent
from delivery mechanisms.

```
cmd/
├── server/     # API server
└── seed/       # Database seeding

internal/
├── domain/           # Business entities
├── application/      # Service logic
├── infrastructure/   # Database repositories
└── ports/           # HTTP handlers

sql/              # Database migrations
```


## Setup Code Repository

1. Create a github/bitbucket/gitlab repository and push all this code as-is.
2. Create a new branch, and provide a pull-request against the main branch with your changes. Instructions to follow.

## Application Setup

- Ensure you have Go installed on your machine.
- Ensure you have Docker installed on your machine.
- Important makefile targets:
  - `make tidy`: Install all dependencies
  - `make docker-up`: Start infrastructure services via docker
  - `make seed`: ⚠️ Destroy and re-create database tables
  - `make run`: Start the application
  - `make test`: Run tests with coverage
  - `make validate`: Run comprehensive validation (format, vet, staticcheck, tests, deps, security)
  - `make docker-down`: Stop docker containers

Follow up for the assignemnt here: [ASSIGNMENT.md](ASSIGNMENT.md)

## Testing

The test suite focuses on:

- business rules in use cases
- HTTP handler behavior
- pagination edge cases
- invalid user inputs


## Further improvements

This API could be improved with:

- Structured logging (e.g. zerolog)
- Metrics and observability (e.g. Prometheus)
- Health check endpoints for monitoring and CI/CD pipelines
- End-to-end tests (e.g. Ginkgo)