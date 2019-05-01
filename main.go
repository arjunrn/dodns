package main

import (
	"context"
	"os"

	"github.com/digitalocean/godo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	doTokenEnvName = "DO_TOKEN"
)

type tokenSource struct {
	StaticToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: t.StaticToken,
	}, nil
}

func getToken() string {
	return os.Getenv(doTokenEnvName)
}

func main() {
	staticToken := getToken()
	if staticToken == "" {
		log.Fatalf("Environment variable %s is not set with Digitalocean token", doTokenEnvName)
	}
	ts := &tokenSource{StaticToken: staticToken}
	oauthClient := oauth2.NewClient(context.Background(), ts)
	doclient := godo.NewClient(oauthClient)
	a, _, err := doclient.Account.Get(context.Background())
	if err != nil {
		log.Errorf("Failed to get account")
	}
	log.Infof("Account: %s", a)
}
