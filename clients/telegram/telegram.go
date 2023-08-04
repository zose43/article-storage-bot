package telegram

import (
	"article-storage-bot/lib/e"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const UPDATE = "getUpdates"
const SEND = "sendMessage"

var EmptyUpdates = errors.New("no one updateResponse")

type Client struct {
	host     string
	basePath string
	client   *http.Client
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	const msg = "Can't get updates"
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	data, err := c.doRequest(UPDATE, q)
	if err != nil {
		return nil, e.Wrap(msg, err)
	}

	var res UpdateResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap(msg, err)
	}
	if !res.Issue {
		return nil, EmptyUpdates
	}
	return res.Results, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(SEND, q)
	if err != nil {
		return e.Wrap("Can't send message", err)
	}
	return nil
}

func (c *Client) doRequest(method string, q url.Values) ([]byte, error) {
	const msg = "Can't do request"
	u := url.URL{
		Host:   c.host,
		Scheme: "https",
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(msg, err)
	}

	req.URL.RawQuery = q.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(msg, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(msg, err)
	}
	return body, err
}

func NewClient(token, host string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.DefaultClient,
	}
}

func newBasePath(token string) string {
	return "bot" + token
}
