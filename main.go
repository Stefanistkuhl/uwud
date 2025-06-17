package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	General  GeneralSettings `toml:"general"`
	Sections []Section       `toml:"sections"`
}

type GeneralSettings struct {
	Title   string `toml:"title"`
	Desc    string `toml:"desc"`
	Icon    string `toml:"icon"`
	Columns int    `toml:"columns"`
}

type Section struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Icon        string `toml:"icon"`
	Items       []Item `toml:"items"`
}

type Item struct {
	Name   string `toml:"name"`
	URL    string `toml:"url"`
	Desc   string `toml:"desc"`
	Health string `toml:"health"`
	Icon   string `toml:"icon"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("views/*.html")
	conf, err := loadConfig("config.toml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	r.GET("/healthcheck", func(c *gin.Context) {
		rawURL := c.Query("url")
		decodedURL, err := url.QueryUnescape(rawURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
			return
		}

		parsedURL, err := url.Parse(decodedURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
			return
		}

		var client *http.Client
		timeout := 15 * time.Second
		if parsedURL.Scheme == "https" {
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
				Timeout: timeout,
			}
		} else {
			client = &http.Client{
				Timeout: timeout,
			}
		}

		req, err := http.NewRequest("GET", decodedURL, nil)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			if os.IsTimeout(err) || (err != nil && err.Error() == "context deadline exceeded") {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "timeout", "message": "Health check timed out"})
			} else {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			}
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{"status": "up"})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "code": resp.StatusCode})
		}
	})

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"Title":       conf.General.Title,
			"Description": conf.General.Desc,
			"Icon":        conf.General.Icon,
			"Columns":     conf.General.Columns,
			"Sections":    conf.Sections,
		})
	})

	log.Println("Server running on :8080")
	r.Run(":8080")
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TOML: %w", err)
	}
	return &config, nil
}
