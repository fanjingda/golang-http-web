package main

import (
	"ghw/ghw"
	"net/http"
)

func main() {
	r := ghw.New()
	r.GET("/", func(c *ghw.Context) {
		c.JSON(http.StatusOK, ghw.H{
			"msg":     "200",
			"context": "hello",
		})
	})

	r.Run(":9999")
}