---
page_title: "oneuptime_rule_exception Data Source - oneuptime"
subcategory: "Security"
description: |-
  What a detection rule skips: an entity, a value or a time window, with who added it, why, and when it expires.
---

# oneuptime_rule_exception (Data Source)

What a detection rule skips: an entity, a value or a time window, with who added it, why, and when it expires.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one rule exception may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_rule_exception" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_rule_exception" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who added the exception. The ID of a `oneuptime_user` (see the data source).
- `detection_rule_id` (String) ID of the detection rule that skips what this names. The ID of a `oneuptime_detection_rule`.
- `entity_type` (String) For an Entity: User, Host or IP Address. The rule skips every event that names it, as who acted, as who was acted on, or among its observables.
- `field_name` (String) For a Value: the event field to match, named as in a Sigma rule - principalUser, processName, an attribute key such as http.user_agent.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `kind` (String) What the rule skips: Entity (a user, host or IP address), Value (one event field's value) or Time Window (every event between two times).
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `reason` (String) Why it is safe for the rule to skip this, for whoever reads the list later.
- `value` (String) For an Entity or a Value: what the rule skips, matched whatever the case. * stands for any run of characters (svc-* skips every service account), and an IP address may be a range such as 10.20.0.0/16.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ends_at` (String) For a Time Window: when the rule sees events again. A window is at most 30 days.
- `expires_at` (String) For an Entity or a Value: when the rule sees it again. Empty: until someone deletes the exception.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this exception belongs to. The ID of a `oneuptime_project`.
- `starts_at` (String) For a Time Window: the first moment the rule skips.
- `updated_at` (String) Date and Time when the object was updated.
