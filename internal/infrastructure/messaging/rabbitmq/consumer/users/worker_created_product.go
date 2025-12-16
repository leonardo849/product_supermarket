package consumer



import (
    "encoding/json"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"
    eventsUser "github.com/leonardo849/product_supermarket/internal/domain/events/user"
    userApplication "github.com/leonardo849/product_supermarket/internal/application/user"
)

type UserCreatedProductConsumer struct {
    channel   *amqp.Channel
    queueName string
    useCase   *userApplication.CreateUserUseCase
    exchange string
}

func NewUserCreatedProductConsumer(
    ch *amqp.Channel,
    queue string,
    useCase *userApplication.CreateUserUseCase,
) *UserCreatedProductConsumer {
    return &UserCreatedProductConsumer{
        channel:   ch,
        queueName: queue,
        useCase:   useCase,
        exchange: "auth_topic",
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
    c.channel.QueueBind(c.queueName, "user.auth.*", c.exchange, false, nil)
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
        for msg := range msgs {
            var event eventsUser.UserCreated

            if err := json.Unmarshal(msg.Body, &event); err != nil {
                log.Println(err.Error())
                msg.Nack(false, false) 
                continue
            }
            log.Print(event)

            body := userApplication.CreateUserInput{
                ID: event.ID,
                AuthUpdatedAt: event.AuthUpdatedAt,
                Role: event.Role,
            }

            _, err := c.useCase.Execute(body)

            if err != nil {
                log.Print(err.Error())
                msg.Nack(false, false) // retry
                continue
            }

            msg.Ack(false)
        }
    }()

    return nil
}
