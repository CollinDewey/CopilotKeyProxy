package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var sessions = make(map[string]string)

func notExpired(token string) bool {
	for _, part := range strings.Split(token, ";") {
		if strings.HasPrefix(part, "exp=") {
			expStr := strings.TrimPrefix(part, "exp=")
			exp, err := strconv.ParseInt(expStr, 10, 64)
			return err == nil && exp >= time.Now().Unix()
		}
	}
	return false
}

func authenticate(key string) (string, error) {
	// Relay if key is valid
	token, exists := sessions[key]
	if exists && notExpired(token) {
		return sessions[key], nil
	}

	// Renew key
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/copilot_internal/v2/token", nil)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	setHeaders(key, req)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %v", err)
	}

	var respJSON map[string]interface{}
	if err := json.Unmarshal(body, &respJSON); err != nil {
		return "", fmt.Errorf("Error parsing JSON: %v", err)
	}

	if token, exists := respJSON["token"].(string); exists {
		sessions[key] = token
		return token, nil
	} else {
		return "", fmt.Errorf("Error making request: %v", err)
	}
}

func setHeaders(key string, req *http.Request) {
	req.Header.Set("Authorization", key)
	req.Header.Set("Editor-Version", "vscode/1.97.2")
	req.Header.Set("User-Agent", "GithubCopilot/1.155.0")
}

func main() {
	var listenAddr string = ":8080"

	if len(os.Args) > 1 {
		listenAddr = os.Args[1]
	}

	target, _ := url.Parse("https://api.individual.githubcopilot.com")
	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		req.Host = target.Host
		token, err := authenticate(req.Header.Get("Authorization"))
		if err != nil {
			return
		}
		setHeaders("Bearer "+token, req)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	fmt.Printf("Starting reverse proxy server on %s\n", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
