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
		packages, err := addUserRepos(h, user)
		if err != nil {
			log.Println(err)
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}
		var up []bson.ObjectId
		for _, p := range packages {
			up = append(up, p.Id)
		}
		user.Packages = up
		_, err = h.back.Model.User.UpsertUser(user)
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


func addUserRepos(h *handler, u *model.UserRow) ([]*model.PackageRow, error) {

	tokenSource := &TokenSource{
		AccessToken: u.Token,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, err
	}

	repoList, _, err := client.Repositories.List(*user.Login, &github.RepositoryListOptions{Type: "owner"})
	if err != nil {
		return nil, err
	}

	packages := make([]*model.PackageRow, 0)
	for _, repo := range repoList {
		if repo.Language != nil && *repo.Language == "Go" {
			pr := &model.PackageRow{
				//Name:   *repo.Name,
				Name:          strings.Replace(*repo.HTMLURL, "https://github.com/", "", -1),
				Author:        u.Login,
				Url:           *repo.HTMLURL,
				Description:   *repo.Description,
				RepositoryUrl: "https://github.com",
				Engine:        model.Git,
				Tags:          getRepoTags(*user.Login, *repo.Name, client),
			}
			newPr, err := h.back.Model.Package.Add(pr)
			if err != nil {
				return nil, err
			}
			packages = append(packages, newPr)
		}
	}

	return packages, nil
}

func getRepoTags(user string, repoName string, client *github.Client) []model.RepositoryTag {
	tags := make([]model.RepositoryTag,0,1)
	githubTags,_,_ := client.Repositories.ListTags(user,repoName,nil)
	for _,v := range githubTags {
		tags = append(tags, model.RepositoryTag{Name: *v.Name, Zip: *v.ZipballURL, Tar: *v.TarballURL, Commit: *v.Commit.URL})
	}
	return tags
}