from prometheus_client import start_http_server, Counter
import time
import json
import yaml
import requests
from confluent_kafka import Producer, Consumer
from syscallid import syscall_name

# 定义 Counter
event_counter = Counter(
    'syscall_event_counter',
    'Counter for syscall events',
    ['name', 'syscall']
)
config_file = '/home/tzx/songshanhu/algorithm算法处理模块/config.yml'

def read_syscall_anomalies_from_file(filename):
    with open(filename, 'r') as file:
        anomalies = [int(line.strip()) for line in file]
    return anomalies

def fetch_events_from_kafka(consumer):
    consumer.subscribe(['event'])
    while True:
        msg = consumer.poll(1.0)  # 等待1秒钟，获取消息
        if msg is None:
            continue
        if msg.error():
            print(f"Kafka error: {msg.error()}")
            continue
        yield json.loads(msg.value().decode('utf-8'))

def update_metrics(event):
    if event.get('Type') == 'Syscall_event' and event.get('Flag') == 2:
        labels = {
            'name': event['ContainerName'],
            'syscall': syscall_name(event['Syscall'])
        }
        ret_value = event['Ret']
        event_counter.labels(**labels).inc(ret_value)

def kafka_produce(message):
    producer = Producer({'bootstrap.servers': 'localhost:19092'})
    producer.produce('alarm', value=message)
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
        container_ip = container_info['ip'].split('/')[0]
        if ip == container_ip:
            return container_name
    return None

def handler(event):
    try:
        config = load_config(config_file)
        anomaly_syscalls = config.get('syscall_id', [])
        high_risk_paths = config.get('high_risk_paths', [])
    except Exception as e:
        return
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

def main():
    # 启动 Prometheus HTTP 服务器，监听 2024 端口
    start_http_server(2024)
    print("Prometheus HTTP server started on port 2024")

    # 初始化 Kafka Consumer
    consumer = Consumer({
        'bootstrap.servers': 'localhost:19092',
        'group.id': 'my_group',
        'auto.offset.reset': 'earliest'
    })

    for event in fetch_events_from_kafka(consumer):
        update_metrics(event)
        handler(event)

    consumer.close()

if __name__ == "__main__":
    main()
