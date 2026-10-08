---
page_title: "oneuptime_incident_template Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident templates for your project
---

# oneuptime_incident_template (Data Source)

Manage incident templates for your project

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_template" "example" {
  title = "example-title"
}

# Or by id:
data "oneuptime_incident_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `change_monitor_status_to_id` (String) Relation to Monitor Status Object ID. All monitors connected to this incident will be changed to this status when the incident is created. The ID of a `oneuptime_monitor_status`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Short description of this incident. This is in markdown and will be visible on the status page.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_severity_id` (String) Incident Severity ID. The ID of a `oneuptime_incident_severity`.
- `initial_incident_state_id` (String) Relation to Incident State Object ID. Incidents created from this template will start in this state. The ID of a `oneuptime_incident_state`.
- `is_scoped_to_status_pages` (Boolean) Whether incidents declared from this template are limited to the status pages in Status Pages. Derived from Status Pages; any value sent for it is ignored.
- `slug` (String) Friendly globally unique name for your object.
- `template_description` (String) Description of the Incident Template.
- `template_name` (String) Name of the Incident Template.
- `title` (String) Title of this incident.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `custom_field_settings` (String) How the Declare Incident form treats each incident custom field when an incident is declared from this template, keyed by the field's template variable key (variableKey). Each value is Required (asked, must be filled in), Optional (asked, may be left empty), Hidden (not asked; the field keeps this template's value) or Default. A field that is not listed, or is Default, follows its own Show on Create and Required on Create settings. Only the dashboard's Declare Incident form applies these settings: incidents created through the API are not checked against them. A JSON value: write it with `jsonencode()`.
- `custom_fields` (String) The custom field values incidents declared from this template start with, keyed by each incident custom field's name. They are merged one field at a time under the values the request or the Declare Incident form supplies. A JSON value: write it with `jsonencode()`.
- `docker_hosts` (Set of String) List of Docker hosts to pre-populate on incidents created from this template. IDs of `oneuptime_docker_host` resources.
- `hosts` (Set of String) List of hosts to pre-populate on incidents created from this template. IDs of `oneuptime_host` resources.
- `kubernetes_clusters` (Set of String) List of Kubernetes clusters to pre-populate on incidents created from this template. IDs of `oneuptime_kubernetes_cluster` resources.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) List of monitors affected by this incident. IDs of `oneuptime_monitor` resources.
- `on_call_duty_policies` (Set of String) List of on-call duty policies affected by this incident template. IDs of `oneuptime_on_call_policy` resources.
- `podman_hosts` (Set of String) List of Podman hosts to pre-populate on incidents created from this template. IDs of `oneuptime_podman_host` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `services` (Set of String) List of services to pre-populate on incidents created from this template. IDs of `oneuptime_service` resources.
- `status_pages` (Set of String) Limit incidents declared from this template to these status pages. Leave empty to reach every status page that lists the incident's monitors. IDs of `oneuptime_status_page` resources.
- `updated_at` (String) Date and Time when the object was updated.
