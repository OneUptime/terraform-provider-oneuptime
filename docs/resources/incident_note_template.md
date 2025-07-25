---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_incident_note_template Resource - oneuptime"
subcategory: ""
description: |-
  Incident note template resource
---

# oneuptime_incident_note_template (Resource)

Incident note template resource

## Example Usage

```terraform
resource "oneuptime_incident_note_template" "example" {
  template_name = "example-template_name"
  template_description = "example-template_description"
}
```

## Schema

- `id` (String) Unique identifier for the resource. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `note` (String) Note. Optional.
- `template_name` (String) Name. Required.
- `template_description` (String) Template Description. Required.
- `created_at` (Map) A date time object.. Computed.
- `updated_at` (Map) A date time object.. Computed.
- `deleted_at` (Map) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_note_template.example <id>
```
