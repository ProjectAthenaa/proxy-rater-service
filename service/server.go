package service

import (
	"context"
	"fmt"
	protos "github.com/ProjectAthenaa/sonic-core/protos/proxy-rater"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent/product"
	"time"
)

type Server struct {
	protos.UnimplementedProxyRaterServer
}

func init() {
	go rater.listen()

	go func() {
		for range time.Tick(time.Second) {
			for _, proxies := range rater.proxies {
				for _, proxy := range proxies {
					fmt.Println(proxy.proxy, proxy.latency)
				}
			}
		}
	}()
}

func (s Server) GetProxy(_ context.Context, site *protos.Site) (*protos.Proxy, error) {
	return rater.GetEntry(product.Site(site.Value)), nil
}
