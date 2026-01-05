package users

import (
	"context"
	"encoding/json"
	"log"

	commom "github.com/leonardo849/product_supermarket/internal/application/common"
	userApplication "github.com/leonardo849/product_supermarket/internal/application/user"
	domainError "github.com/leonardo849/product_supermarket/internal/domain/error"
	eventsUser "github.com/leonardo849/product_supermarket/internal/domain/events/user"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserCreatedProductConsumer struct {
    channel   *amqp.Channel
    queueName string
    useCase   *userApplication.CreateUserUseCase
    exchange string
    publish commom.EventPublisher 
    errorCache domainError.ErrorCache
    useCaseDeleteUser *userApplication.DeleteUserUseCase
}

func NewUserCreatedProductConsumer(
    ch *amqp.Channel,
    queue string,
    useCase *userApplication.CreateUserUseCase,
    useCaseDeleteUser *userApplication.DeleteUserUseCase,
    publish commom.EventPublisher,
    errorCache domainError.ErrorCache,
) *UserCreatedProductConsumer {
    return &UserCreatedProductConsumer{
        channel:   ch,
        queueName: queue,
        useCase:   useCase,
        exchange: "auth_topic",
        publish: publish,
        useCaseDeleteUser: useCaseDeleteUser,
        errorCache: errorCache,
    }
}

func (c *UserCreatedProductConsumer) createExchange()  {
    err := c.channel.ExchangeDeclare(
        c.exchange,
        "topic",     
        true,        
        false,        
        false,        
        false,        
        nil)
    if err != nil {
        log.Fatal(err)
    }
}

func (c *UserCreatedProductConsumer) createQueue() {
    _, err := c.channel.QueueDeclare(
        c.queueName, 
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

func (c *UserCreatedProductConsumer) bindQueue() {
    c.channel.QueueBind(c.queueName, "user.auth.created_worker", c.exchange, false, nil)
}

func (c *UserCreatedProductConsumer) Start() error {
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
        defer func() {
			if r := recover(); r != nil {
				log.Println("panic recovered in consumer:", r)
			}
		}()
        for msg := range msgs {

            var event eventsUser.UserCreated
            
            if err := json.Unmarshal(msg.Body, &event); err != nil {
                log.Println(err.Error())
                msg.Nack(false, false) 
                continue
            }
            has, err := c.errorCache.HasAuthError(context.Background(), event.ID)
            if err != nil {
                log.Println(err.Error())
            }
            if has {
                if err := c.useCaseDeleteUser.Execute(event.ID); err != nil {
                    log.Println(err.Error())    
                }
                msg.Nack(false, false)
                continue
            }
            body := userApplication.CreateUserInput{
                ID: event.ID,
                AuthUpdatedAt: event.AuthUpdatedAt,
                Role: event.Role,
            }
            _, err = c.useCase.Execute(body)
            if err != nil {
                log.Print(err.Error())
                msg.Nack(false, false)  
                body := eventsUser.EmitUserCreatedError{
                    ID: body.ID,
                }
                c.publish.Publish(body)
                continue
            }

            msg.Ack(false)
            createdUser := eventsUser.EmitUserCreated{
                ID: body.ID,
            }
            c.publish.Publish(createdUser)
        }
    }()

    return nil
}
