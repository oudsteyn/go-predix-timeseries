package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/oudsteyn/go-predix-timeseries/api"
	"github.com/oudsteyn/go-predix-timeseries/dataquality"
	"github.com/oudsteyn/go-predix-timeseries/measurement"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var ingestUrl = flag.String("ingestUrl", "", "ingestion endpoint")
var zoneId = flag.String("zoneId", "", "Prefix-Zone-Id")
var uaaIssuerId = flag.String("uaaIssuerId", "", "UAA IssuerId")
var clientId = flag.String("clientId", "", "Client ID")
var clientSecret = flag.String("clientSecret", "", "Client Secret")

func main() {
	flag.Parse()
	if flag.NFlag() < 5 {
		flag.Usage()
		os.Exit(1)
	}
	conf := clientcredentials.Config{
		ClientID:     *clientId,
		ClientSecret: *clientSecret,
		TokenURL:     *uaaIssuerId,
	}
	t, err := conf.Token(oauth2.NoContext)
	if err != nil {
		fmt.Printf("Auth failed: %s", err)
		os.Exit(1)
	}
	a := api.Ingest(*ingestUrl, t.AccessToken, *zoneId)
	if a == nil {
		fmt.Println("Connection to Time Series service failed")
		os.Exit(1)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	for _ = range time.Tick(time.Second) {
		m := a.IngestMessage()
		if m != nil {
			t, _ := m.AddTag("test_tag")
			v := randInt(0, 100)
			t.AddDatapoint(measurement.Int(v), dataquality.Good)
			e := m.Send()
			if e != nil {
				fmt.Printf("Ingesting failed %s\n", e)
			} else {
				fmt.Printf("Send %d\n", v)
			}
		} else {
			fmt.Println("Creating ingest message failed")
			os.Exit(1)
		}
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
