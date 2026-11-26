---
page_title: "oneuptime_blocklist_feed Data Source - oneuptime"
subcategory: "Other"
description: |-
  A secret URL a firewall or a proxy reads the project's blocklist from, on a schedule: Palo Alto, FortiGate, Check Point, Cisco, OPNsense, pfSense, nftables, nginx and others.
---

# oneuptime_blocklist_feed (Data Source)

A secret URL a firewall or a proxy reads the project's blocklist from, on a schedule: Palo Alto, FortiGate, Check Point, Cisco, OPNsense, pfSense, nftables, nginx and others.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one blocklist feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_blocklist_feed" "example" {
  name = "Example blocklist feed"
}

# Or by id:
data "oneuptime_blocklist_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `allowed_source_addresses` (String) The only addresses and ranges the feed answers, one a line: the firewall's own. Empty to answer any address that has the secret.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `family` (String) What the feed serves: IPv4 and IPv6 addresses (Both), or one of them, for a firewall that keeps them apart.
- `firewall` (String) The firewall or proxy that reads the feed, whose setup steps its page opens on. Any of them can read any feed.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_fetch_source_ip` (String) The address the feed was last fetched from.
- `last_fetch_user_agent` (String) The user agent the feed was last fetched with: the firewall or tool that read it.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) What the feed is called where it is enforced: the firewall or proxy that reads it, such as Palo Alto edge firewall.
- `secret_hint` (String) The last four characters of the feed's secret, shown in its place so a person can tell which secret a firewall holds. Never enough to use it.
- `served_count` (Number) How many addresses and ranges the last list a firewall fetched held.
- `usual_fetch_interval_in_seconds` (Number) How often the firewall usually fetches the feed: the median gap between its latest fetches. Three times it without a fetch, and the firewall has stopped reading the feed.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_fetched_at` (String) When a firewall last fetched the feed.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `previous_secret_expires_at` (String) When the secret the last rotation replaced stops working: a day after the rotation, or at once when the feed's secret is revoked.
- `project_id` (String) ID of the project this feed belongs to. The ID of a `oneuptime_project`.
- `revoked_at` (String) When the feed's secret was revoked. A revoked feed answers no firewall until a new secret is made.
- `secret_made_at` (String) When the feed's current secret was made.
- `served_confirmed_at` (String) When a firewall last fetched the feed and was found to hold its last list: the time its blocks were last confirmed there.
- `stopped_reading_told_at` (String) When the Project Owners were told the firewall stopped reading the feed. Cleared when it reads it again, so they are told once each time.
- `updated_at` (String) Date and Time when the object was updated.
