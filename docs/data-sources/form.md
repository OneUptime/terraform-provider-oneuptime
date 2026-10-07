---
page_title: "oneuptime_form Data Source - oneuptime"
subcategory: "Other"
description: |-
  Forms anyone with the link can fill in, without a OneUptime account. Each submission creates an incident or a scheduled maintenance event in this project.
---

# oneuptime_form (Data Source)

Forms anyone with the link can fill in, without a OneUptime account. Each submission creates an incident or a scheduled maintenance event in this project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_form" "by_name" {
  name = "example-form"
}

data "oneuptime_form" "by_id" {
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
- `description` (String) Shown at the top of the form's public page, above the questions: what the form is for and what happens after it is sent. Markdown... Computed.
- `is_enabled` (Bool) Whether the form's link works. While it is off, the public page shows a not-available message and nothing can be submitted... Computed.
- `share_key` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `target_type` (String) What each submission creates: Incident, or ScheduledMaintenance (a scheduled maintenance event)... Computed.
- `fields` (String) The questions the form asks, in order. Each has an id, a source (Question: one of the form's own, answered by type; TargetField: a built-in field of what the form creates, by targetField; TargetCustomField: one of its custom fields, by customFieldId; Submitter: the submitter's Name or Email), a label, optional help text, isRequired and isHidden (not shown on the public form, and answered only from the template a submission started from; never required, and never a field the target cannot be created without). A new form starts with a title, a description and the submitter's name and email... Computed.
- `templates` (String) Named sets of answers a submission can start from, in the order the form lists them. Each has an id, a name (unique within the form), isDefault (the form opens with it; at most one template) and answers: an object keyed by question id, each answer as a submission sends it - text, a number, true or false, an option's value, or a list of values for a multi-select. The public form lists the templates above its questions and fills in a template's answers when one is chosen, or when its link names one (?template=<id>). Hidden questions are answered only from the template a submission started from... Computed.
- `target_settings` (String) What every submission starts with besides the answers. For incidents: defaultTitle, incidentSeverityId, incidentTemplateId, monitorIds, labelIds, onCallDutyPolicyIds, ownerUserIds and ownerTeamIds. For scheduled maintenance events: defaultTitle, monitorIds, statusPageIds, labelIds, ownerUserIds, ownerTeamIds, showOnStatusPages and notifySubscribers... Computed.
- `success_message` (String) Shown after the form is submitted, together with the number of what the submission created. Markdown... Computed.
- `ip_whitelist` (String) The networks the form can be opened and submitted from: one IPv4 or IPv6 address, or one IPv4 range in CIDR notation (such as 10.0.0.0/8), per line. IPv6 ranges are not supported. Leave it empty to allow any network... Computed.
- `logo_file_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `logo_alt_text` (String) What the logo says, read out by screen readers: usually your organization's name. Leave it empty and screen readers skip the logo... Computed.
- `favicon_file_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
