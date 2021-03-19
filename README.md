## mgd-check
使用 `Dockerfile` 构建Docker镜像运行即可
```
$ docker build -t mgd-check .
$ docker run -d -p 8090:8080 mgd-check
```