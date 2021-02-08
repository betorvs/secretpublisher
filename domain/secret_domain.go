package domain

import "github.com/betorvs/secretpublisher/appcontext"

// Secret struct
type Secret struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Checksum  string            `json:"checksum"`
	Data      map[string]string `json:"data"`
	Labels    map[string]string `json:"labels"`
}

// Repository interface
type Repository interface {
	appcontext.Component
	GetSecretByName(secret string, namespace string) (string, error)
	PostOrPUTSecret(method string, secret string, body []byte) error
	DeleteSecretK8S(secret string, namespace string) error
}

// GetRepository func return Repository interface
func GetRepository() Repository {
	return appcontext.Current.Get(appcontext.Repository).(Repository)
}
