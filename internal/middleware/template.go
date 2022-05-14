package mw

import (
	"github.com/labstack/echo/v4"
	"zenith/internal/view"
)

func TemplateMW(debug bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if debug {
				c.Echo().Renderer = view.NewTemplate()
			}
			return next(c)
		}
	}
}
