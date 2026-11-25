// Apier examples in Go. Standard library only.
//
//	go run . verify     registration status and signals (keyless public sandbox)
//	go run . authority  coded signing rules, signaturrett and prokura (keyless public sandbox)
//	go run . proof      an Authority answer in a signed certificate
//
// verify and authority call the public sandbox for the fictional test company
// 999999999 and need no key.
//
// proof without APIER_API_KEY calls the bearer sandbox. The bearer is invented
// here, needs no signup, and the certificate is signed with the public sandbox
// key, so it is for rehearsing verification code only. With APIER_API_KEY and
// APIER_ORG_NUMBER set, proof calls production. That call needs a real key and
// is metered: https://www.apier.no/pricing
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var orgNumber = regexp.MustCompile(`^[0-9]{9}$`)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: go run . verify|authority|proof")
	}
	base, err := baseURL()
	if err != nil {
		return err
	}

	var target, bearer, field string
	switch args[0] {
	case "verify":
		target = base + "/api/v1/sandbox/public/company/999999999/verify"
		field = "verification_status"
	case "authority":
		target = base + "/api/v1/sandbox/public/company/999999999/authority"
		field = "classification"
	case "proof":
		field = "signature"
		if key := os.Getenv("APIER_API_KEY"); key != "" {
			org := os.Getenv("APIER_ORG_NUMBER")
			if !orgNumber.MatchString(org) {
				return fmt.Errorf("set APIER_ORG_NUMBER to a 9-digit organisasjonsnummer when APIER_API_KEY is set")
			}
			target = base + "/api/v1/company/" + org + "/authority?certificate=true"
			bearer = key
		} else {
			target = base + "/api/v1/sandbox/company/818000006/authority?certificate=true"
			bearer = fmt.Sprintf("apier_sandbox_test_%d", time.Now().UnixNano())
		}
	default:
		return fmt.Errorf("unknown example %q; use verify, authority or proof", args[0])
	}

	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return err
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from %s", resp.StatusCode, target)
	}

	var body struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return err
	}
	if _, ok := body.Data[field]; !ok {
		return fmt.Errorf("no data.%s in response", field)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, raw, "", "  "); err != nil {
		return err
	}
	fmt.Println(pretty.String())
	return nil
}

// baseURL returns APIER_BASE_URL without a trailing slash. Only https, or
// http on localhost, so a key is never sent to an untrusted host.
func baseURL() (string, error) {
	raw := os.Getenv("APIER_BASE_URL")
	if raw == "" {
		raw = "https://www.apier.no"
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid APIER_BASE_URL: %w", err)
	}
	host := u.Hostname()
	localHTTP := u.Scheme == "http" && (host == "localhost" || host == "127.0.0.1")
	if u.Scheme != "https" && !localHTTP {
		return "", fmt.Errorf("APIER_BASE_URL must use https:// (or http://localhost for local dev), got %q", raw)
	}
	return strings.TrimRight(raw, "/"), nil
}
