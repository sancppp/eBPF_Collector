from flask import Flask, render_template, request, redirect, url_for, flash
import yaml
import os

app = Flask(__name__)
app.secret_key = 'your_secret_key'

config_file_path = '/home/tzx/rqyc/aigorithm算法处理模块/config.yml'

def load_config():
    with open(config_file_path, 'r') as file:
        config = yaml.safe_load(file)
    return config

def save_config(config):
    with open(config_file_path, 'w') as file:
        yaml.safe_dump(config, file, default_flow_style=False, sort_keys=False, indent=2, width=80)

@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        config = {
            'containers': {},
            'rules': {
                'deny': []
            },
            'high_risk_paths': request.form.getlist('high_risk_paths'),
            'syscall_id': [int(x) for x in request.form.getlist('syscall_id')]
        }

        container_names = request.form.getlist('container_name')
        for container in container_names:
            ip = request.form.get(f'{container}_ip')
            config['containers'][container] = {'ip': ip}

        deny_from = request.form.getlist('deny_from')
        deny_to = request.form.getlist('deny_to')
        for i in range(len(deny_from)):
            config['rules']['deny'].append({'from': deny_from[i], 'to': deny_to[i]})

        save_config(config)
        flash('Configuration updated successfully!', 'success')
        return redirect(url_for('index'))

    config = load_config()
    return render_template('index.html', config=config)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)
