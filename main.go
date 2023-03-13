package main

import (
	"ghw/ghw"
	"net/http"
)

func main() {
	r := ghw.New()
	v1 := r.Group("/v1")
	v1.Use()
	v1.GET("/", func(c *ghw.Context) {
		c.JSON(http.StatusOK, ghw.H{
			"msg":     "200",
			"context": "hello",
		})
	})

	r.Run(":9999")
}
