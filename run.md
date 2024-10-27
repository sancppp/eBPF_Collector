
## demo部署：
anolis
第一次运行时mysql没初始化好，可能会panic，重启一下即可
```shell
cd demo被观测系统
docker-compose up
```

## prometheus+grafana+kafka部署：

ubuntu

```shell
cd algorithm
mkdir prometheus_data
sudo chmod 777 prometheus_data
mkdir grafana_data
sudo chmod 777 grafana_data
修改docker-compose.yml文件：ip地址改为本机ip地址
docker-compose up # 输出日志，ctrl+c停止全部服务
```

## ebpf模块部署：

anolis

```shell
cd ebpfXxx
go generate ./...
go build
sudo ./ebpf_exporter

./cadivisor -port 8000
```

## Algorithm模块部署：

ubuntu

python虚拟环境：
```shell
cd algorithm算法处理模块
python3 -m venv rqyc_venv
source rqyc_venv/bin/activate
cd algorithm
pip install prometheus_client pyyaml requests confluent_kafka
python3 ./main.py
```

## 配置修改

ubuntu

```shell
source rqyc_venv/bin/activate
cd change_config
pip install Flask pyyaml requests
python3 ./change_config.py
```

## 展示前端

ubuntu

```shell
cd rqyc-front
npm install
npm run serve
```

## websocket模块部署：

ubuntu

```shell
