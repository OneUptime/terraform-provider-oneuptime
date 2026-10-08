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
var _ datasource.DataSource = &LlmCostBudgetDataSource{}

func NewLlmCostBudgetDataSource() datasource.DataSource {
    return &LlmCostBudgetDataSource{}
}

// LlmCostBudgetDataSource defines the data source implementation.
type LlmCostBudgetDataSource struct {
    client *Client
}

// LlmCostBudgetDataSourceModel describes the data source data model.
type LlmCostBudgetDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    DailyBudgetInUsd types.Number `tfsdk:"daily_budget_in_usd"`
    ServiceId types.String `tfsdk:"service_id"`
    LlmSystem types.String `tfsdk:"llm_system"`
    LlmModel types.String `tfsdk:"llm_model"`
    CurrentDaySpendInUsd types.Number `tfsdk:"current_day_spend_in_usd"`
    SpendLastEvaluatedAt types.String `tfsdk:"spend_last_evaluated_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *LlmCostBudgetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_llm_cost_budget"
}

func (d *LlmCostBudgetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Daily USD spend budgets for LLM / GenAI telemetry. A worker sums the day's LLM span cost and publishes it as the oneuptime.llm.budget.* metrics, so Metrics monitors, dashboards and anomaly detection can act on spend. Look up an existing llm cost budget by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `current_day_spend_in_usd`, ...): each one set must match, and exactly one llm cost budget may match them all.",

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
                MarkdownDescription: "ID of the project this LLM cost budget belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this budget.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this budget covers.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this budget is evaluated.",
                Optional: true,
                Computed: true,
            },
            "daily_budget_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Daily LLM spend budget in USD, evaluated over the UTC day. Spend and percent-used are published as metrics for monitors to alert on.",
                Optional: true,
                Computed: true,
            },
            "service_id": schema.StringAttribute{
                MarkdownDescription: "Optional telemetry service ID scoping this budget. Null means the budget is project-wide. The ID of a `oneuptime_service`.",
                Optional: true,
                Computed: true,
            },
            "llm_system": schema.StringAttribute{
                MarkdownDescription: "Optional provider filter, e.g. openai, anthropic, aws.bedrock (matches the span's gen_ai provider). Leave empty to count every provider.",
                Optional: true,
                Computed: true,
            },
            "llm_model": schema.StringAttribute{
                MarkdownDescription: "Optional model filter, e.g. gpt-4o (matches the span's requested model exactly). Leave empty to count every model.",
                Optional: true,
                Computed: true,
            },
            "current_day_spend_in_usd": schema.NumberAttribute{
                MarkdownDescription: "LLM spend accrued so far in the current UTC day, in USD. Computed by the worker.",
                Optional: true,
                Computed: true,
            },
            "spend_last_evaluated_at": schema.StringAttribute{
                MarkdownDescription: "The last time the worker evaluated this budget. Computed by the worker.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this budget. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *LlmCostBudgetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LlmCostBudgetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LlmCostBudgetDataSourceModel

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
    if !data.DailyBudgetInUsd.IsNull() && !data.DailyBudgetInUsd.IsUnknown() {
        filters["dailyBudgetInUSD"] = lookupNumber(data.DailyBudgetInUsd)
        filterNames = append(filterNames, "daily_budget_in_usd = "+data.DailyBudgetInUsd.ValueBigFloat().String())
    }
    if !data.ServiceId.IsNull() && !data.ServiceId.IsUnknown() {
        filters["serviceId"] = data.ServiceId.ValueString()
        filterNames = append(filterNames, "service_id = "+fmt.Sprintf("%q", data.ServiceId.ValueString()))
    }
    if !data.LlmSystem.IsNull() && !data.LlmSystem.IsUnknown() {
        filters["llmSystem"] = data.LlmSystem.ValueString()
        filterNames = append(filterNames, "llm_system = "+fmt.Sprintf("%q", data.LlmSystem.ValueString()))
    }
    if !data.LlmModel.IsNull() && !data.LlmModel.IsUnknown() {
        filters["llmModel"] = data.LlmModel.ValueString()
        filterNames = append(filterNames, "llm_model = "+fmt.Sprintf("%q", data.LlmModel.ValueString()))
    }
    if !data.CurrentDaySpendInUsd.IsNull() && !data.CurrentDaySpendInUsd.IsUnknown() {
        filters["currentDaySpendInUSD"] = lookupNumber(data.CurrentDaySpendInUsd)
        filterNames = append(filterNames, "current_day_spend_in_usd = "+data.CurrentDaySpendInUsd.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the llm cost budget up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the llm cost budget up by.",
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
        "dailyBudgetInUSD": true,
        "serviceId": true,
        "llmSystem": true,
        "llmModel": true,
        "currentDaySpendInUSD": true,
        "spendLastEvaluatedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/llm-cost-budget/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read llm_cost_budget, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No llm cost budget found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read llm_cost_budget: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/llm-cost-budget/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list llm_cost_budget, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list llm_cost_budget: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No llm cost budget matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one llm cost budget matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for llm_cost_budget.")
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
    if val, ok := item["dailyBudgetInUSD"].(float64); ok {
        data.DailyBudgetInUsd = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["dailyBudgetInUSD"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DailyBudgetInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.DailyBudgetInUsd = types.NumberNull()
        }
    } else {
        data.DailyBudgetInUsd = types.NumberNull()
    }
    if obj, ok := item["serviceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := item["serviceId"].(string); ok {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := item["llmSystem"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmSystem = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmSystem = types.StringValue(string(jsonBytes))
        } else {
            data.LlmSystem = types.StringNull()
        }
    } else if val, ok := item["llmSystem"].(string); ok {
        data.LlmSystem = types.StringValue(val)
    } else {
        data.LlmSystem = types.StringNull()
    }
    if obj, ok := item["llmModel"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmModel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmModel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmModel = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmModel = types.StringValue(string(jsonBytes))
        } else {
            data.LlmModel = types.StringNull()
        }
    } else if val, ok := item["llmModel"].(string); ok {
        data.LlmModel = types.StringValue(val)
    } else {
        data.LlmModel = types.StringNull()
    }
    if val, ok := item["currentDaySpendInUSD"].(float64); ok {
        data.CurrentDaySpendInUsd = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["currentDaySpendInUSD"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CurrentDaySpendInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.CurrentDaySpendInUsd = types.NumberNull()
        }
    } else {
        data.CurrentDaySpendInUsd = types.NumberNull()
    }
    if obj, ok := item["spendLastEvaluatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpendLastEvaluatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SpendLastEvaluatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SpendLastEvaluatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SpendLastEvaluatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.SpendLastEvaluatedAt = types.StringNull()
        }
    } else if val, ok := item["spendLastEvaluatedAt"].(string); ok {
        data.SpendLastEvaluatedAt = types.StringValue(val)
    } else {
        data.SpendLastEvaluatedAt = types.StringNull()
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
