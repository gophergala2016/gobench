package handler

import (
	"fmt"
	"time"
	//	"github.com/gorilla/context"
	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
	"net/http"
)

// GithubConfig holds GitHub app credentials
type githubConfig struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type HandlerConfig struct {
	Github githubConfig `json:"github"`
}

type handler struct {
	cfg   *HandlerConfig
	store session.CookieStore
}

func New(cfg *HandlerConfig, store session.CookieStore) handler {
	return handler{cfg: cfg, store: store}
}

func (h *handler) NotFoundHandler(err error, c *echo.Context) {

	if he, ok := err.(*echo.HTTPError); ok {
		if he.Code() == http.StatusNotFound {
			http.Error(c.Response(), "Error 404. Page not found", http.StatusNotFound)
			return
		}
	}

	return
}

func (h *handler) IndexGetHandler(c *echo.Context) error {
	s := session.Default(c)
	fmt.Println(s.Get("user"))
	s.Set("user", time.Now().String())
	s.Save()
	return c.Render(http.StatusOK, "index.html", nil)
}
