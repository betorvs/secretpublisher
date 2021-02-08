package usecase

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"

	"github.com/betorvs/secretpublisher/config"
	"github.com/betorvs/secretpublisher/domain"
	"github.com/betorvs/secretpublisher/utils"
)

// ManageSecret func
func ManageSecret(secretName string) error {
	// check if secret exist
	var checksum string
	for _, v := range config.StringData {
		checksum += v
	}
	test := createCheckSum(checksum)
	res, err := CheckSecret(secretName)
	if err != nil {
		return err
	}
	if res != "notFound" {
		if config.Debug {
			fmt.Printf("[DEBUG] Checking %s, %s, %s ", secretName, config.SecretNamespace, res)
		}
		parsedRes := utils.RemoveQuotes(res)
		check := test == parsedRes
		if check {
			fmt.Printf("[OK] Secret %s already exist\n", secretName)
		} else {
			if config.Debug {
				fmt.Println("[DEBUG] Updating")
			}
			errUpdate := UpdateSecret(secretName)
			if errUpdate != nil {
				return errUpdate
			}
			fmt.Println("[OK] Updated")
		}
	} else {
		if config.Debug {
			fmt.Println("[DEBUG] Creating")
		}
		errCreate := CreateSecret(secretName)
		if errCreate != nil {
			return errCreate
		}
		fmt.Println("[OK] Created")
	}
	return nil

}

// CreateSecret func
func CreateSecret(secretName string) error {
	secretClient := domain.GetRepository()
	var checksum string
	for _, v := range config.StringData {
		checksum += v
	}
	secret := &domain.Secret{
		Name:      secretName,
		Namespace: config.SecretNamespace,
		Checksum:  createCheckSum(checksum),
		Data:      config.StringData,
		Labels:    config.Labels,
	}
	bodymarshal, err := json.Marshal(&secret)
	if err != nil {
		errlocal := utils.ErrorHandler(err)
		return errlocal
	}
	method := "POST"
	errGateway := secretClient.PostOrPUTSecret(method, secretName, bodymarshal)
	if errGateway != nil {
		errlocal := utils.ErrorHandler(errGateway)
		return errlocal
	}
	return nil
}

// UpdateSecret func
func UpdateSecret(secretName string) error {
	secretClient := domain.GetRepository()
	var checksum string
	for _, v := range config.StringData {
		checksum += v
	}
	secret := &domain.Secret{
		Name:      secretName,
		Namespace: config.SecretNamespace,
		Checksum:  createCheckSum(checksum),
		Data:      config.StringData,
		Labels:    config.Labels,
	}
	bodymarshal, err := json.Marshal(&secret)
	if err != nil {
		errlocal := utils.ErrorHandler(err)
		return errlocal
	}
	method := "PUT"
	errGateway := secretClient.PostOrPUTSecret(method, secretName, bodymarshal)
	if errGateway != nil {
		errlocal := utils.ErrorHandler(errGateway)
		return errlocal
	}
	return nil
}

// CheckSecret func
func CheckSecret(secretName string) (string, error) {
	secretClient := domain.GetRepository()
	res, errGateway := secretClient.GetSecretByName(secretName, config.SecretNamespace)
	if errGateway != nil {
		errlocal := utils.ErrorHandler(errGateway)
		return "", errlocal
	}
	// fmt.Printf("%s", res)
	return res, nil
}

// DeleteSecret func
func DeleteSecret(secretName string) error {
	secretClient := domain.GetRepository()
	errGateway := secretClient.DeleteSecretK8S(secretName, config.SecretNamespace)
	if errGateway != nil {
		errlocal := utils.ErrorHandler(errGateway)
		return errlocal
	}
	return nil
}

// createCheckSum func
// Create a shasum hash similar to
// echo -n "value" | shasum -a 512
func createCheckSum(value string) string {
	shasum := sha512.New()
	_, err := shasum.Write([]byte(value))
	if err != nil {
		return ""
	}
	checksum := fmt.Sprintf("%x", shasum.Sum(nil))
	return checksum
}
