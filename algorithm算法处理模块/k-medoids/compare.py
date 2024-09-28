import json

# 从txt文件中读取系统调用日志
def read_syscall_log(file_path):
    with open(file_path, 'r') as file:
        data = file.read().strip().split(' ')
    return data

# 读取第一个文件（正常数据）
file1_path = 'correct_1.txt'
file1_syscall_ids = read_syscall_log(file1_path)

# 读取第二个文件（异常数据）
file2_path = 'anomaly_1.txt'
file2_syscall_ids = read_syscall_log(file2_path)

# 比对两个文件，找出第一个文件中的异常数据
anomalies = set(file1_syscall_ids).intersection(set(file2_syscall_ids))

# 输出异常数据到新的txt文件
output_file_path = 'detected_anomalies.txt'
with open(output_file_path, 'w') as file:
    for anomaly in anomalies:
        file.write(f"{anomaly}\n")

# 打印异常数据
print("Detected Anomalies:")
for anomaly in anomalies:
    print(anomaly)
