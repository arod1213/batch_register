package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Pretty(v any) {
	data, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return
	}
	fmt.Println("", string(data))
}

func getModel[T any](endpoint string, auth *auth) (*T, error) {
	body, err := fetchBody(endpoint, auth.AccessToken, "GET")
	if err != nil {
		return nil, err
	}

	var x T
	err = json.Unmarshal(body, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func fetchBody(endpoint string, accessToken string, requestType string) ([]byte, error) {
	data := url.Values{}
	url := fmt.Sprint(endpoint)

	req, err := http.NewRequest(requestType, url, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Request creation error:", err)
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

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
	return body, err
}
