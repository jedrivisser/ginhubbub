package ginhubbub

import (
	"encoding/base64"
	"net/url"
	"path"
)

// TopicToID convert a topic url to a base64 url safe string that can be used as an ID
func TopicToID(topic string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(topic))
}

// TopicFromID convert a previously encoded topic back into a topic url
func TopicFromID(id string) string {
	topic, _ := base64.RawURLEncoding.DecodeString(id)
	return string(topic)
}

// CallbackForTopic creates a topic specific callbackURL from a base callbackURL and a topic
func CallbackForTopic(callbackURL string, topic string) string {
	u, _ := url.Parse(callbackURL)
	u.Path = path.Join(u.Path, TopicToID(topic))
	return u.String()
}
