package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	// "testapp/db"
	// "testapp/utils"
	cmd "github.com/IAmRiteshKoushik/pulse/cmd"
)

func GithubLoginHandler(c *gin.Context) {
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		cmd.EnvVars.GithubClientId, cmd.EnvVars.GithubRedirectUri)

	fmt.Println(redirectURL)

	c.Redirect(http.StatusSeeOther, redirectURL)
}

func GithubCallbackHandler(c *gin.Context) {
	code := c.Query("code")

	githubAccessToken := getGithubAccessToken(code)

	githubData := getGithubData(githubAccessToken)

	loggedinHandler(c, githubData)
}

func loggedinHandler(c *gin.Context, githubData string) {
	if githubData == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED!"})
		return
	}

	var ghUser struct {
		Login string `json:"login"`
	}
	if err := json.Unmarshal([]byte(githubData), &ghUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse GitHub data"})
		return
	}

	// var email string
	// err := db.DB.QueryRow("SELECT email FROM users WHERE github_username = $1", ghUser.Login).Scan(&email)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "GitHub user not registered"})
	// 	return
	// }

	// Return the GitHub data as a pretty-printed JSON response
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	c.Data(http.StatusOK, "application/json", prettyJSON.Bytes())
}

func getGithubAccessToken(code string) string {
	clientID := cmd.EnvVars.GithubClientId
	clientSecret := cmd.EnvVars.GithubClientSecret
	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var ghresp struct {
		AccessToken string `json:"access_token"`
	}
	json.Unmarshal(respbody, &ghresp)

	return ghresp.AccessToken
}

func getGithubData(accessToken string) string {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	// Github data here
	return string(respbody)
}
