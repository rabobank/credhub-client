package credhub

import "time"

type Info struct {
	AuthServer struct {
		Url string `json:"url"`
	} `json:"auth-server"`
	App struct {
		Name string `json:"name"`
	} `json:"app"`
}

type CredentialNames struct {
	Credentials []struct {
		Name             string    `json:"name"`
		VersionCreatedAt time.Time `json:"version_created_at"`
	} `json:"credentials"`
}

type CredentialRequest struct {
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Value    any            `json:"value"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type Credential[T any] struct {
	Id               string         `json:"id"`
	Name             string         `json:"name"`
	Type             string         `json:"type"`
	VersionCreatedAt time.Time      `json:"version_created_at"`
	Value            T              `json:"value"`
	Metadata         map[string]any `json:"metadata,omitempty"`
}

type Credentials[T any] struct {
	Data []Credential[T] `json:"data"`
}
