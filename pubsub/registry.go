package pubsub

import "fmt"

// Register channels to topics
type Register interface {
	addChannel(topic string, ch chan interface{})
	sendMessage(topic string, message interface{})
	removeTopic(topic string)
	removeChannel(ch chan interface{})
}

type registry struct {
	topics   map[string]map[chan interface{}]bool
	channels map[chan interface{}]map[string]bool
}

func (reg *registry) addChannel(topic string, ch chan interface{}) {
	if reg.topics[topic] == nil {
		reg.topics[topic] = make(map[chan interface{}]bool)
	}
	reg.topics[topic][ch] = true

	if reg.channels[ch] == nil {
		reg.channels[ch] = make(map[string]bool)
	}
	reg.channels[ch][topic] = true
}
func (reg *registry) sendMessage(topic string, message interface{}) {
	fmt.Printf("--->A0 %v", message)
	for ch := range reg.topics[topic] {
		fmt.Printf("--->A1 ")
		ch <- message
	}
	fmt.Printf("--->A3 ")
}
func (reg *registry) removeTopic(topic string) {
	delete(reg.topics, topic)
	for ch := range reg.channels {
		delete(reg.channels[ch], topic)
	}
}
func (reg *registry) removeChannel(ch chan interface{}) {
	delete(reg.channels, ch)
	for topic := range reg.topics {
		delete(reg.topics[topic], ch)
	}
}
