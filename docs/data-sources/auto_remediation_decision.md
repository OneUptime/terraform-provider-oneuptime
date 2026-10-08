---
page_title: "oneuptime_auto_remediation_decision Data Source - oneuptime"
subcategory: "Other"
description: |-
  Why auto-remediation did, or did not, act on an incident or alert: what each fix path did, or why it did nothing, each time the rule engine evaluated the signal.
---

# oneuptime_auto_remediation_decision (Data Source)

Why auto-remediation did, or did not, act on an incident or alert: what each fix path did, or why it did nothing, each time the rule engine evaluated the signal.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one auto remediation decision may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_auto_remediation_decision" "example" {
  incident_id = oneuptime_incident.example.id
}

# Or by id:
data "oneuptime_auto_remediation_decision" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of the alert this decision is about. The ID of a `oneuptime_alert`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this decision is about. The ID of a `oneuptime_incident`.
- `stage` (String) WaitingForInvestigation while remediation waits for the signal's AI investigation to finish, Evaluated once every fix path was evaluated.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `entries` (String) One entry per fix path - a Kubernetes cluster, an infrastructure resource, an Auto Remediation Rule, or the project - with a reason code saying what it did or why it did nothing. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
