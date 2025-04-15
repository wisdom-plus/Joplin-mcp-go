package main

import (
	"fmt"
	"io"
	"net/http"
)

func findAvailablePort() int {
  var port int = 41184
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

func main() {
	port := findAvailablePort()
	fmt.Printf("見つかったポート: %d\n", port)
}
