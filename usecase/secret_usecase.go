package usecase

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"

	"github.com/betorvs/secretpublisher/config"
	"github.com/betorvs/secretpublisher/domain"
	"github.com/betorvs/secretpublisher/gateway/kubeclient"
	"github.com/betorvs/secretpublisher/utils"
)

// GenerateSecret func uses generates a secret struct from flags
func GenerateSecret(secretName string) *domain.Secret {
	var checksum string
	for _, v := range config.StringData {
		checksum += v
	}
	secret := &domain.Secret{
		Name:        secretName,
		Namespace:   config.SecretNamespace,
		Checksum:    createCheckSum(checksum),
		Data:        config.StringData,
		Labels:      config.Labels,
		Annotations: config.Annotations,
	}
	return secret
}

// ManageSecret func
func ManageSecret(secretName string, secret *domain.Secret) error {
	// check if secret exist
	test := secret.Checksum
	res, err := CheckSecret(secretName, secret.Namespace)
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
			errUpdate := UpdateSecret(secretName, secret)
			if errUpdate != nil {
				return errUpdate
			}
			fmt.Println("[OK] Updated")
		}
	} else {
		if config.Debug {
			fmt.Println("[DEBUG] Creating")
		}
		errCreate := CreateSecret(secretName, secret)
		if errCreate != nil {
			return errCreate
		}
		fmt.Println("[OK] Created")
	}
	return nil

}

// CreateSecret func
func CreateSecret(secretName string, secret *domain.Secret) error {
	secretClient := domain.GetRepository()
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
func UpdateSecret(secretName string, secret *domain.Secret) error {
	secretClient := domain.GetRepository()
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
func CheckSecret(secretName, namespace string) (string, error) {
	secretClient := domain.GetRepository()
	res, errGateway := secretClient.GetSecretByName(secretName, namespace)
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

// ScanSecret func
func ScanSecret(labels string) (string, error) {
	res, errGateway := kubeclient.GetSecrets(config.SecretNamespace, labels)
	if errGateway != nil {
		errlocal := utils.ErrorHandler(errGateway)
		return "", errlocal
	}
	// fmt.Printf("%#v\n", res)
	// create a loop to check using manage secret
	if len(res.Items) == 0 {
		// fmt.Println(len(res.Items))
		return fmt.Sprintf("Secrets with label %s not found\n", labels), nil
	}
	var countErrors int
	var countErrorsNames []string
	for _, item := range res.Items {
		// fmt.Println(item.Name)

		data := make(map[string]string)
		for k, v := range item.Data {
			data[k] = string(v)
		}
		destination := config.DestinationNamespace
		if destination == "" {
			destination = item.Namespace
		}
		name := item.Name
		if config.NameSuffix != "" {
			name = fmt.Sprintf("%s-%s", item.Name, config.NameSuffix)
		}
		newSecret := rewriteSecret(name, destination, data, item.Labels, item.Annotations)
		err := ManageSecret(name, newSecret)
		if err != nil {
			countErrors++
			countErrorsNames = append(countErrorsNames, item.Name)
		}
	}
	if countErrors != 0 {
		return "NOK", fmt.Errorf("Cannot process these secrets: %v", countErrorsNames)
	}
	return "OK", nil
}

// ScanConfigMap func
func ScanConfigMap(labels string) (string, error) {
	res, errGateway := kubeclient.GetConfigMaps(config.SecretNamespace, labels)
	if errGateway != nil {
		errlocal := utils.ErrorHandler(errGateway)
		return "", errlocal
	}
	var countErrors int
	var countErrorsNames []string
	for _, item := range res.Items {
		data := make(map[string]string)
		for k, v := range item.Data {
			data[k] = string(v)
		}
		destination := config.DestinationNamespace
		if destination == "" {
			destination = item.Namespace
		}
		name := item.Name
		if config.NameSuffix != "" {
			name = fmt.Sprintf("%s-%s", item.Name, config.NameSuffix)
		}
		newSecret := rewriteSecret(name, destination, data, item.Labels, item.Annotations)
		err := ManageSecret(name, newSecret)
		if err != nil {
			countErrors++
			countErrorsNames = append(countErrorsNames, item.Name)
		}
	}
	if countErrors != 0 {
		return "NOK", fmt.Errorf("Cannot process these config maps: %v", countErrorsNames)
	}
	return "OK", nil
}

// local rewrite func to rewrite secret and config map from K8S
func rewriteSecret(secretName, namespace string, data, labels, annotations map[string]string) *domain.Secret {
	var checksum string
	for _, v := range data {
		checksum += v
	}
	secret := &domain.Secret{
		Name:        secretName,
		Namespace:   namespace,
		Checksum:    createCheckSum(checksum),
		Data:        data,
		Labels:      labels,
		Annotations: annotations,
	}
	return secret
}
