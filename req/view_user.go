package req

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	auth "github.com/dombartenope/filterfinder/auth"
	user "github.com/dombartenope/filterfinder/user"
)

func ViewUser() user.User {
	client := &http.Client{}
	app_id, api_key, onesignal_id := auth.CheckForAuth()

	base_url := "https://onesignal.com/api/v1/apps"

	url := fmt.Sprintf("%s/%s/users/by/onesignal_id/%s", base_url, app_id, onesignal_id)
	fmt.Printf("url: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	apiKey := fmt.Sprintf("Basic %s", api_key)

	req.Header.Add("Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading the response body: %v", err)
	}
	// Unmarshal to the User > Property struct
	var apiResp user.User
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Fatalf("Error parsing the JSON response: %v", err)
	}

	// // Write the response to a JSON file
	outputFile := "user.json"
	jsonData, err := json.MarshalIndent(apiResp, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling the JSON data: %v", err)
	}
	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		log.Fatalf("Error writing to the JSON file: %v", err)
	}
	fmt.Sprintf("Response written to JSON file. Resp is \n%s\n ", apiResp)

	return apiResp
}
