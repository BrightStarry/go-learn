#### 日志监控

- 监听nginx日志的实时变动,获取写入的每一行日志,经过解析后,写入influxDB数据库
- 并用grafana展示日志信息图标

- 并且,对程序运行状况进行响应监控
    - 总处理日志行数
    - 系统吞吐量
    - read channel 长度
    - write channel 长度
    - 运行总时间
    - 错误数

#### grafana
- 监控图表网站
- 安装
```
docker pull grafana/grafana

docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  -e "GF_SERVER_ROOT_URL=http://grafana.server.name" \
  -e "GF_SECURITY_ADMIN_PASSWORD=123456" \
  grafana/grafana
 
打开 ip:3000 即可进入(默认用户名:admin,密码则是GF_SECURITY_ADMIN_PASSWORD配置)
配置好influxDB的数据源
```

#### influxDB 
- 时序型数据库,对于时序型数据,性能比mysql快很多
- 时序型数据是基于时间的一些列数据
- 安装
```
docker pull influxdb
```
- [官方文档](http://docs.influxdata.com/influxdb/v1.5/introduction/getting-started/)
- [go客户端](https://github.com/influxdata/influxdb/blob/master/client/README.md)
