package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"backend-sarpras/internal/config"
)

func UploadPDFToSupabase(objectPath string, fileBytes []byte) error {
	cfg := config.GetSupabaseConfig()
	if cfg.URL == "" || cfg.ServiceKey == "" || cfg.Bucket == "" {
		return fmt.Errorf("supabase config incomplete")
	}

	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", cfg.URL, cfg.Bucket, objectPath)

	req, err := http.NewRequest("POST", url, bytes.NewReader(fileBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ServiceKey)
	req.Header.Set("Content-Type", "application/pdf")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload gagal: %s", string(body))
	}

	return nil
}

func GenerateSignedURL(objectPath string) (string, error) {
	cfg := config.GetSupabaseConfig()
	if cfg.URL == "" || cfg.ServiceKey == "" || cfg.Bucket == "" {
		return "", fmt.Errorf("supabase config incomplete")
	}

	url := fmt.Sprintf("%s/storage/v1/object/sign/%s/%s", cfg.URL, cfg.Bucket, objectPath)

	payload := fmt.Sprintf(`{"expiresIn": %d}`, cfg.Expire)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ServiceKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gagal generate signed URL: %s", string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		SignedURL string `json:"signedURL"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("invalid signed url response: %s", string(body))
	}
	if result.SignedURL == "" {
		return "", fmt.Errorf("signed url empty")
	}

	if strings.HasPrefix(result.SignedURL, "http") {
		return result.SignedURL, nil
	}

	signedPath := result.SignedURL
	if !strings.HasPrefix(signedPath, "/") {
		signedPath = "/" + signedPath
	}
	if !strings.HasPrefix(signedPath, "/storage/v1/") {
		signedPath = "/storage/v1" + signedPath
	}

	return fmt.Sprintf("%s%s", cfg.URL, signedPath), nil
}
