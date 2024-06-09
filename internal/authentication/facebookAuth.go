package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the Facebook authentication page
	FacebookAppID, _ := loadEnvVariables("FacebookAppID")
	FaceBookLogInURL := fmt.Sprintf("https://www.facebook.com/v14.0/dialog/oauth?client_id=%s&redirect_uri=http://localhost:8080/facebook-callback&state=random_state_string", FacebookAppID)
	http.Redirect(w, r, FaceBookLogInURL, http.StatusSeeOther)
}
func HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	// Extract the authorization code from the query parameters
	code := r.URL.Query().Get("code")
	// state := r.URL.Query().Get("state")
	done, err := exchangeFacebookCodeForToken(code)
	if err != nil {
		// Handle error
		return
	}

	userinfo, err := GetFacebookUserInfo(done.AccessToken)
	if err != nil {
		fmt.Print(err)
		return
	}

	signUsingAuth(userinfo.Picture.Data.URL, userinfo.Email, userinfo.ID, userinfo.Name, done.AccessToken, w)

	// Redirect or respond as needed
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func exchangeFacebookCodeForToken(code string) (FacebookTokenResponse, error) {
	// Prepare the token request payload
	FacebookAppID, _ := loadEnvVariables("FacebookAppID")
	FacebookAppSecret, _ := loadEnvVariables("FacebookAppSecret")
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", FacebookAppID)
	data.Set("client_secret", FacebookAppSecret)
	data.Set("redirect_uri", "http://localhost:8080/facebook-callback")
	// Send the token request
	resp, err := http.PostForm("https://graph.facebook.com/v14.0/oauth/access_token", data)
	if err != nil {
		return FacebookTokenResponse{}, err
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return FacebookTokenResponse{}, fmt.Errorf("failed to exchange Facebook code for token. Status: %s", resp.Status)
	}
	// Parse the token response
	var tokenResp FacebookTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return FacebookTokenResponse{}, err
	}
	return tokenResp, nil
}

type FacebookTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type FacebookUserInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

func GetFacebookUserInfo(token string) (FacebookUserInfo, error) {
	// Prepare the request to Facebook userinfo endpoint
	req, err := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email,picture.type(large)", nil)
	if err != nil {
		return FacebookUserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return FacebookUserInfo{}, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {

		return FacebookUserInfo{}, fmt.Errorf("Facebook API returned non-200 status code: %s", resp.StatusCode)
	}

	// Decode the response body
	var userInfo FacebookUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return FacebookUserInfo{}, err
	}

	return userInfo, nil
}

func PrintFacebookUserInfo(userInfo FacebookUserInfo) {
	fmt.Println("Facebook User Information:")
	fmt.Println("ID:", userInfo.ID)
	fmt.Println("Name:", userInfo.Name)
	fmt.Println("Email:", userInfo.Email)
}
