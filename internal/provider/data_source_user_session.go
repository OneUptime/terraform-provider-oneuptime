package provider

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserSessionDataSource{}

func NewUserSessionDataSource() datasource.DataSource {
    return &UserSessionDataSource{}
}

// UserSessionDataSource defines the data source implementation.
type UserSessionDataSource struct {
    client *Client
}

// UserSessionDataSourceModel describes the data source data model.
type UserSessionDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    UserId types.String `tfsdk:"user_id"`
    LastActiveAt types.String `tfsdk:"last_active_at"`
    DeviceName types.String `tfsdk:"device_name"`
    DeviceType types.String `tfsdk:"device_type"`
    DeviceOs types.String `tfsdk:"device_os"`
    DeviceBrowser types.String `tfsdk:"device_browser"`
    IpAddress types.String `tfsdk:"ip_address"`
    UserAgent types.String `tfsdk:"user_agent"`
    IsRevoked types.Bool `tfsdk:"is_revoked"`
    RevokedAt types.String `tfsdk:"revoked_at"`
    RevokedReason types.String `tfsdk:"revoked_reason"`
    AdditionalInfo types.String `tfsdk:"additional_info"`
}

func (d *UserSessionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_user_session"
}

func (d *UserSessionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Active user sessions with refresh tokens and device metadata for enhanced authentication security. Look up an existing user session by `id`, or by any of its other arguments (`device_browser`, `device_name`, `device_os`, ...): each one set must match, and exactly one user session may match them all.",

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
            "user_id": schema.StringAttribute{
                MarkdownDescription: "Identifier for the user that owns this session. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "last_active_at": schema.StringAttribute{
                MarkdownDescription: "Last time this session was used.",
                Computed: true,
            },
            "device_name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for the device used to sign in.",
                Optional: true,
                Computed: true,
            },
            "device_type": schema.StringAttribute{
                MarkdownDescription: "Type of device (e.g., desktop, mobile).",
                Optional: true,
                Computed: true,
            },
            "device_os": schema.StringAttribute{
                MarkdownDescription: "Operating system reported for this session.",
                Optional: true,
                Computed: true,
            },
            "device_browser": schema.StringAttribute{
                MarkdownDescription: "Browser or client application used for this session.",
                Optional: true,
                Computed: true,
            },
            "ip_address": schema.StringAttribute{
                MarkdownDescription: "IP address observed for this session.",
                Optional: true,
                Computed: true,
            },
            "user_agent": schema.StringAttribute{
                MarkdownDescription: "Complete user agent string supplied by the client.",
                Optional: true,
                Computed: true,
            },
            "is_revoked": schema.BoolAttribute{
                MarkdownDescription: "Marks whether the session has been explicitly revoked.",
                Optional: true,
                Computed: true,
            },
            "revoked_at": schema.StringAttribute{
                MarkdownDescription: "Timestamp when the session was revoked, if applicable.",
                Computed: true,
            },
            "revoked_reason": schema.StringAttribute{
                MarkdownDescription: "Optional reason describing why the session was revoked.",
                Optional: true,
                Computed: true,
            },
            "additional_info": schema.StringAttribute{
                MarkdownDescription: "Flexible JSON payload for storing structured session metadata. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
        },
    }
}

func (d *UserSessionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserSessionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data UserSessionDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }
    if !data.DeviceName.IsNull() && !data.DeviceName.IsUnknown() {
        filters["deviceName"] = data.DeviceName.ValueString()
        filterNames = append(filterNames, "device_name = "+fmt.Sprintf("%q", data.DeviceName.ValueString()))
    }
    if !data.DeviceType.IsNull() && !data.DeviceType.IsUnknown() {
        filters["deviceType"] = data.DeviceType.ValueString()
        filterNames = append(filterNames, "device_type = "+fmt.Sprintf("%q", data.DeviceType.ValueString()))
    }
    if !data.DeviceOs.IsNull() && !data.DeviceOs.IsUnknown() {
        filters["deviceOS"] = data.DeviceOs.ValueString()
        filterNames = append(filterNames, "device_os = "+fmt.Sprintf("%q", data.DeviceOs.ValueString()))
    }
    if !data.DeviceBrowser.IsNull() && !data.DeviceBrowser.IsUnknown() {
        filters["deviceBrowser"] = data.DeviceBrowser.ValueString()
        filterNames = append(filterNames, "device_browser = "+fmt.Sprintf("%q", data.DeviceBrowser.ValueString()))
    }
    if !data.IpAddress.IsNull() && !data.IpAddress.IsUnknown() {
        filters["ipAddress"] = data.IpAddress.ValueString()
        filterNames = append(filterNames, "ip_address = "+fmt.Sprintf("%q", data.IpAddress.ValueString()))
    }
    if !data.UserAgent.IsNull() && !data.UserAgent.IsUnknown() {
        filters["userAgent"] = data.UserAgent.ValueString()
        filterNames = append(filterNames, "user_agent = "+fmt.Sprintf("%q", data.UserAgent.ValueString()))
    }
    if !data.IsRevoked.IsNull() && !data.IsRevoked.IsUnknown() {
        filters["isRevoked"] = data.IsRevoked.ValueBool()
        filterNames = append(filterNames, "is_revoked = "+fmt.Sprintf("%t", data.IsRevoked.ValueBool()))
    }
    if !data.RevokedReason.IsNull() && !data.RevokedReason.IsUnknown() {
        filters["revokedReason"] = data.RevokedReason.ValueString()
        filterNames = append(filterNames, "revoked_reason = "+fmt.Sprintf("%q", data.RevokedReason.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the user session up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the user session up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "userId": true,
        "lastActiveAt": true,
        "deviceName": true,
        "deviceType": true,
        "deviceOS": true,
        "deviceBrowser": true,
        "ipAddress": true,
        "userAgent": true,
        "isRevoked": true,
        "revokedAt": true,
        "revokedReason": true,
        "additionalInfo": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        // No get endpoint: find it in the list by id.
        filters["_id"] = data.Id.ValueString()
        filterNames = append(filterNames, fmt.Sprintf("id = %q", data.Id.ValueString()))
    }
    if item == nil {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/user-session/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list user_session, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list user_session: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No user session matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one user session matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for user_session.")
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
    if obj, ok := item["userId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserId = types.StringValue(string(jsonBytes))
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := item["userId"].(string); ok {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if obj, ok := item["lastActiveAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastActiveAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastActiveAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastActiveAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastActiveAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastActiveAt = types.StringNull()
        }
    } else if val, ok := item["lastActiveAt"].(string); ok {
        data.LastActiveAt = types.StringValue(val)
    } else {
        data.LastActiveAt = types.StringNull()
    }
    if obj, ok := item["deviceName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeviceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeviceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeviceName = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceName = types.StringNull()
        }
    } else if val, ok := item["deviceName"].(string); ok {
        data.DeviceName = types.StringValue(val)
    } else {
        data.DeviceName = types.StringNull()
    }
    if obj, ok := item["deviceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeviceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeviceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeviceType = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceType = types.StringNull()
        }
    } else if val, ok := item["deviceType"].(string); ok {
        data.DeviceType = types.StringValue(val)
    } else {
        data.DeviceType = types.StringNull()
    }
    if obj, ok := item["deviceOS"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceOs = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeviceOs = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeviceOs = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeviceOs = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceOs = types.StringNull()
        }
    } else if val, ok := item["deviceOS"].(string); ok {
        data.DeviceOs = types.StringValue(val)
    } else {
        data.DeviceOs = types.StringNull()
    }
    if obj, ok := item["deviceBrowser"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeviceBrowser = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeviceBrowser = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeviceBrowser = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeviceBrowser = types.StringValue(string(jsonBytes))
        } else {
            data.DeviceBrowser = types.StringNull()
        }
    } else if val, ok := item["deviceBrowser"].(string); ok {
        data.DeviceBrowser = types.StringValue(val)
    } else {
        data.DeviceBrowser = types.StringNull()
    }
    if obj, ok := item["ipAddress"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IpAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IpAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IpAddress = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IpAddress = types.StringValue(string(jsonBytes))
        } else {
            data.IpAddress = types.StringNull()
        }
    } else if val, ok := item["ipAddress"].(string); ok {
        data.IpAddress = types.StringValue(val)
    } else {
        data.IpAddress = types.StringNull()
    }
    if obj, ok := item["userAgent"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserAgent = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserAgent = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserAgent = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserAgent = types.StringValue(string(jsonBytes))
        } else {
            data.UserAgent = types.StringNull()
        }
    } else if val, ok := item["userAgent"].(string); ok {
        data.UserAgent = types.StringValue(val)
    } else {
        data.UserAgent = types.StringNull()
    }
    if val, ok := item["isRevoked"].(bool); ok {
        data.IsRevoked = types.BoolValue(val)
    } else {
        data.IsRevoked = types.BoolNull()
    }
    if obj, ok := item["revokedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RevokedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RevokedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RevokedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RevokedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RevokedAt = types.StringNull()
        }
    } else if val, ok := item["revokedAt"].(string); ok {
        data.RevokedAt = types.StringValue(val)
    } else {
        data.RevokedAt = types.StringNull()
    }
    if obj, ok := item["revokedReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RevokedReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RevokedReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RevokedReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RevokedReason = types.StringValue(string(jsonBytes))
        } else {
            data.RevokedReason = types.StringNull()
        }
    } else if val, ok := item["revokedReason"].(string); ok {
        data.RevokedReason = types.StringValue(val)
    } else {
        data.RevokedReason = types.StringNull()
    }
    if obj, ok := item["additionalInfo"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AdditionalInfo = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AdditionalInfo = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AdditionalInfo = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AdditionalInfo = types.StringValue(string(jsonBytes))
        } else {
            data.AdditionalInfo = types.StringNull()
        }
    } else if val, ok := item["additionalInfo"].(string); ok {
        data.AdditionalInfo = types.StringValue(val)
    } else {
        data.AdditionalInfo = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
