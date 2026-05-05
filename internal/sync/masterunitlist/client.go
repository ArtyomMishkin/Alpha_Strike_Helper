package masterunitlist

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) QuickList(ctx context.Context, params map[string]string) ([]Unit, error) {
	u, err := url.Parse(c.baseURL + "/Unit/QuickList")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := u.Query()
	for k, v := range params {
		if strings.TrimSpace(v) == "" {
			continue
		}
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "alpha-strike-helper/masterunitlist-sync")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	var out quickListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return out.Units, nil
}

func (c *Client) FactionAutocomplete(ctx context.Context, term string) ([]LabelValue, error) {
	u, err := url.Parse(c.baseURL + "/Faction/Autocomplete")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	q := u.Query()
	q.Set("term", term)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "alpha-strike-helper/masterunitlist-sync")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	var out []LabelValue
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return out, nil
}

func (c *Client) FactionGroups(ctx context.Context, factions []LabelValue) (map[string]string, error) {
	u, err := url.Parse(c.baseURL + "/Faction/Index")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "alpha-strike-helper/masterunitlist-sync")
	req.Header.Set("Accept", "text/html")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	html := string(body)

	h3Re := regexp.MustCompile(`(?is)<h3>\s*([^<]+?)\s*</h3>`)
	stripRe := regexp.MustCompile(`(?is)<[^>]+>`)
	spaceRe := regexp.MustCompile(`\s+`)

	sections := h3Re.FindAllStringSubmatchIndex(html, -1)
	out := make(map[string]string, 128)
	for i, sec := range sections {
		group := strings.TrimSpace(html[sec[2]:sec[3]])
		start := sec[1]
		end := len(html)
		if i+1 < len(sections) {
			end = sections[i+1][0]
		}
		chunk := html[start:end]
		chunkText := strings.ToLower(strings.TrimSpace(spaceRe.ReplaceAllString(stripRe.ReplaceAllString(chunk, " "), " ")))
		for _, f := range factions {
			name := strings.TrimSpace(f.Label)
			if name == "" {
				continue
			}
			if strings.Contains(chunkText, strings.ToLower(name)) {
				out[name] = group
			}
		}
	}
	return out, nil
}

func DefaultEras() []Era {
	return []Era{
		{ID: 10, Name: "Star League"},
		{ID: 11, Name: "Early Succession War"},
		{ID: 255, Name: "Late Succession War - LosTech"},
		{ID: 256, Name: "Late Succession War - Renaissance"},
		{ID: 13, Name: "Clan Invasion"},
		{ID: 247, Name: "Civil War"},
		{ID: 14, Name: "Jihad"},
		{ID: 15, Name: "Early Republic"},
		{ID: 254, Name: "Late Republic"},
		{ID: 16, Name: "Dark Age"},
		{ID: 257, Name: "ilClan"},
	}
}

func eraIDString(id int) string {
	return strconv.Itoa(id)
}
