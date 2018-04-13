package ginhubbub

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
)

// Client allows you to make PubSubHubbub subscriptions and register callback
// handlers that will be executed when an update is received.
type Client struct {
	// URL of the PubSubHubbub Hub to make requests to.
	hubURL string
	// callbackURL that the client will be served from, should be accessible by the hub
	callbackURL string
	from        string       // String passed in the "From" header.
	httpClient  *http.Client // e.g. http.Client{}.
}

// NewClient creates a pubsubhubbub client
func NewClient(hubURL string, callbackURL string, from string) *Client {
	return &Client{
		hubURL:      hubURL,
		callbackURL: callbackURL,
		from:        from,
		httpClient:  &http.Client{},
	}
}

// Subscribe sends a subscribe notrequestification to the hub
// https://pubsubhubbub.github.io/PubSubHubbub/pubsubhubbub-core-0.4.html#rfc.section.5.1
func (client *Client) Subscribe(topic string) {
	log.Println("Subscribing to", topic)
	client.subscriptionRequest(topic, "subscribe")
}

// Unsubscribe sends an unsubscribe request to the hub
// https://pubsubhubbub.github.io/PubSubHubbub/pubsubhubbub-core-0.4.html#rfc.section.5.1
func (client *Client) Unsubscribe(topic string) {
	log.Println("Unsubscribing from", topic)
	client.subscriptionRequest(topic, "unsubscribe")
}

// SubscriptionRequest subscribes or unsubscribes to a topic depending on the mdoe
// https://pubsubhubbub.github.io/PubSubHubbub/pubsubhubbub-core-0.4.html#rfc.section.5.1
func (client *Client) subscriptionRequest(topic string, mode string) {
	body := url.Values{
		"hub.callback": {CallbackForTopic(client.callbackURL, topic)},
		"hub.topic":    {topic},
		"hub.mode":     {mode},
	}

	req, _ := http.NewRequest(http.MethodPost, client.hubURL, bytes.NewBufferString(body.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("From", client.from)

	resp, err := client.httpClient.Do(req)

	if err != nil {
		log.Printf("%s failed, %s, %s", mode, topic, err)

	} else if resp.StatusCode != http.StatusAccepted {
		log.Printf("%s failed, %s status = %s", mode, topic, resp.Status)
	}
}
