package main

import (
	"net/http"
	"time"

	"github.com/lxc/incus/client"
	"github.com/lxc/incus/shared/util"
)

type srcManual struct{}

func (s *srcManual) present() bool {
	if !util.PathExists("/var/lib/lxd") {
		return false
	}

	return true
}

func (s *srcManual) name() string {
	return "manual installation"
}

func (s *srcManual) stop() error {
	d, err := s.connect()
	if err != nil {
		return err
	}

	httpClient, err := d.GetHTTPClient()
	if err != nil {
		return err
	}

	// Request shutdown, this shouldn't return until daemon has stopped so use a large request timeout.
	httpTransport := httpClient.Transport.(*http.Transport)
	httpTransport.ResponseHeaderTimeout = 3600 * time.Second
	_, _, err = d.RawQuery("PUT", "/internal/shutdown", nil, "")
	if err != nil {
		return err
	}

	return nil
}

func (s *srcManual) start() error {
	return nil
}

func (s *srcManual) purge() error {
	return nil
}

func (s *srcManual) connect() (incus.InstanceServer, error) {
	return incus.ConnectIncusUnix("/var/lib/lxd/unix.socket", &incus.ConnectionArgs{SkipGetServer: true})
}

func (s *srcManual) paths() (*daemonPaths, error) {
	return &daemonPaths{
		daemon: "/var/lib/lxd",
		logs:   "/var/log/lxd",
		cache:  "/var/cache/lxd",
	}, nil
}
