package common


type EventPublisher interface {
    Publish(event any) error
}