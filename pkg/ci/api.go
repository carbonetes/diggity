package ci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/carbonetes/diggity/internal/log"
)

const (
	tokenURL = "https://tent-api.carbonetes.com/personal-access-token/is-expired"
	saveURL  = "https://tent-api.carbonetes.com/integrations/bom/plugin/save"
)

var tokenId = "0"
var permitted = false

func PersonalAccessToken(token string, pluginType string) {
	// Payload
	payload := map[string]string{
		"token":      token,
		"pluginType": pluginType,
	}

	// Perform HTTP POST request
	resp, body := apiRequest(payload, tokenURL)
	// ---------------

	if resp.StatusCode != 200 {
		var appError ApplicationErrorResponse
		if err := json.Unmarshal(body, &appError); err != nil {
			log.Fatal("Failed to parse response:", err)
			os.Exit(1)
		}
		log.Print("Error: ", appError.Message)
		os.Exit(1)
	}
	// Unmarshal the body into the struct
	var result TokenCheckResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal("Failed to parse response:", err)
		os.Exit(1)
	}

	for _, p := range result.Permissions {
		if p.Label == "Pipelines" {
			for _, lp := range p.Permissions {
				if lp == "write" {
					permitted = true
				}
			}
		}
	}

	if !permitted {
		log.Fatal("Error: You do not have pipeline write permission.")
		os.Exit(1)
	}

	tokenId = result.PersonalAccessTokenId
	if result.PersonalAccessTokenId == "" {
		log.Fatal("Status Code:", resp.StatusCode)
		log.Fatal("Error: Unable to fetch token id.")
		os.Exit(1)
	}

}

func SavePluginRepository(bom *cyclonedx.BOM, repoName string, pluginName string, start time.Time, secrets interface{}) {

	// Secrets
	secretsBytes, err := json.Marshal(secrets)
	if err != nil {
		log.Fatal("Failed to marshal secrets:", err)
		os.Exit(1)
	}
	secretsJSONString := string(secretsBytes)

	var bomJSONString string
	if bom == nil || bom.Components == nil {
		// Empty State
		bomJSONString = ""

	} else {
		// BOM
		bomBytes, err := json.Marshal(bom)
		if err != nil {
			log.Fatal("Failed to marshal components:", err)
			os.Exit(1)
		}
		bomJSONString = string(bomBytes)

	}
	// Payload
	payload := map[string]interface{}{
		"repoName":              repoName,
		"personalAccessTokenId": tokenId,
		"pluginName":            pluginName,
		"bom":                   bomJSONString,
		"duration":              fmt.Sprintf("%.2f", time.Since(start).Seconds()),
		"secrets":               secretsJSONString,
	}
	// Perform HTTP POST request
	resp, body := apiRequest(payload, saveURL)
	// ---------------

	var result PluginRepo

	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal("Failed to parse response:", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Status Code:", resp.StatusCode)
		log.Fatal("Response Body:", string(body))
		os.Exit(1)
	}

}

func apiRequest(payload any, url string) (*http.Response, []byte) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Perform HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body (modern way)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return resp, body

}
