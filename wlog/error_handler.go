package wlog

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func MiddlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		var sb strings.Builder

		sb.WriteString("method: ")
		sb.WriteString(c.Request().Method)
		sb.WriteString(" | ")
		sb.WriteString("uri: ")
		sb.WriteString(c.Request().URL.String())
		sb.WriteString(" | ")
		sb.WriteString("ip: ")
		sb.WriteString(c.Request().RemoteAddr)

		I("middlewareLogging", "incomming request", sb.String())

		return next(c)
	}
}

func ErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if ok {
		report.Message = fmt.Sprintf("http error %d - %v", report.Code, report.Message)
	} else {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	D("errorHandler", "report", report.Message.(string))
	c.JSON(report.Code, InternalError())
}
