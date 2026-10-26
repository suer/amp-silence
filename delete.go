package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var silenceIdPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func DeleteSilence(endpoint string, silenceId string) ([]byte, error) {
	if !silenceIdPattern.MatchString(silenceId) {
		return nil, fmt.Errorf("invalid silence id: %q", silenceId)
	}

	url, err := url.JoinPath(endpoint, "/alertmanager/api/v2/silence", silenceId)
	if err != nil {
		return nil, err
	}

	return RequestWithSigv4(http.MethodDelete, url, []byte{})
}
