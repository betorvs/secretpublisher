package usecase

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/betorvs/secretpublisher/config"
	"github.com/betorvs/secretpublisher/domain"
	"github.com/betorvs/secretpublisher/gateway/kubeclient"
	"github.com/betorvs/secretpublisher/utils"
	"gopkg.in/yaml.v2"
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
		if test == parsedRes {
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
	// create a loop to check using manage secret
	if len(res.Items) == 0 {
		return fmt.Sprintf("Secrets with label %s not found\n", labels), nil
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

// ScanSubvalueSecret func
func ScanSubvalueSecret(labels string) (string, error) {
	res, errGateway := kubeclient.GetSecrets(config.SecretNamespace, labels)
	if errGateway != nil {
		errlocal := utils.ErrorHandler(errGateway)
		return "", errlocal
	}
	// create a loop to check using manage secret
	if len(res.Items) == 0 {
		return fmt.Sprintf("Secrets with label %s not found\n", labels), nil
	}
	var countErrors int
	var countErrorsNames []string
	for _, item := range res.Items {
		if config.DisabledLabel != "" {
			if searchLabels(config.DisabledLabel, item.Labels) {
				fmt.Printf("Skiping secret %s \n", item.Name)
				continue
			}
		}
		data := make(map[string]string)
		var suffixName, key, subkey string
		if strings.Contains(config.MatchKey, ".") {
			splited := strings.Split(config.MatchKey, ".")
			key = splited[0]
			subkey = splited[1]
		}
		for k, v := range item.Data {
			if k == config.NameSuffix {
				suffixName = string(v)
			}
			if k == key {
				// fmt.Println(k)
				temp := make(map[string]string)
				err := yaml.Unmarshal(v, &temp)
				if err != nil {
					fmt.Println("fail in Unmarshal")
					countErrors++
				}
				// fmt.Println(temp)
				// if subkey is not empty
				if subkey != "" {
					data[subkey] = temp[subkey]
				}
			}

		}
		destination := config.DestinationNamespace
		if destination == "" {
			destination = item.Namespace
		}
		labels := make(map[string]string)
		if config.NewLabels != "" {
			if strings.Contains(config.NewLabels, "=") {
				splited := strings.Split(config.NewLabels, "=")
				labels[splited[0]] = splited[1]
			}
		}
		annotations := make(map[string]string)
		if config.NewAnnotations != "" {
			if strings.Contains(config.NewAnnotations, "=") {
				splited := strings.Split(config.NewAnnotations, "=")
				annotations[splited[0]] = splited[1]
			}
		}
		name := fmt.Sprintf("%s-%s-%s", item.Name, subkey, suffixName)
		localName := fmt.Sprintf("%s-%s", subkey, suffixName)
		if config.MiddleName != "" {
			localName = fmt.Sprintf("%s-%s-%s", subkey, config.MiddleName, suffixName)
		}
		// fmt.Println(localName)
		localData := map[string]string{
			localName: data[subkey],
		}
		// fmt.Println(localData)
		newSecret := rewriteSecret(name, destination, localData, labels, annotations)
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

func searchLabels(label string, labels map[string]string) bool {
	var key, value string
	if strings.Contains(label, "=") {
		splited := strings.Split(label, "=")
		key = splited[0]
		value = splited[1]
	} else {
		key = label
	}
	if len(labels) == 0 {
		return false
	}
	for k, v := range labels {
		if value != "" && k == key && v == value {
			return true
		}
		if value == "" {
			if k == key {
				return true
			}
		}

	}
	return false
}
