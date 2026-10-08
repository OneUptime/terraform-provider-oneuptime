---
page_title: "oneuptime_incident_template Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident templates for your project
---

# oneuptime_incident_template (Resource)

Manage incident templates for your project

## Example Usage

```terraform
resource "oneuptime_incident_template" "example" {
  title                = "This is an example of longer text content that might be stored in this field."
  template_name        = "Example short text"
  template_description = "This is an example of longer text content that might be stored in this field."
  description          = "Managed by Terraform"
}
```

## Schema

### Required

- `template_description` (String) Description of the Incident Template.
- `template_name` (String) Name of the Incident Template.
- `title` (String) Title of this incident.

### Optional

- `change_monitor_status_to_id` (String) Relation to Monitor Status Object ID. All monitors connected to this incident will be changed to this status when the incident is created. The ID of a `oneuptime_monitor_status`.
- `custom_field_settings` (String) How the Declare Incident form treats each incident custom field when an incident is declared from this template, keyed by the field's template variable key (variableKey). Each value is Required (asked, must be filled in), Optional (asked, may be left empty), Hidden (not asked; the field keeps this template's value) or Default. A field that is not listed, or is Default, follows its own Show on Create and Required on Create settings. Only the dashboard's Declare Incident form applies these settings: incidents created through the API are not checked against them. A JSON value: write it with `jsonencode()`.
- `custom_fields` (String) The custom field values incidents declared from this template start with, keyed by each incident custom field's name. They are merged one field at a time under the values the request or the Declare Incident form supplies. A JSON value: write it with `jsonencode()`.
- `description` (String) Short description of this incident. This is in markdown and will be visible on the status page.
- `docker_hosts` (Set of String) List of Docker hosts to pre-populate on incidents created from this template. IDs of `oneuptime_docker_host` resources.
- `hosts` (Set of String) List of hosts to pre-populate on incidents created from this template. IDs of `oneuptime_host` resources.
- `incident_severity_id` (String) Incident Severity ID. The ID of a `oneuptime_incident_severity`.
- `initial_incident_state_id` (String) Relation to Incident State Object ID. Incidents created from this template will start in this state. The ID of a `oneuptime_incident_state`.
- `kubernetes_clusters` (Set of String) List of Kubernetes clusters to pre-populate on incidents created from this template. IDs of `oneuptime_kubernetes_cluster` resources.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) List of monitors affected by this incident. IDs of `oneuptime_monitor` resources.
- `on_call_duty_policies` (Set of String) List of on-call duty policies affected by this incident template. IDs of `oneuptime_on_call_policy` resources.
- `podman_hosts` (Set of String) List of Podman hosts to pre-populate on incidents created from this template. IDs of `oneuptime_podman_host` resources.
- `services` (Set of String) List of services to pre-populate on incidents created from this template. IDs of `oneuptime_service` resources.
- `status_pages` (Set of String) Limit incidents declared from this template to these status pages. Leave empty to reach every status page that lists the incident's monitors. IDs of `oneuptime_status_page` resources.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_scoped_to_status_pages` (Boolean) Whether incidents declared from this template are limited to the status pages in Status Pages. Derived from Status Pages; any value sent for it is ignored.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_template.example <id>
```
