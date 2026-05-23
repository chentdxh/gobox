package box

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// JSON formats a JSON string with indentation.
func JSON(input string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// JSONCompact compacts a JSON string.
func JSONCompact(input string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	out, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Base64Encode encodes a string to base64.
func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// Base64Decode decodes a base64 string.
func Base64Decode(input string) (string, error) {
	out, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// URLEncode encodes a string for URL.
func URLEncode(input string) string {
	return url.QueryEscape(input)
}

// URLDecode decodes a URL-encoded string.
func URLDecode(input string) (string, error) {
	return url.QueryUnescape(input)
}

// Hash computes various hashes of input.
func Hash(input, algo string) (string, error) {
	data := []byte(input)
	switch strings.ToLower(algo) {
	case "md5":
		h := md5.Sum(data)
		return hex.EncodeToString(h[:]), nil
	case "sha1":
		h := sha1.Sum(data)
		return hex.EncodeToString(h[:]), nil
	case "sha256":
		h := sha256.Sum256(data)
		return hex.EncodeToString(h[:]), nil
	case "sha512":
		h := sha512.Sum512(data)
		return hex.EncodeToString(h[:]), nil
	default:
		return "", fmt.Errorf("unknown algorithm: %s (use md5, sha1, sha256, sha512)", algo)
	}
}

// TimestampToDate converts a unix timestamp to human-readable date.
func TimestampToDate(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format("2006-01-02 15:04:05 MST")
}

// DateToTimestamp converts a date string to unix timestamp.
func DateToTimestamp(dateStr string) (int64, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006-01-02T15:04:05",
	}
	for _, f := range formats {
		t, err := time.Parse(f, dateStr)
		if err == nil {
			return t.Unix(), nil
		}
	}
	return 0, fmt.Errorf("cannot parse date: %s", dateStr)
}

// HexEncode encodes a string to hex.
func HexEncode(input string) string {
	return hex.EncodeToString([]byte(input))
}

// HexDecode decodes a hex string.
func HexDecode(input string) (string, error) {
	out, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Count returns character, word, and line counts.
func Count(input string) (chars, words, lines int) {
	chars = len(input)
	lines = strings.Count(input, "\n") + 1
	words = len(strings.Fields(input))
	return
}
