import asyncio
import websockets
from aiokafka import AIOKafkaConsumer
import logging

# 配置日志
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Kafka 配置
KAFKA_BROKER = 'localhost:19092'
KAFKA_TOPIC = 'alarm'

# 全局 WebSocket 连接列表
connected_clients = set()

class KafkaConsumerSingleton:
    _instance = None
    _lock = asyncio.Lock()

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super(KafkaConsumerSingleton, cls).__new__(cls)
            cls._instance.consumer = None
        return cls._instance

    async def start_consumer(self):
        if self.consumer is None:
            self.consumer = AIOKafkaConsumer(
                KAFKA_TOPIC,
                bootstrap_servers=KAFKA_BROKER,
                group_id='websocket-group',
                value_deserializer=lambda x: x.decode('utf-8')
            )
            await self.consumer.start()
            asyncio.create_task(self.consume_messages())

    async def consume_messages(self):
        try:
            async for message in self.consumer:
                message_value = message.value
                logger.info(f"Received message from Kafka: {message_value}")
                await broadcast_message(message_value)
        except Exception as e:
            logger.error(f"Error consuming Kafka message: {e}")
        finally:
            await self.consumer.stop()

async def broadcast_message(message_value):
    if connected_clients:
        logger.info(f"Broadcasting message to {len(connected_clients)} clients")
        await asyncio.gather(*[client.send(message_value) for client in connected_clients])

async def websocket_handler(websocket, path):
    logger.info("Client connected")
    connected_clients.add(websocket)
    try:
        async for message in websocket:
            logger.info(f"Received message from client: {message}")
            # 这里可以处理客户端发来的消息
    except websockets.exceptions.ConnectionClosed:
        logger.info("Client connection closed")
    except Exception as e:
        logger.error(f"Error in websocket connection: {e}")
    finally:
        connected_clients.remove(websocket)
        logger.info("Client disconnected")

async def main():
    # 启动 Kafka 消费者单例
    kafka_consumer = KafkaConsumerSingleton()
    await kafka_consumer.start_consumer()

    # 启动 WebSocket 服务器
    websocket_server = await websockets.serve(websocket_handler, "192.168.0.249", 8090)
    logger.info("WebSocket server is running on ws://192.168.0.249:8090")

    await websocket_server.wait_closed()  # 保持 WebSocket 服务器运行

if __name__ == "__main__":
    asyncio.run(main())
