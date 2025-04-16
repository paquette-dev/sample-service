# Sample Service

## Overview

A simple REST API service built with Go and Echo framework that serves as a backend for the sample-client Angular application.

## Features

- RESTful API endpoints
- Built with Echo framework for high performance
- Database integration with SQLite
- User management endpoints
- Error handling and validation
- Testing with Ginkgo and Gomega

## Prerequisites

- Go 1.22.0 or higher

## Installation

Clone the repository:

```bash
git clone https://github.com/paquette-dev/sample-service.git
```

Navigate to the project directory:

```bash
cd sample-service
```

Install dependencies:

```bash
go mod download
```

## Running the service

Start the server:

```bash
go run cmd/server/main.go
```

## Testing

Run the tests:

```bash
go test -v ./...
```
