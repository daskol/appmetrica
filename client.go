package appmetrica

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

// Client binds HTTP API to simple function calls.
type Client struct {
	apikey     []byte
	apikeyPost []byte
	limiters   [3]rate.Limiter
}

func NewClient(token string) *Client {
	c := new(Client)
	c.apikey = []byte("OAuth " + token)
	return c
}

func (c *Client) Close() error {
	return nil
}

// Application возвращает информацию об указанном приложении.
func (c *Client) Application(id int) (*Application, error) {
	var req = fasthttp.AcquireRequest()
	var uri = req.URI()

	uri.SetScheme(`https`)
	uri.SetHost(`api.appmetrica.yandex.ru`)
	uri.SetPath(`/management/v1/application/` + strconv.Itoa(id))
	req.Header.SetBytesV("Authorization", c.apikey)

	var res = fasthttp.AcquireResponse()
	var cli = fasthttp.Client{Name: "appmetrica-go/0.0.0"}

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	if err := cli.Do(req, res); err != nil {
		return nil, err
	}

	println(string(res.Body()))

	var body Response
	var buf = bytes.NewBuffer(res.Body())
	var dec = json.NewDecoder(buf)

	if err := dec.Decode(&body); err != nil {
		return nil, err
	}

	if body.ErrorCode != 0 {
		return nil, NewError(body.ErrorCode, body.ErrorMessage)
	}

	return body.Application, nil
}

// Applications возвращает информацию о приложениях, доступных пользователю.
func (c *Client) Applications() ([]Application, error) {
	var req = fasthttp.AcquireRequest()
	var uri = req.URI()

	uri.SetScheme(`https`)
	uri.SetHost(`api.appmetrica.yandex.ru`)
	uri.SetPath(`/management/v1/applications`)
	req.Header.SetBytesV("Authorization", c.apikey)

	var res = fasthttp.AcquireResponse()
	var cli = fasthttp.Client{Name: "appmetrica-go/0.0.0"}

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	if err := cli.Do(req, res); err != nil {
		return nil, err
	}

	println(string(res.Body()))
	println(req.String())

	var body Response
	var buf = bytes.NewBuffer(res.Body())
	var dec = json.NewDecoder(buf)

	if err := dec.Decode(&body); err != nil {
		return nil, err
	}

	if body.ErrorCode != 0 {
		return nil, NewError(body.ErrorCode, body.ErrorMessage)
	}

	return body.Applications, nil
}

// ModifyApplication добавляет приложение в AppMetrica.
func (c *Client) CreateApplication(name, tz string) (*Application, error) {
	return nil, nil
}

// ModifyApplication изменяет настройки приложения.
func (c *Client) ModifyApplication(id int, name, tz string) (*Application, error) {
	return nil, nil
}

// DeleteApplication удаляет приложение.
func (c *Client) DeleteApplication(id int) error {
	return nil
}

// ImportEvent загружает информацию о событии.
func (c *Client) ImportEvent() error {
}

// ImportEvents загружает информацию о событиях.
func (c *Client) ImportEvents() error {
	return nil
}

func (c *Client) SetAPIKey(token []byte) {
	c.apikey = token
}

func (c *Client) SetPostAPIKey(token []byte) {
	c.apikeyPost = token
}
