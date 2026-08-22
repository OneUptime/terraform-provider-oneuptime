---
page_title: "oneuptime_alert_measurement_value Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  The computed value of one alert measurement for one alert, recomputed from the alert's timeline rather than accumulated
---

# oneuptime_alert_measurement_value (Data Source)

The computed value of one alert measurement for one alert, recomputed from the alert's timeline rather than accumulated Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_alert_measurement_value" "by_name" {
  name = "example-alert_measurement_value"
}

data "oneuptime_alert_measurement_value" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_measurement_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `started_at` (String) A date time object.. Computed.
- `ended_at` (String) A date time object.. Computed.
- `value_in_seconds` (Number) The measured duration in seconds. Only set when the status is Recorded - a measurement that could not be computed is left blank rather than written as zero... Computed.
- `status` (String) The outcome of evaluating this measurement: Recorded, Pending, Not Applicable or Invalid... Computed.
- `status_message` (String) Why this measurement has the status it has, in plain words - for example which anchor has not been reached yet, or by how much the end precedes the start... Computed.
- `start_alert_state_timeline_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `end_alert_state_timeline_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `computed_at` (String) A date time object.. Computed.
