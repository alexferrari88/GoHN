package setup

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/alexferrari88/gohn/pkg/gohn"
)

const (
	baseTestURLPath = "/hn-v0"
)

func Init() (client *gohn.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseTestURLPath+"/", http.StripPrefix(baseTestURLPath, mux))

	server := httptest.NewServer(apiHandler)

	client = gohn.NewClient(nil)
	url, _ := url.Parse(server.URL + baseTestURLPath + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}
