import asyncio
import websockets
from aiokafka import AIOKafkaConsumer
from threading import Lock

# Kafka 消费者配置
KAFKA_BROKER = 'localhost:19092'
KAFKA_TOPIC = 'alarm'

# 全局WebSocket连接列表
connected_clients = set()

class KafkaConsumerSingleton:
    _instance = None
    _lock = Lock()

    def __new__(cls, loop):
        with cls._lock:
            if cls._instance is None:
                cls._instance = super(KafkaConsumerSingleton, cls).__new__(cls)
                cls._instance.loop = loop
                cls._instance.consumer = None
            return cls._instance

    async def start_consumer(self):
        if self.consumer is None:
            self.consumer = AIOKafkaConsumer(
                KAFKA_TOPIC,
                loop=self.loop,
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
                print(f"Received message: {message_value}")
                await broadcast_message(message_value)
        except Exception as e:
            print(f"Error consuming Kafka message: {e}")
        finally:
            await self.consumer.stop()

# 广播消息到所有连接的WebSocket客户端
async def broadcast_message(message_value):
    if connected_clients:  # 仅在有连接的客户端时广播
        print(f"Broadcasting message to {len(connected_clients)} clients")
        await asyncio.wait([client.send(message_value) for client in connected_clients])

# WebSocket 处理函数
async def websocket_handler(websocket, path):
    print("Client connected")
    connected_clients.add(websocket)
    try:
        async for message in websocket:
            pass  # 处理收到的消息（如果需要）
    except websockets.exceptions.ConnectionClosed:
        pass
    finally:
        print("Client disconnected")
        connected_clients.remove(websocket)

# 主函数
async def main():
    loop = asyncio.get_event_loop()

    # 启动 Kafka 消费者单例
    kafka_consumer = KafkaConsumerSingleton(loop)
    await kafka_consumer.start_consumer()

    # 启动 WebSocket 服务器
    websocket_server = await websockets.serve(websocket_handler, "localhost", 8090)
    print("WebSocket server is running on ws://localhost:8090")

    await websocket_server.wait_closed()  # 保持 WebSocket 服务器运行

# 运行主函数
if __name__ == "__main__":
    asyncio.run(main())