package handler

import (
	"github.com/gophergala2016/gobench/backend"
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
	cfg     *HandlerConfig
	store   session.CookieStore
	back    *backend.Backend
}

func New(cfg *HandlerConfig, store session.CookieStore, b *backend.Backend) handler {
	return handler{cfg: cfg, store: store, back: b}
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
	return c.Render(http.StatusOK, "index.html", nil)
}
