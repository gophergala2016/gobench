package frontend

import (
	"fmt"
	"github.com/gophergala2016/gobench/backend"
	"github.com/gophergala2016/gobench/frontend/handler"
	"github.com/gorilla/context"
	"github.com/hydrogen18/stoppableListener"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/syntaqx/echo-middleware/session"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

// Config holds frontend configuration params
type Config struct {
	Host             string                `json:"host"`
	Port             int                   `json:"port"`
	TemplateFolder   string                `json:"templateFolder`
	AssetFolder      string                `json:"assetFolder`
	SessionSecretKey string                `json:"sessionSecretKey`
	HandlerCfg       handler.HandlerConfig `json:"handler"`
}

//  Frontend provides single point of access to fronend layer
type Frontend struct {
	log    *log.Logger
	back   *backend.Backend
	cfg    *Config
	router *echo.Echo
	store  session.CookieStore
}

// New creates Frontend instance, initialises routers
func New(cfg *Config, l *log.Logger, b *backend.Backend) (*Frontend, error) {
	f := &Frontend{log: l, back: b, cfg: cfg, router: echo.New(), store: session.NewCookieStore([]byte(cfg.SessionSecretKey))}

	f.router.SetRenderer(f)
	f.router.HTTP2(false)

	// TODO, write logs into file (probable logwriter can be used)
	f.router.SetLogOutput(os.Stdout)

	f.router.Use(mw.Logger())
	f.router.Use(mw.Recover())
//	f.router.Use(frontendMiddleware())
	f.router.Use(session.Sessions("ESESSION", f.store))

	f.router.Static("/css", path.Join(f.cfg.AssetFolder, "/css"))
	f.router.Static("/img", path.Join(f.cfg.AssetFolder, "/img"))
	f.router.Static("/js", path.Join(f.cfg.AssetFolder, "/js"))
	f.router.Static("/fonts", path.Join(f.cfg.AssetFolder, "/fonts"))
	f.router.Favicon(path.Join(f.cfg.AssetFolder, "/img/favicon.ico"))

	h := handler.New(&f.cfg.HandlerCfg, f.store, f.back)
	f.router.SetHTTPErrorHandler(h.NotFoundHandler)
	f.router.Get("/", h.IndexGetHandler)
	f.router.Get("/search",h.SearchbackHandler)
	f.router.Get("/oauth", h.OauthRequestHandler)
	f.router.Get("/oauth/callback", h.OauthCallbackHandler)
	f.router.Get("/dashboard", h.DashboardGetHandler)

	return f, nil
}

// Start launches gracefull HTTP listener
func (f *Frontend) Start() error {

	listenAddr := fmt.Sprintf("%s:%d", f.cfg.Host, f.cfg.Port)
	originalListener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	sl, err := stoppableListener.New(originalListener)
	if err != nil {
		return err
	}

	server := http.Server{Handler: context.ClearHandler(f.router)}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		server.Serve(sl)
	}()

	f.log.Println("Start serving HTTP requests at ", listenAddr)
	select {
	case signal := <-stop:
		f.log.Println("Got signal: ", signal)
	}
	f.log.Println("Stopping listener")
	sl.Stop()
	f.log.Println("Waiting on server")
	wg.Wait()

	return nil
}

// Render implements Echo's Renderer interface
func (f *Frontend) Render(w io.Writer, name string, data interface{}) error {

	// TODO: implement templates caching
	tmpl, err := template.ParseFiles(path.Join(f.cfg.TemplateFolder, "layout.html"), path.Join(f.cfg.TemplateFolder, name))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Template reading error. Details: %s", name, err.Error()))
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Template rendering error. Details: %s", name, err.Error()))
	}

	return nil
}

//func frontendMiddleware() echo.MiddlewareFunc {
//	return func(h echo.HandlerFunc) echo.HandlerFunc {
//		return func(c *echo.Context) error {
//			s := session.Default(c)
//			user := s.Get("user")
//			c.Set("user", user)
//
//			return nil
//		}
//	}
//}