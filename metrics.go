package metrics

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Backends holds a map of registered metrics backends.
var Backends = map[string]Interface{}

// current is the current selected backend.
var current Interface = nil

// servicePrefix is prefixed to all event service names.
var servicePrefix = ""

// Event represents an event which can be published to
// a backend.
type Event struct {
	State      string
	Host       string
	Service    string
	HttpStatus int
	Metric     int64
	Ttl        float32
	Tags       []string
	Transient  bool
	Attributes map[string]interface{}
}

// SetAttr sets an attribute with the given key and value.
// The value must be a string or a boolean.
func (e *Event) SetAttr(k string, v interface{}) {
	if e.Attributes == nil {
		e.Attributes = map[string]interface{}{}
	}
	e.Attributes[k] = v
}

// Interface implements the Publish method.
// It is the type of all metrics backends.
type Interface interface {
	Publish(e ...*Event) error
}

// Register registers a backend with the metrics library.
func Register(name string, m Interface) {
	Backends[name] = m
}

// Use selects the metrics backend to use for Publishing.
func Use(name string) error {
	if _, ok := Backends[name]; !ok {
		return errors.New("backend not found")
	}
	current = Backends[name]
	return nil
}

// Publish publishes one or more events to the current backend.
func Publish(evs ...*Event) error {
	if current != nil {
		for _, e := range evs {
			// TODO Make copy of event
			e.Service = servicePrefix + e.Service

			err := current.Publish(e)
			if err != nil {
				return fmt.Errorf("error publishing metric '%s': %s", e.Service, err)
			}
		}
	}
	return nil
}

// PublishHttpAccess publishes an HTTP access to the current backend.
func PublishHttpAccess(r *http.Request, d time.Duration, status int) error {
	return Publish(
		&Event{
			Service:   "inbound.timings",
			Tags:      []string{"http", "inbound", "percentiles"},
			Metric:    int64(d / time.Millisecond),
			Transient: true,
		},
		&Event{
			Service:    "outbound",
			HttpStatus: status,
			Tags:       []string{"http", "outbound", "rate"},
			Metric:     1,
			Transient:  true,
		},
	)
}

// SetPrefix sets the event service prefix.
// Use this to prefix all events with the name of the service
// generating the events.
func SetPrefix(pre string) {
	servicePrefix = pre
}
