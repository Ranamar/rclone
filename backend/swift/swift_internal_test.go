package swift

import (
	"testing"
	"time"

	"github.com/ncw/swift"
	"github.com/stretchr/testify/assert"
)

func TestInternalUrlEncode(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string
	}{
		{"", ""},
		{"abcdefghijklmopqrstuvwxyz", "abcdefghijklmopqrstuvwxyz"},
		{"ABCDEFGHIJKLMOPQRSTUVWXYZ", "ABCDEFGHIJKLMOPQRSTUVWXYZ"},
		{"0123456789", "0123456789"},
		{"abc/ABC/123", "abc/ABC/123"},
		{"   ", "%20%20%20"},
		{"&", "%26"},
		{"ß£", "%C3%9F%C2%A3"},
		{"Vidéo Potato Sausage?&£.mkv", "Vid%C3%A9o%20Potato%20Sausage%3F%26%C2%A3.mkv"},
	} {
		got := urlEncode(test.in)
		if got != test.want {
			t.Logf("%q: want %q got %q", test.in, test.want, got)
		}
	}
}

func TestInternalShouldRetryHeaders(t *testing.T) {
	headers := swift.Headers{
		"Content-Length": "64",
		"Content-Type":   "text/html; charset=UTF-8",
		"Date":           "Mon: 18 Mar 2019 12:11:23 GMT",
		"Retry-After":    "1",
	}
	err := &swift.Error{
		StatusCode: 429,
		Text:       "Too Many Requests",
	}
	start := time.Now()
	retry, gotErr := shouldRetryHeaders(headers, err)
	dt := time.Since(start)
	assert.True(t, retry)
	assert.Equal(t, err, gotErr)
	assert.True(t, dt > time.Second/2)
}
