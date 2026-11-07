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
var _ datasource.DataSource = &CloudBudgetDataSource{}

func NewCloudBudgetDataSource() datasource.DataSource {
    return &CloudBudgetDataSource{}
}

// CloudBudgetDataSource defines the data source implementation.
type CloudBudgetDataSource struct {
    client *Client
}

// CloudBudgetDataSourceModel describes the data source data model.
type CloudBudgetDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    CloudAccountId types.String `tfsdk:"cloud_account_id"`
    MonthlyAmount types.Number `tfsdk:"monthly_amount"`
    Currency types.String `tfsdk:"currency"`
    AlertThresholds types.String `tfsdk:"alert_thresholds"`
    AlertOnProjection types.Bool `tfsdk:"alert_on_projection"`
    AlertOnAnomaly types.Bool `tfsdk:"alert_on_anomaly"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    SpentThisMonth types.Number `tfsdk:"spent_this_month"`
    ProjectedThisMonth types.Number `tfsdk:"projected_this_month"`
    LastCheckedAt types.String `tfsdk:"last_checked_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *CloudBudgetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_cloud_budget"
}

func (d *CloudBudgetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Monthly budgets for what your connected cloud accounts cost. OneUptime alerts the first time in a month the spend passes each threshold, and when the month is on course to pass the budget. Look up an existing cloud budget by `id`, or by any of its other arguments (`name`, `alert_on_anomaly`, `alert_on_projection`, ...): each one set must match, and exactly one cloud budget may match them all.",

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
                MarkdownDescription: "ID of the project this budget belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "What the budget is called.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the budget is for, in your team's words.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether the budget is watched and raises alerts.",
                Optional: true,
                Computed: true,
            },
            "cloud_account_id": schema.StringAttribute{
                MarkdownDescription: "ID of the account the budget is for. Empty: every account of the project. The ID of a `oneuptime_cloud_account`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_cloud_account"),
            },
            "monthly_amount": schema.NumberAttribute{
                MarkdownDescription: "What the accounts may cost in a calendar month (UTC), in the budget's currency.",
                Optional: true,
                Computed: true,
            },
            "currency": schema.StringAttribute{
                MarkdownDescription: "The currency of the amount, as the clouds bill the accounts: an account billed in another is not counted.",
                Optional: true,
                Computed: true,
            },
            "alert_thresholds": schema.StringAttribute{
                MarkdownDescription: "The shares of the budget, in percent, an alert is raised at, comma separated: the first time in a month the spend passes each.",
                Optional: true,
                Computed: true,
            },
            "alert_on_projection": schema.BoolAttribute{
                MarkdownDescription: "Also alert when the month is on course to pass the budget, at the pace of the last full week: before the money is spent.",
                Optional: true,
                Computed: true,
            },
            "alert_on_anomaly": schema.BoolAttribute{
                MarkdownDescription: "Also alert when a day costs far more than usual: a service at twice its usual day, or the whole account at one and a half times.",
                Optional: true,
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the severity of the alerts the budget raises. Empty: the project's least severe. The ID of a `oneuptime_alert_severity`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_alert_severity"),
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "The on-call policies the budget's alerts page. IDs of `oneuptime_on_call_policy` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "spent_this_month": schema.NumberAttribute{
                MarkdownDescription: "What the accounts cost this month so far, when last checked.",
                Optional: true,
                Computed: true,
            },
            "projected_this_month": schema.NumberAttribute{
                MarkdownDescription: "What the accounts are on course to cost this month, at the pace of the last full week, when last checked.",
                Optional: true,
                Computed: true,
            },
            "last_checked_at": schema.StringAttribute{
                MarkdownDescription: "When the budget was last compared with what the accounts cost.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this budget. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *CloudBudgetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudBudgetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CloudBudgetDataSourceModel

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
    if !data.CloudAccountId.IsNull() && !data.CloudAccountId.IsUnknown() {
        filters["cloudAccountId"] = data.CloudAccountId.ValueString()
        filterNames = append(filterNames, "cloud_account_id = "+fmt.Sprintf("%q", data.CloudAccountId.ValueString()))
    }
    if !data.MonthlyAmount.IsNull() && !data.MonthlyAmount.IsUnknown() {
        filters["monthlyAmount"] = lookupNumber(data.MonthlyAmount)
        filterNames = append(filterNames, "monthly_amount = "+data.MonthlyAmount.ValueBigFloat().String())
    }
    if !data.Currency.IsNull() && !data.Currency.IsUnknown() {
        filters["currency"] = data.Currency.ValueString()
        filterNames = append(filterNames, "currency = "+fmt.Sprintf("%q", data.Currency.ValueString()))
    }
    if !data.AlertThresholds.IsNull() && !data.AlertThresholds.IsUnknown() {
        filters["alertThresholds"] = data.AlertThresholds.ValueString()
        filterNames = append(filterNames, "alert_thresholds = "+fmt.Sprintf("%q", data.AlertThresholds.ValueString()))
    }
    if !data.AlertOnProjection.IsNull() && !data.AlertOnProjection.IsUnknown() {
        filters["alertOnProjection"] = data.AlertOnProjection.ValueBool()
        filterNames = append(filterNames, "alert_on_projection = "+fmt.Sprintf("%t", data.AlertOnProjection.ValueBool()))
    }
    if !data.AlertOnAnomaly.IsNull() && !data.AlertOnAnomaly.IsUnknown() {
        filters["alertOnAnomaly"] = data.AlertOnAnomaly.ValueBool()
        filterNames = append(filterNames, "alert_on_anomaly = "+fmt.Sprintf("%t", data.AlertOnAnomaly.ValueBool()))
    }
    if !data.AlertSeverityId.IsNull() && !data.AlertSeverityId.IsUnknown() {
        filters["alertSeverityId"] = data.AlertSeverityId.ValueString()
        filterNames = append(filterNames, "alert_severity_id = "+fmt.Sprintf("%q", data.AlertSeverityId.ValueString()))
    }
    if !data.SpentThisMonth.IsNull() && !data.SpentThisMonth.IsUnknown() {
        filters["spentThisMonth"] = lookupNumber(data.SpentThisMonth)
        filterNames = append(filterNames, "spent_this_month = "+data.SpentThisMonth.ValueBigFloat().String())
    }
    if !data.ProjectedThisMonth.IsNull() && !data.ProjectedThisMonth.IsUnknown() {
        filters["projectedThisMonth"] = lookupNumber(data.ProjectedThisMonth)
        filterNames = append(filterNames, "projected_this_month = "+data.ProjectedThisMonth.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the cloud budget up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the cloud budget up by.",
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
        "cloudAccountId": true,
        "monthlyAmount": true,
        "currency": true,
        "alertThresholds": true,
        "alertOnProjection": true,
        "alertOnAnomaly": true,
        "alertSeverityId": true,
        "onCallDutyPolicies": true,
        "spentThisMonth": true,
        "projectedThisMonth": true,
        "lastCheckedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/cloud-budget/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cloud_budget, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud budget found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read cloud_budget: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/cloud-budget/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list cloud_budget, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list cloud_budget: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud budget matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one cloud budget matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for cloud_budget.")
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
    data.MonthlyAmount = numberFromAPI(data.MonthlyAmount, item["monthlyAmount"])
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
    if obj, ok := item["alertThresholds"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertThresholds = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertThresholds = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertThresholds = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertThresholds = types.StringValue(string(jsonBytes))
        } else {
            data.AlertThresholds = types.StringNull()
        }
    } else if val, ok := item["alertThresholds"].(string); ok {
        data.AlertThresholds = types.StringValue(val)
    } else {
        data.AlertThresholds = types.StringNull()
    }
    if val, ok := item["alertOnProjection"].(bool); ok {
        data.AlertOnProjection = types.BoolValue(val)
    } else {
        data.AlertOnProjection = types.BoolNull()
    }
    if val, ok := item["alertOnAnomaly"].(bool); ok {
        data.AlertOnAnomaly = types.BoolValue(val)
    } else {
        data.AlertOnAnomaly = types.BoolNull()
    }
    if obj, ok := item["alertSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := item["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if val, ok := item["onCallDutyPolicies"].([]interface{}); ok {
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
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        data.OnCallDutyPolicies = types.SetNull(types.StringType)
    }
    data.SpentThisMonth = numberFromAPI(data.SpentThisMonth, item["spentThisMonth"])
    data.ProjectedThisMonth = numberFromAPI(data.ProjectedThisMonth, item["projectedThisMonth"])
    if obj, ok := item["lastCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastCheckedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastCheckedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastCheckedAt = types.StringNull()
        }
    } else if val, ok := item["lastCheckedAt"].(string); ok {
        data.LastCheckedAt = types.StringValue(val)
    } else {
        data.LastCheckedAt = types.StringNull()
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
