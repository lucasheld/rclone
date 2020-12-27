package client

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	host       = "127.0.0.1"
	port       = 1119
	encryption = false
	group      = "misc.test"
	messageId  = "<i.am.an.article@example.com>"
	user       = "\"Demo User\" <nobody@example.net>"
)

func buildHead() string {
	return "Message-Id: " + messageId + "\n" +
		"From: " + user + "\n" +
		"Newsgroups: " + group + "\n" +
		"Subject: I am just a test article" + "\n"

}

func buildBody() string {
	return "This is just a test article.\n"
}

func buildArticle() string {
	return buildHead() + "\n" + buildBody()
}

func TestConnect(t *testing.T) {
	_, err := NewClient(host, port, encryption)
	assert.NoError(t, err)
}

func TestGroup(t *testing.T) {
	client, _ := NewClient(host, port, encryption)

	g, err := client.Group("misc.test")
	assert.NoError(t, err)

	assert.Equal(t, g.Number, int64(0))
	assert.Equal(t, g.Low, int64(0))
	assert.Equal(t, g.High, int64(0))
	assert.Equal(t, g.Group, "misc.test")
}

func TestAuth(t *testing.T) {
	client, _ := NewClient(host, port, encryption)

	err := client.Auth("testuser", "testpass")
	assert.NoError(t, err)
}

func TestPost(t *testing.T) {
	client, _ := NewClient(host, port, encryption)
	_, _ = client.Group(group)

	buffer := strings.NewReader(buildArticle())
	err := client.Post(buffer)
	assert.NoError(t, err)

	testArticle(t, client)
	testHeader(t, client)
	testBody(t, client)
	testStat(t, client)
}

func testArticle(t *testing.T, client *Client) {
	article, err := client.Article(messageId)
	assert.NoError(t, err)

	assert.Equal(t, article.Number, int64(1))
	assert.Equal(t, article.MessageId, messageId)

	buffer, err := ioutil.ReadAll(article.Reader)
	assert.NoError(t, err)
	assert.Equal(t, string(buffer), buildArticle())
}

func testHeader(t *testing.T, client *Client) {
	head, err := client.Head(messageId)
	assert.NoError(t, err)

	assert.Equal(t, head.Number, int64(1))
	assert.Equal(t, head.MessageId, messageId)

	buffer, err := ioutil.ReadAll(head.Reader)
	assert.NoError(t, err)
	// TODO: ignore line order
	assert.Equal(t, string(buffer), buildHead())
}

func testBody(t *testing.T, client *Client) {
	body, err := client.Body(messageId)
	assert.NoError(t, err)

	assert.Equal(t, body.Number, int64(1))
	assert.Equal(t, body.MessageId, messageId)

	buffer, err := ioutil.ReadAll(body.Reader)
	assert.NoError(t, err)
	assert.Equal(t, string(buffer), buildBody())
}

func testStat(t *testing.T, client *Client) {
	err := client.Stat(messageId)
	assert.NoError(t, err)
}

func TestLast(t *testing.T) {
	client, _ := NewClient(host, port, encryption)
	_, _ = client.Group(group)

	article, err := client.Last()
	assert.NoError(t, err)
	assert.Equal(t, article.Number, int64(1))
	assert.Equal(t, article.MessageId, messageId)
}

func TestNext(t *testing.T) {
	client, _ := NewClient(host, port, encryption)
	_, _ = client.Group(group)

	article, err := client.Next()
	assert.NoError(t, err)
	assert.Equal(t, article.Number, int64(1))
	assert.Equal(t, article.MessageId, messageId)
}
