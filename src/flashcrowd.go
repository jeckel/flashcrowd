package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"log"
	"net/url"
	"os"
	"time"
)

// TwitterConnect is a structure containing identifier to connect to the Twitter API
type TwitterConnect struct {
	consumerKey, consumerSecret, twitterApiKey, twitterApiKeySecret string
}

const (
	LEVEL_LIMIT = 15
)

var api *anaconda.TwitterApi

func getTwitterApiClient(twitterAuthKey TwitterConnect) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(twitterAuthKey.consumerKey)
	anaconda.SetConsumerSecret(twitterAuthKey.consumerSecret)
	return anaconda.NewTwitterApi(twitterAuthKey.twitterApiKey, twitterAuthKey.twitterApiKeySecret)
}

func init() {
	log.Print("Loading setup...")

	twitterConnect := TwitterConnect{
		consumerKey:         os.Getenv("CONSUMER_KEY"),
		consumerSecret:      os.Getenv("CONSULER_KEY_SECRET"),
		twitterApiKey:       os.Getenv("ACCESS_TOKEN"),
		twitterApiKeySecret: os.Getenv("ACCESS_TOKEN_SECRET"),
	}
	api = getTwitterApiClient(twitterConnect)
}

func main() {
	log.Print("Starting...")

	currentLevel := 0
	// @todo : use go routines, loop over an array, be carefull of ratelimit
	currentLevel += getBuzzLevel("#RERC", 30*time.Minute)
	currentLevel += getBuzzLevel("@RERC_SNCF", 30*time.Minute)

	log.Printf("Current overall level: %d\n", currentLevel)

	if currentLevel > LEVEL_LIMIT {
		log.Printf("Current level is over %d, publishing alert on Twitter", LEVEL_LIMIT)
		v := url.Values{}
		api.PostTweet(fmt.Sprintf("Current RER C Buzz level: %d", currentLevel), v)
	}

	log.Print("Exiting...")
}

// getBuzzLevel return the current buzz/crowd level for the query string
func getBuzzLevel(query string, interval time.Duration) int {
	now := time.Now()
	currentLevel := 0

	v := url.Values{}
	v.Set("count", "30")
	v.Set("result_type", "recent")
	searchResult, err := api.GetSearch(query, v)
	if err != nil {
		log.Fatal(err)
	}

	for i, tweet := range searchResult.Statuses {
		createdAt, err := tweet.CreatedAtTime()
		if err != nil {
			log.Printf("Error with createdAt for tweet #%d", i)
			continue
		}
		delta := now.Sub(createdAt)
		if delta < interval {
			// add one point for tweet during the interval
			currentLevel++
			// add the retweet count as additionnal points
			currentLevel += tweet.RetweetCount
			//log.Printf("%ds ago: (%s) RT:%d ==> %s", (delta / time.Second), tweet.CreatedAt, tweet.RetweetCount, tweet.Text)
		}
	}
	log.Printf("Current level for \"%s\": %d\n", query, currentLevel)
	return currentLevel
}
