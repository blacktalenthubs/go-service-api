package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8081/api"

// Models that match our API
type Skill struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Consultant struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	SkillIDs []int  `json:"skill_ids"`
}

// APIResponse wraps the API responses
type APIResponse struct {
	URL     string
	Method  string
	Status  int
	Body    []byte
	Error   error
	Message string
}

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// API request function with proper error handling
func (c *Client) makeRequest(method, url string, payload interface{}) APIResponse {
	var reqBody io.Reader = nil
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return APIResponse{
				URL:     url,
				Method:  method,
				Error:   err,
				Message: "Failed to marshal JSON",
			}
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return APIResponse{
			URL:     url,
			Method:  method,
			Error:   err,
			Message: "Failed to create request",
		}
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return APIResponse{
			URL:     url,
			Method:  method,
			Error:   err,
			Message: "Failed to send request",
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{
			URL:     url,
			Method:  method,
			Status:  resp.StatusCode,
			Error:   err,
			Message: "Failed to read response body",
		}
	}

	return APIResponse{
		URL:    url,
		Method: method,
		Status: resp.StatusCode,
		Body:   body,
	}
}

// Specific API operations
func (c *Client) CreateSkill(skill Skill) APIResponse {
	return c.makeRequest("POST", baseURL+"/skills", skill)
}

func (c *Client) GetSkills() APIResponse {
	return c.makeRequest("GET", baseURL+"/skills", nil)
}

func (c *Client) CreateConsultant(consultant Consultant) APIResponse {
	return c.makeRequest("POST", baseURL+"/consultants", consultant)
}

func (c *Client) GetConsultants() APIResponse {
	return c.makeRequest("GET", baseURL+"/consultants", nil)
}

func (c *Client) GetConsultantsBySkill(skillID int) APIResponse {
	return c.makeRequest("GET", fmt.Sprintf("%s/consultants/skills/%d", baseURL, skillID), nil)
}

func (c *Client) UpdateConsultant(id int, consultant Consultant) APIResponse {
	return c.makeRequest("PUT", fmt.Sprintf("%s/consultants/%d", baseURL, id), consultant)
}

func (c *Client) DeleteConsultant(id int) APIResponse {
	return c.makeRequest("DELETE", fmt.Sprintf("%s/consultants/%d", baseURL, id), nil)
}
