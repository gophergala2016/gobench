package handler

import (
	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
	"net/http"
	"fmt"
)

// GithubConfig holds GitHub app credentials
type githubConfig struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type HandlerConfig struct  {
	Github githubConfig `json:"github"`
}

type handler struct {
	cfg *HandlerConfig
}

func New(cfg *HandlerConfig) handler  {
	return handler{cfg: cfg}
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
	return c.Render(http.StatusOK, "index.html", nil)
}
