# pubstore
A publication store which can be associated with an lcp-server

**Note: This project is for demonstration purposes only and should not be used in a production environment.**

## Quick Start

1. Clone the repository:

```shell
   git clone https://github.com/edrlab/pubstore.git
```

2. Compile
```shell
make build
```

or 

```shell
GOPATH=$PWD/build go install cmd/pubstore/pubstore.go
```

3. run

```shell
make run
```

or 

```shell
./build/bin/pubstore
```

4. Run the PubStore server:

Access the PubStore API at http://localhost:8080 (or the appropriate base URL according to your configuration).



## Deployment

it's currently deployed on Google Cloud Platform Cloud Run and Cloud SQL (postgresql:14)

https://pubstore.edrlab.org

### ENV Variable

The following environment variables are used in the Dockerfile:

- `BASE_URL`: The base URL for the pubstore application. Default value: `http://localhost:8080`.
- `PORT`: The port on which the HTTP server will listen. Default value: `8080`.
- `LCP_SERVER_URL`: The URL of the LCP server. Default value: `https://front-prod.edrlab.org/lcpserver`.
- `LCP_SERVER_USERNAME`: The username for the LCP server. Default value: `adm_username`.
- `LCP_SERVER_PASSWORD`: The password for the LCP server. Default value: `adm_password`.
- `DSN`: The database connection string. Default value: `""` (empty string).

You can modify these environment variables according to your requirements. Make sure to set the appropriate values based on your deployment environment.

Note: The environment variables are used during the Docker image build process and are set as defaults in the resulting image. You can override these defaults by providing custom values when running the container.

**Note: Be careful when using the LCP_SERVER_URL environment variable. The default value points to https://front-prod.edrlab.org, which is intended for test purposes only. Make sure to use a secure and production-ready LCP server URL in a real production environment.**

### Docker 

```
docker run --name my-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres

docker build -t pubstore .

docker run -p 8080:8080 -e DSN="host=host.docker.internal user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable" pubstore
```


## API

### Swagger

Swagger documentation : https://pubstore.edrlab.org/api/v1/swagger/index.html

to compile the swagger documentation

you need to install https://github.com/swaggo/swag


```shell
make docs
make build
make run
```

### User API

The User API provides endpoints to manage user data.

#### Create User

Create a new user.

- **URL:** `/user`
- **Method:** `POST`
- **Request Body:** User object
- **Response:** Created User object
- **Error Responses:**
  - `400 Bad Request`: Invalid request payload or validation errors
  - `500 Internal Server Error`: Failed to create user

#### Get User

Retrieve a user by ID.

- **URL:** `/user/{id}`
- **Method:** `GET`
- **Response:** User object
- **Error Responses:**
  - `500 Internal Server Error`: Internal server error

#### Update User

Update a user by ID.

- **URL:** `/user/{id}`
- **Method:** `PUT`
- **Request Body:** Updated User object
- **Response:** Updated User object
- **Error Responses:**
  - `400 Bad Request`: Invalid request payload or validation errors
  - `500 Internal Server Error`: Failed to update user

#### Delete User

Delete a user by ID.

- **URL:** `/user/{id}`
- **Method:** `DELETE`
- **Response:** Success message
- **Error Responses:**
  - `500 Internal Server Error`: Failed to delete user

### Publication API

The Publication API provides endpoints to manage publication data.

#### Create Publication

Create a new publication.

- **URL:** `/publications`
- **Method:** `POST`
- **Request Body:** Publication object
- **Response:** Created Publication object
- **Error Responses:**
  - `400 Bad Request`: Invalid request payload or validation errors
  - `500 Internal Server Error`: Failed to create publication

#### Get Publication

Retrieve a publication by ID.

- **URL:** `/publications/{id}`
- **Method:** `GET`
- **Response:** Publication object
- **Error Responses:**
  - `500 Internal Server Error`: Internal server error

#### Update Publication

Update a publication by ID.

- **URL:** `/publications/{id}`
- **Method:** `PUT`
- **Request Body:** Updated Publication object
- **Response:** Updated Publication object
- **Error Responses:**
  - `400 Bad Request`: Invalid request payload or validation errors
  - `500 Internal Server Error`: Failed to update publication

#### Delete Publication

Delete a publication by ID.

- **URL:** `/publications/{id}`
- **Method:** `DELETE`
- **Response:** Success message
- **Error Responses:**
  - `500 Internal Server Error`: Failed to delete publication
