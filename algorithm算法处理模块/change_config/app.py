from flask import Flask, render_template, request, redirect, url_for, flash
import yaml
import requests
import os

app = Flask(__name__)
app.secret_key = 'your_secret_key'

config_file_path = '/home/tzx/rqyc/aigorithm算法处理模块/config.yml'

def load_container_info():
    url = 'http://192.168.252.131:8888/containerinfo'
    response = requests.get(url)
    if response.status_code == 200:
        container_info = response.json()
        config = load_config()
        config['containers'] = {container['Name']: {'cid': container['CID'], 'ip': container['IP']} for container in container_info}
        save_config(config)
        print("容器信息已更新到配置文件")
    else:
        print(f"无法获取容器信息，状态码: {response.status_code}")

def load_config():
    if os.path.exists(config_file_path):
        with open(config_file_path, 'r') as file:
            config = yaml.safe_load(file)
            if config is None:
                config = {}
            return config
    return {}

def save_config(config):
    with open(config_file_path, 'w') as file:
        yaml.safe_dump(config, file, default_flow_style=False, sort_keys=False, indent=2, width=80)

def clear_network_blacklist():
    url = 'http://192.168.252.131:8889/cnetworkclean'
    response = requests.get(url)
    if response.status_code == 200:
        print("网络黑名单已清空")
    else:
        print(f"无法清空网络黑名单，状态码: {response.status_code}")

def update_network_config(deny_rules):
    base_url = 'http://192.168.252.131:8889/cnetworkconfig'
    for rule in deny_rules:
        from_ip = rule['from']
        to_ip = rule['to']
        url = f'{base_url}?sip={from_ip}&dip={to_ip}'
        response = requests.get(url)
        if response.status_code == 200:
            print(f"已更新网络配置：{from_ip} -> {to_ip}")
        else:
            print(f"无法更新网络配置：{from_ip} -> {to_ip}，状态码: {response.status_code}")

@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        config = load_config()
        config['containers'] = {}
        config['rules'] = {
            'deny': []
        }
        config['high_risk_paths'] = request.form.getlist('high_risk_paths')
        config['syscall_id'] = [int(x) for x in request.form.getlist('syscall_id')]

        container_names = request.form.getlist('container_name')
        for container in container_names:
            ip = request.form.get(f'{container}_ip')
            config['containers'][container] = {'ip': ip}

        deny_from = request.form.getlist('deny_from')
        deny_to = request.form.getlist('deny_to')
        new_deny_rules = []
        for i in range(len(deny_from)):
            new_deny_rules.append({'from': deny_from[i], 'to': deny_to[i]})

        # 清空网络黑名单
        clear_network_blacklist()

        # 更新网络配置
        update_network_config(new_deny_rules)

        config['rules']['deny'] = new_deny_rules

        save_config(config)
        flash('Configuration updated successfully!', 'success')
        return redirect(url_for('index'))
    
    # 从8888端口获取容器信息，更新配置文件
    load_container_info()
    config = load_config()
    return render_template('index.html', config=config)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)
