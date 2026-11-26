---
page_title: "oneuptime_cloud_waf Resource - oneuptime"
subcategory: "Cloud"
description: |-
  The WAFs of connected cloud accounts that enforce the project's blocklist: AWS WAF web ACLs, Azure WAF policies and Google Cloud Armor policies, each holding OneUptime's own rules.
---

# oneuptime_cloud_waf (Resource)

The WAFs of connected cloud accounts that enforce the project's blocklist: AWS WAF web ACLs, Azure WAF policies and Google Cloud Armor policies, each holding OneUptime's own rules.

## Example Usage

```terraform
resource "oneuptime_cloud_waf" "example" {

}
```

## Schema

### Optional

- `client_address` (String) Where the WAF reads the client's address: the request's source address, or a header a proxy or CDN in front of it writes. Allowed values: `SourceIp`, `XForwardedFor`, `CfConnectingIp`, `TrueClientIp`, `ForwardedClient`, `UserIp`. Defaults to `SourceIp`.
- `cloud_resource_id` (String) The WAF's row in the cloud inventory: a web ACL, Azure WAF policy or Cloud Armor policy of a connection whose blocking setting is on. The ID of a `oneuptime_cloud_resource`.

### Read-Only

- `checked_at` (String) When OneUptime last read the WAF.
- `cloud_connection_id` (String) ID of the cloud connection OneUptime writes the WAF through. The ID of a `oneuptime_cloud_connection`.
- `cloud_provider` (String) The cloud the WAF is in: aws, azure or gcp.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who picked the WAF. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) The WAF's name: prod-alb.
- `project_id` (String) ID of the project this WAF enforces the blocklist of. The ID of a `oneuptime_project`.
- `region` (String) The Region the WAF is in: us-east-1 for a CloudFront distribution's.
- `rule_details` (String) What OneUptime keeps in the WAF: its rules' names and priorities, what they hold or reference, and the metric whose sampled requests are evidence. A JSON value: read it with `jsondecode()`.
- `status` (String) Where the WAF stands, as OneUptime last read it: Checking, Enforcing, Waiting for your rule, No free priority, Applying, Changed by hand, or Failed.
- `status_message` (String) Why the WAF is not enforcing the blocklist, and what to do: empty when it is.
- `status_reason` (String) Why the WAF is not enforcing the blocklist, as a sentence and its values that the Dashboard words in the reader's language: empty when it is. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
- `waf_id` (String) The WAF as its cloud names it: a web ACL's ARN, an Azure WAF policy's resource ID, or a Cloud Armor policy's name with its project.

## Import

Import an existing cloud waf by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_cloud_waf.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_cloud_waf.example <id>
```
