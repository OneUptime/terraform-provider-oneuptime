---
page_title: "oneuptime_collector_config Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  The OpenTelemetry Collector config OneUptime sends to the hosts it manages, by system and label.
---

# oneuptime_collector_config (Resource)

The OpenTelemetry Collector config OneUptime sends to the hosts it manages, by system and label.

## Example Usage

```terraform
resource "oneuptime_collector_config" "example" {
  name        = "Example collector config"
  description = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this collector config's place among the collector configs of its project, lowest first. Setting it places this collector config and moves no other, so a configuration that numbers every collector config of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.collector_configs`:

```terraform
order = index(local.collector_configs, each.key) + 1
```

- Give each collector config a number of its own: two collector configs at one number are shown the older first.
- A collector config created without `order` goes after the last one of its project.
- Deleting a collector config leaves a gap in the numbers, which changes nothing: collector configs go by their numbers, not by counting.
- Moving a collector config in the OneUptime Dashboard numbers the project's collector configs 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `name` (String) What the hosts this config is for have in common.

### Optional

- `config_yaml` (String) The collector's config.yaml. Left empty on create, it starts as OneUptime's default for the config's system. Every saved change is checked and becomes a new version.
- `description` (String) What this config collects, for the people who edit it.
- `labels` (Set of String) Only hosts with at least one of these labels get this config. Leave empty for every host of its system. IDs of `oneuptime_label` resources.
- `order` (Number) Where this config sits in the list. A host runs the first config in the list that applies to it. A new config is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).
- `os_type` (String) The system of the hosts this config is for: linux, windows, darwin (macOS) or any. Allowed values: `any`, `linux`, `windows`, `darwin`. Defaults to `linux`.

### Read-Only

- `config_hash` (String) The SHA-256 of the current version's config, which a host reports back once it runs it.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `current_version_number` (Number) The number of the config's current version. Every saved change to the config adds one.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing collector config by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_collector_config.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_collector_config.example <id>
```
