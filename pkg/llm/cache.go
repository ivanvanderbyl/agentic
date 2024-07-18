package llm

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CacheTransport struct {
	Transport       http.RoundTripper
	CacheDomains    []string
	CachePath       string
	CacheExpiration time.Duration
}

func (t *CacheTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.shouldCache(req.URL.Hostname()) {
		return t.Transport.RoundTrip(req)
	}

	cacheKey, err := t.GetCacheKey(req)
	if err != nil {
		return nil, err
	}

	// Check if we have a cached response
	// If we do, and it's not expired, return it
	if t.requestHasCachedResponse(cacheKey) {
		cachedResp, err := t.getCachedResponse(cacheKey)
		if err == nil && !t.isExpired(cachedResp) {
			return cachedResp, nil
		}
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		err = t.cacheResponse(cacheKey, resp)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (t *CacheTransport) shouldCache(hostname string) bool {
	if len(t.CacheDomains) == 0 {
		return true
	}

	for _, domain := range t.CacheDomains {
		if strings.HasSuffix(hostname, domain) {
			return true
		}
	}
	return false
}

func (t *CacheTransport) GetCacheKey(req *http.Request) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	buf := bytes.NewBuffer(nil)
	buf.WriteString(req.Method)
	headerKeys := make([]string, 0, len(req.Header))

	for k := range req.Header {
		headerKeys = append(headerKeys, k)
	}
	slices.Sort(headerKeys)
	for _, k := range headerKeys {
		buf.WriteString(k)
		sortedValues := req.Header[k]
		slices.Sort(sortedValues)
		buf.WriteString(strings.Join(sortedValues, ","))
		buf.WriteRune(';')
	}

	buf.WriteString(req.URL.String())
	buf.Write(body)

	return uuid.NewSHA1(uuid.NameSpaceOID, buf.Bytes()).String(), nil
}

func (t *CacheTransport) requestHasCachedResponse(cacheKey string) bool {
	cacheFile := filepath.Join(t.CachePath, cacheKey)

	_, err := os.Stat(cacheFile)
	if err != nil {
		return false
	}

	return true

}

func (t *CacheTransport) getCachedResponse(cacheKey string) (*http.Response, error) {
	cacheFile := filepath.Join(t.CachePath, cacheKey)
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *CacheTransport) isExpired(resp *http.Response) bool {
	if t.CacheExpiration == 0 {
		return false
	}

	cacheTime := resp.Header.Get("X-Cache-Time")
	if cacheTime == "" {
		return true
	}

	cacheTimestamp, err := time.Parse(time.RFC3339, cacheTime)
	if err != nil {
		return true
	}

	return time.Since(cacheTimestamp) > t.CacheExpiration
}

func (t *CacheTransport) cacheResponse(cacheKey string, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	finalBuffer := bytes.NewBuffer(body)
	resp.Body = io.NopCloser(finalBuffer)
	resp.Header.Set("X-Cache-Time", time.Now().Format(time.RFC3339))

	otherResp := http.Response{
		StatusCode:    resp.StatusCode,
		Proto:         resp.Proto,
		ProtoMajor:    resp.ProtoMajor,
		ProtoMinor:    resp.ProtoMinor,
		Header:        resp.Header,
		ContentLength: int64(finalBuffer.Len()),
		Body:          io.NopCloser(bytes.NewReader(finalBuffer.Bytes())),
	}

	buf := bytes.NewBuffer(nil)
	err = otherResp.Write(buf)
	if err != nil {
		return errors.Wrap(err, "failed to write response to buffer")
	}

	cacheFile := filepath.Join(t.CachePath, cacheKey)
	err = os.MkdirAll(filepath.Dir(cacheFile), 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create cache directory")
	}

	err = os.WriteFile(cacheFile, buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write cache file")
	}

	return nil
}

func NewCacheTransport(transport http.RoundTripper, cacheDomains []string, cachePath string, cacheExpiration time.Duration) *CacheTransport {
	return &CacheTransport{
		Transport:       transport,
		CacheDomains:    cacheDomains,
		CachePath:       cachePath,
		CacheExpiration: cacheExpiration,
	}
}
