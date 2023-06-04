package utils

import (
	events "broker-service/internals/event"

	"github.com/gofiber/fiber/v2"
)

func LogRequest(c *fiber.Ctx, event events.Config) {
	var query interface{}
	c.QueryParser(&query)
	payload := events.Payload{
		Name: events.REQUEST,
		Data: struct {
			Method string
			Path   string
			Body   string
			User   string
			Query  interface{}
		}{
			Method: c.Method(),
			Path:   c.Path(),
			Body:   string(c.Body()),
			Query:  query,
		},
	}

	event.LogEventViaRabbit(&payload)
}
