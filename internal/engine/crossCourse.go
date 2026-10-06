package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CrossCourse struct {
	Client  *http.Client
	BaseURL string
}

func NewCrossCourse() *CrossCourse {
	return &CrossCourse{
		Client:  &http.Client{Timeout: 10 * time.Second},
		BaseURL: "https://api.nbrb.by/exrates/rates/",
	}
}

func (c *CrossCourse) Get(p CrossCourseParams) (*Rate, error) {
	if c.Client == nil {
		c.Client = http.DefaultClient
	}
	if c.BaseURL == "" {
		c.BaseURL = "https://api.nbrb.by/exrates/rates/"
	}
	if p.CurrencyID == "" {
		return nil, fmt.Errorf("currency id must not be empty")
	}

	path := strings.TrimRight(c.BaseURL, "/") + "/" + url.PathEscape(p.CurrencyID)
	reqURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	if p.ParamMode != 0 {
		values.Set("parammode", strconv.Itoa(p.ParamMode))
	}
	if p.Periodicity != 0 {
		values.Set("periodicity", strconv.Itoa(p.Periodicity))
	}
	if !p.OnDate.IsZero() {
		values.Set("ondate", p.OnDate.Format("2006-01-02"))
	}
	if len(values) > 0 {
		reqURL.RawQuery = values.Encode()
	}

	resp, err := c.Client.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<10))
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rate Rate
	if err := json.Unmarshal(body, &rate); err != nil {
		return nil, err
	}

	return &rate, nil
}
