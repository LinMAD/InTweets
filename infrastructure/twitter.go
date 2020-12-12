package infrastructure

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/LinMAD/InTweets/core"
	"github.com/LinMAD/InTweets/domain"
)

// ClientTwitter ...
type ClientTwitter struct {
	owner      string
	twitAPI    *anaconda.TwitterApi
	twitStream *anaconda.Stream
	log        *core.Logger
}

// InitClientTwitter ...
func InitClientTwitter(owner string, tc *domain.TwitterCredential, l *core.Logger) *ClientTwitter {
	anaconda.SetConsumerKey(tc.ConsumerKey)
	anaconda.SetConsumerSecret(tc.ConsumerSecret)

	c := &ClientTwitter{
		twitAPI: anaconda.NewTwitterApi(tc.AccessToken, tc.AccessTokenSecret),
		log:     l,
		owner:   owner,
	}
	c.twitAPI.SetLogger(l)

	l.Infof("Twitter client initialized for owner: %s...", owner)

	return c
}

// StartStream of tweets for interested track
func (t *ClientTwitter) StartStream(track string) bool {
	if t.twitStream != nil {
		t.StopStream()
	}

	t.twitStream = t.twitAPI.PublicStreamFilter(
		url.Values{"track": []string{track}},
	)

	t.log.Infof("Owner %s, started Twitter stream...", t.owner)

	return t.twitStream != nil
}

// StopStream ...
func (t *ClientTwitter) StopStream() {
	if t.twitStream != nil {
		t.log.Infof("Owner %s, stopped Twitter stream...", t.owner)
		t.twitStream.Stop()
	}
}

// FetchTweet from stream
func (t *ClientTwitter) FetchTweet(tweetChan chan string, exitChan chan bool) {
	defer close(tweetChan)

	for {
		select {
		case v := <-t.twitStream.C:
			tweet, ok := v.(anaconda.Tweet)
			if !ok {
				t.log.Errorf(
					"Owner %s, has issue with tweet, unable to cast data: unexpected type of %T must be type of %s",
					t.owner,
					v,
					"anaconda.Tweet",
				)
				return
			}

			t.log.Debugf("%v\n", tweet.Text)
			tweetChan <- tweet.Text
		case <-exitChan:
			return
		}
	}
}
