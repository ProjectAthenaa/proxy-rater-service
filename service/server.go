package service

import (
	"context"
	protos "github.com/ProjectAthenaa/sonic-core/protos/proxy-rater"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent/product"
)

type Server struct {
	protos.UnimplementedProxyRaterServer
}

func init() {
	go rater.listen()
}

func (s Server) GetProxy(_ context.Context, site *protos.Site) (*protos.Proxy, error) {
	return rater.GetEntry(product.Site(site.Value)), nil
}
