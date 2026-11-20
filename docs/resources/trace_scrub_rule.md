---
page_title: "oneuptime_trace_scrub_rule Resource - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Configure rules to automatically detect and scrub sensitive data (PII) from spans at ingest time.
---

# oneuptime_trace_scrub_rule (Resource)

Configure rules to automatically detect and scrub sensitive data (PII) from spans at ingest time.

## Example Usage

```terraform
resource "oneuptime_trace_scrub_rule" "example" {
  name         = "Example trace scrub rule"
  pattern_type = "email"
  description  = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`sort_order` is this trace scrub rule's place among the trace scrub rules of its project, lowest first. Setting it places this trace scrub rule and moves no other, so a configuration that numbers every trace scrub rule of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.trace_scrub_rules`:

```terraform
sort_order = index(local.trace_scrub_rules, each.key) + 1
```

- Give each trace scrub rule a number of its own: two trace scrub rules at one number are shown the older first.
- A trace scrub rule created without `sort_order` goes after the last one of its project.
- Deleting a trace scrub rule leaves a gap in the numbers, which changes nothing: trace scrub rules go by their numbers, not by counting.
- Moving a trace scrub rule in the OneUptime Dashboard numbers the project's trace scrub rules 1..n again; the next plan shows that as changes to `sort_order`.


## Schema

### Required

- `name` (String) Friendly name for this scrub rule.
- `pattern_type` (String) The type of sensitive data to detect: email, creditCard, ssn, phoneNumber, ipAddress, sensitiveKeys (the whole value of every attribute whose key looks sensitive, such as password or token), or custom (the regular expression in customRegex). Any other value is refused: it would scrub nothing. Allowed values: `email`, `creditCard`, `ssn`, `phoneNumber`, `ipAddress`, `sensitiveKeys`, `custom`.

### Optional

- `custom_regex` (String) The regular expression a 'custom' rule scrubs, written without slashes or flags and matched case-sensitively. Required when patternType is 'custom': a pattern that is empty, does not compile, or matches empty text is refused. Ignored for the other pattern types.
- `description` (String) Description of what this scrub rule does.
- `fields_to_scrub` (String) Which span fields to scrub: 'name' (the span name), 'attributes' (attribute values), 'events' (span event attributes), or 'all', the default. A sensitiveKeys rule always scrubs attribute and event attribute values, whatever this says. Defaults to `all`.
- `is_enabled` (Boolean) Whether this scrub rule is active. Defaults to `true`.
- `scrub_action` (String) How to scrub matched data: 'redact' replaces it with [REDACTED] (the default), 'mask' partially hides it, 'hash' replaces it with a short hash of the value. Allowed values: `mask`, `hash`, `redact`. Defaults to `redact`.
- `sort_order` (Number) Where this rule is applied among the project's span scrub rules, lowest number first. A new rule is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this trace scrub rule. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this trace scrub rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing trace scrub rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_trace_scrub_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_trace_scrub_rule.example <id>
```
