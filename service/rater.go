package service

import (
	"context"
	"encoding/base64"
	"fmt"
	protos "github.com/ProjectAthenaa/sonic-core/protos/proxy-rater"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent/product"
	"github.com/bradfitz/slice"
	"github.com/prometheus/common/log"
	"github.com/tcnksm/go-httpstat"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var rater = NewRater()

type Rater struct {
	ctx     context.Context
	proxies map[product.Site][]*ratedProxy
	locker  sync.Mutex
}

type ratedProxy struct {
	latency       time.Duration
	proxy         string
	authorization string
}

func NewRater() *Rater {
	r := &Rater{
		ctx:     context.Background(),
		proxies: map[product.Site][]*ratedProxy{},
		locker:  sync.Mutex{},
	}

	return r
}

func (r *Rater) Rate(proxy string, site product.Site) {
	var authorization string

	if v := strings.Split(proxy, ":"); len(v) == 4 {
		authorization = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", v[0], v[1])))
	}

	req, _ := http.NewRequest("GET", siteMap[site], nil)

	var result httpstat.Result

	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := http.DefaultClient

	if authorization != "" {
		pr, _ := url.Parse(fmt.Sprintf("http://%s", proxy))
		req.Header.Set("Proxy-Authorization", fmt.Sprintf("Basic %s", authorization))
		client = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(pr)},
		}
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error("do req: ", err)
		return
	}

	if _, err = io.Copy(ioutil.Discard, res.Body); err != nil {
		log.Error("io copy: ", err)
		return
	}
	res.Body.Close()
	result.End(time.Now())

	length := result.ContentTransfer(time.Now())

	r.locker.Lock()
	defer r.locker.Unlock()

	rProxy := &ratedProxy{
		latency:       length,
		proxy:         proxy,
		authorization: authorization,
	}

	r.addProxy(rProxy, site)
}

func (r *Rater) addProxy(proxy *ratedProxy, site product.Site) {
	var availableIndex = -1

	for i, v := range r.proxies[site] {
		if v == nil {
			availableIndex = i
			break
		}
	}

	if availableIndex != -1 {
		r.proxies[site][availableIndex] = proxy
	}

	r.proxies[site] = append(r.proxies[site], proxy)

	slice.SortInterface(r.proxies[site][:], func(i, j int) bool {
		return r.proxies[site][i].latency < r.proxies[site][i].latency
	})

	if len(r.proxies[site]) > 500 {
		r.proxies[site] = r.proxies[site][:500]
	}
}

func (r *Rater) GetEntry(site product.Site) *protos.Proxy {
	rand.Seed(time.Now().UnixNano())

	p := r.proxies[site][rand.Intn(len(r.proxies[site])-1)]

	return &protos.Proxy{
		Value:         p.proxy,
		Authorization: p.authorization,
	}
}
