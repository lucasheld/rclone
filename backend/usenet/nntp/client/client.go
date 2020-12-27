package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

/*
https://tools.ietf.org/html/rfc2980
https://tools.ietf.org/html/rfc3977
https://tools.ietf.org/html/rfc4643
https://tools.ietf.org/html/rfc4644
https://tools.ietf.org/html/rfc6048
*/
type Client struct {
	conn *textproto.Conn
}

func NewClient(host string, port int, encryption bool) (*Client, error) {
	net := "tcp"
	addr := fmt.Sprintf("%s:%d", host, port)

	var conn *textproto.Conn
	if encryption {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		conf := &tls.Config{RootCAs: certPool}
		tlsConn, err := tls.Dial(net, addr, conf)
		if err != nil {
			return nil, err
		}
		conn = textproto.NewConn(tlsConn)
	} else {
		var err error
		conn, err = textproto.Dial(net, addr)
		if err != nil {
			return nil, err
		}
	}

	c := &Client{conn}
	err := c.connect()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c Client) connect() error {
	_, message, err := c.conn.ReadCodeLine(200)
	fmt.Printf("Connect message: %s\n", message)
	return err
}

func (c Client) Auth(username string, password string) error {
	err := c.conn.PrintfLine("AUTHINFO USER %s", username)
	if err != nil {
		return err
	}
	_, message, err := c.conn.ReadCodeLine(381)
	if err != nil {
		return err
	}
	fmt.Printf("Authinfo user message: %s\n", message)

	err = c.conn.PrintfLine("AUTHINFO PASS %s", password)
	if err != nil {
		return err
	}
	_, message, err = c.conn.ReadCodeLine(281)
	if err != nil {
		return err
	}
	fmt.Printf("Authinfo pass message: %s\n", message)
	return nil
}

type Group struct {
	Group  string
	Number int64
	Low    int64
	High   int64
}

func (c Client) Group(groupName string) (*Group, error) {
	err := c.conn.PrintfLine("GROUP %s", groupName)
	if err != nil {
		return nil, err
	}
	_, message, err := c.conn.ReadCodeLine(211)
	if err != nil {
		return nil, err
	}

	parameters := strings.Split(message, " ")
	group := &Group{}
	group.Number, err = strconv.ParseInt(parameters[0], 10, 64)
	if err != nil {
		return nil, err
	}
	group.Low, err = strconv.ParseInt(parameters[1], 10, 64)
	if err != nil {
		return nil, err
	}
	group.High, err = strconv.ParseInt(parameters[2], 10, 64)
	if err != nil {
		return nil, err
	}
	group.Group = parameters[3]
	return group, nil
}

//func (c Client) Listgroup() {
//
//}

type Article struct {
	Number    int64
	MessageId string
}

func (c Client) execLastOrNext(cmd string) (*Article, error) {
	err := c.conn.PrintfLine(cmd)
	if err != nil {
		return nil, err
	}
	_, message, err := c.conn.ReadCodeLine(223)
	if err != nil {
		return nil, err
	}
	parameters := strings.Split(message, " ")

	article := &Article{}
	article.Number, err = strconv.ParseInt(parameters[0], 10, 64)
	article.MessageId = parameters[1]
	return article, nil
}

func (c Client) Last() (*Article, error) {
	return c.execLastOrNext("LAST")
}

func (c Client) Next() (*Article, error) {
	return c.execLastOrNext("NEXT")
}

type ArticleReader struct {
	MessageId string
	Number    int64
	Reader    io.Reader
}

func (c Client) newArticleReader(expectCode int) (*ArticleReader, error) {
	_, message, err := c.conn.ReadCodeLine(expectCode)
	if err != nil {
		return nil, err
	}
	parameters := strings.Split(message, " ")
	number, err := strconv.ParseInt(parameters[0], 10, 64)
	if err != nil {
		return nil, err
	}
	messageId := parameters[1]

	articleReader := &ArticleReader{
		MessageId: messageId,
		Number:    number,
		Reader:    c.conn.DotReader(),
	}
	return articleReader, nil
}

func formatMessageId(messageId string) string {
	if messageId[0] != '<' {
		messageId = "<" + messageId + ">"
	}
	return messageId
}

func (c Client) Article(messageId string) (*ArticleReader, error) {
	messageId = formatMessageId(messageId)
	err := c.conn.PrintfLine("ARTICLE %s", messageId)
	if err != nil {
		return nil, err
	}
	return c.newArticleReader(220)
}

func (c Client) Head(messageId string) (*ArticleReader, error) {
	messageId = formatMessageId(messageId)
	err := c.conn.PrintfLine("HEAD %s", messageId)
	if err != nil {
		return nil, err
	}
	return c.newArticleReader(221)
}

func (c Client) Body(messageId string) (*ArticleReader, error) {
	messageId = formatMessageId(messageId)
	err := c.conn.PrintfLine("BODY %s", messageId)
	if err != nil {
		return nil, err
	}
	return c.newArticleReader(222)
}

func (c Client) Stat(messageId string) error {
	messageId = formatMessageId(messageId)

	err := c.conn.PrintfLine("STAT %s", messageId)
	if err != nil {
		return err
	}
	_, _, err = c.conn.ReadCodeLine(223)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) Post(reader io.Reader) error {
	err := c.conn.PrintfLine("POST")
	if err != nil {
		return err
	}
	_, _, err = c.conn.ReadCodeLine(340)
	if err != nil {
		return err
	}

	writer := c.conn.DotWriter()
	_, err = io.Copy(writer, reader)
	writer.Close()
	if err != nil {
		return err
	}

	_, _, err = c.conn.ReadCodeLine(240)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
