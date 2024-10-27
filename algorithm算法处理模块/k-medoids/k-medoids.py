import numpy as np
import matplotlib.pyplot as plt
from sklearn.cluster import KMeans
from sklearn.preprocessing import OneHotEncoder
from sklearn.preprocessing import StandardScaler
import json
import requests


url = 'http://192.168.252.131:8089/events'
response = requests.get(url)
data = response.json()

# 从txt文件中读取系统调用日志
def read_syscall_log(file_path):
    with open(file_path, 'r') as file:
        data = file.read().strip().split(' ')
    return data

# 读取正确的数据和异常的数据
# correct_file_path = 'correct_1.txt'
# anomaly_file_path = 'anomaly_1.txt'
# correct_systemcall_ids = read_syscall_log(correct_file_path)
# anomaly_systemcall_ids = read_syscall_log(anomaly_file_path)
#
# # 将正确的数据和异常的数据合并
# systemcall_ids = correct_systemcall_ids + anomaly_systemcall_ids
#
# # 将systemcall_id转换为数值特征（One-Hot编码）
# encoder = OneHotEncoder()
# data = encoder.fit_transform(np.array(systemcall_ids).reshape(-1, 1)).toarray()
#
# # 使用K均值算法进行聚类
# kmeans = KMeans(n_clusters=3, random_state=42)
# kmeans.fit(data[:len(correct_systemcall_ids)])  # 只使用正确的数据进行聚类
# labels = kmeans.predict(data)  # 使用聚类模型预测所有数据点的标签
# centroids = kmeans.cluster_centers_
#
# # 计算每个点到其簇中心的距离
# distances = np.linalg.norm(data - centroids[labels], axis=1)
#
# # 确定异常点（设定一个距离阈值）
# threshold = np.percentile(distances[:len(correct_systemcall_ids)], 95)  # 使用正确数据的95百分位数作为阈值
# anomalies = np.array(systemcall_ids)[distances > threshold]
#
# # 输出异常点到txt文件
# with open('detected_anomalies.txt', 'w') as file:
#     for anomaly in anomalies:
#         file.write(f"{anomaly}\n")
#
# # 可视化（如果需要）
# plt.scatter(data[:, 0], data[:, 1], c='blue', marker='o', label='Normal Data')
# for anomaly in anomalies:
#     idx = systemcall_ids.index(anomaly)
#     plt.scatter(data[idx, 0], data[idx, 1], c='red', marker='x', label='Anomalies' if anomaly == anomalies[0] else "")
#
# plt.scatter(centroids[:, 0], centroids[:, 1], c='green', marker='*', s=200, label='Centroids')
# plt.title('K-Means Anomaly Detection in System Call Logs')
# plt.legend()
# plt.show()
#
#
# # 打印异常点
# print("Detected Anomalies:")
# for anomaly in anomalies:
#     print(anomaly)

with open('anomaly_1.txt', 'r') as file:
    anomaly_syscalls = {int(line.strip()) for line in file}

# 解析数据并提取特征
features = []
for item in data:
    if item['Type'] == 'Syscall_event':
        features.append([item['Syscall'], item['Flag']])

# 转换为NumPy数组
features = np.array(features)

# 标准化数据
scaler = StandardScaler()
features_scaled = scaler.fit_transform(features)

# 执行K均值聚类
kmeans = KMeans(n_clusters=2, random_state=0).fit(features_scaled)
labels = kmeans.labels_
centers = kmeans.cluster_centers_

# 计算样本与聚类中心的距离
distances = np.linalg.norm(features_scaled - centers[labels], axis=1)

# 设定阈值来检测异常值，这里假设距离超过一定阈值的样本为异常值
threshold = np.percentile(distances, 95)
potential_anomalies = [i for i in range(len(distances)) if distances[i] > threshold]

# 确定异常syscall
anomalies = []
for i in potential_anomalies:
    if features[i, 0] in anomaly_syscalls:
        anomalies.append(data[i])

# 输出异常值
print("检测到的异常值:")
for anomaly in anomalies:
    print(anomaly)