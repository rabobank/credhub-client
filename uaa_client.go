package credhub

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/cloudfoundry-community/go-uaa"
	"golang.org/x/oauth2"
)

func newUaaClient(options *Options) (HttpClient, error) {
	result := &uaaAuthenticatedClient{httpClient: &http.Client{}}
	if options != nil && options.IgnoreSsl {
		result.httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	// let's get the authentication server url
	info, e := getJson[Info](result.httpClient, options.Url+"/info")
	if e != nil {
		return nil, e
	}

	if result.uaaClient, e = uaa.New(info.AuthServer.Url, uaa.WithClientCredentials(options.Client, options.Secret, uaa.JSONWebToken), uaa.WithSkipSSLValidation(options != nil && options.IgnoreSsl)); e != nil {
		return nil, e
	}

	if result.token, e = result.uaaClient.Token(context.Background()); e != nil {
		return nil, e
	}

	return result, nil
}

type uaaAuthenticatedClient struct {
	httpClient *http.Client
	uaaClient  *uaa.API
	token      *oauth2.Token
}

func (uac *uaaAuthenticatedClient) Do(request *http.Request) (response *http.Response, e error) {
	if !uac.token.Valid() {
		if uac.token, e = uac.uaaClient.Token(context.Background()); e != nil {
			return nil, e
		}
	}

	uac.token.SetAuthHeader(request)
	return uac.httpClient.Do(request)
}
