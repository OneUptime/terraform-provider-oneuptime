---
page_title: "oneuptime_rule_exception Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  What a detection rule skips: an entity, a value or a time window, with who added it, why, and when it expires.
---

# oneuptime_rule_exception (Resource)

What a detection rule skips: an entity, a value or a time window, with who added it, why, and when it expires.

## Example Usage

```terraform
resource "oneuptime_rule_exception" "example" {
  detection_rule_id = oneuptime_detection_rule.example.id
  kind              = "Example short text"
  reason            = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `detection_rule_id` (String) ID of the detection rule that skips what this names. The ID of a `oneuptime_detection_rule`.
- `kind` (String) What the rule skips: Entity (a user, host or IP address), Value (one event field's value) or Time Window (every event between two times).
- `reason` (String) Why it is safe for the rule to skip this, for whoever reads the list later.

### Optional

- `ends_at` (String) For a Time Window: when the rule sees events again. A window is at most 30 days.
- `entity_type` (String) For an Entity: User, Host or IP Address. The rule skips every event that names it, as who acted, as who was acted on, or among its observables.
- `expires_at` (String) For an Entity or a Value: when the rule sees it again. Empty: until someone deletes the exception.
- `field_name` (String) For a Value: the event field to match, named as in a Sigma rule - principalUser, processName, an attribute key such as http.user_agent.
- `starts_at` (String) For a Time Window: the first moment the rule skips.
- `value` (String) For an Entity or a Value: what the rule skips, matched whatever the case. * stands for any run of characters (svc-* skips every service account), and an IP address may be a range such as 10.20.0.0/16.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who added the exception. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this exception belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing rule exception by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_rule_exception.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_rule_exception.example <id>
```
