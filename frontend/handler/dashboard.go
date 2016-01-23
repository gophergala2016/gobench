package handler

import (
	"github.com/google/go-github/github"
	"github.com/gophergala2016/gobench/backend/model"
	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strings"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func (h *handler) DashboardGetHandler(c *echo.Context) error {
	s := session.Default(c)
	if s.Get("user") == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	user, err := h.back.Model.User.GetByLogin(s.Get("user").(string))
	if (err != nil || user == &model.UserRow{}) {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	var repos []model.RepositoryRow
	if c.Query("just_signup") != "" {
		tokenSource := &TokenSource{
			AccessToken: user.Token,
		}
		oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
		client := github.NewClient(oauthClient)
		user, _, err := client.Users.Get("")
		if err != nil {
			log.Println(err)
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		repoList, _, err := client.Repositories.List(*user.Login, &github.RepositoryListOptions{Type: "owner"})
		if err != nil {
			log.Println(err)
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		for _, repo := range repoList {
			if repo.Language != nil && *repo.Language == "Go" {
				r := model.RepositoryRow{
					Name:   *repo.Name,
					Url:    strings.Replace(*repo.HTMLURL, "https://", "", -1),
					Engine: model.Git,
				}
				repos = append(repos, r)
			}
		}
	} else {
		repos, err = h.back.Model.Repository.Items(user.Repos)
		if err != nil {
			log.Println(err)
		}
	}

	data := struct {
		Repos []model.RepositoryRow
	}{
		repos,
	}
	return c.Render(http.StatusOK, "dashboard.html", data)

}
