package pubsub

type Subscriber interface {
	Subscribe(topic string) <-chan interface{}
	// Unsubscribe(ch chan interface{})
}

type Publisher interface {
	Publish(topic string, message interface{})
	AsSubscriber() Subscriber
	Shutdown()
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

func (p *pubSub) Shutdown() {
	for topic := range p.register.getTopics() {
		p.register.removeTopic(topic)
	}
}

func (p *pubSub) Subscribe(topic string) <-chan interface{} {
	ch := make(chan interface{}, 3)
	p.register.addChannel(topic, ch)
	return ch
}
