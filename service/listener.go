package service

import (
	"github.com/ProjectAthenaa/sonic-core/sonic/core"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent/product"
	"strings"
)

func (r *Rater) listen() {
	pubSub := core.Base.GetRedis("cache").PSubscribe(r.ctx, "proxies:*")

	for msg := range pubSub.Channel() {
		site := product.Site(strings.Split(msg.Channel, ":")[1])
		go r.Rate(msg.Payload, site)
	}
}
