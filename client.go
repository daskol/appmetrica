package appmetrica

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

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
	c.client = &fasthttp.Client{Name: "appmetrica-go/0.0.0"}
	c.SetAPIKey(token)
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
func (c *Client) ImportEvent(event ImportEvent) error {
	return ErrNotImplemented
}

// ImportEvents загружает информацию о событиях. Функция в качестве аргумента
// принимает интерфейс Reader. Предполагается, что пользователь самостоятельно
// подготовил тело запроса в формате CSV, как того требует спецификация
// AppMetrica API. Задачу может упростить реализация интерфейса Reader тип
// EventImporter, который фильтрует и форматирует список событий.
func (c *Client) ImportEvents(reader io.Reader) error {
	req, res := c.prepare()
	req.Header.SetMethod("POST")
	req.Header.SetContentType(`text/csv; charset=UTF-8`)
	io.Copy(req.BodyWriter(), reader)

	uri := req.URI()
	uri.SetPath(`/logs/v1/import/events.csv`)

	args := uri.QueryArgs()
	args.SetBytesV(`post_api_key`, c.apikeyPost)

	var _, err = c.do(req, res, nil)
	return err
}

func (c *Client) SetAPIKey(token string) {
	c.apikey = []byte(`OAuth ` + token)
}

func (c *Client) SetPostAPIKey(token string) {
	c.apikeyPost = []byte(token)
}

func (c *Client) do(req *fasthttp.Request, res *fasthttp.Response, msg interface{}) (*Response, error) {
	var err error
	var obj Response

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	// Encoder JSON message.
	if msg != nil {
		enc := json.NewEncoder(req.BodyWriter())
		req.Header.SetContentType(`application/json; charset=UTF-8`)

		if err = enc.Encode(msg); err != nil {
			return &obj, err
		}
	}

	// Print prepared request.
	println("\033[1mRequest dump string\033[0m")
	println(req.String())

	// Make request.
	if err = c.client.Do(req, res); err != nil {
		return &obj, err
	}

	// Print received response.
	println("\033[1mResponse dump string\033[0m")
	println(res.String())

	contentType := string(res.Header.Peek(`Content-Type`))
	contentType = strings.Split(contentType, ";")[0]

	switch contentType {
	case "application/json", "application/x-yametrika+json":
		return c.processJSON(res)
	case "text/plain":
		return c.processPlainText(res)
	default:
		var status = res.StatusCode()
		var message = "unexpected content type: " + contentType
		return &obj, NewError(status, message)
	}
}

func (c *Client) processPlainText(res *fasthttp.Response) (*Response, error) {
	if status := res.StatusCode(); status != http.StatusOK {
		var message = string(res.Body())
		return nil, NewError(status, message)
	} else {
		return nil, nil
	}
}

func (c *Client) processJSON(res *fasthttp.Response) (*Response, error) {
	var buf = bytes.NewBuffer(res.Body())
	var dec = json.NewDecoder(buf)
	var obj Response

	if err := dec.Decode(&obj); err != nil {
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
