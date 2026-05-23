package box

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// JWTClaims represents decoded JWT payload.
type JWTClaims struct {
	Raw     string
	Subject string
	Issuer  string
	Audience string
	Expires time.Time
	Issued  time.Time
	NotBefore time.Time
	ID      string
	Claims  map[string]any
}

// JWTHeader represents decoded JWT header.
type JWTHeader struct {
	Algorithm string
	Type      string
	Raw       map[string]any
}

// DecodeJWT decodes a JWT token (header + payload) without verification.
func DecodeJWT(token string) (*JWTHeader, *JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, nil, fmt.Errorf("invalid JWT: expected 3 parts, got %d", len(parts))
	}

	header, err := decodeJWTPart(parts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("header: %w", err)
	}

	payload, err := decodeJWTClaims(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("payload: %w", err)
	}

	h := &JWTHeader{
		Raw: header,
	}
	if alg, ok := header["alg"].(string); ok {
		h.Algorithm = alg
	}
	if typ, ok := header["typ"].(string); ok {
		h.Type = typ
	}

	c := &JWTClaims{
		Raw:    parts[1],
		Claims: payload,
	}
	if sub, ok := payload["sub"].(string); ok {
		c.Subject = sub
	}
	if iss, ok := payload["iss"].(string); ok {
		c.Issuer = iss
	}
	if aud, ok := payload["aud"].(string); ok {
		c.Audience = aud
	}
	if jti, ok := payload["jti"].(string); ok {
		c.ID = jti
	}
	if exp, ok := payload["exp"].(float64); ok {
		c.Expires = time.Unix(int64(exp), 0)
	}
	if iat, ok := payload["iat"].(float64); ok {
		c.Issued = time.Unix(int64(iat), 0)
	}
	if nbf, ok := payload["nbf"].(float64); ok {
		c.NotBefore = time.Unix(int64(nbf), 0)
	}

	return h, c, nil
}

func decodeJWTPart(part string) (map[string]any, error) {
	// Add padding
	part = padBase64(part)
	decoded, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		// Try standard base64
		decoded, err = base64.RawStdEncoding.DecodeString(part)
		if err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
	}
	var result map[string]any
	if err := json.Unmarshal(decoded, &result); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return result, nil
}

func decodeJWTClaims(part string) (map[string]any, error) {
	return decodeJWTPart(part)
}

func padBase64(s string) string {
	switch len(s) % 4 {
	case 2:
		return s + "=="
	case 3:
		return s + "="
	}
	return s
}
