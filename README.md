## mgd-check

某不可描述丁自动打卡

运行后访问 `host:port/register` 录入信息即可。每天 7 点上班，17 点下班

### TODO

- [ ] 邮件提醒
- [ ] 数据持久化
- [ ] 失败重试
- [x] Web录入界面

### 运行
使用 `Dockerfile` 构建Docker镜像运行即可
```
$ docker build -t mgd-check .
$ docker run -d -p 8090:8080 mgd-check
```