package metrics

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var Backends = map[string]Interface{}
var current Interface = nil
var prefix = ""

type Event struct {
	State      string
	Host       string
	Service    string
	HttpStatus int
	Metric     int64
	Ttl        float32
	Tags       []string
	Attributes map[string]string
}

func (e *Event) SetAttr(k string, v string) {
	if e.Attributes == nil {
		e.Attributes = map[string]string{}
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

func Publish(evs ...*Event) error {
	if current != nil {
		for _, e := range evs {
			// TODO Make copy of event
			e.Service = prefix + e.Service

			err := current.Publish(e)
			if err != nil {
				return fmt.Errorf("error publishing metric '%s': %s", e.Service, err)
			}
		}
	}
	return nil
}

func PublishHttpAccess(r *http.Request, d time.Duration, status int) error {
	return Publish(
		&Event{
			Service: "inbound.request.timings",
			Tags:    []string{"http", "inbound", "percentiles"},
			Metric:  int64(d / time.Millisecond),
		},
		&Event{
			Service:    "outbound",
			HttpStatus: status,
			Tags:       []string{"http", "outbound", "rate"},
			Metric:     1,
		},
	)
}

func SetPrefix(pre string) {
	prefix = pre
}
