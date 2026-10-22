---
page_title: "oneuptime_collector_config_version Data Source - oneuptime"
subcategory: "Other"
description: |-
  Every saved version of a Collector Config, so each host can say which version it runs.
---

# oneuptime_collector_config_version (Data Source)

Every saved version of a Collector Config, so each host can say which version it runs.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one collector config version may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_collector_config_version" "example" {
  host_collector_config_id = oneuptime_collector_config.example.id
}

# Or by id:
data "oneuptime_collector_config_version" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `change_note` (String) What this version did: created the config, changed it, or put an older version back.
- `config_hash` (String) The SHA-256 of this version's config, which a host reports back once it runs it.
- `config_yaml` (String) The collector's config.yaml as this version saved it.
- `created_by_user_id` (String) ID of the user who saved this version. The ID of a `oneuptime_user` (see the data source).
- `host_collector_config_id` (String) ID of the Collector Config this is a version of. The ID of a `oneuptime_collector_config`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `version_number` (Number) The version's number: 1 for the config as created, then one more for every saved change.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
