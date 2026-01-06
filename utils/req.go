package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func FetchBody(endpoint string, accessToken string, requestType string, params url.Values) ([]byte, error) {
	url := fmt.Sprint(endpoint)

	req, err := http.NewRequest(requestType, url, strings.NewReader(params.Encode()))
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

	statusCode := resp.StatusCode
	if statusCode < 200 || statusCode > 200 {
		errStr := fmt.Sprintln("Bad request: received non 200 err code")
		return []byte{}, errors.New(errStr)
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return nil, err
	}

	return body, err
}
