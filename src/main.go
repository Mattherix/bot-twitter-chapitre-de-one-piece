// The majority of the code is from https://tutorialedge.net/golang/writing-a-twitter-bot-golang/
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func main() {
	rand.Seed(time.Now().Unix())
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	client, user, err := getClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
		return
	}

	for chapterNumber := 1; chapterNumber < 1017; chapterNumber++ {
		haveATweet, err := chapterHaveATweet(client, user, chapterNumber)
		if err != nil {
			log.Printf("Error will trying to check if chapter %d have a tweet", chapterNumber)
			log.Println(err)
			return
		}
		if !haveATweet {
			tweet, _, err := tweet(client, chapterNumber)
			if err != nil {
				log.Printf("Error while tweeting on chapter %d\n", chapterNumber)
				log.Println(err)
				return
			} else {
				log.Printf("Tweet: %s\n", tweet.Text)
			}
		}
	}

}

func getClient(creds *Credentials) (*twitter.Client, *twitter.User, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Connected to twitter")
	log.Printf("User's ACCOUNT: %+s\n", user.Name)
	return client, user, nil
}

func chapterHaveATweet(client *twitter.Client, user *twitter.User, chapterNumber int) (bool, error) {
	var search twitter.SearchTweetParams
	search.Query = fmt.Sprintf("from:%s #onepiece%d", user.ScreenName, chapterNumber)
	search.Count = 1
	tweets, _, err := client.Search.Tweets(&search)
	if err != nil {
		return false, err
	}
	if len(tweets.Statuses) != 0 {
		log.Printf("Chapter %d already have a tweet\n", chapterNumber)
		return true, nil
	} else {
		log.Printf("Chapter %d have no tweet", chapterNumber)
		return false, nil
	}
}

func tweet(client *twitter.Client, chapterNumber int) (*twitter.Tweet, *http.Response, error) {
	image := os.Getenv("IMAGE_LINK")
	messages := []string{
		"Cette semaine encore gOda à frapper",
		"Juste exceptionnel",
		"Cette semaine encore je croie qu'on peu le dire",
		"On a véccue un chapitre incroyable",
		"Quel Masterclass",
		"Là juste respect",
		"Merci Oda",
	}
	message := fmt.Sprintf("%s #onepiece%d %s", messages[rand.Intn(len(messages))], chapterNumber, image)
	tweet, resp, err := client.Statuses.Update(message, nil)
	return tweet, resp, err
}
