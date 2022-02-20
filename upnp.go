package nymo

import (
	"context"
	"errors"

	"github.com/huin/goupnp/dcps/internetgateway2"
	"golang.org/x/sync/errgroup"
)

type routerClient interface {
	AddPortMappingCtx(
		ctx context.Context,
		NewRemoteHost string,
		NewExternalPort uint16,
		NewProtocol string,
		NewInternalPort uint16,
		NewInternalClient string,
		NewEnabled bool,
		NewPortMappingDescription string,
		NewLeaseDuration uint32,
	) (err error)

	DeletePortMappingCtx(
		ctx context.Context,
		NewRemoteHost string,
		NewExternalPort uint16,
		NewProtocol string,
	) (err error)

	GetExternalIPAddressCtx(
		ctx context.Context,
	) (NewExternalIPAddress string, err error)
}

func pickRouterClient(ctx context.Context) (routerClient, error) {
	tasks, _ := errgroup.WithContext(ctx)

	var ip1Clients []*internetgateway2.WANIPConnection1
	tasks.Go(func() (err error) {
		ip1Clients, _, err = internetgateway2.NewWANIPConnection1Clients()
		return
	})
	var ip2Clients []*internetgateway2.WANIPConnection2
	tasks.Go(func() (err error) {
		ip2Clients, _, err = internetgateway2.NewWANIPConnection2Clients()
		return
	})
	var ppp1Clients []*internetgateway2.WANPPPConnection1
	tasks.Go(func() (err error) {
		ppp1Clients, _, err = internetgateway2.NewWANPPPConnection1Clients()
		return
	})

	if err := tasks.Wait(); err != nil {
		return nil, err
	}

	switch {
	case len(ip2Clients) >= 1:
		return ip2Clients[0], nil
	case len(ip1Clients) >= 1:
		return ip1Clients[0], nil
	case len(ppp1Clients) >= 1:
		return ppp1Clients[0], nil
	default:
		return nil, errors.New("no services found")
	}
}
