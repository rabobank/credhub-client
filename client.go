package credhub

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	FindByPathEndpoint   = "%s/api/v1/data?path=%s"
	FindByNameEndpoint   = "%s/api/v1/data?name-like=%s"
	GetByNameEndpoint    = "%s/api/v1/data?name=%s&current=true"
	GetByIdEndpoint      = "%s/api/v1/data/%s"
	SetByNameEndpoint    = "%s/api/v1/data"
	DeleteByNameEndpoint = "%s/api/v1/data?name=%s"
)

type Client interface {
	FindByPath(pathPrefix string) (*CredentialNames, error)
	FindByName(namePrefix string) (*CredentialNames, error)
	GetByName(name string) (*Credential[any], error)
	GetJsonCredentialByName(name string) (*Credential[map[string]any], error)
	GetJsonByName(name string) (map[string]any, error)
	GetById(id string) (*Credential[any], error)
	GetJsonCredentialById(id string) (*Credential[map[string]any], error)
	GetJsonById(id string) (map[string]any, error)
	SetByName(credentialType string, name string, value any) (*Credential[any], error)
	SetJsonByName(name string, value map[string]any) (*Credential[map[string]any], error)
	DeleteByName(name string) error
}

type Options struct {
	Url    string
	Client string
	Secret string
}

func New(options *Options) (Client, error) {
	if options == nil {
		options = &Options{}
	}
	if len(options.Url) == 0 {
		options.Url = "https://credhub.service.cf.internal:8844"
	} else if options.Url[len(options.Url)-1] == '/' {
		// remove trailing slash if present
		options.Url = options.Url[:len(options.Url)-1]
	}
	result := &client{url: options.Url}
	var e error
	if len(options.Client) == 0 || len(options.Secret) == 0 {
		// no credentials are provided, assume it should run in a cf environment and use mtls to authenticate
		result.httpClient = newMtlsClient()
	} else {
		// authenticate with the uaa authentication provider returned by credhub/info
		if result.httpClient, e = newUaaClient(options); e != nil {
			return nil, e
		}
	}
	return result, nil
}

type client struct {
	url        string
	httpClient HttpClient
}

func (c *client) FindByPath(pathPrefix string) (*CredentialNames, error) {
	return getJson[CredentialNames](c.httpClient, fmt.Sprintf(FindByPathEndpoint, c.url, url.PathEscape(pathPrefix)))
}

func (c *client) FindByName(namePrefix string) (*CredentialNames, error) {
	return getJson[CredentialNames](c.httpClient, fmt.Sprintf(FindByNameEndpoint, c.url, url.PathEscape(namePrefix)))
}

func (c *client) GetByName(name string) (*Credential[any], error) {
	if credentials, e := getJson[Credentials[any]](c.httpClient, fmt.Sprintf(GetByNameEndpoint, c.url, url.PathEscape(name))); e != nil {
		return nil, e
	} else if len(credentials.Data) == 0 {
		return nil, errors.New("no data")
	} else {
		return &credentials.Data[0], nil
	}
}

func (c *client) GetJsonCredentialByName(name string) (*Credential[map[string]any], error) {
	if credentials, e := getJson[Credentials[map[string]any]](c.httpClient, fmt.Sprintf(GetByNameEndpoint, c.url, url.PathEscape(name))); e != nil {
		return nil, e
	} else if len(credentials.Data) == 0 {
		return nil, errors.New("no data")
	} else {
		return &credentials.Data[0], nil
	}
}

func (c *client) GetJsonByName(name string) (map[string]any, error) {
	if credentials, e := c.GetJsonCredentialByName(name); e != nil {
		return nil, e
	} else {
		return credentials.Value, nil
	}
}

func (c *client) GetById(id string) (*Credential[any], error) {
	if credentials, e := getJson[Credential[any]](c.httpClient, fmt.Sprintf(GetByIdEndpoint, c.url, url.PathEscape(id))); e != nil {
		return nil, e
	} else {
		return credentials, nil
	}
}

func (c *client) GetJsonCredentialById(id string) (*Credential[map[string]any], error) {
	if credentials, e := getJson[Credential[map[string]any]](c.httpClient, fmt.Sprintf(GetByIdEndpoint, c.url, url.PathEscape(id))); e != nil {
		return nil, e
	} else {
		return credentials, nil
	}
}

func (c *client) GetJsonById(id string) (map[string]any, error) {
	if credentials, e := c.GetJsonCredentialById(id); e != nil {
		return nil, e
	} else {
		return credentials.Value, nil
	}
}

func (c *client) SetByName(credentialType string, name string, value any) (*Credential[any], error) {
	request := &CredentialRequest{
		Name:  name,
		Type:  credentialType,
		Value: value,
	}
	if credentials, e := putAndGetJsons[Credential[any]](c.httpClient, fmt.Sprintf(SetByNameEndpoint, c.url), request); e != nil {
		return nil, e
	} else {
		return credentials, nil
	}
}

func (c *client) SetJsonByName(name string, value map[string]any) (*Credential[map[string]any], error) {
	request := &CredentialRequest{
		Name:  name,
		Type:  "json",
		Value: value,
	}
	if credentials, e := putAndGetJsons[Credential[map[string]any]](c.httpClient, fmt.Sprintf(SetByNameEndpoint, c.url), request); e != nil {
		return nil, e
	} else {
		return credentials, nil
	}
}

func (c *client) DeleteByName(name string) error {
	request, e := http.NewRequest(http.MethodDelete, fmt.Sprintf(DeleteByNameEndpoint, c.url, url.PathEscape(name)), nil)
	if e != nil {
		return e
	}

	_, e = c.httpClient.Do(request)
	return e
}
