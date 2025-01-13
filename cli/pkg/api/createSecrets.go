package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func CreateSecret(repo, env, secret, value string) error {
	uri := fmt.Sprintf("https://api.github.com/repos/%s/environments/%s/secrets/%s", repo, env, secret)
	// Create a new request using http.NewRequest() and set the Authorization header
	req, err := http.NewRequest("PUT", uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create secrets: %w", err)
	}
	// Set the Authorization header using req.Header.Set()
	req.Header.Add("Authorization", "bearer "+setHeader())
	req.Body = ioutil.NopCloser(
		strings.NewReader(fmt.Sprintf(`{"key_id=%s": "encrypted_value=%s"}`, secret, value)),
		// strings.NewReader(fmt.Sprintf(`{"value": "%s"}`, value)),
	)
	// Send the request using http.DefaultClient.Do() and check the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to secret API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to create secrets: %s", body)
	}

	fmt.Printf("Creating secret for %s in %s\n", secret, repo)
	return nil
}
