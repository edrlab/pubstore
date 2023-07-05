# pubstore
A publication store which can be associated with an lcp-server


### Docker 

```

docker run --name my-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres

docker build -t pubstore .

docker run -p 8080:8080 -e DSN="host=host.docker.internal user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable" pubstore

```
