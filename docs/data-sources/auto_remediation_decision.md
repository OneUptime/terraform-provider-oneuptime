---
page_title: "oneuptime_auto_remediation_decision Data Source - oneuptime"
subcategory: "Other"
description: |-
  Why auto-remediation did, or did not, act on an incident or alert: what each fix path did, or why it did nothing, each time the rule engine evaluated the signal.
---

# oneuptime_auto_remediation_decision (Data Source)

Why auto-remediation did, or did not, act on an incident or alert: what each fix path did, or why it did nothing, each time the rule engine evaluated the signal. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_auto_remediation_decision" "by_name" {
  name = "example-auto_remediation_decision"
}

data "oneuptime_auto_remediation_decision" "by_id" {
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
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `stage` (String) WaitingForInvestigation while remediation waits for the signal's AI investigation to finish, Evaluated once every fix path was evaluated... Computed.
- `entries` (String) One entry per fix path - a Kubernetes cluster, an infrastructure resource, an Auto Remediation Rule, or the project - with a reason code saying what it did or why it did nothing... Computed.
