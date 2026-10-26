package main

import "testing"

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
