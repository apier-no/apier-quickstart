#!/usr/bin/env bash
# Proof: an Authority answer wrapped in a signed certificate (?certificate=true).
#
# Without APIER_API_KEY this calls the bearer sandbox. The bearer is invented
# here, needs no signup, and the certificate is signed with the public sandbox
# key, so it is for rehearsing verification code only.
#
# With APIER_API_KEY and APIER_ORG_NUMBER set it calls production. That call
# needs a real key and is metered: https://www.apier.no/pricing
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

if [ -n "${APIER_API_KEY:-}" ]; then
  ORG="${APIER_ORG_NUMBER:-}"
  if ! [[ "$ORG" =~ ^[0-9]{9}$ ]]; then
    echo "Set APIER_ORG_NUMBER to a 9-digit organisasjonsnummer when APIER_API_KEY is set." >&2
    exit 1
  fi
  curl -sS --fail-with-body -H "Authorization: Bearer $APIER_API_KEY" \
    "$BASE/api/v1/company/$ORG/authority?certificate=true"
else
  SUFFIX="$(date +%s)$RANDOM"
  curl -sS --fail-with-body -H "Authorization: Bearer apier_sandbox_test_$SUFFIX" \
    "$BASE/api/v1/sandbox/company/818000006/authority?certificate=true"
fi
echo
