### Minimal Go HTTP Server for redirecting url aliases to full urls

This code is recommended for use in a Docker container with environment variables set.
This is a minimal Go microservice that can be used to redirect url aliases to full urls using a Postgres database and Redis cache.

### Environment Variables

| Variable  | Description                | Default                                                   |
| --------- | -------------------------- | --------------------------------------------------------- |
| DB_URL    | Postgres connection string | postgres://username:password@localhost:5432/database_name |
| RDB_URL   | Redis connection string    | redis://localhost:6379                                    |
| REDIS_TTL | Redis cache TTL in seconds | 300                                                       |

### Usage

1. Build the Docker image: `docker build -t zipit-redirect .`
2. Copy `.env.example` and set the environment variables for your docker container.
3. Run the Docker image: `docker run -p 8080:8080 zipit-redirect`
4. Access the service at `http://localhost:8080`
