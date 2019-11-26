package rabbitmq

import (
    "encoding/json"
    "github.com/streadway/amqp"
)

type RabbitMQ struct {
    _conn    *amqp.Connection
    channel  *amqp.Channel
    Name     string
    exchange string
}

// New(address)用于创建RabbitMQ结构体
func New(s string) *RabbitMQ {
    conn, err := amqp.Dial(s)
    if err != nil {
        panic(err)
    }

    ch, err := conn.Channel()
    if err != nil {
        panic(err)
    }

    q, err := ch.QueueDeclare(
        "",    // name
        false, // durable
        true,  // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        panic(err)
    }
    mq := new(RabbitMQ)
    mq.channel = ch
    mq.Name = q.Name
    return mq
}

// Bind(exchange) 将自己的消息队列和一个exchange绑定
// 所有发往该exchange的消息都能在自己的消息队列中接收到
func (q *RabbitMQ) Bind(exchange string) {
    err := q.channel.QueueBind(
        q.Name,   // queue name
        "",       // routing key
        exchange, // exchange
        false,
        nil,
    )
    if err != nil {
        panic(err)
    }
    q.exchange = exchange
}

// Send(queue,body)可以往某个消息队列发送消息
func (q *RabbitMQ) Send(queue string, body interface{}) {
    str, err := json.Marshal(body)
    if err != nil {
        panic(err)
    }
    err = q.channel.Publish(
        "",
        queue,
        false,
        false,
        amqp.Publishing{
            ReplyTo: q.Name,
            Body:    []byte(str),
        })
    if err != nil {
        panic(err)
    }
}

// Publish(exchange,body)可以往某个exchange发送消息
func (q *RabbitMQ) Publish(exchange string, body interface{}) {
    str, err := json.Marshal(body)
    if err != nil {
        panic(err)
    }
    err = q.channel.Publish(
        exchange,
        "",
        false,
        false,
        amqp.Publishing{
            ReplyTo: q.Name,
            Body:    []byte(str),
        },
    )
    if err != nil {
        panic(err)
    }
}

// Consume用于生成一个接收消息的go channel
func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
    c, err := q.channel.Consume(
        q.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        panic(err)
    }
    return c
}

// 关闭消息队列
func (q *RabbitMQ) Close() {
    q.channel.Close()
    // q._conn.Close()
}
