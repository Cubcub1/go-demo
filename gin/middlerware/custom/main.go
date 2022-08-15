package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("example", "12345")

		c.Next()

		latency := time.Since(t)
		log.Print(latency)

		status := c.Writer.Status()
		log.Panicln(status)
	}
}

func main() {
	r := gin.New()
	r.Use(Logger())

	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)

		log.Panicln(example)

	})

	r.Run(":8080")
}
