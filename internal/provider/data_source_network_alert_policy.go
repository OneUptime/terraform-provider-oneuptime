package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NetworkAlertPolicyDataSource{}

func NewNetworkAlertPolicyDataSource() datasource.DataSource {
    return &NetworkAlertPolicyDataSource{}
}

// NetworkAlertPolicyDataSource defines the data source implementation.
type NetworkAlertPolicyDataSource struct {
    client *Client
}

// NetworkAlertPolicyDataSourceModel describes the data source data model.
type NetworkAlertPolicyDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    MonitorTemplateId types.String `tfsdk:"monitor_template_id"`
    Scope types.String `tfsdk:"scope"`
    LastSyncAt types.String `tfsdk:"last_sync_at"`
    LastSyncError types.String `tfsdk:"last_sync_error"`
    CoveredDeviceCount types.Number `tfsdk:"covered_device_count"`
    TemplateSyncedAt types.String `tfsdk:"template_synced_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *NetworkAlertPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_alert_policy"
}

func (d *NetworkAlertPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Alert on a set of network devices at once: every device matching the policy's sites, roles and labels gets a Network Device monitor provisioned from the policy's monitor template, and kept as devices come and go. Look up an existing network alert policy by `id`, or by any of its other arguments (`name`, `covered_device_count`, `created_by_user_id`, ...): each one set must match, and exactly one network alert policy may match them all.",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
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
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this policy is active. Disable it to stop provisioning monitors for matching devices without deleting the policy.",
                Optional: true,
                Computed: true,
            },
            "monitor_template_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Network Device monitor template every matching device gets a monitor cloned from. Null only after the template was deleted, which disables the policy. The ID of a `oneuptime_monitor_template`.",
                Optional: true,
                Computed: true,
            },
            "scope": schema.StringAttribute{
                MarkdownDescription: "Which devices this policy covers: site ids, device role ids and label ids. A device must match every kind that is listed (AND) and any id within a kind (OR); a kind left empty matches every device. Empty altogether means every device in the project. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "last_sync_at": schema.StringAttribute{
                MarkdownDescription: "When the engine last reconciled this policy's monitors against its matching devices. Managed by the engine.",
                Computed: true,
            },
            "last_sync_error": schema.StringAttribute{
                MarkdownDescription: "Why the engine's last reconciliation of this policy failed, if it did. Cleared by the next successful pass. Managed by the engine.",
                Optional: true,
                Computed: true,
            },
            "covered_device_count": schema.NumberAttribute{
                MarkdownDescription: "How many devices matched this policy's scope at the engine's last reconciliation. Managed by the engine.",
                Optional: true,
                Computed: true,
            },
            "template_synced_at": schema.StringAttribute{
                MarkdownDescription: "When this policy's provisioned monitors were last re-synced from the monitor template. Managed by the engine.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *NetworkAlertPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkAlertPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkAlertPolicyDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.MonitorTemplateId.IsNull() && !data.MonitorTemplateId.IsUnknown() {
        filters["monitorTemplateId"] = data.MonitorTemplateId.ValueString()
        filterNames = append(filterNames, "monitor_template_id = "+fmt.Sprintf("%q", data.MonitorTemplateId.ValueString()))
    }
    if !data.LastSyncError.IsNull() && !data.LastSyncError.IsUnknown() {
        filters["lastSyncError"] = data.LastSyncError.ValueString()
        filterNames = append(filterNames, "last_sync_error = "+fmt.Sprintf("%q", data.LastSyncError.ValueString()))
    }
    if !data.CoveredDeviceCount.IsNull() && !data.CoveredDeviceCount.IsUnknown() {
        filters["coveredDeviceCount"] = lookupNumber(data.CoveredDeviceCount)
        filterNames = append(filterNames, "covered_device_count = "+data.CoveredDeviceCount.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the network alert policy up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the network alert policy up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "isEnabled": true,
        "monitorTemplateId": true,
        "scope": true,
        "lastSyncAt": true,
        "lastSyncError": true,
        "coveredDeviceCount": true,
        "templateSyncedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/network-alert-policy/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_alert_policy, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network alert policy found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_alert_policy: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/network-alert-policy/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network_alert_policy, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list network_alert_policy: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network alert policy matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network alert policy matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for network_alert_policy.")
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
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["monitorTemplateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorTemplateId = types.StringNull()
        }
    } else if val, ok := item["monitorTemplateId"].(string); ok {
        data.MonitorTemplateId = types.StringValue(val)
    } else {
        data.MonitorTemplateId = types.StringNull()
    }
    if obj, ok := item["scope"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Scope = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Scope = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Scope = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Scope = types.StringValue(string(jsonBytes))
        } else {
            data.Scope = types.StringNull()
        }
    } else if val, ok := item["scope"].(string); ok {
        data.Scope = types.StringValue(val)
    } else {
        data.Scope = types.StringNull()
    }
    if obj, ok := item["lastSyncAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSyncAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSyncAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSyncAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSyncAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSyncAt = types.StringNull()
        }
    } else if val, ok := item["lastSyncAt"].(string); ok {
        data.LastSyncAt = types.StringValue(val)
    } else {
        data.LastSyncAt = types.StringNull()
    }
    if obj, ok := item["lastSyncError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSyncError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSyncError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSyncError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSyncError = types.StringValue(string(jsonBytes))
        } else {
            data.LastSyncError = types.StringNull()
        }
    } else if val, ok := item["lastSyncError"].(string); ok {
        data.LastSyncError = types.StringValue(val)
    } else {
        data.LastSyncError = types.StringNull()
    }
    if val, ok := item["coveredDeviceCount"].(float64); ok {
        data.CoveredDeviceCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["coveredDeviceCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CoveredDeviceCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.CoveredDeviceCount = types.NumberNull()
        }
    } else {
        data.CoveredDeviceCount = types.NumberNull()
    }
    if obj, ok := item["templateSyncedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TemplateSyncedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TemplateSyncedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TemplateSyncedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TemplateSyncedAt = types.StringValue(string(jsonBytes))
        } else {
            data.TemplateSyncedAt = types.StringNull()
        }
    } else if val, ok := item["templateSyncedAt"].(string); ok {
        data.TemplateSyncedAt = types.StringValue(val)
    } else {
        data.TemplateSyncedAt = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
