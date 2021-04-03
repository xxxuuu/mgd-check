## mgd-check

某不可描述丁自动打卡

运行后访问 `host:port/register` 录入信息即可。每天 7 点上班，17 点下班

### TODO

- [x] 邮件提醒
- [x] 数据持久化
- [ ] 基于代理池打卡
- [ ] 失败重试
- [x] Web录入界面

### 运行
使用 `Dockerfile` 构建Docker镜像
```
$ docker build -t mgd-check .
```
数据存储在内存，直接运行：
```
$ docker run --name mgd-check -d -p 8090:8080 mgd-check
```
数据存储在MongoDB，并开启邮件提醒：
```
$ docker network create -d bridge mgd
$ docker run -d --name mgd-mongo --network mgd mongo
$ docker run -d --name mgd-check -p 8090:8080 --network mgd \
    -e DB_TYPE=mongodb \
    -e DB_HOST=mgd-mongo:27017 \
    -e EMAIL_ENABLE=true \
    -e EMAIL_HOST=smtp.qq.com \
    -e EMAIL_PORT=587 \
    -e EMAIL_USERNAME=邮箱账号 \
    -e EMAIL_PASSWORD=邮箱密码 \
    mgd-check
```
