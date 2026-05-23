package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chentdxh/gobox/pkg/box"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "json", "jq":
		cmdJSON(args)
	case "jsonc":
		cmdJSONCompact(args)
	case "b64e", "b64enc":
		cmdB64Encode(args)
	case "b64d", "b64dec":
		cmdB64Decode(args)
	case "urle", "urlenc":
		cmdURLEncode(args)
	case "urld", "urldec":
		cmdURLDecode(args)
	case "hash":
		cmdHash(args)
	case "ts", "timestamp":
		cmdTimestamp(args)
	case "d2t", "date2ts":
		cmdDateToTimestamp(args)
	case "hexe", "hexenc":
		cmdHexEncode(args)
	case "hexd", "hexdec":
		cmdHexDecode(args)
	case "jwt":
		cmdJWT(args)
	case "count", "wc":
		cmdCount(args)
	case "all":
		cmdAll(args)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`gobox — dev toolbox for your terminal

Usage: gobox <command> [args]

Commands:
  json, jq     Format and indent JSON
  jsonc        Compact JSON to single line
  b64e, b64enc Base64 encode
  b64d, b64dec Base64 decode
  urle, urlenc URL encode
  urld, urldec URL decode
  hash <algo>  Hash string (md5, sha1, sha256, sha512)
  ts           Unix timestamp → readable date
  d2t          Date string → unix timestamp
  hexe, hexenc Hex encode
  hexd, hexdec Hex decode
  jwt          Decode JWT header + payload
  count, wc    Count chars, words, lines
  all          Show all transforms at once

All commands accept piped input: echo '{"a":1}' | gobox json`)
}

func readInput(args []string) string {
	if len(args) > 0 {
		return strings.Join(args, " ")
	}
	// Read from pipe
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)
		var b strings.Builder
		for {
			line, err := reader.ReadString('\n')
			b.WriteString(line)
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
		}
		return strings.TrimSuffix(b.String(), "\n")
	}
	return ""
}

func cmdJSON(args []string) {
	input := readInput(args)
	if input == "" {
		fmt.Fprintln(os.Stderr, "usage: gobox json <json-string> | echo '...' | gobox json")
		os.Exit(1)
	}
	out, err := box.JSON(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
}

func cmdJSONCompact(args []string) {
	input := readInput(args)
	out, err := box.JSONCompact(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
}

func cmdB64Encode(args []string) {
	input := readInput(args)
	fmt.Println(box.Base64Encode(input))
}

func cmdB64Decode(args []string) {
	input := readInput(args)
	out, err := box.Base64Decode(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
}

func cmdURLEncode(args []string) {
	input := readInput(args)
	fmt.Println(box.URLEncode(input))
}

func cmdURLDecode(args []string) {
	input := readInput(args)
	out, err := box.URLDecode(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
}

func cmdHash(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: gobox hash <algo> <string>\nalgo: md5, sha1, sha256, sha512")
		os.Exit(1)
	}
	algo := args[0]
	input := readInput(args[1:])
	out, err := box.Hash(input, algo)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
}

func cmdTimestamp(args []string) {
	input := readInput(args)
	ts, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid timestamp:", input)
		os.Exit(1)
	}
	fmt.Println(box.TimestampToDate(ts))
}

func cmdDateToTimestamp(args []string) {
	input := readInput(args)
	ts, err := box.DateToTimestamp(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(ts)
}

func cmdHexEncode(args []string) {
	input := readInput(args)
	fmt.Println(box.HexEncode(input))
}

func cmdHexDecode(args []string) {
	input := readInput(args)
	out, err := box.HexDecode(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
}

func cmdJWT(args []string) {
	input := readInput(args)
	if input == "" {
		fmt.Fprintln(os.Stderr, "usage: gobox jwt <token>")
		os.Exit(1)
	}
	header, claims, err := box.DecodeJWT(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Header:\n")
	prettyPrint(header.Raw)
	fmt.Printf("\nPayload:\n")
	prettyPrint(claims.Claims)
	fmt.Printf("\nSubject   : %s\n", claims.Subject)
	fmt.Printf("Issuer    : %s\n", claims.Issuer)
	if !claims.Expires.IsZero() {
		fmt.Printf("Expires   : %s\n", claims.Expires.Format("2006-01-02 15:04:05"))
		if claims.Expires.Before(time.Now()) {
			fmt.Println("             ⚠ EXPIRED")
		}
	}
	if !claims.Issued.IsZero() {
		fmt.Printf("Issued    : %s\n", claims.Issued.Format("2006-01-02 15:04:05"))
	}
}

func cmdCount(args []string) {
	input := readInput(args)
	chars, words, lines := box.Count(input)
	fmt.Printf("Chars : %d\nWords : %d\nLines : %d\n", chars, words, lines)
}

func cmdAll(args []string) {
	input := readInput(args)
	fmt.Printf("Input: %s\n\n", input)
	fmt.Printf("Base64 : %s\n", box.Base64Encode(input))
	fmt.Printf("URL    : %s\n", box.URLEncode(input))
	fmt.Printf("Hex    : %s\n", box.HexEncode(input))

	md5h, _ := box.Hash(input, "md5")
	sha256h, _ := box.Hash(input, "sha256")
	fmt.Printf("MD5    : %s\n", md5h)
	fmt.Printf("SHA256 : %s\n", sha256h)

	chars, words, lines := box.Count(input)
	fmt.Printf("\nChars  : %d\nWords  : %d\nLines  : %d\n", chars, words, lines)
}

func prettyPrint(v any) {
	out, _ := json.MarshalIndent(v, "  ", "  ")
	fmt.Println(" ", string(out))
}
