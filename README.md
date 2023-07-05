# pubstore
A publication store which can be associated with an lcp-server


### Docker 

```

docker build -t pubstore .

docker run -p 8080:8080 -e DSN="host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai" pubstore


```
