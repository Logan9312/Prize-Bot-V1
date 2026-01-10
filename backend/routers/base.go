package routers

import (
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/logan9312/discord-auction-bot/api"
)

func HealthCheck() {
	StartServer(nil)
}

func StartServer(botSession *discordgo.Session) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	// Register API routes
	api.RegisterRoutes(e, botSession)

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
