---
page_title: "oneuptime_huntress_incident_report Data Source - oneuptime"
subcategory: "Other"
description: |-
  Huntress incident reports received through a Huntress connection, each with the incident it opened or the reason it opened none.
---

# oneuptime_huntress_incident_report (Data Source)

Huntress incident reports received through a Huntress connection, each with the incident it opened or the reason it opened none.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one huntress incident report may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_huntress_incident_report" "example" {
  huntress_connection_id = oneuptime_huntress_connection.example.id
}

# Or by id:
data "oneuptime_huntress_incident_report" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `affected_name` (String) The host or identity the report is about, as Huntress names it in the report's subject.
- `huntress_account_id` (String) The id of the Huntress account the report belongs to, or empty when Huntress did not say.
- `huntress_connection_id` (String) ID of the connection that received this report, or empty once it is deleted. The ID of a `oneuptime_huntress_connection`.
- `huntress_incident_report_id` (String) The incident report's id in Huntress.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this report opened. Empty when it opened none, or when the incident was deleted. The ID of a `oneuptime_incident`.
- `last_event_type` (String) The last event Huntress sent about this report.
- `organization_id` (String) The id of the Huntress organization the report is about.
- `organization_name` (String) The name of the Huntress organization the report is about.
- `outcome` (String) What was done with the report: Opening, IncidentOpened, IncidentResolved, OrganizationNotWatched or ClosedBeforeReceived.
- `paged_on_call` (Boolean) Whether the incident was opened with the connection's on-call policies, which pages them.
- `severity` (String) The report's severity in Huntress: critical, high or low, as Huntress last sent it.
- `status` (String) The report's status in Huntress as it last sent it, such as sent, closed or dismissed.
- `subject` (String) The report's subject in Huntress.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_event_received_at` (String) When Huntress last sent an event about this report.
- `project_id` (String) ID of the project this report was received in. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
