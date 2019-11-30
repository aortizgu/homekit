package jobs

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/grandcat/zeroconf"
)

func init() {
	go func() {
		server, err := zeroconf.Register("homekit.", "_http._tcp.", "local.", 9000, nil, nil)
		if err != nil {
			panic(err)
		}
		defer server.Shutdown()

		// Clean exit.
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		select {
		case <-sig:
		}
	}()
}
