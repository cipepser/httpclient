package sdk

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

const (
	timeout = 10
)

// Client have common http client for some api.
type Client struct {
	URL                *url.URL
	HTTPClient         *http.Client
	Username, Password string
	Logger             *log.Logger
}

// NewClient is a constructor of Client.
func NewClient(urlStr, username, password string, logger *log.Logger) (*Client, error) {
	if len(username) == 0 {
		return nil, errors.New("missing username")
	}
	if len(password) == 0 {
		return nil, errors.New("missing password")
	}

	var discardLogger = log.New(ioutil.Discard, "", log.LstdFlags)
	if logger == nil {
		logger = discardLogger
	}

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url: %s", urlStr)
	}

	c := &Client{
		URL:        parsedURL,
		HTTPClient: &http.Client{},
		Username:   username,
		Password:   password,
		Logger:     logger,
	}

	return c, err
}

// NewRequest is a wrapper of http.NewRequest which has timeout by context package.
func (c *Client) NewRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	return req, nil
}

// DecodeBody decode http responce to json format which specified by out.
func DecodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	return dec.Decode(out)
}
