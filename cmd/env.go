package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"slices"

	"github.com/joho/godotenv"
)

var EnvVars *EnvConfig

type EnvConfig struct {
	Environment        string
	Port               int
	Domain             string
	CookieSecure       bool
	DBUrl              string

	EmailAppPassword   string

	GithubClientId string
	GithubClientSecret string
	GithubRedirectUri  string
}

func NewEnvConfig() (*EnvConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf(".env file not found")
	}

	cfg := &EnvConfig{}
	validEnvs := []string{"development", "testing", "production"}

	environment := os.Getenv("ENVIRONMENT")
	port := os.Getenv("PORT")
	domain := os.Getenv("DOMAIN")
	cookieSecure := os.Getenv("COOKIE_SECURE")
	dbUrl := os.Getenv("DATABASE_URL")
	emailAppPassword := os.Getenv("EMAIL_APP_PASSWORD")
	githubClientId := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	githubRedirectUri := os.Getenv("GITHUB_REDIRECT_URI")

	environment = strings.ToLower(environment)
	isValid := slices.Contains(validEnvs, environment)
	if !isValid {
		return nil, fmt.Errorf("Invalid ENVIRONMENT value: %s", environment)
	}
	cfg.Environment = environment
	if port == "" {
		return nil, fmt.Errorf("PORT environment variable is missing.")
	}
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("Invalid PORT value: %w", err)
	}
	cfg.Port = portNum
	if domain == "" {
		return nil, fmt.Errorf("DOMAIN environment variable is missing.")
	}
	cfg.Domain = domain
	if cookieSecure == "" {
		return nil, fmt.Errorf("COOKIE_SECURE environment variable is missing.")
	}
	cookieSecurityBool, err := strconv.ParseBool(cookieSecure)
	if err != nil {
		return nil, fmt.Errorf("Invalid COOKIE_SECURE value %s", cookieSecure)
	}
	cfg.CookieSecure = cookieSecurityBool
	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is missing.")
	}
	cfg.DBUrl = dbUrl
	if emailAppPassword == "" {
		return nil, fmt.Errorf("EMAIL_APP_PASSWORD environment variable is missing.")
	}
	cfg.EmailAppPassword = emailAppPassword
	if githubClientId == "" {
		return nil, fmt.Errorf("GITHUB_CLIENT_ID environment variable is missing.")
	}
	cfg.GithubClientId = githubClientId
	if githubClientSecret == "" {
		return nil, fmt.Errorf("GITHUB_CLIENT_SECRET environment variable is missing.")
	}
	cfg.GithubClientSecret = githubClientSecret
	if githubRedirectUri == "" {
		return nil, fmt.Errorf("GITHUB_REDIRECT_URI environment variable is missing.")
	}
	cfg.GithubRedirectUri = githubRedirectUri

	return cfg, nil
}
