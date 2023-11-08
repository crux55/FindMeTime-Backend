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
	resp.JSON().Array()
}

func TestCreateTagHandler200(t *testing.T) {
	router := httprouter.New()
	router.POST("/api/v1/tags", CreateTagHandler)

	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	requestPayload := map[string]interface{}{
		"Name":        "Test Tag",
		"Description": "This is a test tag",
		"TimeSlots": []map[string]interface{}{
			{
				"DayIndex":  1,
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
}

func TestCreateTagHandler400(t *testing.T) {
	router := httprouter.New()
	router.POST("/api/v1/tags", CreateTagHandler)

	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	requestPayloadNoName := Tag{
		Description: "This is a test tag",
		TimeSlots: []TimeSlot{
			{
				DayIndex:  1,
				StartTime: 9,
				EndTime:   12,
			},
		},
	}
	resp := e.POST("/api/v1/tags").
		WithJSON(requestPayloadNoName).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains("Name is required")

	requestPayloadNoTimeSlot := Tag{
		Name:        "Test Tag",
		Description: "This is a test tag",
	}

	resp = e.POST("/api/v1/tags").
		WithJSON(requestPayloadNoTimeSlot).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains("TimeSlots is required")

	requestPayloadNoTimeSlotOrName := Tag{
		Description: "This is a test tag",
	}

	resp = e.POST("/api/v1/tags").
		WithJSON(requestPayloadNoTimeSlotOrName).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains("Name is required")

	requestPayloadTimeSlotNoDayIndex := Tag{
		Name:        "Test Tag",
		Description: "This is a test tag",
		TimeSlots: []TimeSlot{
			{
				StartTime: 1,
				EndTime:   12,
			},
		},
	}

	resp = e.POST("/api/v1/tags").
		WithJSON(requestPayloadTimeSlotNoDayIndex).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains("DayIndex is required")

	requestPayloadNoStartTime := Tag{
		Name:        "Test Tag",
		Description: "This is a test tag",
		TimeSlots: []TimeSlot{
			{
				DayIndex: 1,
				EndTime:  12,
			},
		},
	}

	resp = e.POST("/api/v1/tags").
		WithJSON(requestPayloadNoStartTime).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains("StartTime is required")

	requestPayloadNoEndTime := Tag{
		Name:        "Test Tag",
		Description: "This is a test tag",
		TimeSlots: []TimeSlot{
			{
				DayIndex:  1,
				StartTime: 12,
			},
		},
	}

	resp = e.POST("/api/v1/tags").
		WithJSON(requestPayloadNoEndTime).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains("EndTime is required")

	requestPayloadNonArrayTimeSlot := map[string]interface{}{
		"Name":        "Test Tag",
		"Description": "This is a test tag",
		"TimeSlots": map[string]interface{}{
			"DayIndex":  1,
			"StartTime": 9,
			"EndTime":   12,
		},
	}
	resp = e.POST("/api/v1/tags").
		WithJSON(requestPayloadNonArrayTimeSlot).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains("TimeSlots is required")

	requestPayloadLargeDayIndex := map[string]interface{}{
		"Name":        "Test Tag",
		"Description": "This is a test tag",
		"TimeSlots": []map[string]interface{}{
			{
				"DayIndex":  7,
				"StartTime": 9,
				"EndTime":   12,
			},
		},
	}

	resp = e.POST("/api/v1/tags").
		WithJSON(requestPayloadLargeDayIndex).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains("DayIndex max value is 6")

	requestPayloadLargeStartTime := map[string]interface{}{
		"Name":        "Test Tag",
		"Description": "This is a test tag",
		"TimeSlots": []map[string]interface{}{
			{
				"DayIndex":  6,
				"StartTime": 24,
				"EndTime":   12,
			},
		},
	}

	resp = e.POST("/api/v1/tags").
		WithJSON(requestPayloadLargeStartTime).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains("StartTime max value is 23")

	requestPayloadLargeEndTime := map[string]interface{}{
		"Name":        "Test Tag",
		"Description": "This is a test tag",
		"TimeSlots": []map[string]interface{}{
			{
				"DayIndex":  6,
				"StartTime": 23,
				"EndTime":   24,
			},
		},
	}

	resp = e.POST("/api/v1/tags").
		WithJSON(requestPayloadLargeEndTime).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains("EndTime max value is 23")

}
