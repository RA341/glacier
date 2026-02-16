package config

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ra341/glacier/pkg/fileutil"
)

func checkGlacierUrl(url string) error {
	url = fmt.Sprintf("%s/%s", strings.TrimSuffix(url, "/"), "api/server/public/ping")

	reqCtx, cancel := context.WithTimeout(
		context.Background(),
		1*time.Second,
	)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer fileutil.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}

	return nil
}

// returns a tested url otherwise err
func testGlacierUrls(input string) (string, error) {
	var er error

	candidates := generateUrlCandidates(input)
	for _, candidate := range candidates {
		err := checkGlacierUrl(candidate)
		if err != nil {
			er = errors.Join(er, fmt.Errorf("%s: %w\n", candidate, err))
			continue
		}
		return candidate, nil
	}

	return "", er
}

// generates potential url candidates based on user input
// "example.com",           // No scheme, no port -> 4 candidates
// "http://example.com",    // Has scheme, no port -> 2 candidates
// "example.com:8080",      // No scheme, has port -> 2 candidates
// "https://example.com:8080", // Has scheme, has port -> 1 candidate
func generateUrlCandidates(input string) []string {
	var candidates []string

	hasScheme := strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")

	var baseInputs []string
	if !hasScheme {
		// Generate candidates with both http and https
		baseInputs = []string{
			"https://" + input,
			"http://" + input,
		}
	} else {
		baseInputs = []string{input}
	}

	// For each base input, check if port is present and generate candidates
	for _, baseInput := range baseInputs {
		parsedURL, err := url.Parse(baseInput)
		if err != nil {
			continue
		}

		hasPort := strings.Contains(parsedURL.Host, ":")
		if !hasPort {
			candidates = append(candidates, baseInput)
			// candidate with default port 6699
			urlWithPort := parsedURL.Scheme + "://" + parsedURL.Host + ":6699"
			if parsedURL.Path != "" {
				urlWithPort += parsedURL.Path
			}
			if parsedURL.RawQuery != "" {
				urlWithPort += "?" + parsedURL.RawQuery
			}
			candidates = append(candidates, urlWithPort)
		} else {
			// Port is already present
			candidates = append(candidates, baseInput)
		}
	}

	return candidates
}
