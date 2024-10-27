import asyncio
import aiohttp
import random
import string
import logging
import argparse

# 设置日志配置
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

PUT_URL = 'http://192.168.252.131:8080/put'
GET_URL = 'http://192.168.252.131:8080/get'

# 从 keys.txt 文件中读取键
def load_keys(filename='keys.txt'):
    with open(filename, 'r') as f:
        keys = [line.strip() for line in f if line.strip()]
    return keys

# 生成随机字符串
def random_string(length=100):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

async def put_request(session, key):
    value = random_string()
    params = {'key': key, 'value': value}
    async with session.get(PUT_URL, params=params) as response:
        if response.status != 200:
            logger.error(f"PUT request failed: {response.status} - {await response.text()}")
        else:
            await response.text()

async def get_request(session, key):
    params = {'key': key}
    async with session.get(GET_URL, params=params) as response:
        if response.status != 200:
            logger.error(f"GET request failed: {response.status} - {await response.text()}")
        else:
            await response.text()

async def worker(session, keys, requests_per_second):
    while True:
        key = random.choice(keys)
        if random.choice([True, False]):
            await put_request(session, key)
        else:
            await get_request(session, key)
        await asyncio.sleep(1 / requests_per_second)

async def main(requests_per_second):
    keys = load_keys()

    async with aiohttp.ClientSession() as session:
        # 启动混合 PUT 和 GET 请求的 workers
        workers = [asyncio.create_task(worker(session, keys, requests_per_second)) for _ in range(requests_per_second)]

        # 一直运行
        await asyncio.gather(*workers)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description="HTTP Pressure Testing Tool")
    parser.add_argument('requests_per_second', type=int, help='Number of requests per second')
    args = parser.parse_args()
    asyncio.run(main(args.requests_per_second))
