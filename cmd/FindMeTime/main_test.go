package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/julienschmidt/httprouter"
)

func TestGetTagsHandler(t *testing.T) {
	router := httprouter.New()
	router.GET("/api/v1/tags", GetTagsHandler)
	server := httptest.NewServer(router)
	defer server.Close()
	e := httpexpect.Default(t, server.URL)
	resp := e.GET("/api/v1/tags").
		Expect()

	resp.Status(http.StatusOK)
}

func TestCreateTagHandler(t *testing.T) {
	router := httprouter.New()
	router.POST("/api/v1/tags", GetTagsHandler)

	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	requestPayload := map[string]interface{}{
		"Name":        "Test Tag",
		"Description": "This is a test tag",
		"TimeSlots": []map[string]interface{}{
			{
				"DayIndex":  0,
				"StartTime": 9,
				"EndTime":   12,
			},
		},
	}
	resp := e.POST("/api/v1/tags").
		WithJSON(requestPayload).
		Expect()
	resp.Status(http.StatusOK)
	resp.JSON().NotNull()
	resp.JSON().Array()
}
