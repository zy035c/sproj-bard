import amqp

host = 'localhost'
with amqp.Connection(host) as c:
    ch = c.channel()

    def on_message(message):
        print('Received message (delivery tag: {}): {}'.format(message.delivery_tag, message.body))
    ch.basic_consume(queue='test', callback=on_message, no_ack=True)
    while True:
        c.drain_events()
