---
page_title: "oneuptime_status_page_footer_link Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage footer links on your status page
---

# oneuptime_status_page_footer_link (Resource)

Manage footer links on your status page

## Example Usage

```terraform
resource "oneuptime_status_page_footer_link" "example" {
  status_page_id = oneuptime_status_page.example.id
  title          = "Privacy Policy"
  link           = "https://short.url/abc123"
}
```

<a id="order"></a>
## Order

`order` is this status page footer link's place among the status page footer links of its status page, lowest first. Setting it places this status page footer link and moves no other, so a configuration that numbers every status page footer link of a status page is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.status_page_footer_links`:

```terraform
order = index(local.status_page_footer_links, each.key) + 1
```

- Give each status page footer link a number of its own: two status page footer links at one number are shown the older first.
- A status page footer link created without `order` goes after the last one of its status page.
- Deleting a status page footer link leaves a gap in the numbers, which changes nothing: status page footer links go by their numbers, not by counting.
- Moving a status page footer link in the OneUptime Dashboard numbers the status page's status page footer links 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `link` (String) URL to a website or any other resource on the internet.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.
- `title` (String) Title of this resource.

### Optional

- `order` (Number) Where this link appears among the status page's footer links, lowest number first. A new link is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page footer link by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_footer_link.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_footer_link.example <id>
```
