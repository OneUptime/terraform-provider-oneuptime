---
page_title: "oneuptime_log_scrub_rule Data Source - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Configure rules to automatically detect and scrub sensitive data (PII) from logs at ingest time.
---

# oneuptime_log_scrub_rule (Data Source)

Configure rules to automatically detect and scrub sensitive data (PII) from logs at ingest time.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one log scrub rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_log_scrub_rule" "example" {
  name = "Example log scrub rule"
}

# Or by id:
data "oneuptime_log_scrub_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this log scrub rule. The ID of a `oneuptime_user` (see the data source).
- `custom_regex` (String) The regular expression a 'custom' rule scrubs, written without slashes or flags and matched case-sensitively. Required when patternType is 'custom': a pattern that is empty, does not compile, or matches empty text is refused. Ignored for the other pattern types.
- `description` (String) Description of what this scrub rule does.
- `fields_to_scrub` (String) Which log fields to scrub: 'body' (the log message), 'attributes' (attribute values), or 'both', the default. A sensitiveKeys rule always scrubs attribute values, whatever this says.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this scrub rule is active.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Friendly name for this scrub rule.
- `pattern_type` (String) The type of sensitive data to detect: email, creditCard, ssn, phoneNumber, ipAddress, sensitiveKeys (the whole value of every attribute whose key looks sensitive, such as password or token), or custom (the regular expression in customRegex). Any other value is refused: it would scrub nothing.
- `scrub_action` (String) How to scrub matched data: 'redact' replaces it with [REDACTED] (the default), 'mask' partially hides it, 'hash' replaces it with a short hash of the value.
- `service_id` (String) ID of the service whose logs this rule scrubs: a rule that names a service scrubs only the logs that service sends. Null scrubs the logs of every service in the project. The ID of a `oneuptime_service`.
- `sort_order` (Number) Where this rule is applied among the project's log scrub rules, lowest number first. A new rule is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them.

### Read-Only

- `attribute_keys` (String) The attributes this rule reads, each by its key or a dotted prefix of it (customer_name, request.headers). A Named Attributes rule hides the whole value of each; any other pattern looks for its matches only in their values. Empty for every attribute. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this log scrub rule belongs to. The ID of a `oneuptime_project`.
- `proof` (String) Written by OneUptime: how this rule was proven to hide what Data Security found, on the very logs it was found in - the records read, the values that held the data and the share of other values it changed. It counts only while the rule is unchanged. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
