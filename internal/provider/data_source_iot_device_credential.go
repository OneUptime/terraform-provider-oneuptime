package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IotDeviceCredentialDataSource{}

func NewIotDeviceCredentialDataSource() datasource.DataSource {
    return &IotDeviceCredentialDataSource{}
}

// NewIotDeviceCredentialLegacyDataSource registers this data source under the
// name it had before type names kept mixed-case words whole,
// oneuptime_io_t_device_credential. Deprecated.
func NewIotDeviceCredentialLegacyDataSource() datasource.DataSource {
    return &IotDeviceCredentialDataSource{isLegacyAlias: true}
}

// IotDeviceCredentialDataSource defines the data source implementation.
type IotDeviceCredentialDataSource struct {
    client *Client
    // Registered under the name this data source had before.
    isLegacyAlias bool
}

// IotDeviceCredentialDataSourceModel describes the data source data model.
type IotDeviceCredentialDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    IotFleetId types.String `tfsdk:"iot_fleet_id"`
    ExternalId types.String `tfsdk:"external_id"`
    Name types.String `tfsdk:"name"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    LastConnectedAt types.String `tfsdk:"last_connected_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    SecretKey types.String `tfsdk:"secret_key"`
}

func (d *IotDeviceCredentialDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    if d.isLegacyAlias {
        resp.TypeName = req.ProviderTypeName + "_io_t_device_credential"
        return
    }
    resp.TypeName = req.ProviderTypeName + "_iot_device_credential"
}

func (d *IotDeviceCredentialDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Registered IoT devices and their per-device MQTT credentials. Registered devices get individual authentication, topic isolation, revocation, and silent-death offline detection. Look up an existing iot device credential by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `external_id`, ...): each one set must match, and exactly one iot device credential may match them all.",

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
            "iot_fleet_id": schema.StringAttribute{
                MarkdownDescription: "ID of the IoT Fleet this device belongs to. The ID of a `oneuptime_iot_fleet`.",
                Optional: true,
                Computed: true,
            },
            "external_id": schema.StringAttribute{
                MarkdownDescription: "The device id — must match the device.id label the device stamps on its datapoints. It is also the <device> segment of the device's MQTT topics, so a device that reports directly over MQTT cannot use an id containing '/', '+', or '#' (such devices can still report through a gateway).",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this device.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Disabled credentials are rejected at MQTT CONNECT and stop the device's silent-death offline detection.",
                Optional: true,
                Computed: true,
            },
            "last_connected_at": schema.StringAttribute{
                MarkdownDescription: "When this credential last authenticated an MQTT connection.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "secret_key": schema.StringAttribute{
                MarkdownDescription: "Secret this device presents as the MQTT password (with this credential's ID as the username).",
                Optional: true,
                Computed: true,
            },
        },
    }
    if d.isLegacyAlias {
        resp.Schema.DeprecationMessage = "oneuptime_io_t_device_credential has been renamed to oneuptime_iot_device_credential. Use the new name; the old one keeps working until then."
    }
}

func (d *IotDeviceCredentialDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IotDeviceCredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IotDeviceCredentialDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.IotFleetId.IsNull() && !data.IotFleetId.IsUnknown() {
        filters["iotFleetId"] = data.IotFleetId.ValueString()
        filterNames = append(filterNames, "iot_fleet_id = "+fmt.Sprintf("%q", data.IotFleetId.ValueString()))
    }
    if !data.ExternalId.IsNull() && !data.ExternalId.IsUnknown() {
        filters["externalId"] = data.ExternalId.ValueString()
        filterNames = append(filterNames, "external_id = "+fmt.Sprintf("%q", data.ExternalId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.SecretKey.IsNull() && !data.SecretKey.IsUnknown() {
        filters["secretKey"] = data.SecretKey.ValueString()
        filterNames = append(filterNames, "secret_key = "+fmt.Sprintf("%q", data.SecretKey.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the iot device credential up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the iot device credential up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "iotFleetId": true,
        "externalId": true,
        "name": true,
        "isEnabled": true,
        "lastConnectedAt": true,
        "createdByUserId": true,
        "secretKey": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/iot-device-credential/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read iot_device_credential, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No iot device credential found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read iot_device_credential: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/iot-device-credential/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list iot_device_credential, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list iot_device_credential: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No iot device credential matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one iot device credential matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for iot_device_credential.")
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
    if obj, ok := item["iotFleetId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IotFleetId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IotFleetId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IotFleetId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IotFleetId = types.StringValue(string(jsonBytes))
        } else {
            data.IotFleetId = types.StringNull()
        }
    } else if val, ok := item["iotFleetId"].(string); ok {
        data.IotFleetId = types.StringValue(val)
    } else {
        data.IotFleetId = types.StringNull()
    }
    if obj, ok := item["externalId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExternalId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExternalId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExternalId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExternalId = types.StringValue(string(jsonBytes))
        } else {
            data.ExternalId = types.StringNull()
        }
    } else if val, ok := item["externalId"].(string); ok {
        data.ExternalId = types.StringValue(val)
    } else {
        data.ExternalId = types.StringNull()
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
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["lastConnectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastConnectedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastConnectedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastConnectedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastConnectedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastConnectedAt = types.StringNull()
        }
    } else if val, ok := item["lastConnectedAt"].(string); ok {
        data.LastConnectedAt = types.StringValue(val)
    } else {
        data.LastConnectedAt = types.StringNull()
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
    if obj, ok := item["secretKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.SecretKey = types.StringNull()
        }
    } else if val, ok := item["secretKey"].(string); ok {
        data.SecretKey = types.StringValue(val)
    } else {
        data.SecretKey = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
