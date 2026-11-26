package shard22

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProbeDataSource{}

func NewProbeDataSource() datasource.DataSource {
    return &ProbeDataSource{}
}

// ProbeDataSource defines the data source implementation.
type ProbeDataSource struct {
    client *Client
}

// ProbeDataSourceModel describes the data source data model.
type ProbeDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    Key types.String `tfsdk:"key"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    ProbeVersion types.String `tfsdk:"probe_version"`
    ReportedVersion types.String `tfsdk:"reported_version"`
    MajorVersionBehindNotifiedFor types.Number `tfsdk:"major_version_behind_notified_for"`
    EgressIpAddress types.String `tfsdk:"egress_ip_address"`
    InterfaceIpAddresses types.String `tfsdk:"interface_ip_addresses"`
    LastAlive types.String `tfsdk:"last_alive"`
    IconFileId types.String `tfsdk:"icon_file_id"`
    ProjectId types.String `tfsdk:"project_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    ShouldAutoEnableProbeOnNewMonitors types.Bool `tfsdk:"should_auto_enable_probe_on_new_monitors"`
    ConnectionStatus types.String `tfsdk:"connection_status"`
    ChecksLateSince types.String `tfsdk:"checks_late_since"`
    ChecksLastLateAt types.String `tfsdk:"checks_last_late_at"`
    DevicePollsLateSince types.String `tfsdk:"device_polls_late_since"`
    DevicePollsLastLateAt types.String `tfsdk:"device_polls_last_late_at"`
    PacketCaptureCapability types.String `tfsdk:"packet_capture_capability"`
    Labels types.Set `tfsdk:"labels"`
}

func (d *ProbeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_probe"
}

func (d *ProbeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manages custom probes. Deploy probes anywhere in the world and connect it to your project. Look up an existing probe by `id`, or by any of its other arguments (`name`, `connection_status`, `created_by_actor_id`, ...): each one set must match, and exactly one probe may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "key": schema.StringAttribute{
                MarkdownDescription: "The secret the probe connects with, its PROBE_KEY. Leave it out and OneUptime generates one.",
                Computed: true,
                Sensitive: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "probe_version": schema.StringAttribute{
                MarkdownDescription: "A version label for the probe. Leave it out and it is 1.0.0.",
                Computed: true,
            },
            "reported_version": schema.StringAttribute{
                MarkdownDescription: "The OneUptime release this probe last said it runs. Set by OneUptime from the probe's own requests. Empty until the probe connects, and for a probe older than 17.1.0, which does not say.",
                Optional: true,
                Computed: true,
            },
            "major_version_behind_notified_for": schema.NumberAttribute{
                MarkdownDescription: "The major release of OneUptime this probe's owners were last told the probe is behind. Set by OneUptime once a day; empty while the probe runs OneUptime's major release.",
                Optional: true,
                Computed: true,
            },
            "egress_ip_address": schema.StringAttribute{
                MarkdownDescription: "The IP address this probe last checked in from, as OneUptime saw it. Set by OneUptime from the probe's own requests; the blocklist never blocks it.",
                Optional: true,
                Computed: true,
            },
            "interface_ip_addresses": schema.StringAttribute{
                MarkdownDescription: "The IP addresses on this probe's own network interfaces, as the probe last listed them when it checked in, separated by commas. Set by OneUptime from the probe's own requests; Security Events knows the probe's own polls and checks by them, and the blocklist never blocks them.",
                Optional: true,
                Computed: true,
            },
            "last_alive": schema.StringAttribute{
                MarkdownDescription: "When the probe last checked in. Set by OneUptime.",
                Computed: true,
            },
            "icon_file_id": schema.StringAttribute{
                MarkdownDescription: "Probe Page Icon File ID. The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_file"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project the probe belongs to. A global probe belongs to none.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User).",
                Optional: true,
                Computed: true,
            },
            "should_auto_enable_probe_on_new_monitors": schema.BoolAttribute{
                MarkdownDescription: "Auto Enable Probe on New Monitors.",
                Optional: true,
                Computed: true,
            },
            "connection_status": schema.StringAttribute{
                MarkdownDescription: "Connection Status of the Probe.",
                Optional: true,
                Computed: true,
            },
            "checks_late_since": schema.StringAttribute{
                MarkdownDescription: "When this probe's monitors started waiting more than a minute past their check time. Set by OneUptime; empty while they are checked on time.",
                Computed: true,
            },
            "checks_last_late_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last found one of this probe's monitors waiting more than a minute past its check time. Set by OneUptime.",
                Computed: true,
            },
            "device_polls_late_since": schema.StringAttribute{
                MarkdownDescription: "When the network devices this probe polls started waiting more than a minute past their poll time. Set by OneUptime; empty while they are polled on time.",
                Computed: true,
            },
            "device_polls_last_late_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last found one of the network devices this probe polls waiting more than a minute past its poll time. Set by OneUptime.",
                Computed: true,
            },
            "packet_capture_capability": schema.StringAttribute{
                MarkdownDescription: "What the probe last reported about packet capture: whether it is turned on (PROBE_PACKET_CAPTURE_ENABLED on the probe), whether tcpdump is installed, the network interfaces it can capture on, and the limits its operator set. Managed by the probe. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to group and filter by. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *ProbeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    d.client = client
}

func (d *ProbeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ProbeDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.ReportedVersion.IsNull() && !data.ReportedVersion.IsUnknown() {
        filters["reportedVersion"] = data.ReportedVersion.ValueString()
        filterNames = append(filterNames, "reported_version = "+fmt.Sprintf("%q", data.ReportedVersion.ValueString()))
    }
    if !data.MajorVersionBehindNotifiedFor.IsNull() && !data.MajorVersionBehindNotifiedFor.IsUnknown() {
        filters["majorVersionBehindNotifiedFor"] = LookupNumber(data.MajorVersionBehindNotifiedFor)
        filterNames = append(filterNames, "major_version_behind_notified_for = "+data.MajorVersionBehindNotifiedFor.ValueBigFloat().String())
    }
    if !data.EgressIpAddress.IsNull() && !data.EgressIpAddress.IsUnknown() {
        filters["egressIpAddress"] = data.EgressIpAddress.ValueString()
        filterNames = append(filterNames, "egress_ip_address = "+fmt.Sprintf("%q", data.EgressIpAddress.ValueString()))
    }
    if !data.InterfaceIpAddresses.IsNull() && !data.InterfaceIpAddresses.IsUnknown() {
        filters["interfaceIpAddresses"] = data.InterfaceIpAddresses.ValueString()
        filterNames = append(filterNames, "interface_ip_addresses = "+fmt.Sprintf("%q", data.InterfaceIpAddresses.ValueString()))
    }
    if !data.IconFileId.IsNull() && !data.IconFileId.IsUnknown() {
        filters["iconFileId"] = data.IconFileId.ValueString()
        filterNames = append(filterNames, "icon_file_id = "+fmt.Sprintf("%q", data.IconFileId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.ShouldAutoEnableProbeOnNewMonitors.IsNull() && !data.ShouldAutoEnableProbeOnNewMonitors.IsUnknown() {
        filters["shouldAutoEnableProbeOnNewMonitors"] = data.ShouldAutoEnableProbeOnNewMonitors.ValueBool()
        filterNames = append(filterNames, "should_auto_enable_probe_on_new_monitors = "+fmt.Sprintf("%t", data.ShouldAutoEnableProbeOnNewMonitors.ValueBool()))
    }
    if !data.ConnectionStatus.IsNull() && !data.ConnectionStatus.IsUnknown() {
        filters["connectionStatus"] = data.ConnectionStatus.ValueString()
        filterNames = append(filterNames, "connection_status = "+fmt.Sprintf("%q", data.ConnectionStatus.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the probe up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the probe up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "key": true,
        "name": true,
        "description": true,
        "slug": true,
        "probeVersion": true,
        "reportedVersion": true,
        "majorVersionBehindNotifiedFor": true,
        "egressIpAddress": true,
        "interfaceIpAddresses": true,
        "lastAlive": true,
        "iconFileId": true,
        "projectId": true,
        "createdByUserId": true,
        "shouldAutoEnableProbeOnNewMonitors": true,
        "connectionStatus": true,
        "checksLateSince": true,
        "checksLastLateAt": true,
        "devicePollsLateSince": true,
        "devicePollsLastLateAt": true,
        "packetCaptureCapability": true,
        "labels": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/probe/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read probe, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No probe found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read probe: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/probe/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list probe, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list probe: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No probe matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one probe matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for probe.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
    }
    if obj, ok := item["key"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Key = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Key = types.StringValue(string(jsonBytes))
        } else {
            data.Key = types.StringNull()
        }
    } else if arr, ok := item["key"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Key = types.StringValue(string(jsonBytes))
        } else {
            data.Key = types.StringNull()
        }
    } else if val, ok := item["key"].(string); ok {
        data.Key = types.StringValue(val)
    } else {
        data.Key = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := item["name"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := item["description"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := item["slug"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := item["probeVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProbeVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProbeVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProbeVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProbeVersion = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeVersion = types.StringNull()
        }
    } else if arr, ok := item["probeVersion"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProbeVersion = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeVersion = types.StringNull()
        }
    } else if val, ok := item["probeVersion"].(string); ok {
        data.ProbeVersion = types.StringValue(val)
    } else {
        data.ProbeVersion = types.StringNull()
    }
    if obj, ok := item["reportedVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReportedVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReportedVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReportedVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReportedVersion = types.StringValue(string(jsonBytes))
        } else {
            data.ReportedVersion = types.StringNull()
        }
    } else if arr, ok := item["reportedVersion"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ReportedVersion = types.StringValue(string(jsonBytes))
        } else {
            data.ReportedVersion = types.StringNull()
        }
    } else if val, ok := item["reportedVersion"].(string); ok {
        data.ReportedVersion = types.StringValue(val)
    } else {
        data.ReportedVersion = types.StringNull()
    }
    data.MajorVersionBehindNotifiedFor = NumberFromAPI(data.MajorVersionBehindNotifiedFor, item["majorVersionBehindNotifiedFor"])
    if obj, ok := item["egressIpAddress"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EgressIpAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EgressIpAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EgressIpAddress = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EgressIpAddress = types.StringValue(string(jsonBytes))
        } else {
            data.EgressIpAddress = types.StringNull()
        }
    } else if arr, ok := item["egressIpAddress"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EgressIpAddress = types.StringValue(string(jsonBytes))
        } else {
            data.EgressIpAddress = types.StringNull()
        }
    } else if val, ok := item["egressIpAddress"].(string); ok {
        data.EgressIpAddress = types.StringValue(val)
    } else {
        data.EgressIpAddress = types.StringNull()
    }
    if obj, ok := item["interfaceIpAddresses"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InterfaceIpAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.InterfaceIpAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.InterfaceIpAddresses = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.InterfaceIpAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.InterfaceIpAddresses = types.StringNull()
        }
    } else if arr, ok := item["interfaceIpAddresses"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.InterfaceIpAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.InterfaceIpAddresses = types.StringNull()
        }
    } else if val, ok := item["interfaceIpAddresses"].(string); ok {
        data.InterfaceIpAddresses = types.StringValue(val)
    } else {
        data.InterfaceIpAddresses = types.StringNull()
    }
    if obj, ok := item["lastAlive"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastAlive = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastAlive = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastAlive = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastAlive = types.StringValue(string(jsonBytes))
        } else {
            data.LastAlive = types.StringNull()
        }
    } else if arr, ok := item["lastAlive"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastAlive = types.StringValue(string(jsonBytes))
        } else {
            data.LastAlive = types.StringNull()
        }
    } else if val, ok := item["lastAlive"].(string); ok {
        data.LastAlive = types.StringValue(val)
    } else {
        data.LastAlive = types.StringNull()
    }
    if obj, ok := item["iconFileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IconFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IconFileId = types.StringValue(string(jsonBytes))
        } else {
            data.IconFileId = types.StringNull()
        }
    } else if arr, ok := item["iconFileId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.IconFileId = types.StringValue(string(jsonBytes))
        } else {
            data.IconFileId = types.StringNull()
        }
    } else if val, ok := item["iconFileId"].(string); ok {
        data.IconFileId = types.StringValue(val)
    } else {
        data.IconFileId = types.StringNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := item["shouldAutoEnableProbeOnNewMonitors"].(bool); ok {
        data.ShouldAutoEnableProbeOnNewMonitors = types.BoolValue(val)
    } else {
        data.ShouldAutoEnableProbeOnNewMonitors = types.BoolNull()
    }
    if obj, ok := item["connectionStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ConnectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ConnectionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ConnectionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ConnectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.ConnectionStatus = types.StringNull()
        }
    } else if arr, ok := item["connectionStatus"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ConnectionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.ConnectionStatus = types.StringNull()
        }
    } else if val, ok := item["connectionStatus"].(string); ok {
        data.ConnectionStatus = types.StringValue(val)
    } else {
        data.ConnectionStatus = types.StringNull()
    }
    if obj, ok := item["checksLateSince"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChecksLateSince = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChecksLateSince = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChecksLateSince = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChecksLateSince = types.StringValue(string(jsonBytes))
        } else {
            data.ChecksLateSince = types.StringNull()
        }
    } else if arr, ok := item["checksLateSince"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ChecksLateSince = types.StringValue(string(jsonBytes))
        } else {
            data.ChecksLateSince = types.StringNull()
        }
    } else if val, ok := item["checksLateSince"].(string); ok {
        data.ChecksLateSince = types.StringValue(val)
    } else {
        data.ChecksLateSince = types.StringNull()
    }
    if obj, ok := item["checksLastLateAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChecksLastLateAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChecksLastLateAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChecksLastLateAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChecksLastLateAt = types.StringValue(string(jsonBytes))
        } else {
            data.ChecksLastLateAt = types.StringNull()
        }
    } else if arr, ok := item["checksLastLateAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ChecksLastLateAt = types.StringValue(string(jsonBytes))
        } else {
            data.ChecksLastLateAt = types.StringNull()
        }
    } else if val, ok := item["checksLastLateAt"].(string); ok {
        data.ChecksLastLateAt = types.StringValue(val)
    } else {
        data.ChecksLastLateAt = types.StringNull()
    }
    if obj, ok := item["devicePollsLateSince"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DevicePollsLateSince = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DevicePollsLateSince = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DevicePollsLateSince = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DevicePollsLateSince = types.StringValue(string(jsonBytes))
        } else {
            data.DevicePollsLateSince = types.StringNull()
        }
    } else if arr, ok := item["devicePollsLateSince"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DevicePollsLateSince = types.StringValue(string(jsonBytes))
        } else {
            data.DevicePollsLateSince = types.StringNull()
        }
    } else if val, ok := item["devicePollsLateSince"].(string); ok {
        data.DevicePollsLateSince = types.StringValue(val)
    } else {
        data.DevicePollsLateSince = types.StringNull()
    }
    if obj, ok := item["devicePollsLastLateAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DevicePollsLastLateAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DevicePollsLastLateAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DevicePollsLastLateAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DevicePollsLastLateAt = types.StringValue(string(jsonBytes))
        } else {
            data.DevicePollsLastLateAt = types.StringNull()
        }
    } else if arr, ok := item["devicePollsLastLateAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DevicePollsLastLateAt = types.StringValue(string(jsonBytes))
        } else {
            data.DevicePollsLastLateAt = types.StringNull()
        }
    } else if val, ok := item["devicePollsLastLateAt"].(string); ok {
        data.DevicePollsLastLateAt = types.StringValue(val)
    } else {
        data.DevicePollsLastLateAt = types.StringNull()
    }
    if obj, ok := item["packetCaptureCapability"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PacketCaptureCapability = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PacketCaptureCapability = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PacketCaptureCapability = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PacketCaptureCapability = types.StringValue(string(jsonBytes))
        } else {
            data.PacketCaptureCapability = types.StringNull()
        }
    } else if arr, ok := item["packetCaptureCapability"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PacketCaptureCapability = types.StringValue(string(jsonBytes))
        } else {
            data.PacketCaptureCapability = types.StringNull()
        }
    } else if val, ok := item["packetCaptureCapability"].(string); ok {
        data.PacketCaptureCapability = types.StringValue(val)
    } else {
        data.PacketCaptureCapability = types.StringNull()
    }
    if val, ok := item["labels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
