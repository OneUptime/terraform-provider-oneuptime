---
page_title: "oneuptime_investigation_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Choose which new incidents or alerts OneUptime AI investigates on its own. With no rule, every one is investigated; with rules, only those that match one.
---

# oneuptime_investigation_rule (Data Source)

Choose which new incidents or alerts OneUptime AI investigates on its own. With no rule, every one is investigated; with rules, only those that match one.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one investigation rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_investigation_rule" "example" {
  name = "Example investigation rule"
}

# Or by id:
data "oneuptime_investigation_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this investigation rule.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `llm_provider_id` (String) ID of the LLM provider OneUptime AI uses for everything this rule starts. Empty: the project's default provider, else the global one. The ID of a `oneuptime_llm_provider`.
- `name` (String) Name of this investigation rule.
- `only_if_open_after_minutes` (Number) Investigate only incidents or alerts still open this many minutes after they were created: OneUptime AI waits that long, then investigates if it is still open. Leave empty to investigate right away. When several rules match, the shortest wait decides. At most 1440 minutes (a day).
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title.
- `trigger_entity_type` (String) Which kind of new signal this rule decides about: Incident or Alert.

### Read-Only

- `alert_severities` (Set of String) Match only alerts with these severities (alert rules only). Leave empty to match any severity. IDs of `oneuptime_alert_severity` resources.
- `code_repositories` (Set of String) The code repositories OneUptime AI may read while it investigates what this rule matches. Leave empty to let it read every repository the project has connected. IDs of `oneuptime_code_repository` resources.
- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `external_mcp_servers` (Set of String) The external MCP servers OneUptime AI may use while it investigates what this rule matches. Leave empty to let it use all of the project's MCP servers. IDs of `oneuptime_external_mcp_server` resources.
- `incident_severities` (Set of String) Match only incidents with these severities (incident rules only). Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.
- `labels` (Set of String) Match only incidents/alerts that carry at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.
- `monitor_labels` (Set of String) Match only when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Match only incidents/alerts from these monitors. Leave empty to match any monitor. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
