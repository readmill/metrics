package riemann

import (
	"flag"

	"github.com/readmill/metrics"
	"github.com/readmill/raidman"
)

type Riemann struct {
	addr   string
	proto  string
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
		err := r.client.Send(&raidman.Event{
			State:   e.State,
			Service: e.Service,
			Metric:  e.Metric,
			Ttl:     e.Ttl,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Riemann) open() (*raidman.Client, error) {
	return raidman.Dial(r.proto, r.addr)
}

func init() {
	addr := flag.String("riemann.addr", ":5555", "riemann host address")
	proto := flag.String("riemann.proto", "tcp", "riemann network protocol (tcp, udp)")
	metrics.Register("riemann", &Riemann{*addr, *proto, nil})
}
