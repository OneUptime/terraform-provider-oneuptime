---
page_title: "oneuptime_monitor Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Monitor is anything that monitors your API, Websites, IP, Network or more. You can also create static monitor that does not monitor anything.
---

# oneuptime_monitor (Resource)

Monitor is anything that monitors your API, Websites, IP, Network or more. You can also create static monitor that does not monitor anything.

## Example Usage

```terraform
# The project's own statuses and severities, looked up by name.
data "oneuptime_monitor_status" "operational" {
  name = "Operational"
}

data "oneuptime_monitor_status" "offline" {
  name = "Offline"
}

data "oneuptime_incident_severity" "critical" {
  name = "Critical Incident"
}

resource "oneuptime_monitor" "example" {
  name                = "Example website"
  description         = "Checks https://example.com every minute"
  monitor_type        = "Website"
  monitoring_interval = "* * * * *"

  monitor_steps = [{
    monitor_destination      = "https://example.com"
    monitor_destination_type = "URL"
    request_type             = "GET"

    # Evaluated top to bottom; the first that matches wins.
    criteria = [
      {
        name                  = "Offline"
        description           = "The website does not answer"
        filter_condition      = "Any"
        change_monitor_status = true
        monitor_status_id     = data.oneuptime_monitor_status.offline.id
        create_incidents      = true

        filters = [
          { check_on = "Is Online", filter_type = "False" },
        ]

        incidents = [{
          title                 = "Example website is down"
          description           = "The website did not respond to the probe."
          incident_severity_id  = data.oneuptime_incident_severity.critical.id
          auto_resolve_incident = true
        }]
      },
      {
        name                  = "Online"
        description           = "The website answers"
        filter_condition      = "All"
        change_monitor_status = true
        monitor_status_id     = data.oneuptime_monitor_status.operational.id

        filters = [
          { check_on = "Is Online", filter_type = "True" },
        ]
      },
    ]
  }]
}
```

## Schema

### Required

- `monitor_type` (String) What is the type of this monitor? Website? API? etc. Allowed values: `Manual`, `Website`, `API`, `Ping`, `Kubernetes`, `Docker`, `Host`, `Podman`, `Docker Swarm`, `Proxmox`, `VMware`, `Ceph`, `Storage Array`, `IoT Device`, `IP`, `Incoming Request`, `Incoming Email`, `Port`, `Server`, `SSL Certificate`, `SQL Query`, `Database`, `Synthetic Monitor`, `Custom JavaScript Code`, `Logs`, `Metrics`, `Traces`, `Exceptions`, `Profiles`, `Security Events`, `Network Device`, `DNS`, `DNSSEC`, `NTP`, `Domain`, `External Status Page`.
- `name` (String) Any friendly name for this monitor.

### Optional

- `ceph_clusters` (Set of String) Ceph clusters this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_ceph_cluster` resources.
- `current_monitor_status_id` (String) Whats the current status ID of this monitor? The ID of a `oneuptime_monitor_status`.
- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `database_servers` (Set of String) Databases this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_database` resources.
- `depends_on_monitors` (Set of String) Parent monitors this monitor depends on. When a parent is offline (or in one of the configured suppression statuses), alerts and incidents from this monitor are suppressed at creation time — the monitor keeps evaluating and its status timeline still updates. IDs of `oneuptime_monitor` resources.
- `description` (String) Friendly description that will help you remember.
- `disable_active_monitoring` (Boolean) Disable active monitoring for this resource? Defaults to `false`.
- `docker_hosts` (Set of String) Docker hosts this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_docker_host` resources.
- `docker_swarm_clusters` (Set of String) Docker Swarm clusters this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_docker_swarm_cluster` resources.
- `hosts` (Set of String) Hosts this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_host` resources.
- `incoming_email_custom_local_part` (String) This field is for Incoming Email Monitor only. A custom name for this monitor's inbound email address: the part before the @, on the server's inbound email domain. While set, it replaces the generated monitor-{secret key} address. Must be unique across all monitors. Set to null to go back to the generated address.
- `incoming_monitor_request` (String) Incoming Monitor Request for Incoming Request Monitor. A JSON value: write it with `jsonencode()`.
- `incoming_request_monitor_heartbeat_checked_at` (String) This field is for Incoming Request monitor only. When was the last time we checked the heartbeat?
- `iot_fleets` (Set of String) IoT fleets this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_iot_fleet` resources.
- `is_archived` (Boolean) Archived monitors are hidden from monitor lists and status pages, are not checked, and open no incidents or alerts. Unarchiving resumes monitoring. Defaults to `false`.
- `kubernetes_clusters` (Set of String) Kubernetes clusters this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_kubernetes_cluster` resources.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `minimum_probe_agreement` (Number) Minimum number of probes that must agree on a status before the monitor status changes. If null, all enabled and connected probes must agree.
- `monitor_steps` (Attributes List) What would you like to monitor and what is the criteria? (see [below for nested schema](#nestedatt--monitor_steps))
- `monitor_template_id` (String) ID of the Monitor Template this monitor was created from. Null for monitors not created from a template. The ID of a `oneuptime_monitor_template`.
- `monitoring_interval` (String) How often would you like OneUptime to monitor this resource? A 5-field cron expression, not a label: "*/5 * * * *" is every five minutes.
- `podman_hosts` (Set of String) Podman hosts this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_podman_host` resources.
- `proxmox_clusters` (Set of String) Proxmox clusters this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_proxmox_cluster` resources.
- `server_monitor_request_received_at` (String) This field is for Server Monitor only. When was the last time we received a request?
- `server_monitor_response` (String) Server Monitor Response for Server Monitor. A JSON value: write it with `jsonencode()`.
- `services` (Set of String) Services this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_service` resources.
- `storage_arrays` (Set of String) Storage arrays this monitor watches. Incidents and alerts it creates are linked to them. IDs of `oneuptime_storage_array` resources.
- `suppress_alerts_when_parent_monitor_statuses` (Set of String) Parent monitor statuses that suppress this monitor's alerts and incidents. When empty, statuses flagged as offline suppress (the default). Only used when Depends On Monitors is set. IDs of `oneuptime_monitor_status` resources.
- `telemetry_monitor_last_monitor_at` (String) This field is for Telemetry Monitor only. When was the last time we monitored?
- `telemetry_monitor_next_monitor_at` (String) This field is for Telemetry Monitor only. When is the next time we should monitor?
- `vmware_v_centers` (Set of String) VMware vCenters this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_vcenter` resources.

### Read-Only

- `archived_at` (String) When this monitor was archived. Empty while it is not archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `auto_provisioned_network_device_id` (String) ID of the Network Device that caused this monitor to be provisioned automatically. The ID of a `oneuptime_network_device`.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `disable_active_monitoring_because_of_manual_incident` (Boolean) Disable Monitoring because of Incident which is creeated manually by user.
- `disable_active_monitoring_because_of_scheduled_maintenance_event` (Boolean) Disable Monitoring because of Ongoing Scheduled Maintenance Event.
- `id` (String) Unique identifier for the resource.
- `incoming_email_monitor_heartbeat_checked_at` (String) This field is for Incoming Email monitor only. When was the last time we checked the heartbeat?
- `incoming_email_monitor_last_email_received_at` (String) This field is for Incoming Email Monitor only. When was the last email received?
- `incoming_email_monitor_request` (String) This field is for Incoming Email Monitor only. Last email data received. A JSON value: write it with `jsonencode()`.
- `incoming_email_secret_key` (String) This field is for Incoming Email Monitor only. Secret Key used to generate unique email address.
- `incoming_request_secret_key` (String) This field is for Incoming Request Monitor only. Secret Key to authenticate the request.
- `is_all_probes_disconnected_from_this_monitor` (Boolean) All Probes Disconnected From This Monitor. Is this monitor not being monitored?
- `is_no_probe_enabled_on_this_monitor` (Boolean) No Probe Enabled On This Monitor. Is this monitor not being monitored?
- `is_owner_notified_of_resource_creation` (Boolean) Are owners notified of when this resource is created?
- `network_alert_policy_id` (String) ID of the Network Alert Policy that provisioned this monitor, when one did. The ID of a `oneuptime_network_alert_policy`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `server_monitor_secret_key` (String) This field is for Server Monitor only. Secret Key to authenticate the request.
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

- `check_on` (String) What this filter inspects (e.g. `Is Online`, `Response Status Code`, `Response Time (in ms)`). Allowed values: `Response Time (in ms)`, `Packet Loss (in %)`, `Jitter (in ms)`, `Port DNS Lookup Time (in ms)`, `Port TCP Connect Time (in ms)`, `Response Status Code`, `Response Header`, `Response Header Value`, `Response Body`, `Is Online`, `Incoming Request`, `Server Process Name`, `Server Process Command`, `Server Process PID`, `Request Body`, `Request Header`, `Request Header Value`, `JavaScript Expression`, `Disk Usage (in %)`, `CPU Usage (in %)`, `Memory Usage (in %)`, `Load Average (1 minute)`, `Load Average (5 minute)`, `Load Average (15 minute)`, `Swap Usage (in %)`, `CPU IO Wait (in %)`, `Expires In Hours`, `Expires In Days`, `Is Self Signed Certificate`, `Is Expired Certificate`, `Is Valid Certificate`, `Is Not A Valid Certificate`, `Is Request Timeout`, `Result Value`, `Error`, `Execution Time (in ms)`, `Screen Size`, `Browser Type`, `Log Count`, `Security Event Count`, `Span Count`, `Exception Count`, `Profile Count`, `Metric Value`, `Email Subject`, `Email From Address`, `Email Body`, `Email To Address`, `Email Received`, `SNMP OID Value`, `SNMP OID Exists`, `SNMP Response Time (in ms)`, `SNMP Device Is Online`, `SNMP Walk Is Succeeding`, `SNMP Interface Is Down`, `SNMP Trap Received (Trap OID)`, `SNMP Interface Utilization (in %)`, `SNMP Interface Errors (per second)`, `SNMP Table Value`, `SNMP Table Row Count`, `SNMP Table Row Is Unhealthy`, `SNMP Trap Varbind Value`, `SNMP Transceiver Not Detected`, `SNMP Transceiver Past Alarm Threshold`, `SNMP Transceiver Past Warning Threshold`, `SNMP Transceiver Reading`, `SNMP Transceiver RX Power Drop (in dB)`, `DNS Response Time (in ms)`, `DNS Is Online`, `DNS Record Value`, `DNSSEC Is Valid`, `DNS Record Exists`, `Domain Expires In Days`, `Domain Registrar`, `Domain Name Server`, `Domain Status Code`, `Domain Is Expired`, `DNSSEC Chain Is Valid`, `DNSSEC DNSKEY Record Exists`, `DNSSEC DS Record Exists At Parent`, `DNSSEC Signature Expires In Days`, `DNSSEC Resolver Consensus (AD Flag)`, `DNSSEC Nameservers Are Consistent`, `NTP Is Online`, `NTP Is Synchronized`, `NTP Stratum`, `NTP Clock Offset (in ms)`, `NTP Response Time (in ms)`, `NTP Root Dispersion (in ms)`, `SQL Is Online`, `SQL Query Row Count`, `SQL Query Scalar Value`, `SQL Query Execution Time (in ms)`, `SQL Query Error`, `Database Is Online`, `Database Metric`, `Database Collection Error`, `External Status Page Is Online`, `External Status Page Overall Status`, `External Status Page Component Status`, `External Status Page Active Incidents`, `External Status Page Response Time (in ms)`.

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

Import an existing monitor by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_monitor.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_monitor.example <id>
```
