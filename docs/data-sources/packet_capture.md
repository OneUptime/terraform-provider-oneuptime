---
page_title: "oneuptime_packet_capture Data Source - oneuptime"
subcategory: "Other"
description: |-
  A packet capture run on one of the project's probes, and the pcap file it produced. Captures and their files are deleted 7 days after they start.
---

# oneuptime_packet_capture (Data Source)

A packet capture run on one of the project's probes, and the pcap file it produced. Captures and their files are deleted 7 days after they start.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one packet capture may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_packet_capture" "example" {
  name = "Example packet capture"
}

# Or by id:
data "oneuptime_packet_capture" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `bpf_filter` (String) A BPF filter expression - the capture-filter language of tcpdump and Wireshark - that decides which packets are kept: host 10.0.0.5 and tcp port 443. Empty keeps every packet.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `end_reason` (String) Why a completed capture stopped: "DurationReached", "PacketLimitReached", "FileSizeLimitReached", "StoppedFromDashboard" or "CaptureToolStopped". Managed by the probe.
- `file_size_in_bytes` (Number) The size of the pcap file in bytes. Managed by OneUptime.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `interface_name` (String) The network interface the probe captures on, as the probe reported it: "eth0", or "any" for every interface at once.
- `max_duration_in_seconds` (Number) How long the capture runs, in seconds, at most: from 5 seconds to 30 minutes, and no longer than the probe allows. It stops earlier when it reaches its packet or file size limit.
- `max_file_size_in_mb` (Number) The capture stops when its file reaches this many megabytes: from 1 to 25, and no more than the probe allows. The file is cut at the last whole packet.
- `max_packets` (Number) The capture stops after this many packets: from 1 to 1,000,000.
- `name` (String) The interface and the filter of the capture, written by OneUptime: eth0: host 10.0.0.5.
- `network_device_id` (String) ID of the network device the capture was started from, if it was started from a device's page. The ID of a `oneuptime_network_device`.
- `packet_count` (Number) How many packets the file holds, counted by OneUptime from the file. Managed by OneUptime.
- `probe_id` (String) ID of the probe that runs the capture: one of the project's own probes, with packet capture turned on. The ID of a `oneuptime_probe`.
- `status` (String) Where the capture is: "Pending" (waiting for the probe), "Running" (the probe is capturing), "Completed" (the file is ready) or "Failed" (see Status Message). Managed by OneUptime and the probe.
- `status_message` (String) Why a capture failed - the interface is gone, the probe may not capture, the filter did not compile - or what the capture tool said when it stopped by itself. Managed by OneUptime and the probe.

### Read-Only

- `completed_at` (String) When the capture completed or failed. Managed by OneUptime.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `started_at` (String) When the probe started capturing. Managed by OneUptime.
- `stop_requested_at` (String) When someone pressed Stop on the running capture. The probe stops at its next check, within about ten seconds, and uploads what it captured. Managed by OneUptime.
- `updated_at` (String) Date and Time when the object was updated.
