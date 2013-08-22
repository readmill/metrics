package metrics

import (
	"errors"
	"fmt"
)

var Backends = map[string]Interface{}
var current Interface = nil

type Event struct {
	State      string
	Host       string
	Service    string
	Metric     int64
	Ttl        float32
	Tags       []string
	Attributes map[string]interface{}
}

func (e *Event) Set(k string, v interface{}) {
	if e.Attributes == nil {
		e.Attributes = map[string]interface{}{}
	}
	e.Attributes[k] = v
}

type Interface interface {
	Publish(e ...*Event) error
}

func Register(name string, m Interface) {
	Backends[name] = m
}

func Use(name string) error {
	if _, ok := Backends[name]; !ok {
		return errors.New("backend not found")
	}
	current = Backends[name]
	return nil
}

func Publish(e *Event) error {
	if current != nil {
		err := current.Publish(e)
		if err != nil {
			return fmt.Errorf("error publishing metric '%s': %s", e.Service, err)
		}
	}
	return nil
}
