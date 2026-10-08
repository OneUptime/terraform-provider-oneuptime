---
page_title: "oneuptime_threat_intel_indicator Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Threat Intel Indicator
---

# oneuptime_threat_intel_indicator (Data Source)

API endpoints for Threat Intel Indicator

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one threat intel indicator may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_threat_intel_indicator" "example" {
  feed_id = "example-feed-id"
}

# Or by id:
data "oneuptime_threat_intel_indicator" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `confidence` (Number) Confidence.
- `feed_id` (String) Feed ID.
- `feed_name` (String) Feed.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `indicator_name` (String) Name.
- `indicator_type` (String) Indicator Type.
- `indicator_value` (String) Indicator Value.
- `revoked` (Boolean) Revoked.
- `stix_id` (String) STIX ID.
- `valid_from` (String) Valid From.
- `valid_until` (String) Valid Until.
- `version` (String) Version.

### Read-Only

- `project_id` (String) Project ID.
- `retention_date` (String) Retention Date.
- `stix_labels` (Set of String) Labels.
