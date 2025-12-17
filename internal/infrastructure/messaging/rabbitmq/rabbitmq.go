package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	// eventsProduct "github.com/leonardo849/product_supermarket/internal/domain/events/product"
	eventsUser "github.com/leonardo849/product_supermarket/internal/domain/events/user"
	amqp "github.com/rabbitmq/amqp091-go"
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


func (p *Publisher) createExchanges() {
    switch p.exchange {
    case "product_auth_direct":
        p.channel.ExchangeDeclare(p.exchange, "direct", true, false, false, false, nil)
    }
}

func (p *Publisher) Publish(event any) error {
    var (
        routingKey string
        body       []byte
        err        error
    )
    p.createExchanges()
    switch e := event.(type) {
    // case eventsProduct.ProductCreated:
    //     routingKey = "email"
    //     body, err = json.Marshal(e)
    case eventsUser.EmitUserCreated:
        routingKey = "user.product.created"
        body, err = json.Marshal(e)
    case eventsUser.EmitUserCreatedError:
        routingKey = "user.product.created_error"
        body, err = json.Marshal(e)
    default:
        return fmt.Errorf("application doesn't know that event %T", event)
    }

    if err != nil {
        return err
    }
    log.Printf("publishing message to exchange: %s, routing key: %s", p.exchange, routingKey)
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