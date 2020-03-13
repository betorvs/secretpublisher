package gateway

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/betorvs/secretpublisher/appcontext"
	"github.com/betorvs/secretpublisher/config"
	"github.com/betorvs/secretpublisher/utils"
)

// Repository struct
type Repository struct {
	Client *http.Client
}

// GetSecretByName func
func (repo Repository) GetSecretByName(secret string, namespace string) (string, error) {
	url := fmt.Sprintf("%s/%s/%s", config.ReceiverURL, namespace, secret)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", utils.ErrorHandler(err)
	}
	if config.EncodingRequest != "disabled" {
		now := time.Now()
		timestamp := string(now.Unix())
		req.Header.Add("X-SECRET-Request-Timestamp", timestamp)
		basestring := fmt.Sprintf("v1:%s:%s", timestamp, secret)
		signature := createHeaderSignature(timestamp, basestring, config.EncodingRequest)
		req.Header.Add("X-SECRET-Signature", signature)
	}
	resp, err := repo.Client.Do(req)
	if err != nil {
		return "", utils.ErrorHandler(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", utils.ErrorHandler(err)
	}
	if resp.StatusCode == 204 {
		return "notFound", nil
	}
	res := fmt.Sprintln(string(bodyText))
	if config.Debug {
		fmt.Printf("[SECRETRECEIVER] Response Code: %s, Body Response: %s \n", resp.Status, res)
	}
	if resp.StatusCode >= 400 {
		errLocal := fmt.Errorf("%s", resp.Status)
		return "", errLocal
	}
	defer resp.Body.Close()
	return res, nil
}

// PostOrPUTSecret func
func (repo Repository) PostOrPUTSecret(method string, secret string, body []byte) error {
	url := config.ReceiverURL
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return utils.ErrorHandler(err)
	}
	if config.EncodingRequest != "disabled" {
		now := time.Now()
		timestamp := string(now.Unix())
		req.Header.Add("X-SECRET-Request-Timestamp", timestamp)
		basestring := fmt.Sprintf("v1:%s:%s", timestamp, secret)
		signature := createHeaderSignature(timestamp, basestring, config.EncodingRequest)
		req.Header.Add("X-SECRET-Signature", signature)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := repo.Client.Do(req)
	if err != nil {
		return utils.ErrorHandler(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
		return err
	}
	if config.Debug {
		s := string(bodyText)
		fmt.Printf("[SECRETRECEIVER] Response Code: %s, Body Response: %s \n", resp.Status, s)
	}
	if resp.StatusCode > 204 {
		errLocal := fmt.Errorf("%s", resp.Status)
		return errLocal
	}
	defer resp.Body.Close()
	return nil
}

// DeleteSecretK8S func
func (repo Repository) DeleteSecretK8S(secret string, namespace string) error {
	url := fmt.Sprintf("%s/%s/%s", config.ReceiverURL, namespace, secret)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return utils.ErrorHandler(err)
	}
	if config.EncodingRequest != "disabled" {
		now := time.Now()
		timestamp := string(now.Unix())
		req.Header.Add("X-SECRET-Request-Timestamp", timestamp)
		basestring := fmt.Sprintf("v1:%s:%s", timestamp, secret)
		signature := createHeaderSignature(timestamp, basestring, config.EncodingRequest)
		req.Header.Add("X-SECRET-Signature", signature)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := repo.Client.Do(req)
	if err != nil {
		return utils.ErrorHandler(err)
	}
	if config.Debug {
		fmt.Printf("[SECRETRECEIVER] Response: %s \n", resp.Status)
	}
	if resp.StatusCode > 204 {
		errLocal := fmt.Errorf("%s", resp.Status)
		return errLocal
	}
	defer resp.Body.Close()
	return nil
}

// createHeaderSignature
func createHeaderSignature(timestamp string, message string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(message)); err != nil {
		fmt.Printf("mac.Write(%v) failed\n", message)
		return ""
	}
	calculatedMAC := "v0=" + hex.EncodeToString(mac.Sum(nil))
	return calculatedMAC
}

func init() {
	if config.TestRun == "true" {
		return
	}
	client := http.Client{
		Timeout: time.Second * config.PublisherTimeout,
	}
	appcontext.Current.Add(appcontext.Repository, Repository{Client: &client})
	if appcontext.Current.Count() != 0 && config.Debug {
		fmt.Println("[INFO] Using Repository")
	}
}
