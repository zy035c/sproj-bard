import amqp
import kombu

host = 'localhost'


with amqp.Connection(host) as c:
    ch = c.channel()
    ch.basic_publish(amqp.Message('Hello World'), routing_key='test')
