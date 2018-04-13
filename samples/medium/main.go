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
	Status  string  `xml:"status>http"`
	Entries []Entry `xml:"entry"`
}

// Entry object in medium post
type Entry struct {
	URL     string `xml:"id"`
	Title   string `xml:"title"`
	Summary string `xml:"summary"`
	Author  Author `xml:"author"`
}

// Author object in medium post
type Author struct {
	Name string `xml:"name"`
}

var callbackURL = flag.String("callbackURL", "", "Callback URL where the HUB can reach you")

func main() {
	flag.Parse()

	if *callbackURL == "" {
		log.Println("Need to set callbackURL flag")
		os.Exit(1)
	}

	log.Println("Medium Story Watcher Started, callbackURL: " + *callbackURL)

	topic := "https://medium.com/feed/latest"

	subscriptions := make(map[string]func(string, []byte))
	subscriptions[topic] = handleMediumFeed

	server := ginhubbub.NewServer(subscriptions)
	client := ginhubbub.NewClient("http://medium.superfeedr.com", *callbackURL, "Medium Test")

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

func handleMediumFeed(contentType string, body []byte) {
	// Print whole body
	// log.Printf(string(body))
	var feed Feed
	xmlError := xml.Unmarshal(body, &feed)

	if xmlError != nil {
		log.Printf("XML Parse Error %v", xmlError)

	} else {
		for _, entry := range feed.Entries {
			log.Printf("%s by %s (%s)", entry.Title, entry.Author.Name, entry.URL)
		}
	}
}
