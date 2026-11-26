#!/usr/bin/env bash
# Authority: the register's coded signing rules (signaturrett and prokura) for one company.
# Keyless: public sandbox, test company 999999999 (fictional data).
set -euo pipefail

BASE="${APIER_BASE_URL:-https://www.apier.no}"
BASE="${BASE%/}"
# https only, or http://localhost for local development.
case "$BASE" in
  https://* | http://localhost | http://localhost:* | http://127.0.0.1 | http://127.0.0.1:*) ;;
  *)
    echo "APIER_BASE_URL must use https:// (or http://localhost for local dev). Got: $BASE" >&2
    exit 1
    ;;
esac

curl -sS --fail-with-body "$BASE/api/v1/sandbox/public/company/999999999/authority"
echo
