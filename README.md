# pubstore
A publication store, to be associated with an LCP Server.

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


### Configuration

The configuration of the server is kept both in a configuration file and in environment variables. It is possible to mix both sets;  environment variables are expressly recommended for confidential information. 

The server will use the PUBSTORE_CONFIG environment variable to find a configuration file. Its value must be a file path. 

Configuration properties are expressed in snake case in the configuration file, and screaming snake case prefixed by `PUBSTORE` when expressed as environment variables. 
As an example, the `port` conguration property becomes the `PUBSTORE_PORT` environment variable, `public_base_url` becomes `PUBSTORE_PUBLIC_BASE_URL`, 
and the `version` property of the `lcp_server` section becomes `PUBSTORE_LCP_SERVER_VERSION`.

- `port`:tThe port on which the HTTP server will listen. Default value: `8080`.
- `public_base_url`: the base URL for the pubstore server. Default value: `http://localhost:8080`.
- `dsn`: the data source name, i.e. database connection string. Default value: `sqlite3://pubstore.sqlite`.
- `oauth_seed`: a string used as a seed for OAuth2 server authorization. 
- `root_dir`: the path to static files and views used by the web interface. Default value: current directory.
//- `resources`: the path to the cover images used by the Web interface.
- `page_size`: the page size used  in the REST API and Web interface.
- `print_limit`: the print limit set in LCP licenses generated from the associated LCP Server. 
- `copy_limit`: the copy limit set in LCP licenses generated from the associated LCP Server. 
- `username`: the Basic Auth username used to notify Pubstore of a new encrypted publication.
- `password`: the Basic Auth password used to notify Pubstore of a new encrypted publication.
- `lcp_server`: a section relative to the access to the associated LCP Server. 

The `lcp_server` section contains:
- `url`: the URL of the LCP server.
- `version`: the version of the LCP server: "v1" or "v2".
- `username`: the username for the LCP server.
- `password`: the password for the LCP server.


You can modify these environment variables according to your requirements. Make sure to set the appropriate values based on your deployment environment.

Note: The environment variables are used during the Docker image build process and are set as defaults in the resulting image. You can override these defaults by providing custom values when running the container.

**Note: Be careful when using the LCP_SERVER_URL environment variable. The default value points to https://front-prod.edrlab.org, which is intended for test purposes only. Make sure to use a secure and production-ready LCP server URL in a real production environment.**

### Docker 

```
docker run --name my-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres

docker build -t pubstore .

docker run -p 8080:8080 -e DSN="host=host.docker.internal user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable" pubstore
```

## Deployment

it's currently deployed on Google Cloud Platform Cloud Run and Cloud SQL (postgresql:14)

https://pubstore.edrlab.org

## API

### Swagger

Swagger documentation : https://pubstore.edrlab.org/api/swagger/index.html

to compile the swagger documentation, you need to install https://github.com/swaggo/swag


```shell
make docs
make build
make run
```
