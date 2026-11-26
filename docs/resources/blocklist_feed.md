---
page_title: "oneuptime_blocklist_feed Resource - oneuptime"
subcategory: "Other"
description: |-
  A secret URL a firewall or a proxy reads the project's blocklist from, on a schedule: Palo Alto, FortiGate, Check Point, Cisco, OPNsense, pfSense, nftables, nginx and others.
---

# oneuptime_blocklist_feed (Resource)

A secret URL a firewall or a proxy reads the project's blocklist from, on a schedule: Palo Alto, FortiGate, Check Point, Cisco, OPNsense, pfSense, nftables, nginx and others.

## Example Usage

```terraform
resource "oneuptime_blocklist_feed" "example" {
  name = "Example blocklist feed"
}
```

## Schema

### Required

- `name` (String) What the feed is called where it is enforced: the firewall or proxy that reads it, such as Palo Alto edge firewall.

### Optional

- `allowed_source_addresses` (String) The only addresses and ranges the feed answers, one a line: the firewall's own. Empty to answer any address that has the secret.
- `family` (String) What the feed serves: IPv4 and IPv6 addresses (Both), or one of them, for a firewall that keeps them apart. Allowed values: `Both`, `IPv4`, `IPv6`. Defaults to `Both`.
- `firewall` (String) The firewall or proxy that reads the feed, whose setup steps its page opens on. Any of them can read any feed. Allowed values: `PaloAlto`, `FortiGate`, `CheckPoint`, `Cisco`, `OPNsense`, `PfSense`, `Nftables`, `Ipset`, `Nginx`, `HAProxy`, `WindowsFirewall`, `Other`. Defaults to `Other`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `id` (String) Unique identifier for the resource.
- `last_fetch_source_ip` (String) The address the feed was last fetched from.
- `last_fetch_user_agent` (String) The user agent the feed was last fetched with: the firewall or tool that read it.
- `last_fetched_at` (String) When a firewall last fetched the feed.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `previous_secret_expires_at` (String) When the secret the last rotation replaced stops working: a day after the rotation, or at once when the feed's secret is revoked.
- `project_id` (String) ID of the project this feed belongs to. The ID of a `oneuptime_project`.
- `revoked_at` (String) When the feed's secret was revoked. A revoked feed answers no firewall until a new secret is made.
- `secret` (String, Sensitive) The feed's secret, which a firewall sends as a password, as a bearer token or in the URL. Shown once: in the answer to the create that makes the feed, or to a rotation of it. OneUptime keeps only its hash. Kept in state from the answer to the create: OneUptime never reads it back, so a resource imported into Terraform has none. Replace the resource for a new one.
- `secret_hint` (String) The last four characters of the feed's secret, shown in its place so a person can tell which secret a firewall holds. Never enough to use it.
- `secret_made_at` (String) When the feed's current secret was made.
- `served_confirmed_at` (String) When a firewall last fetched the feed and was found to hold its last list: the time its blocks were last confirmed there.
- `served_count` (Number) How many addresses and ranges the last list a firewall fetched held.
- `stopped_reading_told_at` (String) When the Project Owners were told the firewall stopped reading the feed. Cleared when it reads it again, so they are told once each time.
- `updated_at` (String) Date and Time when the object was updated.
- `usual_fetch_interval_in_seconds` (Number) How often the firewall usually fetches the feed: the median gap between its latest fetches. Three times it without a fetch, and the firewall has stopped reading the feed.

## Import

Import an existing blocklist feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_blocklist_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_blocklist_feed.example <id>
```
