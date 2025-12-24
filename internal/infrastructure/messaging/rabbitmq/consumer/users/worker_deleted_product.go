package users

import (
    "encoding/json"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"
    eventsUser "github.com/leonardo849/product_supermarket/internal/domain/events/user"
    userApplication "github.com/leonardo849/product_supermarket/internal/application/user"
    // commom "github.com/leonardo849/product_supermarket/internal/application/common"
)

type UserDeletedProductConsumer struct {
	channel   *amqp.Channel
    queueName string
    useCase   *userApplication.DeleteUserUseCase
    exchange string
}

func NewDeletedUserProductConsumer(
	ch *amqp.Channel,
    queue string,
    useCase *userApplication.DeleteUserUseCase,
)*UserDeletedProductConsumer {
	return  &UserDeletedProductConsumer{
		channel: ch,
		queueName: queue,
		useCase: useCase,
		exchange: "auth_topic",
	}
}

func (c *UserDeletedProductConsumer) createExchange() {
	err := c.channel.ExchangeDeclare(
		c.exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *UserDeletedProductConsumer) createQueue() {
	_, err := c.channel.QueueDeclare(
		c.queueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl": int32(60000),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *UserDeletedProductConsumer) bindQueue() {
	c.channel.QueueBind(c.queueName, "user.auth.deleted", c.exchange, false, nil)
}

func (c *UserDeletedProductConsumer) Start() error {
	c.createExchange()
	c.createQueue()
	c.bindQueue()
	msgs, err := c.channel.Consume(
        c.queueName,
        "",
        false, 
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

	go func() {
		for msg := range msgs {
			var event eventsUser.UserDeleted
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Println(err.Error())
                msg.Nack(false, false) 
                continue
			}
			log.Print(event)

			body := event.ID

			err := c.useCase.Execute(body)
			if err != nil {
				log.Print(err.Error())
				msg.Nack(false, true)
			} 

			msg.Ack(false)
		}
	}()
	return  nil
}