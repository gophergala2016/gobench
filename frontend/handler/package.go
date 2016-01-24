package handler

import (
	"github.com/gophergala2016/gobench/backend/model"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strings"
)

func (h *handler) PackageGetHandler(c *echo.Context) error {

	pName := strings.Replace(c.Request().RequestURI, "/p/", "", -1)
	p, err := h.back.Model.Package.GetItem(pName)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	data := struct {
		Package model.PackageRow
	}{
		p,
	}
	return c.Render(http.StatusOK, "package.html", data)
}
