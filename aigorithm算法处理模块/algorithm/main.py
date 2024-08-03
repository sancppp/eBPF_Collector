from prometheus_client import start_http_server, Counter
import time
from kafka import KafkaProducer, KafkaConsumer
import json
import yaml
import requests
from syscallid import syscall_name

# 定义 Counter
event_counter = Counter(
    'syscall_event_counter',
    'Counter for syscall events',
    ['name', 'syscall']
)
config_file = '/home/tzx/rqyc/aigorithm算法处理模块/config.yml'

def read_syscall_anomalies_from_file(filename):
    with open(filename, 'r') as file:
        anomalies = [int(line.strip()) for line in file]
    return anomalies

def fetch_events_from_kafka():
    consumer = KafkaConsumer(
        'event',
        bootstrap_servers=['localhost:19092'],
        value_deserializer=lambda v: json.loads(v.decode('utf-8'))
    )
    for message in consumer:
        events = message.value
        yield events

def update_metrics(event):
    # 只处理满足条件 "Type" == "Syscall_event" 且 "Flag" == 2 的事件
    if event.get('Type') == 'Syscall_event' and event.get('Flag') == 2:
        labels = {
            'name': event['ContainerName'],
            'syscall': syscall_name(event['Syscall'])
        }
        ret_value = event['Ret']
        event_counter.labels(**labels).inc(ret_value)

def kafka_produce(message):
    producer = KafkaProducer(
        bootstrap_servers=['localhost:19092'],
        value_serializer=lambda v: json.dumps(v).encode('utf-8')
    )

    # 发送消息到 Kafka 的 alarm 主题
    producer.send('alarm', value=message)
    producer.flush()

def load_config(config_file):
    with open(config_file, 'r') as file:
        config = yaml.safe_load(file)
    return config

def save_config(config, config_file):
    with open(config_file, 'w') as file:
        yaml.safe_dump(config, file)

def check_container_network_event(record, config):
    daddr = '.'.join(map(str, record.get('Daddr')))
    saddr = '.'.join(map(str, record.get('Saddr')))
    for rule in config.get('rules', {}).get('deny', []):
        if (daddr == rule['from'] and saddr == rule['to']) or (daddr == rule['to'] and saddr == rule['from']):
            return True
    return False

def get_container_name_by_ip(ip, config):
    for container_name, container_info in config.get('containers', {}).items():
        container_ip = container_info['ip'].split('/')[0]  # 去除子网掩码部分
        if ip == container_ip:
            return container_name
    return None

def handler(event):
    config = load_config(config_file)
    anomaly_syscalls = config.get('syscall_id', [])
    high_risk_paths = config.get('high_risk_paths', [])
    found_anomalies = False
    if event["Type"] == "Syscall_event":
        if event["Flag"] == 2 and event["Syscall"] in anomaly_syscalls:
            found_anomalies = True
            kafka_produce(f'容器{event["ContainerName"]}(cid: {event["Cid"]})中的进程{event["Comm"]}(pid: {event["Pid"]})执行异常系统调用{syscall_name(event["Syscall"])}')
        elif event["Flag"] == 0:
            found_anomalies = True
            kafka_produce(f'容器{event["ContainerName"]}(cid: {event["Cid"]})中的进程{event["Comm"]}(pid: {event["Pid"]})执行高危系统调用{syscall_name(event["Syscall"])},使用了{event["Info"]}')
    elif event["Type"] == "Fileopen_event":
        filepath = event["Filename"]
        if any(high_risk_path in filepath for high_risk_path in high_risk_paths):
            found_anomalies = True
            print("发现异常文件访问: ", event)
            kafka_produce(f'容器{event["ContainerName"]}(cid: {event["Cid"]})中的进程{event["Comm"]}(pid: {event["Pid"]})打开了文件{filepath}')
    elif event["Type"] == "Network_event":
        if event["Flag"] == 1:
            return
        saddr = '.'.join(map(str, event.get('Saddr')))
        daddr = '.'.join(map(str, event.get('Daddr')))
        if event["Pid"] > 30 and check_container_network_event(event, config):
            print("发现异常网络事件:", event)
            found_anomalies = True
            kafka_produce(f'容器{get_container_name_by_ip(saddr, config)}(cid: {event["Cid"]})中的进程{event["Comm"]}(pid: {event["Pid"]})与容器{get_container_name_by_ip(daddr, config)}发生了未经授权的网络事件')
    if not found_anomalies:
        pass

def load_container_info():
    url = 'http://192.168.0.202:8888/containerinfo'
    response = requests.get(url)
    if response.status_code == 200:
        container_info = response.json()
        config = load_config(config_file)
        config['containers'] = {container['Name']: {'cid': container['CID'], 'ip': container['IP']} for container in container_info}
        save_config(config, config_file)
        print("容器信息已更新到配置文件")
    else:
        print(f"无法获取容器信息，状态码: {response.status_code}")

def main():
    # 启动 Prometheus HTTP 服务器，监听 2024 端口
    start_http_server(2024)
    print("Prometheus HTTP server started on port 2024")

    # 从8888端口获取容器信息，更新配置文件
    load_container_info()

    for event in fetch_events_from_kafka():
        update_metrics(event)
        handler(event)

if __name__ == "__main__":
    main()
