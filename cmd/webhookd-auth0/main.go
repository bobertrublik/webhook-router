package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	audience := safeGetEnv("AUTH0_AUDIENCE")
	domain := safeGetEnv("AUTH0_DOMAIN")
	clientId := safeGetEnv("AUTH0_CLIENT_ID")
	clientSecret := safeGetEnv("AUTH0_CLIENT_SECRET")

	url := fmt.Sprintf("https://%s/oauth/token", domain)

	payload := strings.NewReader(fmt.Sprintf("{\"client_id\":\"%s\",\"client_secret\":\"%s\",\"audience\":\"%s\",\"grant_type\":\"client_credentials\"}", clientId, clientSecret, audience))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

func safeGetEnv(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("The environment variable '%s' doesn't exist or is not set", key)
	}
	return os.Getenv(key)
}
