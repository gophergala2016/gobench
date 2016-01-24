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
	if s.Get("just_signup") != nil {
		err = addUserRepos(h, user)
		if err != nil {
			log.Println(err)
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}
		s.Delete("just_signup")
		s.Save()
	}
	packages, err := h.back.Model.Package.GetItemsByIdSlice(user.Repos)
	if err != nil {
		log.Println(err)
	}

	data := struct {
		Packages []model.PackageRow
	}{
		packages,
	}
	return c.Render(http.StatusOK, "dashboard.html", data)

}

func addUserRepos(h *handler, u *model.UserRow) error {

	tokenSource := &TokenSource{
		AccessToken: u.Token,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		return err
	}

	repoList, _, err := client.Repositories.List(*user.Login, &github.RepositoryListOptions{Type: "owner"})
	if err != nil {
		return err
	}

	for _, repo := range repoList {
		if repo.Language != nil && *repo.Language == "Go" {
			pr := &model.PackageRow{
				//Name:   *repo.Name,
				Name:          strings.Replace(*repo.HTMLURL, "https://github.com/", "", -1),
				Url:           *repo.HTMLURL,
				RepositoryUrl: "https://github.com",
				Engine:        model.Git,
			}
			err = h.back.Model.Package.Add(pr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
