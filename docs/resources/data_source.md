---
page_title: "oneuptime_data_source Resource - oneuptime"
subcategory: "Other"
description: |-
  Connect external systems — Prometheus, SQL databases, ClickHouse, Loki, Elasticsearch, or REST APIs — and build dashboards on their data.
---

# oneuptime_data_source (Resource)

Connect external systems — Prometheus, SQL databases, ClickHouse, Loki, Elasticsearch, or REST APIs — and build dashboards on their data.

## Example Usage

```terraform
resource "oneuptime_data_source" "example" {
  name             = "Example data source"
  data_source_type = "Example short text"
  description      = "Managed by Terraform"
}
```

## Schema

### Required

- `data_source_type` (String) The kind of external system this data source connects to (Prometheus, PostgreSQL, MySQL, Microsoft SQL Server, ClickHouse, Loki, Elasticsearch, or REST API).
- `name` (String) Any friendly name of this object.

### Optional

- `additional_options` (String) Per-type options that are not secrets — e.g. { "sslEnabled": true, "elasticsearchIndex": "logs-*" }. A JSON value: write it with `jsonencode()`.
- `api_token` (String) Bearer token for HTTP sources (Prometheus behind an auth proxy, Elasticsearch API key, REST APIs). Encrypted at rest and never returned by the API.
- `custom_headers` (String) Extra HTTP headers sent to HTTP-based sources (e.g. auth headers for a proxy). Values are encrypted at rest and never returned by the API. A JSON value: write it with `jsonencode()`.
- `database_host` (String) Hostname or IP address for database sources (PostgreSQL, MySQL, SQL Server, ClickHouse).
- `database_name` (String) Database (or ClickHouse database) to connect to.
- `database_port` (Number) Port for database sources. Defaults per engine: PostgreSQL 5432, MySQL 3306, SQL Server 1433, ClickHouse 8123.
- `description` (String) Friendly description that will help you remember.
- `password` (String) Password for database sources, or HTTP basic-auth password. Encrypted at rest and never returned by the API.
- `url` (String) Base URL for HTTP-based sources (Prometheus, Loki, Elasticsearch, REST API) — e.g. https://prometheus.example.com.
- `username` (String) Username for database sources, or HTTP basic-auth username for HTTP sources. Use a READ-ONLY account — dashboards only ever read.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing data source by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_data_source.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_data_source.example <id>
```
