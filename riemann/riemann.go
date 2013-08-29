package riemann

import (
	"flag"
	"strconv"

	"github.com/readmill/metrics"
	"github.com/readmill/raidman"
)

var (
	HttpStatusAttr = "status"
	PersistAttr    = "persist"
)

type Riemann struct {
	addr   *string
	net    *string
	client *raidman.Client
}

func (r *Riemann) Publish(evs ...*metrics.Event) error {
	if r.client == nil {
		client, err := r.open()
		if err != nil {
			return err
		}
		r.client = client
	}

	for _, e := range evs {
		ev := &raidman.Event{
			State:      e.State,
			Service:    e.Service,
			Metric:     e.Metric,
			Ttl:        e.Ttl,
			Tags:       e.Tags,
			Attributes: e.Attributes,
		}
		if ev.Attributes == nil {
			ev.Attributes = map[string]interface{}{}
		}
		if e.HttpStatus != 0 {
			ev.Attributes[HttpStatusAttr] = strconv.Itoa(e.HttpStatus)
		}
		if !e.Transient {
			ev.Attributes[PersistAttr] = "true"
		}
		err := r.client.Send(ev)
		if err != nil {
			if err == io.EOF {
				r.client = nil
			}
			return err
		}
	}
	return nil
}

func (r *Riemann) open() (*raidman.Client, error) {
	return raidman.Dial(*r.net, *r.addr)
}

func init() {
	addr := flag.String("riemann.addr", ":5555", "riemann host address")
	netwrk := flag.String("riemann.net", "tcp", "riemann network protocol (tcp, udp)")
	metrics.Register("riemann", &Riemann{addr, netwrk, nil})
}
