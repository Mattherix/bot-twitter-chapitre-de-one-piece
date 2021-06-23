// The majority of the code is from https://tutorialedge.net/golang/writing-a-twitter-bot-golang/
package main

import (
	"fmt"
	"log"
	"math/rand"
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

	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	client, err := getClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
	}

	tweet(client, 1)
}

func getClient(creds *Credentials) (*twitter.Client, error) {
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
		return nil, err
	}

	log.Println("Connected to twitter")
	log.Printf("User's ACCOUNT: %+s\n", user.Name)
	return client, nil
}

func tweet(client *twitter.Client, chapterNumber int) {
	image := os.Getenv("IMAGE_LINK")
	rand.Seed(time.Now().Unix())
	messages := []string{
		"Cette semaine encore gOda à frapper",
		"Juste exceptionnel",
		"Cette semaine encore je croie qu'on peu le dire",
		"On a véccue un chapitre excetionnel",
		"Quel Masterclass",
		"Là juste respect",
		"Merci Oda",
	}
	message := fmt.Sprintf("%s #onepiece%d %s", messages[rand.Intn(len(messages))], chapterNumber, image)
	_, _, err := client.Statuses.Update(message, nil)
	if err != nil {
		log.Printf("Error while tweeting on chapter %d\n", chapterNumber)
		log.Println(err)
	}
}
