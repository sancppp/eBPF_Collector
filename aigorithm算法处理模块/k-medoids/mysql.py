import requests
import pymysql

# 从 HTTP 请求中获取数据
url = 'http://192.168.0.202:8089/events'
response = requests.get(url)

if response.status_code == 200:
    data = response.json()
    # print(data)  # 添加这一行来查看数据格式
else:
    print(f"Failed to fetch data. HTTP status code: {response.status_code}")
    exit()

# 连接到 MySQL 数据库
try:
    connection = pymysql.connect(
        host='192.168.0.249',
        database='rqyc',
        user='root',
        password='123456',
        charset='utf8'
    )

    cursor = connection.cursor()

    # 创建表（如果尚不存在）
    create_table_query = '''
    CREATE TABLE IF NOT EXISTS Anomalylog (
        id INT AUTO_INCREMENT PRIMARY KEY,
        InsertionTime DATETIME,
        Info VARCHAR(255)
    )
    '''
    cursor.execute(create_table_query)

    # 插入数据到表中
    insert_query = '''
        INSERT INTO Systemcall (Type, Timestamp, Flag, Pid, Comm, Syscall, Ret, Cid, Info)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
        '''
    for record in data:
        if record.get('Syscall') == 59:
            info_value = 2
        else:
            info_value = 1

        cursor.execute(insert_query, (
            record.get('Type'),
            record.get('Timestamp'),
            record.get('Flag'),
            record.get('Pid'),
            record.get('Comm'),
            record.get('Syscall'),
            record.get('Ret'),
            record.get('Cid'),
            info_value
        ))

    # 提交事务
    connection.commit()

    print("Data has been inserted successfully.")

except pymysql.MySQLError as e:
    print(f"Error: {e}")
finally:
    if connection:
        cursor.close()
        connection.close()
        print("MySQL connection is closed.")
