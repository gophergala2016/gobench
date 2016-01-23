package handler

import (
	"github.com/google/go-github/github"
	"github.com/gophergala2016/gobench/backend/model"
	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"log"
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
		log.Println(err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		log.Println(err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	userLogin := *user.Login
	userToken := token.AccessToken

	u := model.UserRow{
		Login:     userLogin,
		Token:     userToken,
		AvatarURL: *user.AvatarURL,
	}
	ci, err := h.back.Model.User.CreateUser(&u)
	if err != nil {
		log.Println(err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	if (ci.Updated == 0) {
		s.Set("just_signup", true)
	}
	s.Set("user", userLogin)
	s.Set("token", userToken)
	s.Save()

	return c.Redirect(http.StatusFound, "/dashboard")
}
