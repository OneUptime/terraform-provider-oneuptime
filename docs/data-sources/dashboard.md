---
page_title: "oneuptime_dashboard Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Create and manage Dashboards to visualize your data in a single place
---

# oneuptime_dashboard (Data Source)

Create and manage Dashboards to visualize your data in a single place

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one dashboard may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_dashboard" "example" {
  name = "Example dashboard"
}

# Or by id:
data "oneuptime_dashboard" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `dark_theme_logo_file_id` (String) Dashboard Dark Theme Logo File ID. The ID of a `oneuptime_file`.
- `description` (String) Friendly description that will help you remember.
- `enable_master_password` (Boolean) Require visitors to enter a master password before viewing a public dashboard.
- `favicon_file_id` (String) Dashboard Favicon File ID. The ID of a `oneuptime_file`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `ip_whitelist` (String) IP Whitelist for this dashboard. One IP per line. Only used when the dashboard is public.
- `is_archived` (Boolean) Archived dashboards are hidden from the Dashboards list and their public link stops working. Unarchiving restores them as they were.
- `is_public_dashboard` (Boolean) Is this dashboard public?
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `logo_file_id` (String) Dashboard Logo File ID. The ID of a `oneuptime_file`.
- `name` (String) Any friendly name of this object.
- `page_description` (String) Description of the public dashboard page. This will be used for SEO.
- `page_title` (String) Title of the public dashboard page. This will be used for SEO and the browser tab.
- `require_sign_in` (Boolean) Only this dashboard's private users can open its public link, once they sign in with a password or single sign-on.
- `require_sso_for_login` (Boolean) While the dashboard asks its visitors to sign in, only its single sign-on lets them in: private users' passwords no longer do.
- `slug` (String) Friendly globally unique name for your object.
- `theme` (String) How the public dashboard looks to visitors: Light, Dark, or System (light or dark, as each visitor's device is set). Visitors can switch between light and dark from its toolbar. A dashboard created without a theme starts with System.

### Read-Only

- `archived_at` (String) When this dashboard was archived. Empty while it is not archived.
- `created_at` (String) Date and Time when the object was created.
- `dashboard_view_config` (String) The dashboard's widgets and variables, and the time range (timeRange, such as Past 1 Day) and auto-refresh (refreshInterval, such as 1m) it opens on. See the DashboardViewConfig data type. A JSON value: read it with `jsondecode()`.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `master_password` (String, Sensitive) Password required to unlock a public dashboard. This value is stored as a secure hash.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
