package service

import (
	"github.com/ProjectAthenaa/sonic-core/sonic/core"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent/product"
	"github.com/prometheus/common/log"
	"strings"
)

func (r *Rater) listen() {
	pubSub := core.Base.GetRedis("cache").PSubscribe(r.ctx, "proxies:*")

	for msg := range pubSub.Channel() {
		log.Info(msg.Channel, msg.Payload)
		site := product.Site(strings.Split(msg.Channel, ":")[1])
		go r.Rate(msg.Payload, site)
	}
}
