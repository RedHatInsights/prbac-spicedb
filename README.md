# prbac-spicedb
Implement the PRBAC API with a spicedb backend.

# Development
## Run prbac-spicedb with spicedb (using schema in /schema)
```
docker-compose up --build
```
Test using an endpoint like:
```
curl -I localhost:8080/access/
```
## Docker
```
docker build . -t quay.io/ciam_authz/prbac-spicedb
docker run -p8080:8080 --rm quay.io/ciam_authz/prbac-spicedb
```
## Regenerate server code
`oapi-codegen -config api/server.cfg.yaml api/openapi.json`