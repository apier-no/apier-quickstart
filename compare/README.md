# Raw Brreg vs Apier: signing authority for a Norwegian company

Brønnøysundregistrene (Brreg) is the source of this data. Apier reads it from Brreg's open APIs under NLOD and returns it in a different shape. Brreg is not a competitor to Apier, and for many uses the open APIs are all you need.

This page records what each API returned on 2026-11-25, so you can see what you get from each and decide which to call. Everything below comes from the JSON observed that day. Names and birth dates of people are redacted here; the APIs return them in full.

## The company

TELENOR ASA, organisasjonsnummer 982463718. It was picked because its registration has both prokura (one holder alone and two holders jointly) and a signature rule that includes joint signing.

## 1. Enhetsregisteret roles API

```bash
curl -s https://data.brreg.no/enhetsregisteret/api/enheter/982463718/roller
```

HTTP 200, 8,308 bytes, 0.31 s (single sample).

What it returned: three role groups.

| Group (`rollegrupper[].type.kode`) | Roles in the group |
|---|---|
| `DAGL` Daglig leder | 1 person |
| `STYR` Styre | `LEDE` Styrets leder, 9 `MEDL` Styremedlem, 7 `VARA` Varamedlem |
| `REVI` Revisor | 1 company |

Shape of one group:

```json
{
  "type": { "kode": "DAGL", "beskrivelse": "Daglig leder" },
  "sistEndret": "2024-12-02",
  "roller": [
    {
      "type": { "kode": "DAGL", "beskrivelse": "Daglig leder" },
      "person": {
        "navn": { "fornavn": "<redacted>", "mellomnavn": "<redacted>", "etternavn": "<redacted>" },
        "fodselsdato": "<redacted>",
        "erDoed": false
      },
      "avregistrert": false,
      "rekkefolge": 0
    }
  ]
}
```

What it did not return: any signature or prokura group, and any text or code for the signing rule. The response says who holds which role. It does not say who may sign for the company or in which combination. The role types this API publishes (`/enhetsregisteret/api/roller/rolletyper` and `/roller/rollegruppetyper`) contained no signature or prokura code on the same day.

## 2. Brreg Fullmakt API (signatur and prokura)

Brreg publishes signing rules in a second open API, with one call for signatur and one for prokura.

```bash
curl -s https://data.brreg.no/fullmakt/enheter/982463718/signatur
curl -s https://data.brreg.no/fullmakt/enheter/982463718/prokura
```

signatur: HTTP 200, 3,676 bytes, 0.26 s. prokura: HTTP 200, 2,625 bytes, 0.25 s (single samples).

Top-level keys in both: `_links`, `enhet`, `oppslagsData`, `signeringsGrunnlag`, `signeringsKombinasjon`, `status`. In the excerpts below, `_links` and `enhet` are left out and lists of people are abridged to a count.

### signatur

```json
{
  "oppslagsData": { "oppslagsIdent": "10921515", "oppslagsTidspunkt": "2026-11-25T02:49:47" },
  "signeringsGrunnlag": {
    "kode": "SF",
    "tekstforklaring": "Signeringsgrunnlag finnes",
    "muligeSigneringsRoller": {
      "personRolleGrunnlag": [
        { "navn": "<redacted>", "fodselsdato": "<redacted>", "rolle": { "kode": "DAGL", "tekstforklaring": "Daglig leder" } }
      ]
    },
    "signaturProkuraRoller": {
      "signaturProkuraFritekst": "Styrets leder alene, daglig leder alene eller to aksjonærvalgte styremedlemmer i fellesskap."
    }
  },
  "signeringsKombinasjon": {
    "kombinasjon": [
      { "kode": "STYR", "kombinasjonsId": "1", "tekstforklaring": "Styret i fellesskap", "personRolleKombinasjon": ["10 people, same shape as above"] }
    ]
  },
  "status": {
    "kombinasjonStatus": { "kode": "OK", "tekstforklaring": "Kombinasjon finnes." },
    "regelStatus": { "kode": "RI", "tekstforklaring": "Regel er ikke registrert eller standardisert.", "regelIdent": null },
    "rutineStatus": { "kode": "OK", "tekstforklaring": "Maskinelt oppslag utført uten feil." }
  }
}
```

`muligeSigneringsRoller.personRolleGrunnlag` listed 11 people (`DAGL`, `LEDE` and 9 `MEDL`); one is shown.

Two things to note in this answer:

- The free text names three alternatives (the chair alone, the daglig leder alone, or two shareholder-elected board members jointly). The coded part holds one combination, "Styret i fellesskap".
- `regelStatus.kode` is `RI`: the rule is not registered or standardised, and `regelIdent` is null.

### prokura

```json
{
  "signeringsGrunnlag": {
    "kode": "SF",
    "signaturProkuraRoller": {
      "personRolleGrunnlag": ["3 people: one PROK, two POFE"],
      "signaturProkuraFritekst": "Daglig leder alene."
    }
  },
  "signeringsKombinasjon": {
    "kombinasjon": [
      { "kode": "PROK", "kombinasjonsId": "1", "tekstforklaring": "Prokura", "personRolleKombinasjon": ["1 person"] },
      { "kode": "POFE", "kombinasjonsId": "1", "tekstforklaring": "Prokura i fellesskap", "personRolleKombinasjon": ["2 people"] },
      { "kode": "R0002", "kombinasjonsId": "1", "tekstforklaring": "Daglig leder alene", "personRolleKombinasjon": ["1 person, role DAGL"] }
    ]
  },
  "status": {
    "kombinasjonStatus": { "kode": "OK", "tekstforklaring": "Kombinasjon finnes." },
    "regelStatus": { "kode": "RF", "tekstforklaring": "Regel funnet.", "regelIdent": "R0002" },
    "rutineStatus": { "kode": "OK", "tekstforklaring": "Maskinelt oppslag utført uten feil." }
  }
}
```

### A second company, for contrast

For EQUINOR ASA (923609016) both Fullmakt calls returned HTTP 200 with no `signeringsGrunnlag` and no `signeringsKombinasjon`. The only status was:

```json
{ "rutineStatus": { "kode": "NA", "tekstforklaring": "Generelle unntak, finnes roller som har status det ikke kan foretas maskinelt oppslag mot." } }
```

So an HTTP 200 from this API does not always carry a signing rule.

## 3. What you build on top of the raw APIs

Based only on the responses above, a developer who wants "who may sign for this company" from the open APIs writes code for the following.

- Three requests per company (roles, signatur, prokura) and the join between them.
- The three status objects (`rutineStatus`, `regelStatus`, `kombinasjonStatus`) and their codes. That includes the `NA` case, where the answer is HTTP 200 with no signing data, and the `RI` case, where there is no standardised rule.
- Cases where the free text says more than the coded combinations, as in Telenor's signatur answer.
- Field formats that differ between the two APIs. The roles API returns the name as an object and `erDoed` as a boolean. The Fullmakt API returns the name as one string, the birth date as `dd.mm.yyyy`, and `konkurs` in `enhet` as the string `"false"`.
- Personal data. Every response carries names and birth dates, so your storage and logs need a legal basis for it.
- History. Each response describes the state at `oppslagsTidspunkt` and holds no earlier state. If you need to show later what the register said on a given day, you store the answers yourself.
- Evidence. The responses are plain JSON over HTTPS and contain no signature field.

None of this is a fault in the Brreg APIs. They publish the register. The list is the work that sits between the register and an application that acts on it.

## 4. Apier Authority response (sandbox shape)

```bash
curl -s https://www.apier.no/api/v1/sandbox/public/company/999999999/authority
```

HTTP 200, 1,216 bytes, 0.45 s (single sample). This is the keyless public sandbox. The company is fictional and the answer is a fixture, so the timing is not comparable with a live register lookup. It is listed only so that every call on this page has its measured number.

```json
{
  "success": true,
  "data": {
    "org_number": "999999999",
    "name": "Apier Public Sandbox AS",
    "entity_type": "AS",
    "classification": "sole",
    "summary": "Én rolle kan signere alene for selskapet.",
    "kombinasjon_available": true,
    "combinations": [
      { "kombinasjons_id": "SANDBOX-KOMB-DAGL-ALENE", "description": "Daglig leder signerer alene.", "role_codes": ["DAGL"] }
    ],
    "signaturrett_holders": [{ "role": "DAGL", "name": "Ola Nordmann" }],
    "prokura_holders": [],
    "data_sources": ["brreg", "fullmakttjenesten"],
    "resolved_at": "2026-01-01T09:00:00+01:00",
    "country": "NO",
    "source": "brreg.fullmakt",
    "rule_set": "no.brreg.signatur.v1",
    "rule_set_version": "1.0.0"
  },
  "_meta": {
    "data_freshness": "2026-01-01T00:00:00+01:00",
    "last_verified": "2026-01-01T00:00:00+01:00",
    "source": "apier.no",
    "is_sandbox": true,
    "untrusted_text_fields": ["/data/combinations/0/description", "/data/name", "/data/signaturrett_holders/0/name", "/data/signaturrett_holders/0/role"]
  }
}
```

Trimmed: `software_version` and six `_meta` fields are omitted.

The bearer sandbox has a fixture with joint signing and prokura, org 818000006. Its answer has `"classification": "joint"`, one combination with `"role_codes": ["DAGL", "STYRL"]`, and one with `"role_codes": ["PROKU"]`. With `?certificate=true` the same answer comes wrapped in a signed certificate; see [`examples/`](../examples) for the call.

What the Apier response adds to the raw data, as observed in the sandbox shape:

- One request returns signatur and prokura together.
- `classification` is one of a fixed set of values (`sole`, `joint`, `by_role`, `prokura_only`, `no_authority`, `unknown`).
- `combinations[].role_codes` gives each combination as a list of role codes.
- `data_sources`, `resolved_at` and `_meta.data_freshness` state where the answer came from and when.
- `_meta.untrusted_text_fields` lists the fields whose text was written into the register by a third party, so an agent can treat them as data.
- On request, a signed certificate that can be checked offline against a published key set.

What it does not add: a decision. Apier returns the register's coding of the signing rules, uninterpreted. Whether a given person may sign a given document is for your application to work out.

The production response was not called for this page. It needs an API key and each call is priced; see https://www.apier.no/pricing. Its fields are documented at https://www.apier.no/docs/guides/signature-rights-norwegian-company.

## Summary

| | Enhetsregisteret roles API | Brreg Fullmakt API | Apier Authority |
|---|---|---|---|
| Requests for signatur and prokura | does not return them | 2 | 1 |
| Role holders | yes | those in the signing rule | those in the signing rule |
| Signing rule as codes | no | yes, when `regelStatus` is `RF` | `classification` and `combinations[].role_codes` |
| Signing rule as free text | no | `signaturProkuraFritekst` | `combinations[].description` |
| Source and retrieval time in the body | no | `oppslagsTidspunkt` | `data_sources`, `resolved_at`, `_meta.data_freshness` |
| Signature on the answer | no | no | on request (`?certificate=true`) |
| Key | none | none | none in the sandbox; API key in production |
| Price | free | free | free in the sandbox; per call in production |

## When the raw APIs are the better choice

- You need role holders only, not signing rules.
- You want no third party between you and the register.
- You are prepared to handle the status codes and free text yourself and do not need a signed record of the answer.

Brreg's documentation for both APIs is at https://data.brreg.no.
