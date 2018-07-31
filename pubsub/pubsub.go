package pubsub

type Subscriber interface {
	Subscribe(topic string) chan interface{}
	// Unsubscribe(ch chan interface{})
}

type Publisher interface {
	Publish(topic string, message interface{})
	AsSubscriber() Subscriber
}

type pubSub struct {
	register Register
}

func New() Publisher {
	register := &registry{
		topics:   make(map[string]map[chan interface{}]bool),
		channels: make(map[chan interface{}]map[string]bool),
	}
	return &pubSub{register}
}

func (p *pubSub) Publish(topic string, message interface{}) {
	p.register.sendMessage(topic, message)
}

func (p *pubSub) AsSubscriber() Subscriber {
	return p
}

func (p *pubSub) Subscribe(topic string) chan interface{} {
	ch := make(chan interface{})
	p.register.addChannel(topic, ch)
	return ch
}
