import asyncio
import aiohttp
import random
import string
import logging

# 设置日志配置
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

PUT_URL = 'http://192.168.0.202:8080/put'
GET_URL = 'http://192.168.0.202:8080/get'
REQUESTS_PER_SECOND = 100  # 每秒请求数
NUM_KEYS = 1000  # 预生成的随机键数量

# 生成随机字符串
def random_string(length=100):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

# 生成随机键列表
def generate_keys(num_keys=NUM_KEYS, length=100):
    return [random_string(length) for _ in range(num_keys)]

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

async def put_worker(session, keys):
    while True:
        key = random.choice(keys)
        await put_request(session, key)
        await asyncio.sleep(1 / REQUESTS_PER_SECOND)

async def get_worker(session, keys):
    while True:
        key = random.choice(keys)
        await get_request(session, key)
        await asyncio.sleep(1 / REQUESTS_PER_SECOND)

async def main():
    keys = generate_keys()

    async with aiohttp.ClientSession() as session:
        # 启动 PUT 和 GET 请求的 workers
        put_workers = [asyncio.create_task(put_worker(session, keys)) for _ in range(REQUESTS_PER_SECOND)]
        get_workers = [asyncio.create_task(get_worker(session, keys)) for _ in range(REQUESTS_PER_SECOND)]

        # 一直运行
        await asyncio.gather(*put_workers, *get_workers)

if __name__ == '__main__':
    asyncio.run(main())
