package router

import (
	"io"
	"net/http"
	"net/url"
)

type Proxy struct {
	lb LoadBalancer
}

func NewProxy(lb LoadBalancer) *Proxy {
	return &Proxy{
			lb: lb,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		targetBaseUrl := p.lb.SelectTarget() // Select target server
		targetUrl, err := url.Parse(targetBaseUrl)
		if err != nil {
				http.Error(w, "Bad Gateway", http.StatusBadGateway) // 502 ( :0 ohh ) 	
				return
		}
		proxyURL := targetUrl.ResolveReference(r.URL) // Transform Request URL's base

		proxyReq, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
		if err != nil {
				http.Error(w, "Failed to create request", http.StatusInternalServerError)
				return
		}
		proxyReq.Header = r.Header

		p.lb.RecordRequest() // Record request

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
				http.Error(w, "Failed to fetch the URL", http.StatusInternalServerError)
				return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
				for _, value := range values {
						w.Header().Add(key, value)
				}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		p.lb.RecordResponse(targetBaseUrl) // Record response
}