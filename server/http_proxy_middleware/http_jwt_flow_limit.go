package http_proxy_middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackjie2016/gateway/server/dao"
	"github.com/jackjie2016/gateway/server/middleware"
	"github.com/jackjie2016/gateway/server/public"
	"github.com/pkg/errors"
)

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.App)
		if appInfo.Qps > 0 {
			clientLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowAppPrefix+appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.Qps))
			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				middleware.ResponseError(c, 5002, errors.New(fmt.Sprintf("%v flow limit %v", c.ClientIP(), appInfo.Qps)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
