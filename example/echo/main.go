package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/moonrhythm/session"
	"github.com/moonrhythm/session/store"
)

func main() {
	e := echo.New()

	e.Use(echo.WrapMiddleware(session.Middleware(session.Config{
		Store:    new(store.Memory),
		HTTPOnly: true,
		Path:     "/",
		MaxAge:   time.Minute,
	})))

	e.GET("/", handler)

	e.Logger.Fatal(e.Start(":8080"))
}

func handler(c echo.Context) error {
	sess, _ := session.Get(c.Request().Context(), "sess")
	cnt := sess.GetInt("cnt")
	cnt++
	sess.Set("cnt", cnt)
	return c.String(http.StatusOK, fmt.Sprintf("%d views", cnt))
}
