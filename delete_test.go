package main

import "testing"

func TestDeleteSilenceRejectsInvalidId(t *testing.T) {
	ids := []string{"", "../../x", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "12345678-1234-1234-1234-123456789012/../x"}
	for _, id := range ids {
		if _, err := DeleteSilence("https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-1/", id); err == nil {
			t.Errorf("DeleteSilence(%q) should fail", id)
		}
	}
}
