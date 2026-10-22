---
page_title: "oneuptime_collector_config Data Source - oneuptime"
subcategory: "Other"
description: |-
  The OpenTelemetry Collector config OneUptime sends to the hosts it manages, by system and label.
---

# oneuptime_collector_config (Data Source)

The OpenTelemetry Collector config OneUptime sends to the hosts it manages, by system and label.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one collector config may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_collector_config" "example" {
  name = "Example collector config"
}

# Or by id:
data "oneuptime_collector_config" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `config_hash` (String) The SHA-256 of the current version's config, which a host reports back once it runs it.
- `config_yaml` (String) The collector's config.yaml. Left empty on create, it starts as OneUptime's default for the config's system. Every saved change is checked and becomes a new version.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `current_version_number` (Number) The number of the config's current version. Every saved change to the config adds one.
- `description` (String) What this config collects, for the people who edit it.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) What the hosts this config is for have in common.
- `order` (Number) Where this config sits in the list. A host runs the first config in the list that applies to it. A new config is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along. In the dashboard, drag the rows to reorder them.
- `os_type` (String) The system of the hosts this config is for: linux, windows, darwin (macOS) or any.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Only hosts with at least one of these labels get this config. Leave empty for every host of its system. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
