---
page_title: "oneuptime_threat_intel_indicator Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Threat Intel Indicator
---

# oneuptime_threat_intel_indicator (Data Source)

API endpoints for Threat Intel Indicator Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_threat_intel_indicator" "by_name" {
  name = "example-threat_intel_indicator"
}

data "oneuptime_threat_intel_indicator" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `feed_id` (String) Feed ID. Computed.
- `feed_name` (String) Feed. Computed.
- `stix_id` (String) STIX ID. Computed.
- `indicator_type` (String) Indicator Type. Computed.
- `indicator_value` (String) Indicator Value. Computed.
- `indicator_name` (String) Name. Computed.
- `confidence` (Number) Confidence. Computed.
- `stix_labels` (Set) Labels. Computed.
- `valid_from` (String) Valid From. Computed.
- `valid_until` (String) Valid Until. Computed.
- `revoked` (Bool) Revoked. Computed.
- `version` (String) Version. Computed.
- `retention_date` (String) Retention Date. Computed.
