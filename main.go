package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type apiServer struct {
	addr string
	app  *fiber.App
}

func NewApiServer(addr string) *apiServer {
	return &apiServer{addr: addr}
}

func (a *apiServer) Start() {
	a.app = fiber.New()

	a.registerRoutes()

	a.app.Listen(a.addr)
}

func (a *apiServer) registerRoutes() {
	a.app.Get("/", healthHandler)
	a.app.Post("/check", checkerHandler)
}

func pingUrls(urls []string) []responseItem {
	var response []responseItem

	for _, url := range urls {
		var checkedUrl responseItem

		timeBefore := time.Now()

		res, err := http.Get(url)

		elapsedTime := time.Since(timeBefore).Milliseconds()

		if err != nil {
			checkedUrl.Url = url
			checkedUrl.StatusCode = 404
			checkedUrl.ResponseTime = 0
		} else {
			checkedUrl.Url = url
			checkedUrl.StatusCode = res.StatusCode
			checkedUrl.ResponseTime = elapsedTime
		}

		response = append(response, checkedUrl)
	}

	return response
}

// handlers
func healthHandler(c *fiber.Ctx) error {
	return c.SendString("All is Okay")
}

func checkerHandler(c *fiber.Ctx) error {
	c.Accepts("json", "text")

	requestedJson := c.Body()
	urlsForCheck := &urls{}

	r := bytes.NewReader(requestedJson)
	decoder := json.NewDecoder(r)

	err := decoder.Decode(urlsForCheck)
	if err != nil {
		c.SendString(err.Error())
		return err
	}

	checkedUrls := pingUrls(urlsForCheck.Urls)
	structedUrls := responseBody{checkedUrls}
	encodedUrls, err := json.Marshal(structedUrls)
	if err != nil {
		c.SendString(err.Error())
	}
	return c.Send(encodedUrls)
}

// url struct
type urls struct {
	Urls []string `json:"urls"`
}

type responseBody struct {
	Urls []responseItem `json:"urls"`
}

type responseItem struct {
	Url          string `json:"url"`
	StatusCode   int    `json:"statuscode"`
	ResponseTime int64  `json:"responsetime"`
}


func main() {
	srv := NewApiServer(":3000")
	srv.Start()
}
