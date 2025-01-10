# GO DNS RESOLVER

### build instruction
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
