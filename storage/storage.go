package storage

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// NoCacheStorage keeps cookies only memory
// without persisting data on the disk.
type NoCacheStorage struct {
	jar *cookiejar.Jar
}

// Init initializes NoCacheStorage
func (s *NoCacheStorage) Init() error {
	if s.jar == nil {
		var err error
		s.jar, err = cookiejar.New(nil)
		return err
	}
	return nil
}

// Visited implements Storage.Visited()
func (s *NoCacheStorage) Visited(requestID uint64) error {
	return nil
}

func (s *NoCacheStorage) IsVisited(requestID uint64) (bool, error) {
	return false, nil
}

// Cookies implements Storage.Cookies()
func (s *NoCacheStorage) Cookies(u *url.URL) string {
	return StringifyCookies(s.jar.Cookies(u))
}

// SetCookies implements Storage.SetCookies()
func (s *NoCacheStorage) SetCookies(u *url.URL, cookies string) {
	s.jar.SetCookies(u, UnstringifyCookies(cookies))
}

// Close implements Storage.Close()
func (s *NoCacheStorage) Close() error {
	return nil
}

// StringifyCookies serializes list of http.Cookies to string
func StringifyCookies(cookies []*http.Cookie) string {
	// Stringify cookies.
	cs := make([]string, len(cookies))
	for i, c := range cookies {
		cs[i] = c.String()
	}
	return strings.Join(cs, "\n")
}

// UnstringifyCookies deserializes a cookie string to http.Cookies
func UnstringifyCookies(s string) []*http.Cookie {
	h := http.Header{}
	for _, c := range strings.Split(s, "\n") {
		h.Add("Set-Cookie", c)
	}
	r := http.Response{Header: h}
	return r.Cookies()
}
