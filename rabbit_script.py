import pika

def check_rabbitmq_connection():
    try:
        # 定义RabbitMQ连接参数
        connection_params = pika.ConnectionParameters(
            host='localhost',  # RabbitMQ服务器的主机名或IP地址
            port=5672,         # RabbitMQ的默认端口
            virtual_host='/',   # RabbitMQ虚拟主机（默认是'/'）
            credentials=pika.PlainCredentials('guest', 'guest')  # RabbitMQ默认用户名和密码
        )

        # 尝试连接RabbitMQ
        connection = pika.BlockingConnection(connection_params)
        if connection.is_open:
            print("成功连接到 RabbitMQ 服务器！")
        else:
            print("无法连接到 RabbitMQ 服务器。")

        # 关闭连接
        connection.close()

    except pika.exceptions.AMQPConnectionError as e:
        print(f"连接失败: {e}")

if __name__ == '__main__':
    check_rabbitmq_connection()