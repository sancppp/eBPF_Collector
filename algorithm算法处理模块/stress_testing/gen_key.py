import random
import string

# 生成随机字符串
def random_string(length=100):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

# 生成指定数量的随机键
def generate_keys(num_keys=1000, length=100):
    return [random_string(length) for _ in range(num_keys)]

# 将键写入 keys.txt 文件
def write_keys_to_file(keys, filename='keys.txt'):
    with open(filename, 'w') as f:
        for key in keys:
            f.write(key + '\n')

def main():
    num_keys = 1000  # 要生成的键的数量
    keys = generate_keys(num_keys)
    write_keys_to_file(keys)
    print(f"{num_keys} keys have been written to keys.txt")

if __name__ == '__main__':
    main()
