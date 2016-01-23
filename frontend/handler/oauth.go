package handler

import (
	"github.com/google/go-github/github"
	"github.com/gophergala2016/gobench/backend/model"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/syntaqx/echo-middleware/session"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"net/http"
)

var oauthConf = &oauth2.Config{
	Scopes:   []string{"user:email", "repo"},
	Endpoint: githuboauth.Endpoint,
}

func (h *handler) OauthRequestHandler(c *echo.Context) error {
	s := session.Default(c)
	if s.Get("user") != nil {
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
	s.Set("token", token)
	s.Save()

	u := model.UserRow{
		Login:     *user.Login,
		Token:     token.AccessToken,
		AvatarURL: *user.AvatarURL,
	}
	err = h.back.Model.User.CreateUser(&u)
	if err != nil {
		log.Error(err)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	return c.Redirect(http.StatusFound, "/dashboard")
}
