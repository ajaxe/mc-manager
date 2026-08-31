package models

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/labstack/echo/v5"
)

var client *http.Client
var once sync.Once

func authHttpClient() *http.Client {
	once.Do(func() {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client = &http.Client{Transport: t}
	})
	return client
}

func NewAppContext(c *echo.Context, l *slog.Logger) *AppContext {
	if l == nil {
		l = slog.Default()
	}
	return &AppContext{
		Context: c,
		logger:  l,
	}
}

func GetAppContext(c *echo.Context) *AppContext {
	if v := c.Get("appContext"); v != nil {
		if ac, ok := v.(*AppContext); ok {
			return ac
		}
	}
	return NewAppContext(c, slog.Default())
}

type AppContext struct {
	*echo.Context
	logger  *slog.Logger
	isAdmin *bool
}

func (c *AppContext) IsAdmin() bool {
	if c.isAdmin == nil {
		isAdmin, err := c.checkIdentity()
		if err != nil {
			c.logger.Error("is_admin check failed", "error", err)
		}
		c.isAdmin = &isAdmin
	}
	return *c.isAdmin
}

// checkIdentity invokes the introspect endpoint to check if the authenticated
// user is in the admin user list.
func (c *AppContext) checkIdentity() (isAdmin bool, err error) {
	cfg := config.LoadAppConfig()
	if cfg.Server.AuthServerEnabled == false {
		isAdmin = cfg.IsDev
		c.logger.Info("auth server not enabled: isAdmin is false")
		return
	}
	u, err := cfg.AuthIntrospectURL()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 100*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return
	}

	authCkVal := c.authCookieValueInternal(cfg)
	req.AddCookie(&http.Cookie{
		Name:  cfg.Server.AuthCookieName,
		Value: authCkVal,
	})
	c.logger.Info("setting auth cookie value", "auth_cookie", authCkVal)

	c.logger.Info("invoking auth introspect url", "url", u)
	res, err := authHttpClient().Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		err = echo.NewHTTPError(res.StatusCode, string(b))
		return
	}
	c.logger.Info("introspect response", "response", string(b))
	d := IntrospectResult{}
	err = json.Unmarshal(b, &d)
	if err != nil {
		return
	}

	if d.Active == false {
		isAdmin = false
		return
	}

	isAdmin = slices.Contains(cfg.AdminUsers(), d.Username)

	if isAdmin == false {
		c.logger.Info("user is not in admin list", "username", d.Username, "admin_users", cfg.AdminUsers())
	}

	return
}

func (c *AppContext) AuthCookieValue() (v string) {
	cfg := config.LoadAppConfig()
	return c.authCookieValueInternal(cfg)
}
func (c *AppContext) authCookieValueInternal(cfg config.AppConfig) (v string) {
	ck := c.Request().Cookies()
	slog.Info("checking for login cookie", "auth_cookie_name", cfg.Server.AuthCookieName)
	for _, c := range ck {
		if c.Name == cfg.Server.AuthCookieName {
			v = c.Value
			break
		}
	}
	return
}
