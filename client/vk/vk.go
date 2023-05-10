package vk

import (
	"VKBot/lib/e"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type header struct {
	name  string
	value string
}

type param struct {
	name  string
	value string
}

type Client struct {
	Host     string
	BasePath string
	params   []param
	headers  []header
	client   http.Client
}

const (
	sendMessageMethod = "messages.send"
	getLongPollServer = "groups.getLongPollServer"
)

func New(host, token, version, basePath string) *Client {
	return &Client{
		Host:     host,
		BasePath: basePath,
		params:   newBaseParams(),
		headers:  newBaseHeaders(token),
		client:   http.Client{},
	}
}

func newParam(name, value string) param {
	return param{
		name:  name,
		value: value,
	}
}

func newHeader(name, value string) header {
	return header{
		name:  name,
		value: value,
	}
}

func newBaseParams() []param {
	return []param{newParam("v", "5.131")}
}

func newBaseHeaders(token string) []header {
	return []header{newHeader("Authorization", "Bearer "+token)}
}

func (c *Client) GetLongPollServer(groupID string) (longPollServer Response, err error) {
	defer func() { err = e.WrapIfErr("can't get LongPollServer", err) }()

	q := url.Values{}
	q.Add("group_id", groupID)

	var res LongPollServer

	data, err := c.doRequest(getLongPollServer, q)
	if err != nil {
		return res.Response, err
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return res.Response, err
	}

	return res.Response, nil
}

func (c *Client) GetLongPool(key, ts string) (updates GetResponse, err error) {
	defer func() { err = e.WrapIfErr("can't get LongPool", err) }()

	q := url.Values{}
	q.Add("act", "a_check")
	q.Add("key", key)
	q.Add("ts", ts)
	q.Add("wait", "25")

	var res GetResponse

	data, err := c.doRequest("", q)
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *Client) SendMessage(userID int, text, keyboard string) error {
	q := url.Values{}
	rand.Seed(time.Now().UnixNano())
	q.Add("user_id", strconv.Itoa(userID))
	q.Add("random_id", strconv.Itoa(rand.Int()))
	q.Add("message", text)
	q.Add("keyboard", keyboard)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	for _, val := range c.params {
		query.Add(val.name, val.value)
	}

	u := url.URL{
		Scheme: "https",
		Host:   c.Host,
		Path:   path.Join(c.BasePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	for _, val := range c.headers {
		req.Header.Add(val.name, val.value)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
