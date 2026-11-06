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
var _ datasource.DataSource = &CloudCostDataSource{}

func NewCloudCostDataSource() datasource.DataSource {
    return &CloudCostDataSource{}
}

// CloudCostDataSource defines the data source implementation.
type CloudCostDataSource struct {
    client *Client
}

// CloudCostDataSourceModel describes the data source data model.
type CloudCostDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    CloudAccountId types.String `tfsdk:"cloud_account_id"`
    CloudConnectionId types.String `tfsdk:"cloud_connection_id"`
    ProviderValue types.String `tfsdk:"provider_value"`
    AccountId types.String `tfsdk:"account_id"`
    Day types.String `tfsdk:"day"`
    Service types.String `tfsdk:"service"`
    Region types.String `tfsdk:"region"`
    ResourceId types.String `tfsdk:"resource_id"`
    Cost types.Number `tfsdk:"cost"`
    Currency types.String `tfsdk:"currency"`
    ReadAt types.String `tfsdk:"read_at"`
}

func (d *CloudCostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_cloud_cost"
}

func (d *CloudCostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for Cloud Cost Look up an existing cloud cost by `id`, or by any of its other arguments (`account_id`, `cloud_account_id`, `cloud_connection_id`, ...): each one set must match, and exactly one cloud cost may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID.",
                Computed: true,
            },
            "cloud_account_id": schema.StringAttribute{
                MarkdownDescription: "Cloud Account ID.",
                Optional: true,
                Computed: true,
            },
            "cloud_connection_id": schema.StringAttribute{
                MarkdownDescription: "Cloud Connection ID.",
                Optional: true,
                Computed: true,
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "Cloud Provider.",
                Optional: true,
                Computed: true,
            },
            "account_id": schema.StringAttribute{
                MarkdownDescription: "Account.",
                Optional: true,
                Computed: true,
            },
            "day": schema.StringAttribute{
                MarkdownDescription: "Day.",
                Computed: true,
            },
            "service": schema.StringAttribute{
                MarkdownDescription: "Service.",
                Optional: true,
                Computed: true,
            },
            "region": schema.StringAttribute{
                MarkdownDescription: "Region.",
                Optional: true,
                Computed: true,
            },
            "resource_id": schema.StringAttribute{
                MarkdownDescription: "Resource.",
                Optional: true,
                Computed: true,
            },
            "cost": schema.NumberAttribute{
                MarkdownDescription: "Cost.",
                Optional: true,
                Computed: true,
            },
            "currency": schema.StringAttribute{
                MarkdownDescription: "Currency.",
                Optional: true,
                Computed: true,
            },
            "read_at": schema.StringAttribute{
                MarkdownDescription: "Read At.",
                Computed: true,
            },
        },
    }
}

func (d *CloudCostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudCostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CloudCostDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CloudAccountId.IsNull() && !data.CloudAccountId.IsUnknown() {
        filters["cloudAccountId"] = data.CloudAccountId.ValueString()
        filterNames = append(filterNames, "cloud_account_id = "+fmt.Sprintf("%q", data.CloudAccountId.ValueString()))
    }
    if !data.CloudConnectionId.IsNull() && !data.CloudConnectionId.IsUnknown() {
        filters["cloudConnectionId"] = data.CloudConnectionId.ValueString()
        filterNames = append(filterNames, "cloud_connection_id = "+fmt.Sprintf("%q", data.CloudConnectionId.ValueString()))
    }
    if !data.ProviderValue.IsNull() && !data.ProviderValue.IsUnknown() {
        filters["provider"] = data.ProviderValue.ValueString()
        filterNames = append(filterNames, "provider_value = "+fmt.Sprintf("%q", data.ProviderValue.ValueString()))
    }
    if !data.AccountId.IsNull() && !data.AccountId.IsUnknown() {
        filters["accountId"] = data.AccountId.ValueString()
        filterNames = append(filterNames, "account_id = "+fmt.Sprintf("%q", data.AccountId.ValueString()))
    }
    if !data.Service.IsNull() && !data.Service.IsUnknown() {
        filters["service"] = data.Service.ValueString()
        filterNames = append(filterNames, "service = "+fmt.Sprintf("%q", data.Service.ValueString()))
    }
    if !data.Region.IsNull() && !data.Region.IsUnknown() {
        filters["region"] = data.Region.ValueString()
        filterNames = append(filterNames, "region = "+fmt.Sprintf("%q", data.Region.ValueString()))
    }
    if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
        filters["resourceId"] = data.ResourceId.ValueString()
        filterNames = append(filterNames, "resource_id = "+fmt.Sprintf("%q", data.ResourceId.ValueString()))
    }
    if !data.Cost.IsNull() && !data.Cost.IsUnknown() {
        filters["cost"] = lookupNumber(data.Cost)
        filterNames = append(filterNames, "cost = "+data.Cost.ValueBigFloat().String())
    }
    if !data.Currency.IsNull() && !data.Currency.IsUnknown() {
        filters["currency"] = data.Currency.ValueString()
        filterNames = append(filterNames, "currency = "+fmt.Sprintf("%q", data.Currency.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the cloud cost up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the cloud cost up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "projectId": true,
        "cloudAccountId": true,
        "cloudConnectionId": true,
        "provider": true,
        "accountId": true,
        "day": true,
        "service": true,
        "region": true,
        "resourceId": true,
        "cost": true,
        "currency": true,
        "readAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/cloud-cost/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cloud_cost, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud cost found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read cloud_cost: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/cloud-cost/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list cloud_cost, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list cloud_cost: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud cost matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one cloud cost matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for cloud_cost.")
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
    if obj, ok := item["cloudAccountId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudAccountId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountId = types.StringNull()
        }
    } else if val, ok := item["cloudAccountId"].(string); ok {
        data.CloudAccountId = types.StringValue(val)
    } else {
        data.CloudAccountId = types.StringNull()
    }
    if obj, ok := item["cloudConnectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudConnectionId = types.StringNull()
        }
    } else if val, ok := item["cloudConnectionId"].(string); ok {
        data.CloudConnectionId = types.StringValue(val)
    } else {
        data.CloudConnectionId = types.StringNull()
    }
    if obj, ok := item["provider"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProviderValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProviderValue = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderValue = types.StringNull()
        }
    } else if val, ok := item["provider"].(string); ok {
        data.ProviderValue = types.StringValue(val)
    } else {
        data.ProviderValue = types.StringNull()
    }
    if obj, ok := item["accountId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AccountId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AccountId = types.StringValue(string(jsonBytes))
        } else {
            data.AccountId = types.StringNull()
        }
    } else if val, ok := item["accountId"].(string); ok {
        data.AccountId = types.StringValue(val)
    } else {
        data.AccountId = types.StringNull()
    }
    if obj, ok := item["day"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Day = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Day = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Day = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Day = types.StringValue(string(jsonBytes))
        } else {
            data.Day = types.StringNull()
        }
    } else if val, ok := item["day"].(string); ok {
        data.Day = types.StringValue(val)
    } else {
        data.Day = types.StringNull()
    }
    if obj, ok := item["service"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Service = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Service = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Service = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Service = types.StringValue(string(jsonBytes))
        } else {
            data.Service = types.StringNull()
        }
    } else if val, ok := item["service"].(string); ok {
        data.Service = types.StringValue(val)
    } else {
        data.Service = types.StringNull()
    }
    if obj, ok := item["region"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Region = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Region = types.StringValue(string(jsonBytes))
        } else {
            data.Region = types.StringNull()
        }
    } else if val, ok := item["region"].(string); ok {
        data.Region = types.StringValue(val)
    } else {
        data.Region = types.StringNull()
    }
    if obj, ok := item["resourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceId = types.StringNull()
        }
    } else if val, ok := item["resourceId"].(string); ok {
        data.ResourceId = types.StringValue(val)
    } else {
        data.ResourceId = types.StringNull()
    }
    data.Cost = numberFromAPI(data.Cost, item["cost"])
    if obj, ok := item["currency"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Currency = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Currency = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Currency = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Currency = types.StringValue(string(jsonBytes))
        } else {
            data.Currency = types.StringNull()
        }
    } else if val, ok := item["currency"].(string); ok {
        data.Currency = types.StringValue(val)
    } else {
        data.Currency = types.StringNull()
    }
    if obj, ok := item["readAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReadAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReadAt = types.StringValue(string(jsonBytes))
        } else {
            data.ReadAt = types.StringNull()
        }
    } else if val, ok := item["readAt"].(string); ok {
        data.ReadAt = types.StringValue(val)
    } else {
        data.ReadAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
