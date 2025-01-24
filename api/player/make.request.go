package superplayer

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

func (player *Player) HttpRequest(endpoint, method, payload string) ([]byte, int, error) {
	url := fmt.Sprintf("%s%s", "https://citygame-dev.appspot.com", endpoint)

	req, err := http.NewRequest(method, url, bytes.NewBufferString(payload))
	if err != nil {
		return nil, -1, fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Add("Host", "citygame-dev.appspot.com")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", player.AuthCookie)
	req.Header.Add("User-Agent", "PixelPlex/1 CFNetwork/1335.0.3.4 Darwin/21.6.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "it-IT,it;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("App-Version", "4.4.1")

	res, err := player.HttpClient.Do(req)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("Error making request: %w", err)
	}
	defer res.Body.Close()

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, res.StatusCode, fmt.Errorf("Error creating gzip reader: %w", err)
		}
		defer reader.Close()
	default:
		reader = res.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("Error reading response body: %w", err)
	}

	return body, res.StatusCode, nil
}
