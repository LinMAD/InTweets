package main

import (
	"flag"
	"fmt"

	"github.com/LinMAD/InTweets/core"
	"github.com/LinMAD/InTweets/domain"
	"github.com/LinMAD/InTweets/infrastructure"
	"github.com/sirupsen/logrus"
)

func main() {
	// App configuration
	host := flag.String("http_host", "localhost", "Application host")
	port := flag.String("http_port", "8080", "Application port")
	isDebug := flag.Bool("debug_mode", false, "Debug mode enabled")
	flag.Parse()

	// Prepare dependencies
	apiServer := Init(isDebug)

	// Execute
	apiServer.Run(*host + ":" + *port)
}

// Init prepares dependencies before execution of api server
func Init(isDebug *bool) *infrastructure.ServerAPI {
	log := &core.Logger{Logger: logrus.New()}
	if *isDebug {
		log.Level = logrus.DebugLevel
	}

	c, err := DispatchTwitterCredentials(log)
	if err != nil {
		log.Fatal(err)
	}

	apiServer := infrastructure.InitServerAPI(c, log)
	apiServer.LoadRouteHandlers()

	return apiServer
}

// DispatchTwitterCredentials ...
func DispatchTwitterCredentials(log *core.Logger) (*domain.TwitterCredential, error) {
	twitKeys := map[string]string{
		"TWITTER_CONSUMER_KEY":        "",
		"TWITTER_CONSUMER_SECRET":     "",
		"TWITTER_ACCESS_TOKEN":        "",
		"TWITTER_ACCESS_TOKEN_SECRET": "",
	}

	for k := range twitKeys {
		log.Debugf("Looking for Twitter env variable %s", k)

		v, err := core.GetEnvVar(k)
		if err != nil {
			return nil, fmt.Errorf("problem with Twitter auth key, error: %v", err.Error())
		}

		if len(v) == 0 {
			return nil, fmt.Errorf("empty value in Twitter auth key %s", k)
		}

		twitKeys[k] = v
	}

	twitAuth := &domain.TwitterCredential{
		ConsumerKey:       twitKeys["TWITTER_CONSUMER_KEY"],
		ConsumerSecret:    twitKeys["TWITTER_CONSUMER_SECRET"],
		AccessToken:       twitKeys["TWITTER_ACCESS_TOKEN"],
		AccessTokenSecret: twitKeys["TWITTER_ACCESS_TOKEN_SECRET"],
	}

	return twitAuth, nil
}
