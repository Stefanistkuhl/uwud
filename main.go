package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
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

var configHolder atomic.Value

const configPath = "config/config.toml"

func main() {
	conf, loadConfErr := loadConfig(configPath)
	if loadConfErr != nil {
		log.Fatalf("Failed to load configuration: %v", loadConfErr)
	}
	conf = fixImagePath(conf)
	configHolder.Store(conf)

	var lastModTime time.Time
	if info, err := os.Stat(configPath); err == nil {
		lastModTime = info.ModTime()
	}

	go func() {
		for {
			info, err := os.Stat(configPath)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			currentModTime := info.ModTime()

			isAfter := currentModTime.After(lastModTime)
			if isAfter {
				time.Sleep(100 * time.Millisecond)

				newConf, err := loadConfig(configPath)
				if err != nil {
					log.Printf("[ERROR] Failed to load new configuration from %s: %v\n", configPath, err)
				} else {
					newConf = fixImagePath(newConf)
					configHolder.Store(newConf)
					lastModTime = currentModTime
					log.Printf("[INFO] Configuration reloaded successfully. New lastModTime: %s\n", lastModTime.Format(time.RFC3339Nano))
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()

	r := gin.Default()
	r.LoadHTMLGlob("views/dashboard.html")

	r.GET("/healthcheck", handleHealthcheck)
	r.Static("/static", "./static")
	r.GET("/", handleRoot)

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

func handleRoot(c *gin.Context) {
	conf := configHolder.Load().(*Config)
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title":       conf.General.Title,
		"Description": conf.General.Desc,
		"Icon":        conf.General.Icon,
		"Columns":     conf.General.Columns,
		"Sections":    conf.Sections,
	})
}

func handleHealthcheck(c *gin.Context) {
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
		if os.IsTimeout(err) || strings.Contains(err.Error(), "context deadline exceeded") {
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
}

func IsValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

func fixImagePath(conf *Config) *Config {
	if conf.General.Icon != "" && !IsValidURL(conf.General.Icon) {
		conf.General.Icon = "static/images/custom/" + conf.General.Icon
	}
	for i := range conf.Sections {
		section := &conf.Sections[i]
		if section.Icon != "" && !IsValidURL(section.Icon) {
			section.Icon = "static/images/custom/" + section.Icon
		}
		for j := range section.Items {
			item := &conf.Sections[i].Items[j]
			if item.Icon != "" && !IsValidURL(item.Icon) {
				item.Icon = "static/images/custom/" + item.Icon
			}

		}
	}
	return conf

}
