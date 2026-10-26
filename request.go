package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	defaultRequestTimeout = 30 * time.Second
	timeoutEnvName        = "AMP_SILENCE_TIMEOUT"
)

func requestTimeout() (time.Duration, error) {
	v := os.Getenv(timeoutEnvName)
	if v == "" {
		return defaultRequestTimeout, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", timeoutEnvName, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("invalid %s: must be positive: %s", timeoutEnvName, v)
	}
	return d, nil
}

func validateEndpoint(rawURL string) error {
	u, err := neturl.Parse(rawURL)
	if err != nil {
		return err
	}
	if u.Scheme != "https" {
		return fmt.Errorf("endpoint must use https: %s", rawURL)
	}
	host := u.Hostname()
	isAmazonHost := strings.HasSuffix(host, ".amazonaws.com") || strings.HasSuffix(host, ".amazonaws.com.cn")
	if !strings.HasPrefix(host, "aps-workspaces.") || !isAmazonHost {
		return fmt.Errorf("endpoint host must be aps-workspaces.<region>.amazonaws.com: %s", host)
	}
	return nil
}

func RequestWithSigv4(httpMethod string, url string, body []byte) ([]byte, error) {
	if err := validateEndpoint(url); err != nil {
		return nil, err
	}

	timeout, err := requestTimeout()
	if err != nil {
		return nil, err
	}

	reqBody := strings.NewReader(string(body))
	hash := sha256.Sum256(body)
	hashString := hex.EncodeToString(hash[:])

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(httpMethod, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	signer := v4.NewSigner()
	signer.SignHTTP(ctx, credentials, req, hashString, "aps", cfg.Region, time.Now().UTC())
	httpClient := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to add silence: Status code %d: response: %s", response.StatusCode, string(responseBody))
	}

	return responseBody, nil

}
