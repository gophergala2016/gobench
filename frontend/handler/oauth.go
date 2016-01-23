package handler

import (
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"github.com/google/go-github/github"
	"github.com/syntaqx/echo-middleware/session"
	"net/http"
	"github.com/labstack/gommon/log"
)
var oauthConf = &oauth2.Config{
	Scopes:       []string{"user:email", "repo"},
	Endpoint:     githuboauth.Endpoint,
}

func (h *handler) OauthRequestHandler(c *echo.Context) error {
	s := session.Default(c)
	if (s.Get("user") != nil) {
		return c.Redirect(http.StatusFound, "/dashboard")
	}

	oauthConf.ClientID = h.cfg.Github.ClientId
	oauthConf.ClientSecret = h.cfg.Github.ClientSecret
	oauthStateString := "random_string"
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *handler) OauthCallbackHandler(c *echo.Context) error {
	s := session.Default(c)
	code := c.Query("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Error(err)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		log.Error(err)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	s.Set("user", user)
	s.Save()

	return c.Redirect(http.StatusFound, "/dashboard")
}
