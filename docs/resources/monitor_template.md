---
page_title: "oneuptime_monitor_template Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Reusable monitor template. Use it to create new monitors with the same configuration.
---

# oneuptime_monitor_template (Resource)

Reusable monitor template. Use it to create new monitors with the same configuration.

## Example Usage

```terraform
resource "oneuptime_monitor_template" "example" {
  template_name        = "Example short text"
  template_description = "This is an example of longer text content that might be stored in this field."
  monitor_type         = "Website"
}
```

## Schema

### Required

- `monitor_type` (String) What is the type of monitor created from this template? Allowed values: `Manual`, `Website`, `API`, `Ping`, `Kubernetes`, `Docker`, `Host`, `Podman`, `Docker Swarm`, `Proxmox`, `VMware`, `Ceph`, `Storage Array`, `IoT Device`, `IP`, `Incoming Request`, `Incoming Email`, `Port`, `Server`, `SSL Certificate`, `SQL Query`, `Database`, `Synthetic Monitor`, `Custom JavaScript Code`, `Logs`, `Metrics`, `Traces`, `Exceptions`, `Profiles`, `Security Events`, `AI / LLM`, `Network Device`, `DNS`, `DNSSEC`, `NTP`, `Domain`, `External Status Page`.
- `template_description` (String) Description of the Monitor Template.
- `template_name` (String) Name of the Monitor Template.

### Optional

- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `labels` (Set of String) Default labels applied to monitors created from this template. IDs of `oneuptime_label` resources.
- `minimum_probe_agreement` (Number) Default minimum number of probes that must agree on a status before the monitor status changes.
- `monitor_description` (String) Default description applied to monitors created from this template.
- `monitor_name` (String) Default name applied to monitors created from this template. Users can override on creation. Leave it blank to name each monitor after the resource it watches.
- `monitor_steps` (Attributes List) Monitor steps and criteria copied to monitors created from this template. (see [below for nested schema](#nestedatt--monitor_steps))
- `monitoring_interval` (String) Default monitoring interval for monitors created from this template. A 5-field cron expression, not a label: "*/5 * * * *" is every five minutes.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

<a id="nestedatt--monitor_steps"></a>
### Nested Schema for `monitor_steps`

Required:

- `criteria` (Attributes List) Ordered criteria evaluated after each check. The first matching criteria decides the monitor status and incident/alert actions. (see [below for nested schema](#nestedatt--monitor_steps--criteria))

Optional:

- `allow_self_signed_certificates` (Boolean) Accept self-signed TLS certificates (API and Website monitors).
- `browser_types` (List of String) Browsers for Synthetic monitors: `Chromium`, `Firefox`.
- `ceph_monitor` (String) Raw JSON escape hatch for the Ceph monitor config (clusterIdentifier, ...). Write it with `jsonencode()`.
- `custom_code` (String) JavaScript (Custom Code monitors) or Playwright script (Synthetic monitors) executed by this step.
- `database_monitor` (String) Raw JSON escape hatch for the Database Health monitor config (databaseType, host, port, databaseName, username, password, enabledMetricGroups, ...). Write it with `jsonencode()`.
- `dns_monitor` (String) Raw JSON escape hatch for the DNS monitor config (queryName, recordType, ...). Write it with `jsonencode()`.
- `dnssec_monitor` (String) Raw JSON escape hatch for the DNSSEC monitor config (domainName, resolvers). Write it with `jsonencode()`.
- `do_not_follow_redirects` (Boolean) Do not follow HTTP redirects (API and Website monitors).
- `docker_monitor` (String) Raw JSON escape hatch for the Docker monitor config (hostIdentifier, ...). Write it with `jsonencode()`.
- `docker_swarm_monitor` (String) Raw JSON escape hatch for the Docker Swarm monitor config (clusterIdentifier, ...). Write it with `jsonencode()`.
- `domain_monitor` (String) Raw JSON escape hatch for the Domain monitor config (domainName, lookupMethod, timeout, retries). lookupMethod is one of Auto, RDAP, WHOIS and defaults to Auto. Write it with `jsonencode()`.
- `exception_monitor` (String) Raw JSON escape hatch for the Exceptions monitor query config. Write it with `jsonencode()`.
- `external_status_page_monitor` (String) Raw JSON escape hatch for the External Status Page monitor config (statusPageUrl, providerType). Write it with `jsonencode()`.
- `host_monitor` (String) Raw JSON escape hatch for the Host monitor config (hostIdentifier, ...). Write it with `jsonencode()`.
- `iot_monitor` (String) Raw JSON escape hatch for the IoT monitor config (fleetIdentifier, ...). Write it with `jsonencode()`.
- `kubernetes_monitor` (String) Raw JSON escape hatch for the Kubernetes monitor config (clusterIdentifier, ...). Write it with `jsonencode()`.
- `llm_monitor` (String) Raw JSON escape hatch for the AI / LLM monitor config (issues, slowAnswerSeconds, model, telemetryServiceIds, lastXSecondsOfCalls). Write it with `jsonencode()`.
- `log_monitor` (String) Raw JSON escape hatch for the Logs monitor query config (attributes, body, severityTexts, telemetryServiceIds, lastXSecondsOfLogs). Write it with `jsonencode()`.
- `metric_monitor` (String) Raw JSON escape hatch for the Metrics monitor query config (metricViewConfig, rollingTime). The server normalizes this object, so provide the full shape to avoid drift. Write it with `jsonencode()`.
- `monitor_destination` (String) The URL, IP address or hostname this step probes (e.g. `https://example.com`, `8.8.8.8`, `example.com`). Requires `monitor_destination_type`.
- `monitor_destination_type` (String) Kind of destination: `URL` (Website, API, SSL Certificate monitors), `IP` or `Hostname` (Ping, Port monitors). Allowed values: `URL`, `IP`, `Hostname`.
- `network_device_monitor` (String) Raw JSON escape hatch for the Network Device monitor config (networkDeviceId). Write it with `jsonencode()`.
- `podman_monitor` (String) Raw JSON escape hatch for the Podman monitor config (hostIdentifier, ...). Write it with `jsonencode()`.
- `port` (Number) TCP port to probe (Port monitors).
- `profile_monitor` (String) Raw JSON escape hatch for the Profiles monitor query config. Write it with `jsonencode()`.
- `proxmox_monitor` (String) Raw JSON escape hatch for the Proxmox monitor config (clusterIdentifier, ...). Write it with `jsonencode()`.
- `request_body` (String) HTTP request body sent by API monitors.
- `request_headers` (Map of String) HTTP request headers sent by API and Website monitors. Omit instead of passing an empty map.
- `request_timeout_in_ms` (Number) Per-step request timeout in milliseconds for probe-based monitors. The probe clamps anything above 60000 to 60000.
- `request_type` (String) HTTP method for API monitors (`GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `PATCH`). The server defaults to `GET` when omitted. Allowed values: `GET`, `POST`, `DELETE`, `PUT`, `HEAD`, `PATCH`.
- `retry_count` (Number) Per-step retries for probe-based monitors when a check fails. This counts retries AFTER the first attempt: `0` runs the check once, `3` runs it up to four times. The probe clamps anything above 3 to 3.
- `retry_count_on_error` (Number) Retries on script error (Synthetic monitors). This counts retries AFTER the first attempt: `0` runs the script once, `2` runs it up to three times. Must be `0` or greater; the dashboard caps it at 5.
- `screen_size_types` (List of String) Screen sizes for Synthetic monitors: `Mobile`, `Tablet`, `Desktop`.
- `sql_monitor` (String) Raw JSON escape hatch for the SQL Query monitor config (databaseType, host, port, databaseName, query, ...). Write it with `jsonencode()`.
- `storage_array_monitor` (String) Raw JSON escape hatch for the Storage Array monitor config (arrayIdentifier, storageSystem, resourceFilters, ...). Write it with `jsonencode()`.
- `tls_client_certificate` (String) Client certificate (PEM or `{{monitorSecrets.name}}` reference) for mTLS (API and Website monitors).
- `tls_client_key` (String, Sensitive) Client private key (PEM or `{{monitorSecrets.name}}` reference) for mTLS (API and Website monitors).
- `tls_client_key_passphrase` (String, Sensitive) Passphrase for the client private key.
- `trace_monitor` (String) Raw JSON escape hatch for the Traces monitor query config. Write it with `jsonencode()`.
- `vmware_monitor` (String) Raw JSON escape hatch for the VMware monitor config (vcenterIdentifier, resourceFilters, ...). Write it with `jsonencode()`.

<a id="nestedatt--monitor_steps--criteria"></a>
### Nested Schema for `monitor_steps.criteria`

Required:

- `filter_condition` (String) How the filters combine: `All` (every filter must match) or `Any` (one match is enough). Allowed values: `All`, `Any`.
- `filters` (Attributes List) Conditions evaluated against the probe result. At least one filter is required. (see [below for nested schema](#nestedatt--monitor_steps--criteria--filters))
- `name` (String) Human-readable name of this criteria (e.g. `Check if online`).

Optional:

- `alerts` (Attributes List) Alert templates created when this criteria matches and `create_alerts` is true. Omit instead of passing an empty list. (see [below for nested schema](#nestedatt--monitor_steps--criteria--alerts))
- `change_monitor_status` (Boolean) Change the monitor status to `monitor_status_id` when this criteria matches. Defaults to false server-side.
- `create_alerts` (Boolean) Create the alerts declared in `alerts` when this criteria matches. Defaults to false server-side.
- `create_incidents` (Boolean) Create the incidents declared in `incidents` when this criteria matches. Defaults to false server-side.
- `description` (String) Description of what this criteria checks.
- `incident_grouping` (String) Raw JSON escape hatch for per-criteria incident grouping (Incoming Request monitors only; groupByJSONPath, resolvedWhenJSONPath, resolvedWhenValue). Write it with `jsonencode()`.
- `incidents` (Attributes List) Incident templates created when this criteria matches and `create_incidents` is true. Omit instead of passing an empty list. (see [below for nested schema](#nestedatt--monitor_steps--criteria--incidents))
- `is_enabled` (Boolean) Whether this criteria is evaluated. Defaults to true server-side.
- `monitor_status_id` (String) ID of the monitor status (e.g. Operational, Offline) to set when this criteria matches.

<a id="nestedatt--monitor_steps--criteria--filters"></a>
### Nested Schema for `monitor_steps.criteria.filters`

Required:

- `check_on` (String) What this filter inspects (e.g. `Is Online`, `Response Status Code`, `Response Time (in ms)`). Allowed values: `Response Time (in ms)`, `Packet Loss (in %)`, `Jitter (in ms)`, `Port DNS Lookup Time (in ms)`, `Port TCP Connect Time (in ms)`, `Response Status Code`, `Response Header`, `Response Header Value`, `Response Body`, `Is Online`, `Incoming Request`, `Server Process Name`, `Server Process Command`, `Server Process PID`, `Request Body`, `Request Header`, `Request Header Value`, `JavaScript Expression`, `Disk Usage (in %)`, `CPU Usage (in %)`, `Memory Usage (in %)`, `Load Average (1 minute)`, `Load Average (5 minute)`, `Load Average (15 minute)`, `Swap Usage (in %)`, `CPU IO Wait (in %)`, `Expires In Hours`, `Expires In Days`, `Is Self Signed Certificate`, `Is Expired Certificate`, `Is Valid Certificate`, `Is Not A Valid Certificate`, `Is Request Timeout`, `Result Value`, `Error`, `Execution Time (in ms)`, `Screen Size`, `Browser Type`, `Log Count`, `Security Event Count`, `Span Count`, `Bad AI Answers (in %)`, `Bad AI Answers`, `AI Answers`, `Exception Count`, `Profile Count`, `Metric Value`, `Email Subject`, `Email From Address`, `Email Body`, `Email To Address`, `Email Received`, `SNMP OID Value`, `SNMP OID Exists`, `SNMP Response Time (in ms)`, `SNMP Device Is Online`, `SNMP Walk Is Succeeding`, `SNMP Interface Is Down`, `SNMP Trap Received (Trap OID)`, `SNMP Interface Utilization (in %)`, `SNMP Interface Errors (per second)`, `SNMP Table Value`, `SNMP Table Row Count`, `SNMP Table Row Is Unhealthy`, `SNMP Trap Varbind Value`, `SNMP Transceiver Not Detected`, `SNMP Transceiver Past Alarm Threshold`, `SNMP Transceiver Past Warning Threshold`, `SNMP Transceiver Reading`, `SNMP Transceiver RX Power Drop (in dB)`, `DNS Response Time (in ms)`, `DNS Is Online`, `DNS Record Value`, `DNSSEC Is Valid`, `DNS Record Exists`, `Domain Expires In Days`, `Domain Registrar`, `Domain Name Server`, `Domain Status Code`, `Domain Is Expired`, `DNSSEC Chain Is Valid`, `DNSSEC DNSKEY Record Exists`, `DNSSEC DS Record Exists At Parent`, `DNSSEC Signature Expires In Days`, `DNSSEC Resolver Consensus (AD Flag)`, `DNSSEC Nameservers Are Consistent`, `NTP Is Online`, `NTP Is Synchronized`, `NTP Stratum`, `NTP Clock Offset (in ms)`, `NTP Response Time (in ms)`, `NTP Root Dispersion (in ms)`, `SQL Is Online`, `SQL Query Row Count`, `SQL Query Scalar Value`, `SQL Query Execution Time (in ms)`, `SQL Query Error`, `Database Is Online`, `Database Metric`, `Database Collection Error`, `External Status Page Is Online`, `External Status Page Overall Status`, `External Status Page Component Status`, `External Status Page Active Incidents`, `External Status Page Response Time (in ms)`.

Optional:

- `custom_code_monitor_options` (String) Raw JSON escape hatch for Custom Code and Synthetic monitor filter options (resultValuePath - the field of the returned data a `Result Value` filter compares, e.g. `status` or `data.items[0].value`: dots for nested fields, `[n]` for array items). Omit it to compare the whole value. Write it with `jsonencode()`.
- `database_monitor_options` (String) Raw JSON escape hatch for Database Health filter options (metricType). Required on every `Database Metric` filter — it names the series the threshold applies to. Write it with `jsonencode()`.
- `disk_path` (String) Disk path for `Disk Usage (in %)` filters on Server monitors (e.g. `/` or `C:`).
- `evaluate_over_time` (Boolean) Evaluate this filter over a time window instead of the latest probe result.
- `evaluate_over_time_minutes` (Number) Length of the evaluation window in minutes (used with `evaluate_over_time`).
- `evaluate_over_time_no_data_policy` (String) What the filter does while the evaluation window does not hold enough data to judge it — for example a monitor that has just been created. `Ignore` (the default) does not match, `Trigger` treats the missing data as the failure, and `Treat As Zero` compares the window as a single zero. Allowed values: `Ignore`, `Treat As Zero`, `Trigger`.
- `evaluate_over_time_type` (String) Aggregation applied over the evaluation window (e.g. `Average`, `Sum`, `Any Value`). `All Values` only matches once the window is actually covered by data, so give it a window at least twice the monitoring interval. Allowed values: `Average`, `Sum`, `Maximum Value`, `Minimum Value`, `All Values`, `Any Value`.
- `filter_type` (String) Comparison operator for the filter (e.g. `True`, `Equal To`, `Greater Than`). Allowed values: `Equal To`, `Not Equal To`, `Greater Than`, `Less Than`, `Greater Than Or Equal To`, `Less Than Or Equal To`, `Contains`, `Not Contains`, `Starts With`, `Ends With`, `Is Empty`, `Is Not Empty`, `True`, `False`, `Not Recieved In Minutes`, `Recieved In Minutes`, `Evaluates To True`, `Is Executing`, `Is Not Executing`, `Anomalously High`, `Anomalously Low`, `Anomalous`.
- `metric_monitor_options` (String) Raw JSON escape hatch for metric filter options (metricAlias, metricAggregationType, anomaly detection, ...). Write it with `jsonencode()`.
- `snmp_monitor_options` (String) Raw JSON escape hatch for SNMP filter options (oid, interfaceName, and for SNMP table criteria tableKey, tableColumnOid and tableRow). Write it with `jsonencode()`.
- `value` (String) Threshold or comparison value. Always a string — write numbers as strings (e.g. `"200"`).

<a id="nestedatt--monitor_steps--criteria--alerts"></a>
### Nested Schema for `monitor_steps.criteria.alerts`

Required:

- `description` (String) Description of the alert created when this criteria matches.
- `title` (String) Title of the alert created when this criteria matches.

Optional:

- `alert_severity_id` (String) ID of the alert severity to use.
- `auto_resolve_alert` (Boolean) Automatically resolve the alert when the criteria stops matching.
- `is_private` (Boolean) Mark the alert as private.
- `label_ids` (List of String) IDs of labels to attach to the alert. Omit instead of passing an empty list.
- `on_call_policy_ids` (List of String) IDs of on-call duty policies to execute when the alert is created. Omit instead of passing an empty list.
- `owner_team_ids` (List of String) IDs of teams to add as alert owners. Omit instead of passing an empty list.
- `owner_user_ids` (List of String) IDs of users to add as alert owners. Omit instead of passing an empty list.
- `remediation_notes` (String) Remediation notes attached to the alert.

<a id="nestedatt--monitor_steps--criteria--incidents"></a>
### Nested Schema for `monitor_steps.criteria.incidents`

Required:

- `description` (String) Description of the incident created when this criteria matches.
- `title` (String) Title of the incident created when this criteria matches.

Optional:

- `auto_resolve_incident` (Boolean) Automatically resolve the incident when the criteria stops matching.
- `incident_severity_id` (String) ID of the incident severity to use.
- `is_private` (Boolean) Mark the incident as private.
- `label_ids` (List of String) IDs of labels to attach to the incident. Omit instead of passing an empty list.
- `on_call_policy_ids` (List of String) IDs of on-call duty policies to execute when the incident is created. Omit instead of passing an empty list.
- `owner_team_ids` (List of String) IDs of teams to add as incident owners. Omit instead of passing an empty list.
- `owner_user_ids` (List of String) IDs of users to add as incident owners. Omit instead of passing an empty list.
- `remediation_notes` (String) Remediation notes attached to the incident.
- `show_incident_on_status_page` (Boolean) Show the incident on status pages this monitor is attached to.

## Import

Import an existing monitor template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_monitor_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_monitor_template.example <id>
```
