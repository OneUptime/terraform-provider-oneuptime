---
page_title: "oneuptime_blocklist_feed_fetch Data Source - oneuptime"
subcategory: "Other"
description: |-
  Each time a firewall fetched a blocklist feed: when, from where, with what, and whether the list had changed. Kept 30 days.
---

# oneuptime_blocklist_feed_fetch (Data Source)

Each time a firewall fetched a blocklist feed: when, from where, with what, and whether the list had changed. Kept 30 days.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one blocklist feed fetch may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_blocklist_feed_fetch" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_blocklist_feed_fetch" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `auth_method` (String) How the firewall sent the secret: as a password (Basic), a bearer token (Bearer) or in the URL (Path).
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `entry_count` (Number) How many addresses and ranges the list held.
- `file` (String) What was fetched: the plain list (blocklist.txt, with ?comments=1 when it carried comments), its MD5, the nftables script or the nginx include.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_changed` (Boolean) Whether the list was different from the one the fetch before was answered with.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `security_blocklist_feed_id` (String) ID of the feed that was fetched. The ID of a `oneuptime_blocklist_feed`.
- `source_ip` (String) The address the feed was fetched from.
- `status_code` (Number) What the feed answered: 200 with the list, or 304 when the firewall already held it (If-None-Match).
- `used_previous_secret` (Boolean) Whether the firewall sent the secret a rotation replaced, which stops working a day after the rotation.
- `user_agent` (String) What the feed was fetched with: the firewall or tool that read it.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `fetched_at` (String) When the feed was fetched.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this fetch belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
