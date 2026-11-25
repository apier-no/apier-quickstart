#!/usr/bin/env python3
"""Proof: an Authority answer wrapped in a signed certificate (?certificate=true).

Without APIER_API_KEY this calls the bearer sandbox. The bearer is invented
here, needs no signup, and the certificate is signed with the public sandbox
key, so it is for rehearsing verification code only.

With APIER_API_KEY and APIER_ORG_NUMBER set it calls production. That call
needs a real key and is metered: https://www.apier.no/pricing

Standard library only.
"""
import json
import os
import re
import sys
import urllib.error
import urllib.parse
import urllib.request
import uuid


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
    base = base_url()
    api_key = os.environ.get("APIER_API_KEY", "")
    if api_key:
        org = os.environ.get("APIER_ORG_NUMBER", "")
        if not re.fullmatch(r"\d{9}", org):
            print(
                "Set APIER_ORG_NUMBER to a 9-digit organisasjonsnummer when APIER_API_KEY is set.",
                file=sys.stderr,
            )
            return 1
        url = f"{base}/api/v1/company/{org}/authority?certificate=true"
        bearer = api_key
    else:
        url = f"{base}/api/v1/sandbox/company/818000006/authority?certificate=true"
        bearer = f"apier_sandbox_test_{uuid.uuid4().hex}"

    req = urllib.request.Request(url, headers={"Authorization": f"Bearer {bearer}"})
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            body = json.load(r)
    except urllib.error.HTTPError as e:
        print(f"HTTP {e.code} from {url}", file=sys.stderr)
        return 1
    except urllib.error.URLError as e:
        print(f"Request failed: {e.reason}", file=sys.stderr)
        return 1
    if not body.get("data", {}).get("signature"):
        print("No data.signature in response", file=sys.stderr)
        return 1
    print(json.dumps(body, indent=2, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    sys.exit(main())
