---
page_title: "oneuptime_openstack_cloud Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  OpenStack clouds that are being monitored in this project. A cloud is discovered when the OneUptime OpenStack Agent sends its first metrics, or connected with a read-only credential that one of the project's probes collects it with.
---

# oneuptime_openstack_cloud (Resource)

OpenStack clouds that are being monitored in this project. A cloud is discovered when the OneUptime OpenStack Agent sends its first metrics, or connected with a read-only credential that one of the project's probes collects it with.

~> **Renamed:** this resource was called `oneuptime_open_stack_cloud` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing openstack cloud:

```terraform
moved {
  from = oneuptime_open_stack_cloud.example
  to   = oneuptime_openstack_cloud.example
}
```

## Example Usage

```terraform
resource "oneuptime_openstack_cloud" "example" {
  name        = "Example openstack cloud"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this OpenStack cloud. This is the join key — it must match the openstack.cloud.name OTel resource attribute stamped by the OneUptime OpenStack Agent.

### Optional

- `active_instance_count` (Number) How many instances are ACTIVE.
- `agent_version` (String) Version of the OneUptime OpenStack agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.
- `ai_command_allowlist` (String) OneUptime AI only reads an OpenStack cloud today, and a read needs no approval, so no pattern is valid here yet and this stays empty. Optional JSON array of command patterns that Automatic mode may run on this OpenStack cloud without approval even though they are riskier changes. Each pattern is one command line for this OpenStack cloud's agent (openstack) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this OpenStack cloud is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit AI Rule; anyone who may edit the OpenStack cloud can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.
- `ai_remediation_mode` (String) OneUptime AI only reads an OpenStack cloud today: it has no command that changes one, so this stays Disabled and any other mode is refused. Disabled: AI never proposes or runs a change on this OpenStack cloud. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the OpenStack cloud's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the OpenStack cloud can lower the mode; raising it needs Project Owner, Project Admin or Edit AI Rule. Allowed values: `Disabled`, `RequireApproval`, `Automatic`, `BypassApproval`. Defaults to `Disabled`.
- `application_credential_id` (String) The ID of the Keystone application credential the probe logs in with, the better of the two kinds because it can be limited to reading and can expire. Its secret goes in Credential Secret. Changing it asks for the secret again.
- `auth_url` (String) The HTTPS address of the cloud's Keystone, which the probe logs in at and reads the service catalog from, such as https://keystone.example.com:5000/v3. Stored as https://host[:port][/path], without the version. Changing it asks for the credential again.
- `availability_zone_count` (Number) How many availability zones this cloud reports.
- `collection_interval_in_minutes` (Number) How often the probe collects this cloud, in whole minutes from 1 to 60. Two minutes is the OpenStack agent's own default. Collect large clouds less often to go easier on their APIs. Defaults to `2`.
- `collection_method` (String) How this cloud's data reaches OneUptime. Probe: one of your OneUptime probes reads the cloud's APIs with the credential saved here, with no agent to install. Agent: the OneUptime OpenStack Agent you run sends it. Switching to Agent forgets the saved credential. Allowed values: `Agent`, `Probe`. Defaults to `Agent`.
- `collection_options` (String) What the probe reads beyond the agent's defaults, as JSON. collectPorts reads ports, scoped by portNetworkIds and portProjectIds and paged. collectProjectLimits reads each project's quotas, four calls per project. Both are off unless asked for. A JSON value: write it with `jsonencode()`.
- `collection_probe_id` (String) ID of the OneUptime probe that collects this cloud. Changing it asks for the credential again. The ID of a `oneuptime_probe`.
- `collection_regions` (String) The regions the probe reads, as a JSON array of the names the catalog uses, such as ["RegionOne"]. The connection test lists the ones the catalog names. A JSON value: write it with `jsonencode()`.
- `credential_kind` (String) Which kind of credential the probe logs in with: applicationCredential or password. Worked out from what is written with a new secret, so it can be left out. Allowed values: `applicationCredential`, `password`.
- `credential_secret` (String, Sensitive) The application credential's secret, or the password. Write-only: encrypted at rest, never returned by the API, and sent only to the probe that collects this cloud. Leave it out of an update to keep the saved one. Changing the address or the probe, or trusting another certificate, needs it again.
- `description` (String) Friendly description for this OpenStack cloud.
- `down_service_count` (Number) How many enabled services and agents are down. A service someone disabled is not counted.
- `endpoint_interface` (String) Which of the catalog's addresses the probe uses for each service: public, internal or admin. Most probes sit where public answers. One in the control plane's own network may want internal. Allowed values: `public`, `internal`, `admin`. Defaults to `public`.
- `floating_ip_count` (Number) How many floating IPs this cloud reports.
- `hypervisor_count` (Number) How many hypervisors this cloud reports.
- `instance_count` (Number) How many instances this cloud reports, in every project.
- `instance_metrics_options` (String) Whether real per-instance CPU, memory, disk and network series are kept for this cloud, as JSON: mode (automatic, all, selected or off), and for selected, the projectIds and aggregates. Per-instance series multiply fast, so automatic stops keeping them past a stated number of instances. A JSON value: write it with `jsonencode()`.
- `is_ai_investigation_enabled` (Boolean) When on, OneUptime AI runs read-only commands (openstack server, hypervisor, service and agent reads, each one call to the cloud's API) on this OpenStack cloud, through its OpenStack AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the OpenStack cloud can turn it on or off. Defaults to `true`.
- `is_archived` (Boolean) Is this OpenStack cloud archived? Archived OpenStack clouds are hidden from lists but keep collecting telemetry. Defaults to `false`.
- `keystone_username` (String) The Keystone user the probe logs in as, when it logs in with a password instead of an application credential. The password goes in Credential Secret. Changing it asks for the password again.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When metrics were last received from this cloud.
- `memory_allocated_bytes` (Number) Memory promised to instances. This is what is promised, not what is in use.
- `memory_capacity_bytes` (Number) Memory the cloud can promise to instances, as Placement counts it.
- `network_count` (Number) How many networks this cloud reports, in every project.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).
- `price_list` (String) What this cloud's operator says things cost, per hour, in one currency, as JSON. Optional: usage works without it, and a price of zero means no price has been set. A JSON value: write it with `jsonencode()`.
- `project_count` (Number) How many OpenStack projects this cloud reports.
- `project_domain_name` (String) The Keystone domain of that project. Default, when left out.
- `project_name` (String) The project that user logs in to, when it logs in with a password. An application credential already names its own project.
- `region_count` (Number) How many regions this cloud reports.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this OpenStack cloud. Leave blank to use the project-wide default.
- `router_count` (Number) How many routers this cloud reports, in every project.
- `service_count` (Number) How many services and agents (nova-compute, neutron agents, cinder-volume and so on) this cloud reports.
- `storage_capacity_bytes` (Number) Total size of the cloud's Cinder storage pools.
- `storage_used_bytes` (Number) How much of the cloud's Cinder storage pools is used.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this OpenStack cloud (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the OpenStack cloud default, then the project's retention settings. A JSON value: write it with `jsonencode()`.
- `trusted_certificate_fingerprints` (String) The SHA-256 fingerprints of the certificates the probe accepts from this cloud's APIs, as a JSON array, for APIs whose certificates are not from a public authority. More than one, because each service can have a host and a certificate of its own. Empty: every certificate must be signed by an authority the probe's machine trusts. Verification is never skipped. Trusting another certificate asks for the credential again. A JSON value: write it with `jsonencode()`.
- `user_domain_name` (String) The Keystone domain of that user. Default, when left out.
- `vcpu_allocated` (Number) vCPUs promised to instances. This is what is promised, not what is busy.
- `vcpu_capacity` (Number) vCPUs the cloud can promise to instances, as Placement counts them: (total - reserved) * allocation ratio.
- `volume_count` (Number) How many volumes this cloud reports, in every project.

### Read-Only

- `ai_access_configured_at` (String) When OneUptime AI access to this OpenStack cloud was first configured by anyone saving an AI access setting. Set by the server; never cleared, so an OpenStack AI agent that registers later never overwrites a setting an operator chose.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this OpenStack cloud, kept until the next successful command. Set by the server.
- `ai_access_last_verified_at` (String) When a command from OneUptime AI last succeeded on this OpenStack cloud through its OpenStack AI agent. Set by the server.
- `ai_agent_reset_at` (String) When this OpenStack cloud's AI agent was last reset: its key stopped working at once, and it registers again by itself. Set by the server.
- `archived_at` (String) When was this OpenStack cloud archived?
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `capacity_snapshot` (String) What this cloud's capacity added up to when a worker last read its inventory, as JSON: what is allocated and what can be, by zone, aggregate, region and cloud, how many more of each flavor fit, when each thing fills up, and what no instance explains. Written by the server. A JSON value: read it with `jsondecode()`.
- `collection_error` (String) What went wrong in the last collection, in words. Written by the probe.
- `collection_error_code` (String) Why the last collection failed or was partial, as a code the dashboard turns into a next step (InvalidCredential, UntrustedCertificate, ConnectionTimedOut, ...). Written by the probe.
- `collection_settings_version` (Number) Counts changes to the probe collection settings, so a report from a collection that started before a change is not shown as the current status. Written by the server.
- `collection_status` (String) Where probe collection stands: Pending (the probe has not tried the current settings yet), Succeeded, PartiallySucceeded (some of the cloud could not be read) or Failed. Written by the server and the probe.
- `collection_summary` (String) What the last collection found: the regions, how each service answered, how many hypervisors, instances, networks, volumes and projects it read, and what it left out. Written by the probe. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `credentials_updated_at` (String) When the credential was last changed. Set by the server.
- `id` (String) Unique identifier for the resource.
- `is_credential_set` (Boolean) Whether a credential is saved. Set by the server from the secret itself.
- `last_collection_at` (String) When the probe last tried to collect this cloud, whether it worked or not. Written by the server.
- `last_successful_collection_at` (String) When a collection of this cloud last worked, in full or in part. Written by the server.
- `next_collection_at` (String) When the probe collects this cloud next. Written by the server.
- `presented_certificate` (String) The certificate an API presented when the last collection did not trust it: its SHA-256 fingerprint, subject, issuer, validity and the host that showed it, so a person can check it and trust exactly that certificate. Written by the probe. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `spend_snapshot` (String) What this cloud's usage added up to when a worker last read its inventory, as JSON: this month against last, the largest domains and flavors, the top growers, and what could be reclaimed. Estimates, never an invoice. Written by the server. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
- `usage_read_at` (String) The time this cloud's usage has been counted up to. Written by the server, and how a time is never counted twice.

## Import

Import an existing openstack cloud by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_openstack_cloud.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_openstack_cloud.example <id>
```
