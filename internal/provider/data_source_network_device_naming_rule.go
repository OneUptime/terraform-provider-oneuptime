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
var _ datasource.DataSource = &NetworkDeviceNamingRuleDataSource{}

func NewNetworkDeviceNamingRuleDataSource() datasource.DataSource {
    return &NetworkDeviceNamingRuleDataSource{}
}

// NetworkDeviceNamingRuleDataSource defines the data source implementation.
type NetworkDeviceNamingRuleDataSource struct {
    client *Client
}

// NetworkDeviceNamingRuleDataSourceModel describes the data source data model.
type NetworkDeviceNamingRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    Criteria types.String `tfsdk:"criteria"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    Priority types.Number `tfsdk:"priority"`
    IpMatchTarget types.String `tfsdk:"ip_match_target"`
    DiscoveredNamePattern types.String `tfsdk:"discovered_name_pattern"`
    MacAddressPattern types.String `tfsdk:"mac_address_pattern"`
    VendorPattern types.String `tfsdk:"vendor_pattern"`
    Sites types.Set `tfsdk:"sites"`
    NamePattern types.String `tfsdk:"name_pattern"`
    LookupKey types.String `tfsdk:"lookup_key"`
    LookupEntries types.String `tfsdk:"lookup_entries"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *NetworkDeviceNamingRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_network_device_naming_rule"
}

func (d *NetworkDeviceNamingRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Rules that name network devices by your naming convention: which devices, and a name pattern such as UN{site.code}{lookup}. The first enabled rule that matches a device names it when discovery imports it, and the rules can be applied to the devices you already have. Look up an existing network device naming rule by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one network device naming rule may match them all.",

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
            "criteria": schema.StringAttribute{
                MarkdownDescription: "Versioned conditions that determine whether this rule matches a resource. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "A name for this rule, such as Menu boards.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What this rule is for.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule names devices. A disabled rule is skipped and the next one is tried.",
                Optional: true,
                Computed: true,
            },
            "priority": schema.NumberAttribute{
                MarkdownDescription: "Where this rule sits in the list: the first enabled rule from the top that matches a device, and has every value its name pattern needs, names it. The rule at the top has the highest number. A new rule is added to the end of the list, with the lowest number. Setting a number another rule already has puts this rule in that place, and the rules in the way move one place along to make room. In the dashboard, drag the rows to reorder them.",
                Optional: true,
                Computed: true,
            },
            "ip_match_target": schema.StringAttribute{
                MarkdownDescription: "Only devices whose IP address is in this range: a CIDR such as 10.48.145.0/24, or octet ranges such as 10.0-255.0-255.150-159, where * is any number (10.*.*.51). Leave empty to match any address.",
                Optional: true,
                Computed: true,
            },
            "discovered_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex or * wildcard pattern (case-insensitive) matched against the name discovery found for the device. Leave empty to match any name.",
                Optional: true,
                Computed: true,
            },
            "mac_address_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex or * wildcard pattern (case-insensitive) matched against the device's MAC address, in any spelling: 00:05:0b*, 00-05-0B* and 00050b* match the same devices. Leave empty to match any MAC address.",
                Optional: true,
                Computed: true,
            },
            "vendor_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex or * wildcard pattern (case-insensitive) matched against who made the device. Leave empty to match any vendor.",
                Optional: true,
                Computed: true,
            },
            "sites": schema.SetAttribute{
                MarkdownDescription: "Only devices in one of these sites, or in a site below one of them. Leave empty to match devices in any site, or in none. IDs of `oneuptime_network_site` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "name_pattern": schema.StringAttribute{
                MarkdownDescription: "What matched devices are called. Text, with placeholders in curly braces filled in from each device: {site.code}, {site}, {discoveredName}, {ip}, {ip.1} to {ip.4}, {mac}, {vendor} and {lookup}, the value of the rule's lookup table. {n}, {nn}, {nnn} or {nnnn} is the lowest number that makes the name unique. A placeholder can be followed by filters: {vendor|first:3|upper}, with upper, lower, first:N, last:N, digits and nospaces.",
                Optional: true,
                Computed: true,
            },
            "lookup_key": schema.StringAttribute{
                MarkdownDescription: "What the rows of the lookup table are keyed by: lastIpOctet (the last part of the device's IP address), ipAddress, macAddress or discoveredName.",
                Optional: true,
                Computed: true,
            },
            "lookup_entries": schema.StringAttribute{
                MarkdownDescription: "Rows of a key and the value it stands for, put where the name pattern says {lookup}. A device whose key has no row is left to the next rule. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *NetworkDeviceNamingRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkDeviceNamingRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data NetworkDeviceNamingRuleDataSourceModel

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
    if !data.Priority.IsNull() && !data.Priority.IsUnknown() {
        filters["priority"] = lookupNumber(data.Priority)
        filterNames = append(filterNames, "priority = "+data.Priority.ValueBigFloat().String())
    }
    if !data.IpMatchTarget.IsNull() && !data.IpMatchTarget.IsUnknown() {
        filters["ipMatchTarget"] = data.IpMatchTarget.ValueString()
        filterNames = append(filterNames, "ip_match_target = "+fmt.Sprintf("%q", data.IpMatchTarget.ValueString()))
    }
    if !data.DiscoveredNamePattern.IsNull() && !data.DiscoveredNamePattern.IsUnknown() {
        filters["discoveredNamePattern"] = data.DiscoveredNamePattern.ValueString()
        filterNames = append(filterNames, "discovered_name_pattern = "+fmt.Sprintf("%q", data.DiscoveredNamePattern.ValueString()))
    }
    if !data.MacAddressPattern.IsNull() && !data.MacAddressPattern.IsUnknown() {
        filters["macAddressPattern"] = data.MacAddressPattern.ValueString()
        filterNames = append(filterNames, "mac_address_pattern = "+fmt.Sprintf("%q", data.MacAddressPattern.ValueString()))
    }
    if !data.VendorPattern.IsNull() && !data.VendorPattern.IsUnknown() {
        filters["vendorPattern"] = data.VendorPattern.ValueString()
        filterNames = append(filterNames, "vendor_pattern = "+fmt.Sprintf("%q", data.VendorPattern.ValueString()))
    }
    if !data.NamePattern.IsNull() && !data.NamePattern.IsUnknown() {
        filters["namePattern"] = data.NamePattern.ValueString()
        filterNames = append(filterNames, "name_pattern = "+fmt.Sprintf("%q", data.NamePattern.ValueString()))
    }
    if !data.LookupKey.IsNull() && !data.LookupKey.IsUnknown() {
        filters["lookupKey"] = data.LookupKey.ValueString()
        filterNames = append(filterNames, "lookup_key = "+fmt.Sprintf("%q", data.LookupKey.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the network device naming rule up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the network device naming rule up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "criteria": true,
        "projectId": true,
        "name": true,
        "description": true,
        "isEnabled": true,
        "priority": true,
        "ipMatchTarget": true,
        "discoveredNamePattern": true,
        "macAddressPattern": true,
        "vendorPattern": true,
        "sites": true,
        "namePattern": true,
        "lookupKey": true,
        "lookupEntries": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/network-device-naming-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network_device_naming_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network device naming rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read network_device_naming_rule: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/network-device-naming-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network_device_naming_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list network_device_naming_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No network device naming rule matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one network device naming rule matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for network_device_naming_rule.")
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
    if obj, ok := item["criteria"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Criteria = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Criteria = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Criteria = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Criteria = types.StringValue(string(jsonBytes))
        } else {
            data.Criteria = types.StringNull()
        }
    } else if arr, ok := item["criteria"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Criteria = types.StringValue(string(jsonBytes))
        } else {
            data.Criteria = types.StringNull()
        }
    } else if val, ok := item["criteria"].(string); ok {
        data.Criteria = types.StringValue(val)
    } else {
        data.Criteria = types.StringNull()
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
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    data.Priority = numberFromAPI(data.Priority, item["priority"])
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
    } else if arr, ok := item["ipMatchTarget"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.IpMatchTarget = types.StringValue(string(jsonBytes))
        } else {
            data.IpMatchTarget = types.StringNull()
        }
    } else if val, ok := item["ipMatchTarget"].(string); ok {
        data.IpMatchTarget = types.StringValue(val)
    } else {
        data.IpMatchTarget = types.StringNull()
    }
    if obj, ok := item["discoveredNamePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DiscoveredNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DiscoveredNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DiscoveredNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DiscoveredNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoveredNamePattern = types.StringNull()
        }
    } else if arr, ok := item["discoveredNamePattern"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DiscoveredNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoveredNamePattern = types.StringNull()
        }
    } else if val, ok := item["discoveredNamePattern"].(string); ok {
        data.DiscoveredNamePattern = types.StringValue(val)
    } else {
        data.DiscoveredNamePattern = types.StringNull()
    }
    if obj, ok := item["macAddressPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MacAddressPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MacAddressPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MacAddressPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MacAddressPattern = types.StringValue(string(jsonBytes))
        } else {
            data.MacAddressPattern = types.StringNull()
        }
    } else if arr, ok := item["macAddressPattern"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MacAddressPattern = types.StringValue(string(jsonBytes))
        } else {
            data.MacAddressPattern = types.StringNull()
        }
    } else if val, ok := item["macAddressPattern"].(string); ok {
        data.MacAddressPattern = types.StringValue(val)
    } else {
        data.MacAddressPattern = types.StringNull()
    }
    if obj, ok := item["vendorPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VendorPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VendorPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VendorPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VendorPattern = types.StringValue(string(jsonBytes))
        } else {
            data.VendorPattern = types.StringNull()
        }
    } else if arr, ok := item["vendorPattern"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.VendorPattern = types.StringValue(string(jsonBytes))
        } else {
            data.VendorPattern = types.StringNull()
        }
    } else if val, ok := item["vendorPattern"].(string); ok {
        data.VendorPattern = types.StringValue(val)
    } else {
        data.VendorPattern = types.StringNull()
    }
    if val, ok := item["sites"].([]interface{}); ok {
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
        data.Sites = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Sites = types.SetNull(types.StringType)
    }
    if obj, ok := item["namePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.NamePattern = types.StringNull()
        }
    } else if arr, ok := item["namePattern"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.NamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.NamePattern = types.StringNull()
        }
    } else if val, ok := item["namePattern"].(string); ok {
        data.NamePattern = types.StringValue(val)
    } else {
        data.NamePattern = types.StringNull()
    }
    if obj, ok := item["lookupKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LookupKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LookupKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LookupKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LookupKey = types.StringValue(string(jsonBytes))
        } else {
            data.LookupKey = types.StringNull()
        }
    } else if arr, ok := item["lookupKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LookupKey = types.StringValue(string(jsonBytes))
        } else {
            data.LookupKey = types.StringNull()
        }
    } else if val, ok := item["lookupKey"].(string); ok {
        data.LookupKey = types.StringValue(val)
    } else {
        data.LookupKey = types.StringNull()
    }
    if obj, ok := item["lookupEntries"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LookupEntries = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LookupEntries = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LookupEntries = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LookupEntries = types.StringValue(string(jsonBytes))
        } else {
            data.LookupEntries = types.StringNull()
        }
    } else if arr, ok := item["lookupEntries"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LookupEntries = types.StringValue(string(jsonBytes))
        } else {
            data.LookupEntries = types.StringNull()
        }
    } else if val, ok := item["lookupEntries"].(string); ok {
        data.LookupEntries = types.StringValue(val)
    } else {
        data.LookupEntries = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
