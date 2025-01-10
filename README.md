# GO DNS RESOLVER

### build instruction
first create a .env file with :
```env
SERVER_PORT=3000
DATABASE_URL=postgres://go_dns_user:go_dns_password@postgres:5432/go_dns_db?sslmode=disable

POSTGRES_USER=go_dns_user
POSTGRES_PASSWORD=go_dns_password
POSTGRES_DB=go_dns_db
```
then
```bash
docker compose build
```
### start web api

```bash
docker compose up -d
```

### usage 
```bash
curl http://localhost:3000/resolve?domain=snapp.ir
```
### response exanple
```json
{"domain":"snapp.ir","ip":"185.143.xxx.xxx","dns_provider":"xxxxx.alidns.com.","cached":false}
```

### health check
```bash
curl http://localhost:3000/healthy
```

### reaciness
```bash
curl http://localhost:3000/ready
```

### metrics
```bash
curl http://localhost:3000/metrics
```
# Helm Chart

```bash
helm install web-app ./helm-chart --namespace web-app --create-namespace
```

or
```bash
./deply.sh
```

## Run app locally for development
```bash
./run.dev.sh
```

## run tests
```bash
./run.test.sh
```