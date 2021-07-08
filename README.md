URL Checker
-----------
This microservice provides access to URL Checker API for check the status code of a number of URLs.

Run service:
```
go run main.go
```

Example of request:

```
curl -X POST -d '{"urls": ["http://localhost:10007/", "http://example.com"]}' http://localhost:10007/check_urls
```
