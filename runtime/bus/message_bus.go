//runtime/bus/message_bus.go

package runtime_bus

type MessageBus struct {
	subscribers map[string][]chan Message
}

type Message struct {
	Topic string
	Data  []byte
}

func NewMessageBus() *MessageBus {
	return &MessageBus{
		subscribers: make(map[string][]chan Message),
	}
}

func (b *MessageBus) Publish(msg Message) {
	if subs, ok := b.subscribers[msg.Topic]; ok {
		for _, ch := range subs {
			ch <- msg
		}
	}
}

func (b *MessageBus) Subscribe(topic string) chan Message {
	ch := make(chan Message, 10)
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	return ch
}
