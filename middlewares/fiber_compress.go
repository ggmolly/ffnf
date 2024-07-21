package middlewares

import (
	"github.com/gofiber/fiber/v2/middleware/compress"
)

var (
	FiberCompress = compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
)
