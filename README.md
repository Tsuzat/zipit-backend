## Zipit API Server

Zipit API Server is a RESTful API server built using [Fiber](https://github.com/gofiber/fiber) framework.
It's a example URL shortener API built using Fiber.

## Prerequisites

- Go 1.23.3
- Docker

## Installation

1. Clone the repository:

```bash
git clone -b go-fiber https://github.com/Tsuzat/zipit-backend.git zipit-go-fiber
```

2. Change directory to the project:

```bash
cd zipit-go-fiber
```

3. Build the Docker image:

```bash
docker build -t zipit-api-server .
```

4. Run the Docker container:

```bash
docker run -p 8080:8080 zipit-api-server
```

## Build from source

```sh
go build -ldflags="-s -w" -o apiserver .
```

and provide `.env` file in same directory (use `.env.example`) for reference.
