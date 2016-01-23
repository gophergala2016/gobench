package frontend

import (
	"github.com/labstack/echo"
	"net/http"
)

func notFoundHandler(err error, c *echo.Context) {

	if he, ok := err.(*echo.HTTPError); ok {
		if he.Code() == http.StatusNotFound {
			http.Error(c.Response(), "Error 404. Page not found", http.StatusNotFound)
			return
		}
	}

	return
}

func indexGetHandler(c *echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
