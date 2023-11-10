# credhub-client
Credhub is a service that stores secrets and runs in the Cloud Foundry ecosystem.

Actions are executed through its REST API and authentication is done with OAUTH2 through UAA (cloud foundry's user 
authentication and authorizadion service) or MTLS using the application instance client
certificate which is generated and injected in cf app containers.

As such, this client library is targeted only at applications that will run in a CF
environment, as it only makes sense in that context. It is its intention to provide
a simple and straightforward way to interact with credhub.

## Usage

Creating a client is done by calling the `NewClient` function, either with no options, in which case a client will be
created using MTLS authentication (refreshing the client certificate when it expires, as
cloud foundry generates a new instance client certificate every 23h) and using "https://credhub.service.cf.internal:8844"
as the default credhub endpoint, or providing an Options object which allows the configuration of
a custom credhub endpoint an/or client/secret credentials in which case OAuth2 authentication will be used,
refreshing the bearer token as needed.

```go
package main

import (
    "fmt"

    "github.com/rabobank/credhub-client"
)

func main() {
    // Create a client using MTLS authentication. No error is produced in this case.
    mtlsClient, _ := credhub.New(nil)
    // List all secrets the client has access to
    if secrets, e := mtlsClient.FindByPath("/"); e != nil {
        fmt.Println("error listing credentials using mtls", e)
    } else {
        fmt.Println("secrets found:")
        for _, secret := range secrets.Credentials {
            fmt.Println(secret.Name)
        }
    }

    // Create a client using MTLS authentication specifying the credhub url. No error is produced in this case.
    // Options without either the client or the secret will result in creating an mtls authenticated client.
    mtlsClient, _ = credhub.New(&credhub.Options{Url: "https://credhub.service.cf.internal:8844"})
    if secrets, e := mtlsClient.FindByPath("/"); e != nil {
        fmt.Println("error listing credentials using mtls", e)
    } else {
        fmt.Println("secrets found:")
        for _, secret := range secrets.Credentials {
            fmt.Println(secret.Name)
        }
    }

    // Create a client using OAuth2 authentication
    if uaaClient, e := credhub.New(&credhub.Options{Url: "https://credhub.service.cf.internal:8844", Client: "client", Secret: "secret"}); e != nil {
        // when providing the credentials, the client will try to get the auth-server from the credhub /info endpoint
        // and try to retrieve a bearer token. Any of these two actions can result in an error.
        fmt.Println("error creating credhub oauth2 client", e)
    } else if secrets, e := uaaClient.FindByPath("/"); e != nil {
        fmt.Println("error listing credentials using mtls", e)
    } else {
        fmt.Println("secrets found:")
        for _, secret := range secrets.Credentials {
            fmt.Println(secret.Name)
        }
    }
}
```

## Supported Methods

None of the deprecated methods are supported. Current client implementation implements:

```go
package credhub

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
```

### Find Credentials
Finding credentials queries for credentials either by a partial name or a partial path. Returns lists of existing 
credential names. Returned list will only contain credentials the authenticate client has permissions to

### Get Credentials

Retrieving credentials can be done by either getting the credential object, which has also the other 
credential attributes, like metadata, creation date and UUID, or in case of methods that retrieve specific
credential types, returning the actual credential value. 

Only if the authenticated client has read permissions for the credential will the retrieval succeed.

### Set Credentials

Setting credentials can be done through the generic SetByName method, compatible with all
credential types, or calling the specific credential type methods, where the actual expected value
type is enforced by the method signature.

As in the case of other operations, calls will only succeed if the authenticated actor has write permissions for the 
path the name falls into

### Delete Credentials

Credential deletion is the same for all credential types and is done by identifying the credential name

Only if the authenticated actor has write permissions for the