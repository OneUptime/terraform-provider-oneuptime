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
var _ datasource.DataSource = &NetworkDeviceAutoImportRuleDataSource{}

func NewNetworkDeviceAutoImportRuleDataSource() datasource.DataSource {
    return &NetworkDeviceAutoImportRuleDataSource{}
}

// NetworkDeviceAutoImportRuleDataSource defines the data source implementation.
type NetworkDeviceAutoImportRuleDataSource struct {
    client *Client
}

// NetworkDeviceAutoImportRuleDataSourceModel describes the data source data model.
type NetworkDeviceAutoImportRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    IpMatchTarget types.String `tfsdk:"ip_match_target"`
    SysNamePattern types.String `tfsdk:"sys_name_pattern"`
    SysDescrPattern types.String `tfsdk:"sys_descr_pattern"`
    SysObjectIdPattern types.String `tfsdk:"sys_object_id_pattern"`
    IncludePingOnlyHosts types.Bool `tfsdk:"include_ping_only_hosts"`
    IsExclusion types.Bool `tfsdk:"is_exclusion"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *NetworkDeviceAutoImportRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_device_auto_import_rule"
}

func (d *NetworkDeviceAutoImportRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Automatically import matching hosts from network device discovery scan results as Network Devices, with no manual review step Look up an existing network_device_auto_import_rule by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this network device auto import rule.",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled.",
                Computed: true,
            },
            "ip_match_target": schema.StringAttribute{
                MarkdownDescription: "Only trigger for discovered hosts whose IP is inside this CIDR (192.168.1.0/24) or octet range (10.16-22.0-255.51-66) — the same notations a scan target takes. Leave empty to match any address..",
                Computed: true,
            },
            "sys_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex or * wildcard pattern (case-insensitive) matched against the discovered host's SNMP sysName. Leave empty to match any name..",
                Computed: true,
            },
            "sys_descr_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex or * wildcard pattern (case-insensitive) matched against the discovered host's SNMP sysDescr. Leave empty to match any description..",
                Computed: true,
            },
            "sys_object_id_pattern": schema.StringAttribute{
                MarkdownDescription: "An OID prefix (1.3.6.1.4.1.9) or a '*' wildcard OID pattern with literal dots (1.3.6.1.4.1.9.* for Cisco) matched against the discovered host's SNMP sysObjectID — the vendor's registered enterprise OID. Not regex: dots match dots, so 1.3.6.1.4.1.9.* can never match enterprise 94. Leave empty to match any vendor. Only hosts reported by probes new enough to carry sysObjectID can match..",
                Computed: true,
            },
            "include_ping_only_hosts": schema.BoolAttribute{
                MarkdownDescription: "Also import hosts that answered ping but not SNMP. Off by default: a wrong SNMP credential makes every host on a subnet report as ping-only, and this rule would then import all of them as half-identified devices..",
                Computed: true,
            },
            "is_exclusion": schema.BoolAttribute{
                MarkdownDescription: "Invert this rule: matching hosts are NEVER auto-imported, even when another rule matches them. Use it to carve printers, phones, or other unwanted hosts out of a broader rule..",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *NetworkDeviceAutoImportRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkDeviceAutoImportRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkDeviceAutoImportRuleDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a network_device_auto_import_rule.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "description": true,
        "isEnabled": true,
        "ipMatchTarget": true,
        "sysNamePattern": true,
        "sysDescrPattern": true,
        "sysObjectIdPattern": true,
        "includePingOnlyHosts": true,
        "isExclusion": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/network-device-auto-import-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_device_auto_import_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_device_auto_import_rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_device_auto_import_rule: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.Post(ctx, "/network-device-auto-import-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network_device_auto_import_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list network_device_auto_import_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network_device_auto_import_rule found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network_device_auto_import_rule matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for network_device_auto_import_rule.")
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    if obj, ok := item["ipMatchTarget"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IpMatchTarget = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IpMatchTarget = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IpMatchTarget = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IpMatchTarget = types.StringValue(string(jsonBytes))
        } else {
            data.IpMatchTarget = types.StringNull()
        }
    } else if val, ok := item["ipMatchTarget"].(string); ok {
        data.IpMatchTarget = types.StringValue(val)
    } else {
        data.IpMatchTarget = types.StringNull()
    }
    if obj, ok := item["sysNamePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SysNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SysNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SysNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.SysNamePattern = types.StringNull()
        }
    } else if val, ok := item["sysNamePattern"].(string); ok {
        data.SysNamePattern = types.StringValue(val)
    } else {
        data.SysNamePattern = types.StringNull()
    }
    if obj, ok := item["sysDescrPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysDescrPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SysDescrPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SysDescrPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SysDescrPattern = types.StringValue(string(jsonBytes))
        } else {
            data.SysDescrPattern = types.StringNull()
        }
    } else if val, ok := item["sysDescrPattern"].(string); ok {
        data.SysDescrPattern = types.StringValue(val)
    } else {
        data.SysDescrPattern = types.StringNull()
    }
    if obj, ok := item["sysObjectIdPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SysObjectIdPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SysObjectIdPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SysObjectIdPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SysObjectIdPattern = types.StringValue(string(jsonBytes))
        } else {
            data.SysObjectIdPattern = types.StringNull()
        }
    } else if val, ok := item["sysObjectIdPattern"].(string); ok {
        data.SysObjectIdPattern = types.StringValue(val)
    } else {
        data.SysObjectIdPattern = types.StringNull()
    }
    if val, ok := item["includePingOnlyHosts"].(bool); ok {
        data.IncludePingOnlyHosts = types.BoolValue(val)
    } else {
        data.IncludePingOnlyHosts = types.BoolNull()
    }
    if val, ok := item["isExclusion"].(bool); ok {
        data.IsExclusion = types.BoolValue(val)
    } else {
        data.IsExclusion = types.BoolNull()
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
