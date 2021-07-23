<a href="https://www.pixiv.net/artworks/85627908"><img src="assets/Raika@ごっち.png" align="right" width="200px"/></a>

# ☁️ Raika ![Go](https://github.com/wuhan005/Raika/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/Raika)](https://goreportcard.com/report/github.com/wuhan005/Raika) [![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/wuhan005/Raika)

Hybrid cloud serverless function framework.

## Getting started

### Login to cloud platform

#### Aliyun

```bash
Raika platform login  --platform aliyun --region-id cn-hangzhou --account-id <REDACTED>  --access-key-id <REDACTED> --access-key-secret <REDACTED>
```

#### Tencent cloud

```bash
Raika platform login  --platform tencentcloud --region-id ap-shanghai --secret-id <REDACTED> --secret-key <REDACTED>
```

### List the cloud platform accounts

```bash
Raika platform list
```

### Deploy serverless function

```bash
Raika function create \
    --name hello_unknwon \
    --memory 128 \    # MB
    --init-timeout 10 \   # seconds
    --runtime-timeout 10 \  # seconds
    --binary-file hello_unknwon \
    --env MYENV=here_is_env_var
```

### Internal daemon

Raika provides an internal daemon service which allows you to run the serverless function periodically.

#### Start daemon service

```bash
Raika daemon start  
```

#### Create & Run periodic task

```bash
Raika daemon cron create --name helloworld --duration 5

Raika daemon cron run --name=helloworld
```

## License

MIT License
