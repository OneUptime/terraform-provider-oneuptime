---
page_title: "oneuptime_incident_alert Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Alerts linked to an incident. Link alerts to an incident to track them as part of that incident's response.
---

# oneuptime_incident_alert (Data Source)

Alerts linked to an incident. Link alerts to an incident to track them as part of that incident's response.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident alert may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_alert" "example" {
  incident_id = oneuptime_incident.example.id
}

# Or by id:
data "oneuptime_incident_alert" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of the Alert that is linked to the incident. The ID of a `oneuptime_alert`.
- `created_by_user_id` (String) ID of the User who linked the alert to the incident (if it was linked by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the Incident this alert is linked to. The ID of a `oneuptime_incident`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
