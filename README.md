# GinHubbub

## About

GinHubbub is A [PubPubHubbub 0.4][1] Subscriber written in go using [Gin][5]

Based on [GoHubbub][6]

It implements a Client that can be used to subscribe and unsubscribe from Subscriptions

It implements a [Gin][5] Server or Engine that listens on a Subscriber Callback URL to verify subscriptions and receives Notification

## Usage

Download the library:

```sh
$ go get github.com/jedrivisser/ginhubbub
```

To subscribe to a topic:

```sh
# assume the following codes in a main.go file
$ cat main.go
```

```go
package main

import "github.com/jedrivisser/ginhubbub"

func main() {
  hubURL := "http://medium.superfeedr.com" // The hub that should send you updates
  callbackURL := "https://1fde1b31.ngrok.io" // Publicly accessible URL where your server is running
  from := "owner@example.com" // Set in HTTP From header
  topic := "https://medium.com/feed/latest" // Feed that you want to subscribe to

  client := ginhubbub.NewClient(hubURL, callbackURL, from)
  client.Subscribe(topic)
}
```

```sh
$ go run main.go
```

Note that if you do not have a server running when you subscribe the hub will not be able to verify your subscription

To listen for Notifications and verify subscriptions:

```sh
# assume the following codes in a main.go file
$ cat main.go
```

```go
package main

import "github.com/jedrivisser/ginhubbub"

func main() {
  topic := "https://medium.com/feed/latest" // Feed to verify and listen for Notifications from

  subscriptions := make(map[string]func(string, []byte))
  subscriptions[topic] = func(contentType string, body []byte) {
    log.Printf(string(body)) // Print Notification body as string
  }

  server := ginhubbub.NewServer(subscriptions)
  go server.Engine().Run()
}
```

```sh
$ go run main.go
```

## Samples

There are 2 samples, one that subscribes to any new posts on [Medium][2], this is very nice for testing because there are always things being published

The other one is an example of subscribing to new videos in a YouTube channel, see [Subscribe to Push Notifications][3]

If you do not have a server that is exposed to the internet, you can use [ngrok][4] to forward requests to your machine while testing

[1]: https://pubsubhubbub.github.io/PubSubHubbub/pubsubhubbub-core-0.4.html
[2]: https://medium.com
[3]: https://developers.google.com/youtube/v3/guides/push_notifications
[4]: https://ngrok.com/
[5]: https://github.com/gin-gonic/gin
[6]: https://github.com/dpup/gohubbub