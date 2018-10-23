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
	client     *fasthttp.Client
	limiters   [3]rate.Limiter
}

func NewClient(token string) *Client {
	c := new(Client)
	c.apikey = []byte("OAuth " + token)
	c.client = &fasthttp.Client{Name: "appmetrica-go/0.0.0"}
	return c
}

func (c *Client) Close() error {
	return nil
}

// Application возвращает информацию об указанном приложении.
func (c *Client) Application(id int) (*Application, error) {
	req, res := c.prepare()
	uri := req.URI()
	uri.SetPath(`/management/v1/application/` + strconv.Itoa(id))
	var obj, err = c.do(req, res, nil)
	return obj.Application, err
}

// Applications возвращает информацию о приложениях, доступных пользователю.
func (c *Client) Applications() ([]Application, error) {
	req, res := c.prepare()
	uri := req.URI()
	uri.SetPath(`/management/v1/applications`)
	var obj, err = c.do(req, res, nil)
	return obj.Applications, err
}

// ModifyApplication добавляет приложение в AppMetrica.
func (c *Client) CreateApplication(name, tz string) (*Application, error) {
	req, res := c.prepare()
	req.Header.SetMethod("POST")

	uri := req.URI()
	uri.SetPath(`/management/v1/applications`)

	var msg = Response{Application: &Application{Name: name, TimeZoneName: tz}}
	var obj, err = c.do(req, res, msg)
	return obj.Application, err
}

// ModifyApplication изменяет настройки приложения.
func (c *Client) ModifyApplication(id int, name, tz string) (*Application, error) {
	req, res := c.prepare()
	req.Header.SetMethod("PUT")

	uri := req.URI()
	uri.SetPath(`/management/v1/application/` + strconv.Itoa(id))

	var msg = Response{Application: &Application{Name: name, TimeZoneName: tz}}
	var obj, err = c.do(req, res, &msg)
	return obj.Application, err
}

// DeleteApplication удаляет приложение.
func (c *Client) DeleteApplication(id int) error {
	req, res := c.prepare()
	req.Header.SetMethod("DELETE")

	uri := req.URI()
	uri.SetPath(`/management/v1/application/` + strconv.Itoa(id))

	var _, err = c.do(req, res, nil)
	return err
}

// ImportEvent загружает информацию о событии.
func (c *Client) ImportEvent() error {
	return nil
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

func (c *Client) do(req *fasthttp.Request, res *fasthttp.Response, msg interface{}) (*Response, error) {
	var err error
	var obj Response

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	if msg != nil {
		enc := json.NewEncoder(req.BodyWriter())
		req.Header.SetContentType(`application/json; encoding=utf-8`)

		if err = enc.Encode(msg); err != nil {
			return &obj, err
		}
	}

	println("\033[1mRequest dump string\033[0m")
	println(req.String())

	if err = c.client.Do(req, res); err != nil {
		return &obj, err
	}

	println("\033[1mResponse dump string\033[0m")
	println(res.String())

	var buf = bytes.NewBuffer(res.Body())
	var dec = json.NewDecoder(buf)

	if err = dec.Decode(&obj); err != nil {
		return &obj, err
	}

	if obj.ErrorCode != 0 {
		return &obj, NewError(obj.ErrorCode, obj.ErrorMessage)
	}

	return &obj, nil
}

func (c *Client) prepare() (*fasthttp.Request, *fasthttp.Response) {
	var req = fasthttp.AcquireRequest()
	var res = fasthttp.AcquireResponse()
	var uri = req.URI()
	uri.SetScheme(`https`)
	uri.SetHost(`api.appmetrica.yandex.ru`)
	req.Header.Set("User-Agent", `appmetrica/0.0.0+golang/1.11`)
	req.Header.SetBytesV("Authorization", c.apikey)
	return req, res
}
