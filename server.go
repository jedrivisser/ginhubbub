package ginhubbub

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server holds a *gin.Engine that is a PubSubHubbub 0.4 Subscriber and a map of subscriptions
type Server interface {
	Engine() *gin.Engine
	AddSubsciption(topic string, handler func(string, []byte))
	RemoveSubsciption(topic string)
}

type ginServer struct {
	engine        *gin.Engine
	subscriptions map[string]func(string, []byte) // Topic, Content-Type, ResponseBody
}

// NewServer sets up a gin pubsubhubbub engine
// This is a PubSubHubbub 0.4 Subscriber
// spec: https://pubsubhubbub.github.io/PubSubHubbub/pubsubhubbub-core-0.4.html
// subscriptions should contain all the topics you want to be subscribed for
func NewServer(subscriptions map[string]func(string, []byte)) Server {
	router := gin.Default()

	server := &ginServer{
		router,
		subscriptions,
	}

	router.GET("/:topicID", server.handleValidation)
	router.POST("/:topicID", server.handlePublication)

	return server
}

// Engine returns the *gin.Engine that is a PubSubHubbub 0.4 Subscriber
func (server *ginServer) Engine() *gin.Engine {
	return server.engine
}

// AddSubsciption adds a topic to the expected list of subscriptions
func (server *ginServer) AddSubsciption(topic string, handler func(string, []byte)) {
	server.subscriptions[topic] = handler
}

// RemoveSubsciption removes a topic to the expected list of subscriptions
func (server *ginServer) RemoveSubsciption(topic string) {
	delete(server.subscriptions, topic)
}

func (server *ginServer) handleValidation(context *gin.Context) {
	topicFromID := TopicFromID(context.Param("topicID"))
	topic := context.Query("hub.topic")
	mode := context.Query("hub.mode")

	if topic != topicFromID {
		log.Printf("topic for the url: %s, does not match hub.topic: %s", topicFromID, topic)
		context.AbortWithStatus(http.StatusNotFound)
		return
	}

	switch mode {
	// https://pubsubhubbub.github.io/PubSubHubbub/pubsubhubbub-core-0.4.html#validationsub
	case "denied":
		log.Printf("Subscription denied for %s, reason was %s", topic, context.Query("hub.reason"))
		context.Status(http.StatusOK)
	// https://pubsubhubbub.github.io/PubSubHubbub/pubsubhubbub-core-0.4.html#verifysub
	case "subscribe":
		if _, exists := server.subscriptions[topic]; exists {
			log.Printf("Subscription verified for %s, lease is %s seconds", topic, context.Query("hub.lease_seconds"))
			context.String(http.StatusOK, context.Query("hub.challenge"))
		} else {
			log.Printf("Unexpected subscription for %s", topic)
			context.AbortWithStatus(http.StatusNotFound)
		}
	// https://pubsubhubbub.github.io/PubSubHubbub/pubsubhubbub-core-0.4.html#verifysub
	case "unsubscribe":
		if _, exists := server.subscriptions[topic]; !exists {
			log.Printf("Unsubscribe confirmed for %s", topic)
			context.String(http.StatusOK, context.Query("hub.challenge"))
		} else {
			log.Printf("Unexpected unsubscribe for %s", topic)
			context.AbortWithStatus(http.StatusNotFound)
		}
	// Some hub.mode message we should not have received
	default:
		log.Printf("Unexpected hub.mode: %s", mode)
		context.AbortWithStatus(http.StatusNotFound)
	}
}

func (server *ginServer) handlePublication(context *gin.Context) {
	contentType := context.ContentType()
	body, _ := context.GetRawData()
	topic := TopicFromID(context.Param("topicID"))

	if handler, exists := server.subscriptions[topic]; exists {
		log.Printf("Received publication from: %s", topic)
		handler(contentType, body)
		context.Status(http.StatusOK)
	} else {
		log.Printf("Received unexpected publication from: %s", topic)
		context.AbortWithStatus(http.StatusNotFound)
	}
}
