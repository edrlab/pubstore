# Use the official Go image as the base image
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code to the container
COPY . .

# Build the pubstore executable
RUN CGO_ENABLED=0 go build -o pubstore ./cmd/pubstore/pubstore.go

# Use a minimal base image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the pubstore executable from the build stage
COPY --from=build /app/pubstore .

# Set the environment variables
ENV BASE_URL="http://localhost:8080"
ENV PORT=8080
ENV LCP_SERVER_URL="https://front-prod.edrlab.org/lcpserver"
ENV LCP_SERVER_USERNAME="adm_username"
ENV LCP_SERVER_PASSWORD="adm_password"

# Expose the port on which the HTTP server will listen
EXPOSE 8080

# Run the pubstore when the container starts
CMD ["./pubstore"]
