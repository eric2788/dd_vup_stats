package middlewares

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	def := cors.DefaultConfig()
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"https://ddstats.ericlamm.xyz",
			"https://ddstats.pages.dev",
			os.Getenv("DEV_HOST"),
		},
		AllowWebSockets: true,
		AllowMethods:    def.AllowMethods,
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"Origin",
			"Content-Length",
		},
	})
}
