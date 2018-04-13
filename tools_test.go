package ginhubbub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testTopic = "https://www.youtube.com/xml/feeds/videos.xml?channel_id=UC99lkbVG8I5hRSZa4FD8zgw"
var testID = "aHR0cHM6Ly93d3cueW91dHViZS5jb20veG1sL2ZlZWRzL3ZpZGVvcy54bWw_Y2hhbm5lbF9pZD1VQzk5bGtiVkc4STVoUlNaYTRGRDh6Z3c"
var testCallbackForTopic = "https://localhost/callback/aHR0cHM6Ly93d3cueW91dHViZS5jb20veG1sL2ZlZWRzL3ZpZGVvcy54bWw_Y2hhbm5lbF9pZD1VQzk5bGtiVkc4STVoUlNaYTRGRDh6Z3c"

func TestTopicToID(t *testing.T) {
	ID := TopicToID(testTopic)
	assert.Equal(t, ID, testID)
}

func TestTopicFromID(t *testing.T) {
	topic := TopicFromID(testID)
	assert.Equal(t, topic, testTopic)
}

func TestCallbackForTopicWithoutSlash(t *testing.T) {
	testCallbackURL := "https://localhost/callback"

	callbackForTopic := CallbackForTopic(testCallbackURL, testTopic)
	assert.Equal(t, callbackForTopic, testCallbackForTopic)
}

func TestCallbackForTopicWithSlash(t *testing.T) {
	testCallbackURL := "https://localhost/callback/"

	callbackForTopic := CallbackForTopic(testCallbackURL, testTopic)
	assert.Equal(t, callbackForTopic, testCallbackForTopic)
}