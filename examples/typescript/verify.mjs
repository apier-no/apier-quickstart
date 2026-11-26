#!/usr/bin/env node
// @ts-check
// Verify: registration status and signals for one company.
// Keyless: public sandbox, test company 999999999 (fictional data).
// Node 20+, built-in fetch, no dependencies.

const res = await fetch(`${baseUrl()}/api/v1/sandbox/public/company/999999999/verify`);
if (!res.ok) {
  await res.body?.cancel();
  fail(`HTTP ${res.status} from ${res.url}`);
} else {
  /** @type {{ data?: Record<string, unknown> }} */
  const body = await res.json();
  if (!body?.data?.verification_status) fail('No data.verification_status in response');
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
