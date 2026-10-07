---
page_title: "oneuptime_form Resource - oneuptime"
subcategory: "Other"
description: |-
  Forms anyone with the link can fill in, without a OneUptime account. Each submission creates an incident or a scheduled maintenance event in this project.
---

# oneuptime_form (Resource)

Forms anyone with the link can fill in, without a OneUptime account. Each submission creates an incident or a scheduled maintenance event in this project.

## Example Usage

```terraform
resource "oneuptime_form" "example" {
  name = "Example short text"
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `name` (String) The form's name, shown as the heading of its public page. Unique within the project...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Shown at the top of the form's public page, above the questions: what the form is for and what happens after it is sent. Markdown...
- `is_enabled` (Bool) Whether the form's link works. While it is off, the public page shows a not-available message and nothing can be submitted...
- `target_type` (String) What each submission creates: Incident, or ScheduledMaintenance (a scheduled maintenance event)...
- `fields` (String) The questions the form asks, in order. Each has an id, a source (Question: one of the form's own, answered by type; TargetField: a built-in field of what the form creates, by targetField; TargetCustomField: one of its custom fields, by customFieldId; Submitter: the submitter's Name or Email), a label, optional help text and isRequired. A new form starts with a title, a description and the submitter's name and email...
- `target_settings` (String) What every submission starts with besides the answers. For incidents: defaultTitle, incidentSeverityId, incidentTemplateId, monitorIds, labelIds, onCallDutyPolicyIds, ownerUserIds and ownerTeamIds. For scheduled maintenance events: defaultTitle, monitorIds, statusPageIds, labelIds, ownerUserIds, ownerTeamIds, showOnStatusPages and notifySubscribers...
- `success_message` (String) Shown after the form is submitted, together with the number of what the submission created. Markdown...
- `ip_whitelist` (String) The networks the form can be opened and submitted from: one IPv4 or IPv6 address, or one IPv4 range in CIDR notation (such as 10.0.0.0/8), per line. IPv6 ranges are not supported. Leave it empty to allow any network...
- `logo_file_id` (String) A unique identifier for an object, represented as a UUID..
- `logo_alt_text` (String) What the logo says, read out by screen readers: usually your organization's name. Leave it empty and screen readers skip the logo...
- `favicon_file_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `share_key` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_form.example <id>
```
