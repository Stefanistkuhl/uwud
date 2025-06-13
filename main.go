package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Items []item `toml:item`
}

type item struct {
	Name   string `toml:name`
	Url    string `toml:url`
	Desc   string `toml:desc`
	Health string `toml:health`
	Icon   string `toml:Icon`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("views/*.html")

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"Title": "My Dashboard",
			"Items": "erm",
		})
	})

	log.Println("Server running on :8080")
	r.Run(":8080")
}

func loadConfig(path string) (*Config, error) {

}
