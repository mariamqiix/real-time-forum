package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type GitHubTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// route: /login/github

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	GitHubClientID, _ := loadEnvVariables("GitHubClientID")
	// Redirect the user to the GitHub authentication page
	http.Redirect(w, r, fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=http://localhost:8080/github-callback", GitHubClientID), http.StatusSeeOther)
}

func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Extract the authorization code from the query parameters
	code := r.URL.Query().Get("code")
	done, _ := ExchangeGitHubCodeForToken(code)
	UserInfo, _ := GetGitHubUserInfo(done.AccessToken)

	signUsingAuth(UserInfo.ProfilePicURL, UserInfo.ProfileURL, strconv.Itoa(UserInfo.ID), UserInfo.Username, done.AccessToken, w)
	// Redirect or respond as needed
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ExchangeGitHubCodeForToken(code string) (GitHubTokenResponse, error) {
	GitHubClientID, _ := loadEnvVariables("GitHubClientID")
	GitHubClientSecret, _ := loadEnvVariables("GitHubClientSecret")
	// Prepare the token request payload
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", GitHubClientID)
	data.Set("client_secret", GitHubClientSecret)
	data.Set("redirect_uri", "http://localhost:8080/github-callback")
	data.Set("grant_type", "authorization_code")

	// Send the token request
	resp, err := http.PostForm("https://github.com/login/oauth/access_token", data)
	if err != nil {
		return GitHubTokenResponse{}, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return GitHubTokenResponse{}, fmt.Errorf("failed to exchange GitHub code for token. Status: %s", resp.Status)
	}

	// Parse the token response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GitHubTokenResponse{}, err
	}

	tokenResp, err := parseGitHubTokenResponse(string(body))
	if err != nil {
		return GitHubTokenResponse{}, err
	}

	return tokenResp, nil
}

func parseGitHubTokenResponse(response string) (GitHubTokenResponse, error) {
	values, err := url.ParseQuery(response)
	if err != nil {
		return GitHubTokenResponse{}, err
	}

	accessToken := values.Get("access_token")
	tokenType := values.Get("token_type")
	scope := values.Get("scope")

	return GitHubTokenResponse{
		AccessToken: accessToken,
		TokenType:   tokenType,
		Scope:       scope,
	}, nil
}

// GetGitHubUserInfo retrieves user information from GitHub using an access token
func GetGitHubUserInfo(accessToken string) (GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return GitHubUser{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GitHubUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GitHubUser{}, fmt.Errorf("GitHub API returned non-200 status code: %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return GitHubUser{}, err
	}

	// Fetch the user's profile picture URL
	user.ProfilePicURL, err = getGitHubProfilePicURL(user.Username)
	if err != nil {
		// If there's an error fetching the profile picture URL, we'll just ignore it and continue
		fmt.Printf("Warning: Failed to fetch profile picture URL for user %s: %v\n", user.Username, err)
	}

	return user, nil
}

// Function to fetch the profile picture URL from GitHub
func getGitHubProfilePicURL(username string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/%s", username), nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned non-200 status code: %d", resp.StatusCode)
	}

	var userData struct {
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return "", err
	}

	return userData.AvatarURL, nil
}

type GitHubUser struct {
	ID            int    `json:"id"`
	Username      string `json:"login"`
	Email         string `json:"email"`
	ProfileURL    string `json:"html_url"`
	ProfilePicURL string // Adding profile picture URL field
}
