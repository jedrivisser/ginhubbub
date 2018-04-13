package main

import (
	"encoding/xml"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jedrivisser/ginhubbub"
)

// Feed object in medium post
type Feed struct {
	Entries []Entry `xml:"entry"`
}

// Entry object in medium post
type Entry struct {
	ID    string `xml:"videoId"`
	Title string `xml:"title"`
}

var callbackURL = flag.String("callbackURL", "", "Callback URL where the HUB can reach you")

func main() {
	flag.Parse()

	if *callbackURL == "" {
		log.Println("Need to set callbackURL flag")
		os.Exit(1)
	}

	log.Println("Youtube Watcher Started, callbackURL: " + *callbackURL)

	topic := "https://www.youtube.com/xml/feeds/videos.xml?channel_id=UC99lkbVG8I5hRSZa4FD8zgw"

	subscriptions := make(map[string]func(string, []byte))
	subscriptions[topic] = handleYoutubeFeed

	server := ginhubbub.NewServer(subscriptions)
	client := ginhubbub.NewClient("https://pubsubhubbub.appspot.com", *callbackURL, "Youtube Test")

	for topic := range subscriptions {
		client.Subscribe(topic)
	}

	go server.Engine().Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	for topic := range subscriptions {
		server.RemoveSubsciption(topic)
		client.Unsubscribe(topic)
	}

	time.Sleep(time.Second * 5)
}

func handleYoutubeFeed(contentType string, body []byte) {
	var feed Feed
	xmlError := xml.Unmarshal(body, &feed)

	if xmlError != nil {
		log.Printf("XML Parse Error %v", xmlError)

	} else {
		for _, entry := range feed.Entries {
			log.Printf("%s <https://youtu.be/%s>", entry.Title, entry.ID)
		}
	}
}
