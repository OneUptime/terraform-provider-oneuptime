package provider

import (
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
    Key types.String `tfsdk:"key"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    ProbeVersion types.String `tfsdk:"probe_version"`
    LastAlive types.String `tfsdk:"last_alive"`
    IconFileId types.String `tfsdk:"icon_file_id"`
    ProjectId types.String `tfsdk:"project_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    ShouldAutoEnableProbeOnNewMonitors types.Bool `tfsdk:"should_auto_enable_probe_on_new_monitors"`
    ConnectionStatus types.String `tfsdk:"connection_status"`
    PacketCaptureCapability types.String `tfsdk:"packet_capture_capability"`
    Labels types.Set `tfsdk:"labels"`
}

func (d *ProbeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_probe"
}

func (d *ProbeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manages custom probes. Deploy probes anywhere in the world and connect it to your project. Look up an existing probe by `id`, or by any of its other arguments (`name`, `connection_status`, `created_by_user_id`, ...): each one set must match, and exactly one probe may match them all.",

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
            "key": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Project Owner, Project Admin], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Probe]",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Settings Admin, Settings Member, Settings Viewer, Read Probe, Create Monitor, Edit Monitor, Read Monitor, Create All Operational Resources, Edit All Operational Resources, Read All Operational Resources, Create Monitor Probe, Edit Monitor Probe, Read Monitor Probe, Create Network Device, Edit Network Device, Read Network Device, Create Network Device Discovery Scan, Edit Network Device Discovery Scan, Read Network Device Discovery Scan, Create Network Site, Edit Network Site, Read Network Site], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Probe]",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Settings Admin, Settings Member, Settings Viewer, Read Probe, Create Monitor, Edit Monitor, Read Monitor, Create All Operational Resources, Edit All Operational Resources, Read All Operational Resources, Create Monitor Probe, Edit Monitor Probe, Read Monitor Probe, Create Network Device, Edit Network Device, Read Network Device, Create Network Device Discovery Scan, Edit Network Device Discovery Scan, Read Network Device Discovery Scan, Create Network Site, Edit Network Site, Read Network Site], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Probe]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "probe_version": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Settings Admin, Settings Member, Settings Viewer, Read Probe], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Probe]",
                Computed: true,
            },
            "last_alive": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Settings Admin, Settings Member, Settings Viewer, Read Probe, Create Monitor, Edit Monitor, Read Monitor, Create All Operational Resources, Edit All Operational Resources, Read All Operational Resources, Create Monitor Probe, Edit Monitor Probe, Read Monitor Probe, Create Network Device, Edit Network Device, Read Network Device, Create Network Device Discovery Scan, Edit Network Device Discovery Scan, Read Network Device Discovery Scan, Create Network Site, Edit Network Site, Read Network Site], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "icon_file_id": schema.StringAttribute{
                MarkdownDescription: "Probe Page Icon File ID. The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Probe], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Settings Admin, Settings Member, Settings Viewer, Read Probe, Create Monitor, Edit Monitor, Read Monitor, Create All Operational Resources, Edit All Operational Resources, Read All Operational Resources, Create Monitor Probe, Edit Monitor Probe, Read Monitor Probe, Create Network Device, Edit Network Device, Read Network Device, Create Network Device Discovery Scan, Edit Network Device Discovery Scan, Read Network Device Discovery Scan, Create Network Site, Edit Network Site, Read Network Site], Update: [No access - you don't have permission for this operation]",
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
            "packet_capture_capability": schema.StringAttribute{
                MarkdownDescription: "What the probe last reported about packet capture: whether it is turned on (PROBE_PACKET_CAPTURE_ENABLED on the probe), whether tcpdump is installed, the network interfaces it can capture on, and the limits its operator set. Managed by the probe. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
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
    if !data.Key.IsNull() && !data.Key.IsUnknown() {
        filters["key"] = data.Key.ValueString()
        filterNames = append(filterNames, "key = "+fmt.Sprintf("%q", data.Key.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
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
        "key": true,
        "name": true,
        "description": true,
        "slug": true,
        "probeVersion": true,
        "lastAlive": true,
        "iconFileId": true,
        "projectId": true,
        "createdByUserId": true,
        "shouldAutoEnableProbeOnNewMonitors": true,
        "connectionStatus": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No probe matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one probe matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
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
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
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
    } else if val, ok := item["probeVersion"].(string); ok {
        data.ProbeVersion = types.StringValue(val)
    } else {
        data.ProbeVersion = types.StringNull()
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
    } else if val, ok := item["connectionStatus"].(string); ok {
        data.ConnectionStatus = types.StringValue(val)
    } else {
        data.ConnectionStatus = types.StringNull()
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
