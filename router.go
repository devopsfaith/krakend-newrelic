package metrics

import (
	"fmt"

	"github.com/gin-gonic/gin"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

var errNoApp = fmt.Errorf("No NewRelic app defined")

func Middleware() (gin.HandlerFunc, error) {
	if app == nil {
		return emptyMW, errNoApp
	}

	return nrgin.Middleware(*app), nil
}

func emptyMW(c *gin.Context) {
	c.Next()
}
