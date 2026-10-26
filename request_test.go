package main

import (
	"testing"
	"time"
)

func TestValidateEndpoint(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{"https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-1/alertmanager/api/v2/silences", false},
		{"https://aps-workspaces.cn-north-1.amazonaws.com.cn/workspaces/ws-1/", false},
		{"http://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-1/", true},
		{"https://example.com/workspaces/ws-1/", true},
		{"https://aps-workspaces.ap-northeast-1.amazonaws.com.evil.example/", true},
		{"https://evil.example/aps-workspaces.ap-northeast-1.amazonaws.com", true},
		{"https://aps-workspaces.ap-northeast-1.amazonaws.com@evil.example/", true},
	}
	for _, tt := range tests {
		err := validateEndpoint(tt.url)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateEndpoint(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
		}
	}
}

func TestRequestTimeout(t *testing.T) {
	tests := []struct {
		env     string
		want    time.Duration
		wantErr bool
	}{
		{"", 30 * time.Second, false},
		{"5s", 5 * time.Second, false},
		{"2m", 2 * time.Minute, false},
		{"abc", 0, true},
		{"10", 0, true},
		{"0s", 0, true},
		{"-5s", 0, true},
	}
	for _, tt := range tests {
		t.Setenv(timeoutEnvName, tt.env)
		got, err := requestTimeout()
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("requestTimeout() with %q = %v, %v; want %v, wantErr %v", tt.env, got, err, tt.want, tt.wantErr)
		}
	}
}
