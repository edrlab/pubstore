# Use the official Go image as the base image
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code to the container
COPY . .

# Build the pubstore executable
RUN CGO_ENABLED=0 go build -o pubstore  -tags PGSQL ./cmd/pubstore/pubstore.go

# Use a minimal base image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the pubstore executable from the build stage
COPY --from=build /app/pubstore .
COPY views /app/views
COPY static /app/static

# Set the environment variables
ENV PUBSTORE_PORT=8080
ENV PUBSTORE_PUBLIC_BASE_URL="http://localhost:8080"
ENV PUBSTORE_DSN=""
ENV PUBSTORE_OAUTH_SEED="oauth-seed"
ENV PUBSTORE_RESOURCES=""
ENV PUBSTORE_PAGE_SIZE=""
ENV PUBSTORE_PRINT_LIMIT="20"
ENV PUBSTORE_COPY_LIMIT="2000"
ENV PUBSTORE_USERNAME="adm_username"
ENV PUBSTORE_PASSWORD="adm_password"
ENV PUBSTORE_LCPSERVER_URL="https://front-prod.edrlab.org/lcpserver"
ENV PUBSTORE_LCPSERVER_VERSION="v1"
ENV PUBSTORE_LCPSERVER_PROVIDER="https://edrlab.org"
ENV PUBSTORE_LCPSERVER_USERNAME="adm_username"
ENV PUBSTORE_LCPSERVER_PASSWORD="adm_password"

# Expose the port on which the HTTP server will listen
EXPOSE $PORT

# Run the pubstore when the container starts
CMD ["./pubstore"]
