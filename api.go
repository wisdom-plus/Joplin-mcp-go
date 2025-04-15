package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JoplinClient struct {
	BaseURL string
	Token   string
	Port 		int
}

type Note struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// NewJoplinClient creates a new JoplinClient instance
func NewJoplinClient(baseURL, token string) *JoplinClient {
	return &JoplinClient{
		BaseURL: baseURL,
		Token:   token,
		Port:  0,
	}
}

// GetNote retrieves a note by its ID
func (c *JoplinClient) GetNote(noteID string) (*Note, error) {
	if err := c.ensurePort(); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/notes/%s?token=%s", c.BaseURL, noteID, c.Token)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch note: %s", string(body))
	}

	var note Note
	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &note, nil
}

func(c *JoplinClient) findPort() int {
  var port int
  for portToTest := 41184; portToTest <= 41194; portToTest++ {
		url := fmt.Sprintf("http://127.0.0.1:%d/ping", portToTest)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		if string(body) == "JoplinClipperServer" {
			port = portToTest
			break
		}
  }

	return port
}

func (c *JoplinClient) ensurePort() error {
	if c.Port != 0 {
		return nil
	}
	port := c.findPort()
	if port == 0 {
		return fmt.Errorf("Joplin Clipper Server のポートが見つかりません")
	}
	c.Port = port
	c.BaseURL = fmt.Sprintf("http://127.0.0.1:%d", port)
	return nil
}
