package usecase

import (
	"testing"

	"github.com/betorvs/secretpublisher/appcontext"
	"github.com/stretchr/testify/assert"
)

var (
	RepositoryGetSecretByNameCalls int
	RepositoryPostSecretCalls      int
	RepositoryPUTSecretCalls       int
	RepositoryPostOrPUTSecretCalls int
	RepositoryDeleteCalls          int
)

func TestCreateCheckSum(t *testing.T) {
	value := "value"
	test := createCheckSum(value)
	assert.Contains(t, test, "ec2c83edecb60304")
}

type RepositoryMock struct {
}

func (repo RepositoryMock) GetSecretByName(secret string, namespace string) (string, error) {
	RepositoryGetSecretByNameCalls++
	return "notFound", nil
}

func (repo RepositoryMock) PostOrPUTSecret(method string, secret string, body []byte) error {
	switch method {
	case "PUT":
		RepositoryPUTSecretCalls++
	case "POST":
		RepositoryPostSecretCalls++
	default:
		RepositoryPostOrPUTSecretCalls++
	}
	return nil
}

func (repo RepositoryMock) DeleteSecretK8S(secret string, namespace string) error {
	RepositoryDeleteCalls++
	return nil
}

func TestCheckSecret(t *testing.T) {
	repo := RepositoryMock{}
	appcontext.Current.Add(appcontext.Repository, repo)
	_, err := CheckSecret("foo", "default")
	assert.NoError(t, err)
	expected := 1
	if RepositoryGetSecretByNameCalls != expected {
		t.Fatalf("Invalid 2.1 TestCheckSecret %d", RepositoryGetSecretByNameCalls)
	}
}

func TestCreateSecret(t *testing.T) {
	repo := RepositoryMock{}
	appcontext.Current.Add(appcontext.Repository, repo)
	secret := GenerateSecret("foo")
	err := CreateSecret("foo", secret)
	assert.NoError(t, err)
	expected := 1
	if RepositoryPostSecretCalls != expected {
		t.Fatalf("Invalid 3.1 TestCreateSecret %d", RepositoryPostSecretCalls)
	}
}

func TestUpdateSecret(t *testing.T) {
	repo := RepositoryMock{}
	appcontext.Current.Add(appcontext.Repository, repo)
	secret := GenerateSecret("foo")
	err := UpdateSecret("foo", secret)
	assert.NoError(t, err)
	expected := 1
	if RepositoryPUTSecretCalls != expected {
		t.Fatalf("Invalid 4.1 TestUpdateSecret %d", RepositoryPUTSecretCalls)
	}
}

func TestDeleteSecret(t *testing.T) {
	repo := RepositoryMock{}
	appcontext.Current.Add(appcontext.Repository, repo)
	err := DeleteSecret("foo")
	assert.NoError(t, err)
	expected := 1
	if RepositoryDeleteCalls != expected {
		t.Fatalf("Invalid 2.1 TestDeleteSecret %d", RepositoryDeleteCalls)
	}

}

func TestManageSecret(t *testing.T) {
	repo := RepositoryMock{}
	appcontext.Current.Add(appcontext.Repository, repo)
	secret := GenerateSecret("foo")
	test := ManageSecret("foo", secret)
	assert.NoError(t, test)
}
