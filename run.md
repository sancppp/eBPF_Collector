
## demo部署：
anolis
第一次运行时mysql没初始化好，可能会panic，重启一下即可
```shell
cd demo被观测系统
docker-compose up
```

## ebpf模块部署：

anolis

```shell
cd ebpfXxx
go build
sudo ./ebpf_exporter
```

## Algorithm模块部署：

ubuntu

python虚拟环境：
```shell
cd algorithm
python3 -m venv rqyc_venv
source rqyc_venv/bin/activate
```

kafka+prometheus+grafana
sudo chmod 777 prometheus_data
修改docker-compose.yml文件：ip地址改为本机ip地址
docker-compose up # 输出日志，ctrl+c停止全部服务


