package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Editor-Version", "vscode/1.97.2")
	req.Header.Set("User-Agent", "GithubCopilot/1.155.0")
}

func main() {
	// Request
	reqBody := []byte(`{"client_id":"Iv1.b507a08c87ecfe98","scope":"read:user"}`)
	req, _ := http.NewRequest("POST", "https://github.com/login/device/code", bytes.NewBuffer(reqBody))
	setHeaders(req)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	// Parse the response
	fmt.Println(string(body))
	var respData map[string]interface{}
	json.Unmarshal(body, &respData)

	deviceCode := respData["device_code"].(string)
	userCode := respData["user_code"].(string)
	verificationUri := respData["verification_uri"].(string)

	fmt.Printf("Please visit \033]8;;%s\033\\GitHub OAuth's page\033]8;;\033\\ and enter code %s to authenticate.\n", verificationUri, userCode)

	// Poll for access token
	var accessToken string
	for {
		time.Sleep(5 * time.Second)

		pollRequestBody := fmt.Sprintf(`{"client_id":"Iv1.b507a08c87ecfe98","device_code":"%s","grant_type":"urn:ietf:params:oauth:grant-type:device_code"}`, deviceCode)
		pollReq, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer([]byte(pollRequestBody)))

		setHeaders(pollReq)

		pollResp, _ := client.Do(pollReq)
		pollRespBody, _ := io.ReadAll(pollResp.Body)
		pollResp.Body.Close()

		var pollData map[string]interface{}
		json.Unmarshal(pollRespBody, &pollData)

		if token, ok := pollData["access_token"].(string); ok {
			accessToken = token
			break
		}
	}

	fmt.Println(accessToken)
}
