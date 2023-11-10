package credhub

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"time"
)

func newMtlsClient() HttpClient {
	return &mtlsAuthenticatedClient{httpClient: &http.Client{Transport: http.DefaultTransport.(*http.Transport).Clone()}}
}

type mtlsAuthenticatedClient struct {
	httpClient  *http.Client
	certificate *x509.Certificate
}

func (mac *mtlsAuthenticatedClient) Do(request *http.Request) (response *http.Response, e error) {
	if mac.expired() {
		if e := mac.renew(); e != nil {
			return nil, e
		}
	}

	return mac.httpClient.Do(request)
}

func (mac *mtlsAuthenticatedClient) expired() bool {
	if mac.certificate == nil || time.Now().Add(time.Minute).After(mac.certificate.NotAfter) {
		return true
	}
	return false
}

func (mac *mtlsAuthenticatedClient) renew() error {
	certificate, e := tls.LoadX509KeyPair(
		os.Getenv("CF_INSTANCE_CERT"),
		os.Getenv("CF_INSTANCE_KEY"),
	)
	if e != nil {
		return e
	}

	mac.httpClient.Transport.(*http.Transport).TLSClientConfig.Certificates = []tls.Certificate{certificate}
	mac.certificate, e = x509.ParseCertificate(certificate.Certificate[0])
	if e != nil {
		return e
	}

	return nil
}
