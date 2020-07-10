package analytics

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type AnalyticsHandler interface {
	GetURLAnalytics(w http.ResponseWriter, r *http.Request)
}

type AnalyticsReverseProxyHandler struct {
	reverseProxy *httputil.ReverseProxy
}

func NewAnalyticsReverseProxyHandler(analyticsServer string) (AnalyticsHandler, error) {
	url, err := url.Parse(analyticsServer)
	if err != nil {
		return nil, err
	}
	return &AnalyticsReverseProxyHandler{reverseProxy: httputil.NewSingleHostReverseProxy(url)}, nil
}

func (handler AnalyticsReverseProxyHandler) GetURLAnalytics(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/panel", handler.reverseProxy).ServeHTTP(w, r)
}
