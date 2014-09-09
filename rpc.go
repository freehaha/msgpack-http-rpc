package rpc

import (
	"bytes"
	"errors"
	"github.com/ugorji/go/codec"
	"net/http"
)

type Client struct {
	URL    string
	handle codec.MsgpackHandle
	msgId  int
}

const (
	REQUEST   = 0
	RESPONSE  = 1
	NOTIFY    = 2
	MAX_MSGID = 256
)

type Request []interface{}

type Response []interface{}

func New(URL string) *Client {
	return &Client{
		URL:   URL,
		msgId: 0,
	}
}

func (c *Client) Call(method string, params interface{}) (result interface{}, err error) {
	var r Response
	data := make([]byte, 0, 256)
	buf := bytes.NewBuffer(data)
	enc := codec.NewEncoder(buf, &c.handle)
	enc.Encode(Request{
		REQUEST,
		c.msgId,
		method,
		params,
	})
	c.msgId += 1
	if c.msgId > MAX_MSGID {
		c.msgId = 0
	}
	res, err := http.Post(c.URL, "application/x-msgpack", buf)
	if err != nil {
		return
	}
	defer res.Body.Close()
	dec := codec.NewDecoder(res.Body, &c.handle)
	dec.Decode(&r)
	if t, ok := r[0].(int64); !ok || t != RESPONSE {
		err = errors.New("unknown type")
		return
	}
	if errMsg, ok := r[2].(string); ok && errMsg != "" {
		err = errors.New(errMsg)
		return
	}
	result = r[3]
	return
}
