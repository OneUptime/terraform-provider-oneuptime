---
page_title: "oneuptime_incident_postmortem_template Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage postmortem templates for your incidents
---

# oneuptime_incident_postmortem_template (Data Source)

Manage postmortem templates for your incidents

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident postmortem template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_postmortem_template" "example" {
  postmortem_note = "example-postmortem-note"
}

# Or by id:
data "oneuptime_incident_postmortem_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `postmortem_note` (String) Markdown template used when documenting an incident postmortem.
- `template_description` (String) Description of the Postmortem Template.
- `template_name` (String) Name of the Postmortem Template.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
