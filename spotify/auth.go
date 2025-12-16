package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getAuth() (*auth, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	var auth *auth
	auth, err := getToken(clientID, clientSecret)
	if err != nil {
		return nil, err
	} else if auth.AccessToken == "" {
		err = errors.New("no valid token returned")
		return nil, err
	}
	return auth, nil
}

func getToken(clientID string, clientSecret string) (*auth, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	// Create request
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Request creation error:", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return nil, err
	}

	var auth auth
	err = json.Unmarshal(body, &auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}
