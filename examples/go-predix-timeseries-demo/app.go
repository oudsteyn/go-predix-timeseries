package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/oudsteyn/go-predix-timeseries/api"
	"github.com/oudsteyn/go-predix-timeseries/query"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	var port, authToken string
	var a *api.Api

	if appEnv, err := cfenv.Current(); err != nil {
		port = ":8080"
	} else {
		port = fmt.Sprintf(":%d", appEnv.Port)
		if uaaService, err := appEnv.Services.WithLabel("predix-uaa"); err == nil {
			conf := clientcredentials.Config{
				ClientID:     os.Getenv("CLIENT_ID"),
				ClientSecret: os.Getenv("CLIENT_SECRET"),
				TokenURL:     uaaService[0].Credentials["issuerId"].(string),
			}
			t, err := conf.Token(oauth2.NoContext)
			if err != nil {
				log.Printf("Auth failed: %s", err)
			}
			authToken = t.AccessToken
		}
		if tsService, err := appEnv.Services.WithLabel("predix-timeseries"); err == nil {
			query := tsService[0].Credentials["query"].(map[string]interface{})
			u, _ := url.Parse(query["uri"].(string))
			queryUrl := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
			zoneId := query["zone-http-header-value"].(string)
			a = api.Query(queryUrl, authToken, zoneId)
		}
	}
	r := gin.Default()
	r.LoadHTMLFiles("./static/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/temp", func(c *gin.Context) {
		if a != nil {
			r, e := a.Query().LatestDatapoints(query.Tag("test_tag"))
			if e == nil {
				if r.Tags[0].Stats.RawCount > 0 {
					c.JSON(200, gin.H{"temp": r.Tags[0].Results[0].Values[0].Measure})
				}
			} else {
				c.JSON(500, gin.H{"status": e.Error()})
			}
		} else {
			c.JSON(500, gin.H{"status": "Connection to Time Series service failed"})
		}
	})
	r.Run(port)
}
