package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	os.Setenv("CONFIG_PATH", `C:\Users\games\Code\Go\FindMeTime-Backend\envs\local.yml`)
	// Create a new Tag
	tag := Tag{
		Name: "Test Tag",
		TimeSlots: []TimeSlot{
			{
				StartDayIndex: 1,
				StartTime:     9,
				EndDayIndex:   1,
				EndTime:       12,
			},
		},
	}

	// Convert the Tag to JSON
	jsonTag, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("Failed to convert tag to JSON: %v", err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", "/tags", bytes.NewBuffer(jsonTag))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a handler to test
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateTagHandler(w, r, httprouter.Params{})
	})
	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `Tag: {Id: Name:Test Tag Description: TimeSlots:[{StartDayIndex:1 StartTime:9 EndDayIndex:1 EndTime:12}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateTagHandler400(t *testing.T) {
	os.Setenv("CONFIG_PATH", `C:\Users\games\Code\Go\FindMeTime-Backend\envs\local.yml`)
	router := httprouter.New()
	router.POST("/api/v1/tags", CreateTagHandler)

	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	testCases := []struct {
		name     string
		payload  interface{}
		errorMsg string
	}{
		{
			name: "No Name",
			payload: Tag{
				Description: "This is a test tag",
				TimeSlots: []TimeSlot{
					{
						StartDayIndex: 1,
						StartTime:     9,
						EndTime:       12,
						EndDayIndex:   1,
					},
				},
			},
			errorMsg: "Key: 'Tag.Name' Error:Field validation for 'Name' failed on the 'required' tag\n",
		},
		{
			name: "No TimeSlots",
			payload: Tag{
				Name:        "Test Tag",
				Description: "This is a test tag",
			},
			errorMsg: "Key: 'Tag.TimeSlots' Error:Field validation for 'TimeSlots' failed on the 'required' tag\n",
		},
		{
			name: "No TimeSlots or name",
			payload: Tag{
				Description: "This is a test tag",
			},
			errorMsg: "Key: 'Tag.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'Tag.TimeSlots' Error:Field validation for 'TimeSlots' failed on the 'required' tag\n",
		},
		{
			name: "Missing StartDayIndex",
			payload: Tag{
				Name:        "Test Tag",
				Description: "This is a test tag",
				TimeSlots: []TimeSlot{
					{
						StartTime:   9,
						EndDayIndex: 1,
						EndTime:     12,
					},
				},
			},
			errorMsg: "Key: 'Tag.TimeSlots[0].StartDayIndex' Error:Field validation for 'StartDayIndex' failed on the 'required' tag\n",
		},
		{
			name: "Missing StartTime",
			payload: Tag{
				Name:        "Test Tag",
				Description: "This is a test tag",
				TimeSlots: []TimeSlot{
					{
						StartDayIndex: 1,
						EndDayIndex:   1,
						EndTime:       12,
					},
				},
			},
			errorMsg: "Key: 'Tag.TimeSlots[0].StartTime' Error:Field validation for 'StartTime' failed on the 'required' tag\n",
		},
		{
			name: "Missing EndDayIndex",
			payload: Tag{
				Name:        "Test Tag",
				Description: "This is a test tag",
				TimeSlots: []TimeSlot{
					{
						StartDayIndex: 1,
						StartTime:     9,
						EndTime:       12,
					},
				},
			},
			errorMsg: "Key: 'Tag.TimeSlots[0].EndDayIndex' Error:Field validation for 'EndDayIndex' failed on the 'required' tag\n",
		},
		{
			name: "Missing EndTime",
			payload: Tag{
				Name:        "Test Tag",
				Description: "This is a test tag",
				TimeSlots: []TimeSlot{
					{
						StartDayIndex: 1,
						StartTime:     9,
						EndDayIndex:   1,
					},
				},
			},
			errorMsg: "Key: 'Tag.TimeSlots[0].EndTime' Error:Field validation for 'EndTime' failed on the 'required' tag\n",
		},
		{
			name: "StartDayIndex Out of Range",
			payload: Tag{
				Name:        "Test Tag",
				Description: "This is a test tag",
				TimeSlots: []TimeSlot{
					{
						StartDayIndex: 7,
						StartTime:     9,
						EndDayIndex:   1,
						EndTime:       12,
					},
				},
			},
			errorMsg: "Key: 'Tag.TimeSlots[0].StartDayIndex' Error:Field validation for 'StartDayIndex' failed on the 'max' tag\n",
		},
		{
			name: "StartTime Out of Range",
			payload: Tag{
				Name:        "Test Tag",
				Description: "This is a test tag",
				TimeSlots: []TimeSlot{
					{
						StartDayIndex: 1,
						StartTime:     24,
						EndDayIndex:   1,
						EndTime:       12,
					},
				},
			},
			errorMsg: "Key: 'Tag.TimeSlots[0].StartTime' Error:Field validation for 'StartTime' failed on the 'max' tag\n",
		},
		{
			name: "EndDayIndex Out of Range",
			payload: Tag{
				Name:        "Test Tag",
				Description: "This is a test tag",
				TimeSlots: []TimeSlot{
					{
						StartDayIndex: 1,
						StartTime:     9,
						EndDayIndex:   7,
						EndTime:       12,
					},
				},
			},
			errorMsg: "Key: 'Tag.TimeSlots[0].EndDayIndex' Error:Field validation for 'EndDayIndex' failed on the 'max' tag\n",
		},
		{
			name: "EndTime Out of Range",
			payload: Tag{
				Name:        "Test Tag",
				Description: "This is a test tag",
				TimeSlots: []TimeSlot{
					{
						StartDayIndex: 1,
						StartTime:     9,
						EndDayIndex:   1,
						EndTime:       24,
					},
				},
			},
			errorMsg: "Key: 'Tag.TimeSlots[0].EndTime' Error:Field validation for 'EndTime' failed on the 'max' tag\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := e.POST("/api/v1/tags").
				WithJSON(tc.payload).
				Expect()
			fmt.Print(tc.name)
			resp.Status(http.StatusBadRequest)
			resp.Body().Contains(tc.errorMsg)
		})
	}
}
