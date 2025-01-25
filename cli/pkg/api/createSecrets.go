package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/crypto/nacl/box"
)

func GetGithubPublicKey(org, repo, env string) (interface{}, error) {
	uri := fmt.Sprintf("https://api.github.com/repos/%s/%s/environments/%s/secrets/public-key", org, repo, env)
	// Create a new request using http.NewRequest() and set the Authorization header
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}
	// Set the Authorization header using req.Header.Set()
	req.Header.Add("Authorization", "bearer "+setHeader())
	// Send the request using http.DefaultClient.Do() and check the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to secret API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get public key: %s", body)
	}

	var publicKey map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	fmt.Printf("Getting public key for %s/%s\n", repo, env)
	return publicKey, nil // this is the public key payload which includes key_id and key_value
}

func EncryptValue(publicKey, value string) (string, error) {
	// Decode the base64 public key
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: '%s' \n %w", publicKey, err)
	}

	// Ensure the public key is 32 bytes long
	if len(publicKeyBytes) != 32 {
		return "", fmt.Errorf("public key is not 32 bytes long")
	}

	// convert to an array for use with box.seal
	var publicKeyArray [32]byte
	copy(publicKeyArray[:], publicKeyBytes)

	//Generate a one time ephemeral private key
	// var ephemeralPrivateKey, ephemeralPublicKey [32]byte
	// if _, err := rand.Read(ephemeralPrivateKey[:]); err != nil {
	// 	return "", fmt.Errorf("failed to generate ephemeral private key: %w", err)
	// }

	//compute the ephemeral public key
	// box.Precompute(&ephemeralPublicKey, &ephemeralPrivateKey, &publicKeyArray)

	//Encrypt the value
	nonce := &[24]byte{}
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	encrypted := box.Seal(nil, []byte(value), nonce, &publicKeyArray, new([32]byte))

	fmt.Println("Encrypting value")
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func CreateUpdateSecret(org, repo, env, secret, value, publickey, publickey_id string) error {

	encryptedValue, err := EncryptValue(publickey, value)
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	uri := fmt.Sprintf("https://api.github.com/repos/%s/%s/environments/%s/secrets/%s", org, repo, env, secret)
	// Create a new request using http.NewRequest() and set the Authorization header
	req, err := http.NewRequest("PUT", uri, nil)
	if err != nil {
		return fmt.Errorf("failed to sent PUT API request for create/update secrets: %w", err)
	}
	// Set the Authorization header using req.Header.Set()
	req.Header.Add("Authorization", "bearer "+setHeader())
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-version", "2022-11-28")
	req.Body = io.NopCloser(
		strings.NewReader(fmt.Sprintf(`
			{
				"encrypted_value":"%s",
				"key_id":"%s"
			}
		`, encryptedValue, publickey_id)),
		// strings.NewReader(fmt.Sprintf(`{"value": "%s"}`, value)),
	)
	// Send the request using http.DefaultClient.Do() and check the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to secret API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create secret '%s': %s", secret, body)
	}

	fmt.Printf("Creating secret for %s in %s\n", secret, repo)
	return nil
}
