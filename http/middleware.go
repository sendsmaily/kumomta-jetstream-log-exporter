package http

import "github.com/labstack/echo/v4"

// ErrorLogger is a simple error logging middleware for Echo.
func ErrorLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				c.Logger().Error(err)
			}

			return err
		}
	}
}
