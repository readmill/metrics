package metrics

import (
	"log"
	"os"

	"github.com/readmill/metrics"
)

type StdLogger struct {
	*log.Logger
}

func (l *StdLogger) Publish(e *metrics.Event) error {
	l.Logger.Printf("%s %q (%q): %d", e.Service, e.Tags, e.Attributes, e.Metric)
	return nil
}

func init() {
	metrics.Register("stdout", &StdLogger{log.New(os.Stdout, "[metrics] ", log.LstdFlags)})
	metrics.Register("stderr", &StdLogger{log.New(os.Stderr, "[metrics] ", log.LstdFlags)})
}
