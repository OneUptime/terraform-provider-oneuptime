package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MonitorTemplateDataSource{}

func NewMonitorTemplateDataSource() datasource.DataSource {
    return &MonitorTemplateDataSource{}
}

// MonitorTemplateDataSource defines the data source implementation.
type MonitorTemplateDataSource struct {
    client *Client
}

// MonitorTemplateDataSourceModel describes the data source data model.
type MonitorTemplateDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    TemplateName types.String `tfsdk:"template_name"`
    TemplateDescription types.String `tfsdk:"template_description"`
    Slug types.String `tfsdk:"slug"`
    MonitorName types.String `tfsdk:"monitor_name"`
    MonitorDescription types.String `tfsdk:"monitor_description"`
    MonitorType types.String `tfsdk:"monitor_type"`
    MonitorSteps types.String `tfsdk:"monitor_steps"`
    MonitoringInterval types.String `tfsdk:"monitoring_interval"`
    Labels types.Set `tfsdk:"labels"`
    CustomFields types.String `tfsdk:"custom_fields"`
    MinimumProbeAgreement types.Number `tfsdk:"minimum_probe_agreement"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *MonitorTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_monitor_template"
}

func (d *MonitorTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Reusable monitor template. Use it to create new monitors with the same configuration. Look up an existing monitor template by `id`, or by any of its other arguments (`created_by_user_id`, `minimum_probe_agreement`, `monitor_description`, ...): each one set must match, and exactly one monitor template may match them all.",

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
            "template_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Monitor Template.",
                Optional: true,
                Computed: true,
            },
            "template_description": schema.StringAttribute{
                MarkdownDescription: "Description of the Monitor Template.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "monitor_name": schema.StringAttribute{
                MarkdownDescription: "Default name applied to monitors created from this template. Users can override on creation. Leave it blank to name each monitor after the resource it watches.",
                Optional: true,
                Computed: true,
            },
            "monitor_description": schema.StringAttribute{
                MarkdownDescription: "Default description applied to monitors created from this template.",
                Optional: true,
                Computed: true,
            },
            "monitor_type": schema.StringAttribute{
                MarkdownDescription: "What is the type of monitor created from this template?",
                Optional: true,
                Computed: true,
            },
            "monitor_steps": schema.StringAttribute{
                MarkdownDescription: "Monitor steps and criteria copied to monitors created from this template.",
                Computed: true,
            },
            "monitoring_interval": schema.StringAttribute{
                MarkdownDescription: "Default monitoring interval for monitors created from this template. A 5-field cron expression, not a label: \"*/5 * * * *\" is every five minutes.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Default labels applied to monitors created from this template. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "minimum_probe_agreement": schema.NumberAttribute{
                MarkdownDescription: "Default minimum number of probes that must agree on a status before the monitor status changes.",
                Optional: true,
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

func (d *MonitorTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MonitorTemplateDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.TemplateName.IsNull() && !data.TemplateName.IsUnknown() {
        filters["templateName"] = data.TemplateName.ValueString()
        filterNames = append(filterNames, "template_name = "+fmt.Sprintf("%q", data.TemplateName.ValueString()))
    }
    if !data.TemplateDescription.IsNull() && !data.TemplateDescription.IsUnknown() {
        filters["templateDescription"] = data.TemplateDescription.ValueString()
        filterNames = append(filterNames, "template_description = "+fmt.Sprintf("%q", data.TemplateDescription.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.MonitorName.IsNull() && !data.MonitorName.IsUnknown() {
        filters["monitorName"] = data.MonitorName.ValueString()
        filterNames = append(filterNames, "monitor_name = "+fmt.Sprintf("%q", data.MonitorName.ValueString()))
    }
    if !data.MonitorDescription.IsNull() && !data.MonitorDescription.IsUnknown() {
        filters["monitorDescription"] = data.MonitorDescription.ValueString()
        filterNames = append(filterNames, "monitor_description = "+fmt.Sprintf("%q", data.MonitorDescription.ValueString()))
    }
    if !data.MonitorType.IsNull() && !data.MonitorType.IsUnknown() {
        filters["monitorType"] = data.MonitorType.ValueString()
        filterNames = append(filterNames, "monitor_type = "+fmt.Sprintf("%q", data.MonitorType.ValueString()))
    }
    if !data.MonitoringInterval.IsNull() && !data.MonitoringInterval.IsUnknown() {
        filters["monitoringInterval"] = data.MonitoringInterval.ValueString()
        filterNames = append(filterNames, "monitoring_interval = "+fmt.Sprintf("%q", data.MonitoringInterval.ValueString()))
    }
    if !data.MinimumProbeAgreement.IsNull() && !data.MinimumProbeAgreement.IsUnknown() {
        filters["minimumProbeAgreement"] = lookupNumber(data.MinimumProbeAgreement)
        filterNames = append(filterNames, "minimum_probe_agreement = "+data.MinimumProbeAgreement.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the monitor template up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the monitor template up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "templateName": true,
        "templateDescription": true,
        "slug": true,
        "monitorName": true,
        "monitorDescription": true,
        "monitorType": true,
        "monitorSteps": true,
        "monitoringInterval": true,
        "labels": true,
        "customFields": true,
        "minimumProbeAgreement": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/monitor-template/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor_template, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No monitor template found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read monitor_template: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/monitor-template/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list monitor_template, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list monitor_template: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No monitor template matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one monitor template matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for monitor_template.")
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
    if obj, ok := item["templateName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TemplateName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TemplateName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TemplateName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TemplateName = types.StringValue(string(jsonBytes))
        } else {
            data.TemplateName = types.StringNull()
        }
    } else if val, ok := item["templateName"].(string); ok {
        data.TemplateName = types.StringValue(val)
    } else {
        data.TemplateName = types.StringNull()
    }
    if obj, ok := item["templateDescription"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TemplateDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TemplateDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TemplateDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TemplateDescription = types.StringValue(string(jsonBytes))
        } else {
            data.TemplateDescription = types.StringNull()
        }
    } else if val, ok := item["templateDescription"].(string); ok {
        data.TemplateDescription = types.StringValue(val)
    } else {
        data.TemplateDescription = types.StringNull()
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
    if obj, ok := item["monitorName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorName = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorName = types.StringNull()
        }
    } else if val, ok := item["monitorName"].(string); ok {
        data.MonitorName = types.StringValue(val)
    } else {
        data.MonitorName = types.StringNull()
    }
    if obj, ok := item["monitorDescription"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorDescription = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorDescription = types.StringNull()
        }
    } else if val, ok := item["monitorDescription"].(string); ok {
        data.MonitorDescription = types.StringValue(val)
    } else {
        data.MonitorDescription = types.StringNull()
    }
    if obj, ok := item["monitorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorType = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorType = types.StringNull()
        }
    } else if val, ok := item["monitorType"].(string); ok {
        data.MonitorType = types.StringValue(val)
    } else {
        data.MonitorType = types.StringNull()
    }
    if obj, ok := item["monitorSteps"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorSteps = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorSteps = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorSteps = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorSteps = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorSteps = types.StringNull()
        }
    } else if val, ok := item["monitorSteps"].(string); ok {
        data.MonitorSteps = types.StringValue(val)
    } else {
        data.MonitorSteps = types.StringNull()
    }
    if obj, ok := item["monitoringInterval"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitoringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitoringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitoringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.MonitoringInterval = types.StringNull()
        }
    } else if val, ok := item["monitoringInterval"].(string); ok {
        data.MonitoringInterval = types.StringValue(val)
    } else {
        data.MonitoringInterval = types.StringNull()
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
    if obj, ok := item["customFields"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := item["customFields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
    }
    if val, ok := item["minimumProbeAgreement"].(float64); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["minimumProbeAgreement"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
        } else {
            data.MinimumProbeAgreement = types.NumberNull()
        }
    } else {
        data.MinimumProbeAgreement = types.NumberNull()
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
