package router

import (
	"io"
	"net/http"
	"net/url"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(r.RequestURI)
	if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
	}

	proxyReq, err := http.NewRequest(r.Method, url.String(), r.Body)
	if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
	}
	proxyReq.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
			http.Error(w, "Failed to fetch ", http.StatusInternalServerError)
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
}