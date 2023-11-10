package credhub

import (
	"net/http"
)

type HttpClient interface {
	Do(request *http.Request) (response *http.Response, e error)
}
