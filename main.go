package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ggmolly/ffnf/middlewares"
	"github.com/ggmolly/ffnf/orm"
	"github.com/ggmolly/ffnf/routes"

	"github.com/gofiber/fiber/v2"

	"encoding/json"

	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

var (
	BindAddress = "127.0.0.1"
	Port        = "8000"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v\n", err)
	}
	Port = os.Getenv("PORT")
	BindAddress = os.Getenv("BIND_ADDRESS")
	orm.InitDatabase()
}

func main() {
	engine := html.New("./views", ".html")

	if os.Getenv("MODE") != "production" {
		log.Println("dev mode enabled, reloading templates on each request")
		engine.Reload(true)
	}
	app := fiber.New(fiber.Config{
		Views:        engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		AppName:      "FFNF",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ProxyHeader:  "CF-Connecting-IP",
	})
	app.Use(middlewares.FiberLogger)
	app.Use(middlewares.FiberCompress)

	// Static resources
	app.Static("/static", "./static")

	// Routes
	app.Get("/", routes.Index)

	api := app.Group("/api/v1")
	{
		releasesApi := api.Group("/releases")
		{
			releasesApi.Get("/:n", routes.GetLastReleases)
			releasesApi.Post("/webhook/:secret", routes.ReleaseWebhook)
		}

		// Notices are transformed releases
		noticesApi := api.Group("/notices")
		{
			noticesApi.Get("/after/:id", routes.GetNoticesAfter)
		}
	}

	// Listen on port 8000
	if err := app.Listen(fmt.Sprintf("%s:%s", BindAddress, Port)); err != nil {
		log.Fatalf("Failed to listen on port %s: %v\n", Port, err)
	}
}
