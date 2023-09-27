# prbac-spicedb
Implement the PRBAC API with a spicedb backend.

# Development
## Docker
```
docker build . -t quay.io/ciam_authz/prbac-spicedb
docker run -p8080:8080 --rm quay.io/ciam_authz/prbac-spicedb
```
## Regenerate server code
`oapi-codegen -config api/server.cfg.yaml api/openapi.json`