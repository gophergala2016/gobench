package handler

import (
	"github.com/google/go-github/github"
	"github.com/gophergala2016/gobench/backend/model"
	"github.com/syntaqx/echo-middleware/session"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strings"
	"labix.org/v2/mgo/bson"
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
	packages, err := h.back.Model.Package.GetItemsByIdSlice(user.Packages)
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


func (h *handler) RemoveFromFavPostHandler(c *echo.Context) error {

	s := session.Default(c)
	if s.Get("user") == "" {
		return c.JSON(http.StatusForbidden, "Access denied")
	}
	user, err := h.back.Model.User.GetByLogin(s.Get("user").(string))
	if (err != nil || user == &model.UserRow{}) {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	p, err := h.back.Model.Package.GetItem(c.Form("package"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var packages []bson.ObjectId
	for _, up := range user.Packages {
		if up != p.Id {
			packages = append(packages, up)
		}
	}
	user.Packages = packages
	_, err = h.back.Model.User.UpsertUser(user)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}


func (h *handler) AddToFavPostHandler(c *echo.Context) error {
	s := session.Default(c)
	if s.Get("user") == "" {
		return c.JSON(http.StatusForbidden, "Access denied")
	}
	user, err := h.back.Model.User.GetByLogin(s.Get("user").(string))
	if (err != nil || user == &model.UserRow{}) {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	p, err := h.back.Model.Package.GetItem(c.Form("package"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user.Packages = append(user.Packages, p.Id)
	_, err = h.back.Model.User.UpsertUser(user)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
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
			r := model.RepositoryRow{
				Name:   *repo.Name,
				Url:    strings.Replace(*repo.HTMLURL, "https://", "", -1),
				Engine: model.Git,
			}
			err = h.back.Model.Repository.Add(r)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}