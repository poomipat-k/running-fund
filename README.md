# Running Fund

A Go-based application designed to manage and track running events and their associated funds.

## 🚀 Features

- Event Management: Create, update, and delete running events.

- Fund Tracking: Monitor funds raised for each event.

- User Authentication: Secure login and registration system.

- API Endpoints: RESTful API for integration with other services.

## 🛠️ Technologies Used

- Backend: Go

- Database: PostgreSQL

- Containerization: Docker

- Web Server: Nginx

- Development Tools:
  - Docker Compose for local development
  - Makefile for task automation

## ⚙️ Setup Instructions

## Prerequisites

Ensure you have the following installed:

[Go](https://go.dev/doc/install)

[Docker](https://docs.docker.com/engine/install/)

[Docker Compose](https://docs.docker.com/compose/install/)

[Make](https://makefiletutorial.com/)

## Local Development

Clone the repository:

```
git clone https://github.com/poomipat-k/running-fund.git
cd running-fund
```

Copy the example environment variables:

```
cp .env.example .env
```

Build and start the application using Docker Compose:

```
make build
make up
```

Access the application at `http://localhost:8080`

Running Tests
To run the tests:

```
make test
```

## 📦 Docker Images

- running-fund: The main application image.

- running-fund-dev: Development environment image with debugging tools.

## 📄 API Documentation

API endpoints are documented within the codebase. Refer to the Go files in the main.go and related packages for detailed information.

## 🧪 Testing

Unit and integration tests are located in the pkg directory. To run them:

```
make test
```
