package pubsub

import "sync"

// Register channels to topics
type Register interface {
	addChannel(topic string, ch chan interface{})
	sendMessage(topic string, message interface{})
	removeTopic(topic string)
	removeChannel(ch chan interface{})
	getTopics() map[string]map[chan interface{}]bool
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
	// Added Waitgroup to minimize the probability of wrong order of
	// mesages if sendMessage is called multiple times quickly
	var wg sync.WaitGroup
	for ch := range reg.topics[topic] {
		wg.Add(1)
		go func(ch chan interface{}, wg *sync.WaitGroup) {
			wg.Done()
			ch <- message
		}(ch, &wg)
	}
	wg.Wait()
}
func (reg *registry) removeTopic(topic string) {
	for ch := range reg.topics[topic] {
		reg.remove(topic, ch)
	}
}

func (reg *registry) removeChannel(ch chan interface{}) {
	for topic := range reg.channels[ch] {
		reg.remove(topic, ch)
	}
}

func (reg *registry) remove(topic string, ch chan interface{}) {
	if _, ok := reg.topics[topic]; !ok {
		return
	}

	if _, ok := reg.topics[topic][ch]; !ok {
		return
	}

	delete(reg.topics[topic], ch)
	delete(reg.channels[ch], topic)

	if len(reg.topics[topic]) == 0 {
		delete(reg.topics, topic)
	}

	if len(reg.channels[ch]) == 0 {
		close(ch)
		delete(reg.channels, ch)
	}
}

func (reg *registry) getTopics() map[string]map[chan interface{}]bool {
	return reg.topics
}
