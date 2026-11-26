#!/usr/bin/env python3
"""Verify: registration status and signals for one company.

Keyless: public sandbox, test company 999999999 (fictional data).
Standard library only.
"""
import json
import os
import sys
import urllib.error
import urllib.parse
import urllib.request


def base_url() -> str:
    """APIER_BASE_URL without a trailing slash. Only https, or http on
    localhost, so a key is never sent to an untrusted host."""
    raw = os.environ.get("APIER_BASE_URL", "https://www.apier.no")
    u = urllib.parse.urlparse(raw)
    is_local_http = u.scheme == "http" and u.hostname in ("localhost", "127.0.0.1")
    if u.scheme != "https" and not is_local_http:
        raise ValueError(
            f"APIER_BASE_URL must use https:// (or http://localhost for local dev). Got: {raw}"
        )
    return raw.rstrip("/")


def main() -> int:
    url = f"{base_url()}/api/v1/sandbox/public/company/999999999/verify"
    try:
        with urllib.request.urlopen(url, timeout=30) as r:
            body = json.load(r)
    except urllib.error.HTTPError as e:
        print(f"HTTP {e.code} from {url}", file=sys.stderr)
        return 1
    except urllib.error.URLError as e:
        print(f"Request failed: {e.reason}", file=sys.stderr)
        return 1
    if not body.get("data", {}).get("verification_status"):
        print("No data.verification_status in response", file=sys.stderr)
        return 1
    print(json.dumps(body, indent=2, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    sys.exit(main())
