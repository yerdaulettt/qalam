package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"api-gateway/internal/config"
)

func newProxy(serviceUrl string) (*httputil.ReverseProxy, error) {
	sUrl, err := url.Parse(serviceUrl)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(sUrl), nil
}

func InitRoutes(c config.Config) (*http.ServeMux, error) {
	r := http.NewServeMux()

	for _, s := range c.Services {
		p, err := newProxy(s.Url)
		if err != nil {
			return nil, err
		}

		for _, route := range s.Endpoints {
			r.HandleFunc(route, p.ServeHTTP)
		}
	}

	return r, nil
}
