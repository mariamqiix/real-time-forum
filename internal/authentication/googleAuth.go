package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	// "sandbox/internal/sessionmanager"
)

var GoogleClientID string

const (
	GooglelogInURL = "https://accounts.google.com/o/oauth2/auth?client_id=494333147558-4fdt1969hq590gcuhm9qrpe0sf5c70rg.apps.googleusercontent.com&redirect_uri=http://localhost:8080/login/google/callback&response_type=code&scope=email"
	// GooglelogOutURL   = "https://accounts.google.com/logout"
)

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

// route: /login/google
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the Google authentication page
	http.Redirect(w, r, GooglelogInURL, http.StatusSeeOther)
}

// route: /login/google/callback
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting google callback")
	// Extract the authorization code from the query parameters
	code := r.URL.Query().Get("code")

	// Exchange the authorization code for an access token
	gToken, err := ExchangeGoogleCodeForToken(code)
	if err != nil {
		fmt.Printf("HandleGoogleCallback: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	// Get the user's email from Google

	userInfo, err := GetGoogleUserInfo(gToken.AccessToken)
	if err != nil {
		fmt.Printf("HandleGoogleCallback: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	signUsingAuth(userInfo.Picture, userInfo.Email, userInfo.Sub, userInfo.Name, gToken.AccessToken, w)

	// Redirect or respond as needed
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetGoogleUserEmail(accessToken string) (string, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request to the Google UserInfo endpoint
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return "", err
	}

	// Set the Authorization header with the access token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user info. Status: %s", resp.Status)
	}

	// Parse the response body to extract the email
	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", err
	}

	return userInfo.Email, nil
}

func ExchangeGoogleCodeForToken(code string) (*GoogleTokenResponse, error) {
	// Prepare the token request payload
	GoogleClientID, _ := loadEnvVariables("GoogleClientID")
	GoogleClientSecret, _ := loadEnvVariables("GoogleClientSecret")
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", GoogleClientID)
	data.Set("client_secret", GoogleClientSecret)
	data.Set("redirect_uri", "http://localhost:8080/login/google/callback")
	data.Set("grant_type", "authorization_code")

	// Send the token request
	resp, err := http.PostForm("https://accounts.google.com/o/oauth2/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange Google code for token. Status: %s", resp.Status)
	}
	// Parse the token response
	var tokenResp GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

type GoogleUserInfo struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Picture    string `json:"picture"`
}

func GetGoogleUserInfo(token string) (GoogleUserInfo, error) {
	// Prepare the request to Google userinfo endpoint
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return GoogleUserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GoogleUserInfo{}, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return GoogleUserInfo{}, fmt.Errorf("Google API returned non-200 status code: %d", resp.StatusCode)
	}

	// Decode the response body
	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return GoogleUserInfo{}, err
	}

	return userInfo, nil
}
