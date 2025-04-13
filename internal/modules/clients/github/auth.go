package github_client

import (
	"encoding/json"
	"net/http"
	"net/url"

	"gin-realword-example/internal/modules/core"
	"gin-realword-example/internal/modules/shared"
)

const (
	githubOauthEntryUrl  = "https://github.com/login/oauth/authorize"
	githubAccessTokenUrl = "https://github.com/login/oauth/access_token"
	githubApiUserUrl     = "https://api.github.com/user"
)

var (
	webHost            string
	githubClientID     string
	githubClientSecret string
)

func init() {
	webHost = core.ConfigStore.GetString(shared.ConfigKeyWebHost)
	githubClientID = core.ConfigStore.GetString(shared.ConfigKeyAuthGithubClientID)
	githubClientSecret = core.ConfigStore.GetString(shared.ConfigKeyAuthGithubClientSecret)
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type User struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}

func BuildOauthEntryUrl(callbackUrl string) (*url.URL, error) {
	redirectUri, err := url.JoinPath(webHost, callbackUrl)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(githubOauthEntryUrl)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Add("client_id", githubClientID)
	query.Add("redirect_uri", redirectUri)
	u.RawQuery = query.Encode()
	return u, nil
}

func GetAccessToken(code string) (*AccessTokenResponse, error) {
	u, err := url.Parse(githubAccessTokenUrl)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Add("client_id", githubClientID)
	query.Add("client_secret", githubClientSecret)
	query.Add("code", code)
	u.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var tokenResponse AccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, err
	}
	return &tokenResponse, nil
}

func GetUser(accessToken string) (*User, error) {
	req, err := http.NewRequest(http.MethodGet, githubApiUserUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
