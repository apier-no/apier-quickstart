#!/usr/bin/env node
// @ts-check
// Proof: an Authority answer wrapped in a signed certificate (?certificate=true).
//
// Without APIER_API_KEY this calls the bearer sandbox. The bearer is invented
// here, needs no signup, and the certificate is signed with the public sandbox
// key, so it is for rehearsing verification code only.
//
// With APIER_API_KEY and APIER_ORG_NUMBER set it calls production. That call
// needs a real key and is metered: https://www.apier.no/pricing
//
// Node 20+, built-in fetch, no dependencies.

const base = baseUrl();
const apiKey = process.env.APIER_API_KEY ?? '';
/** @type {string} */
let url;
/** @type {string} */
let bearer;
if (apiKey) {
  const org = process.env.APIER_ORG_NUMBER ?? '';
  if (!/^\d{9}$/.test(org)) {
    console.error('Set APIER_ORG_NUMBER to a 9-digit organisasjonsnummer when APIER_API_KEY is set.');
    process.exit(1);
  }
  url = `${base}/api/v1/company/${org}/authority?certificate=true`;
  bearer = apiKey;
} else {
  url = `${base}/api/v1/sandbox/company/818000006/authority?certificate=true`;
  bearer = `apier_sandbox_test_${crypto.randomUUID()}`;
}

const res = await fetch(url, { headers: { Authorization: `Bearer ${bearer}` } });
if (!res.ok) {
  await res.body?.cancel();
  fail(`HTTP ${res.status} from ${res.url}`);
} else {
  /** @type {{ data?: Record<string, unknown> }} */
  const body = await res.json();
  if (!body?.data?.signature) fail('No data.signature in response');
  else console.log(JSON.stringify(body, null, 2));
}

/** @param {string} message */
function fail(message) {
  console.error(message);
  process.exitCode = 1;
}

/**
 * APIER_BASE_URL without a trailing slash. Only https, or http on localhost,
 * so a key is never sent to an untrusted host.
 * @returns {string}
 */
function baseUrl() {
  const raw = process.env.APIER_BASE_URL ?? 'https://www.apier.no';
  const u = new URL(raw);
  const isLocalHttp = u.protocol === 'http:' && (u.hostname === 'localhost' || u.hostname === '127.0.0.1');
  if (u.protocol !== 'https:' && !isLocalHttp) {
    throw new Error(`APIER_BASE_URL must use https:// (or http://localhost for local dev). Got: ${raw}`);
  }
  return raw.replace(/\/+$/, '');
}
