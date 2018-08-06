package pubsub

// Subscriber allows Subscribe/Unsubscribe to topics
type Subscriber interface {
	Subscribe(topic string) chan interface{}
	Unsubscribe(ch chan interface{})
}

// Publisher emits message to topics
type Publisher interface {
	Publish(topic string, message interface{})
	AsSubscriber() Subscriber
	Shutdown()
}

type pubSub struct {
	register Register
}

// New returns a publisher of topics
func New() Publisher {
	register := &registry{
		topics:   make(map[string]map[chan interface{}]bool),
		channels: make(map[chan interface{}]map[string]bool),
	}
	return &pubSub{register}
}

// Publish a message of a topic
func (p *pubSub) Publish(topic string, message interface{}) {
	p.register.sendMessage(topic, message)
}

// AsSubscriber returns a subscriber
func (p *pubSub) AsSubscriber() Subscriber {
	return p
}

// Shutdown removes all subscriptions and topics
func (p *pubSub) Shutdown() {
	for topic := range p.register.getTopics() {
		p.register.removeTopic(topic)
	}
}

// Subscribe to a topic, returns a channel to receive the messages of the topic
func (p *pubSub) Subscribe(topic string) chan interface{} {
	ch := make(chan interface{})
	p.register.addChannel(topic, ch)
	return ch
}

// Unsubscribe from a topic
func (p *pubSub) Unsubscribe(ch chan interface{}) {
	p.register.removeChannel(ch)
}
