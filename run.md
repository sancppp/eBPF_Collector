Algorithm模块部署：

python虚拟环境：
```shell
cd algorithm模块部署
python3 -m venv rqyc_venv
source rqyc_venv/bin/activate
```

kafka+prometheus+grafana
sudo chmod 777 prometheus_data
修改docker-compose.yml文件：ip地址改为本机ip地址
docker-compose up # 输出日志，ctrl+c停止全部服务

