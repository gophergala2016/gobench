package handler


import (
  //  "github.com/gophergala2016/gobench/backend/model"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
	"github.com/gophergala2016/gobench/backend/model"
)

func (h *handler) SearchbackHandler(c *echo.Context) error {
	searchq := c.Query("search")
	packages,err := h.back.Model.Package.GetItems(searchq)
	  if err != nil {
		  log.Error(err)
		  return c.Redirect(http.StatusTemporaryRedirect, "/")
	  }
	if len(searchq) > 50 && len(searchq) < 4 {
		return c.Redirect(http.StatusBadRequest,"/")
	}
	data := struct {
		Packages []model.PackageRow
	}{
		packages,
	}
	return c.Render(http.StatusOK, "search.html", data)
}
