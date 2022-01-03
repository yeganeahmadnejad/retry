package retry

import (
	"fmt"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin/pkg/parse"
	"github.com/yeganeahmadnejad/fanout"
	"github.com/coredns/coredns/plugin/pkg/transport"

)

func initFanout(c *caddy.Controller) (*fanout.Fanout, error) {
	f := fanout.New()
	from := "."

	if !c.Args(&from) {
		return f, c.ArgErr()
	}

	to := c.RemainingArgs()
	if len(to) == 0 {
		return f, c.ArgErr()
	}

	toHosts, err := parse.HostPortOrFile(to...)
	if err != nil {
		return f, err
	}

	for c.NextBlock() {
		return f, fmt.Errorf("additional parameters not allowed")
	}
	for _, host := range toHosts {
		trans, h := parse.Transport(host)
	    if trans != transport.DNS {
	    return f, fmt.Errorf("only dns transport allowed")
		}
		client:=fanout.NewClient(h, trans)
		f.AddClient(client)

	}
	return f, nil
}
