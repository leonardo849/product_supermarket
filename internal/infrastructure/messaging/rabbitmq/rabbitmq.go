package rabbitmq

import (
    "encoding/json"
    "fmt"

    amqp "github.com/rabbitmq/amqp091-go"
	eventsProduct "github.com/leonardo849/product_supermarket/internal/domain/events/product"
)


func NewConnection(rabbitURL string) (*amqp.Connection, error) {
    return amqp.Dial(rabbitURL)
}


type Publisher struct {
    channel  *amqp.Channel
    exchange string
}

func NewPublisher(ch *amqp.Channel, exchange string) *Publisher {
    return &Publisher{
        channel:  ch,
        exchange: exchange,
    }
}




func (p *Publisher) Publish(event any) error {
    var (
        routingKey string
        body       []byte
        err        error
    )

    switch e := event.(type) {
    case eventsProduct.ProductCreated:
        routingKey = "email"
        body, err = json.Marshal(e)
    default:
        return fmt.Errorf("application doesn't know that event %T", event)
    }

    if err != nil {
        return err
    }

    return p.channel.Publish(
        p.exchange,
        routingKey,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}