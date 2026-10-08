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
var _ datasource.DataSource = &SecurityEventConnectionDataSource{}

func NewSecurityEventConnectionDataSource() datasource.DataSource {
    return &SecurityEventConnectionDataSource{}
}

// SecurityEventConnectionDataSource defines the data source implementation.
type SecurityEventConnectionDataSource struct {
    client *Client
}

// SecurityEventConnectionDataSourceModel describes the data source data model.
type SecurityEventConnectionDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    ProviderValue types.String `tfsdk:"provider_value"`
    Config types.String `tfsdk:"config"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    PollIntervalInMinutes types.Number `tfsdk:"poll_interval_in_minutes"`
    AlertingOnly types.Bool `tfsdk:"alerting_only"`
    LastSuccessfulPollAt types.String `tfsdk:"last_successful_poll_at"`
    LastEventIngestedAt types.String `tfsdk:"last_event_ingested_at"`
    LastPollResult types.String `tfsdk:"last_poll_result"`
    LastPolledAt types.String `tfsdk:"last_polled_at"`
    Cursor types.String `tfsdk:"cursor"`
    LastError types.String `tfsdk:"last_error"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *SecurityEventConnectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_event_connection"
}

func (d *SecurityEventConnectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Managed connections to SIEM, EDR, cloud security and identity products. Records are polled on an interval and ingested as security events. Look up an existing security event connection by `id`, or by any of its other arguments (`name`, `alerting_only`, `created_by_user_id`, ...): each one set must match, and exactly one security event connection may match them all.",

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
                MarkdownDescription: "ID of the project this connection belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this connection.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What this connection imports and why.",
                Optional: true,
                Computed: true,
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "Which security product this connection polls, e.g. 'microsoft-sentinel' or 'crowdstrike-falcon'. Fixed once created.",
                Optional: true,
                Computed: true,
            },
            "config": schema.StringAttribute{
                MarkdownDescription: "Provider-specific, non-secret settings such as tenant, workspace, region or base URL. Keys are defined by the provider catalog. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this connection is polled on its schedule.",
                Optional: true,
                Computed: true,
            },
            "poll_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often new records are polled, in minutes.",
                Optional: true,
                Computed: true,
            },
            "alerting_only": schema.BoolAttribute{
                MarkdownDescription: "For providers that distinguish alerting from non-alerting records: import only the alerting ones.",
                Optional: true,
                Computed: true,
            },
            "last_successful_poll_at": schema.StringAttribute{
                MarkdownDescription: "When a complete poll last finished successfully, including an empty result.",
                Computed: true,
            },
            "last_event_ingested_at": schema.StringAttribute{
                MarkdownDescription: "When a new security event was last imported from this connection.",
                Computed: true,
            },
            "last_poll_result": schema.StringAttribute{
                MarkdownDescription: "The latest scheduled or on-demand poll result with counts and warnings. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "last_polled_at": schema.StringAttribute{
                MarkdownDescription: "When this connection was last polled. Null means it has never run.",
                Computed: true,
            },
            "cursor": schema.StringAttribute{
                MarkdownDescription: "Poll cursor: the end of the last completely processed creation-time window, as an ISO string.",
                Optional: true,
                Computed: true,
            },
            "last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent poll error with credentials redacted, if any. Cleared on the next successful poll.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *SecurityEventConnectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurityEventConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SecurityEventConnectionDataSourceModel

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
    if !data.ProviderValue.IsNull() && !data.ProviderValue.IsUnknown() {
        filters["provider"] = data.ProviderValue.ValueString()
        filterNames = append(filterNames, "provider_value = "+fmt.Sprintf("%q", data.ProviderValue.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.PollIntervalInMinutes.IsNull() && !data.PollIntervalInMinutes.IsUnknown() {
        filters["pollIntervalInMinutes"] = lookupNumber(data.PollIntervalInMinutes)
        filterNames = append(filterNames, "poll_interval_in_minutes = "+data.PollIntervalInMinutes.ValueBigFloat().String())
    }
    if !data.AlertingOnly.IsNull() && !data.AlertingOnly.IsUnknown() {
        filters["alertingOnly"] = data.AlertingOnly.ValueBool()
        filterNames = append(filterNames, "alerting_only = "+fmt.Sprintf("%t", data.AlertingOnly.ValueBool()))
    }
    if !data.Cursor.IsNull() && !data.Cursor.IsUnknown() {
        filters["cursor"] = data.Cursor.ValueString()
        filterNames = append(filterNames, "cursor = "+fmt.Sprintf("%q", data.Cursor.ValueString()))
    }
    if !data.LastError.IsNull() && !data.LastError.IsUnknown() {
        filters["lastError"] = data.LastError.ValueString()
        filterNames = append(filterNames, "last_error = "+fmt.Sprintf("%q", data.LastError.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the security event connection up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the security event connection up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "provider": true,
        "config": true,
        "isEnabled": true,
        "pollIntervalInMinutes": true,
        "alertingOnly": true,
        "lastSuccessfulPollAt": true,
        "lastEventIngestedAt": true,
        "lastPollResult": true,
        "lastPolledAt": true,
        "cursor": true,
        "lastError": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-event-connection/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_event_connection, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security event connection found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_event_connection: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-event-connection/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list security_event_connection, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list security_event_connection: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security event connection matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security event connection matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for security_event_connection.")
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
    if obj, ok := item["config"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Config = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Config = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Config = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Config = types.StringValue(string(jsonBytes))
        } else {
            data.Config = types.StringNull()
        }
    } else if val, ok := item["config"].(string); ok {
        data.Config = types.StringValue(val)
    } else {
        data.Config = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["pollIntervalInMinutes"].(float64); ok {
        data.PollIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["pollIntervalInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PollIntervalInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.PollIntervalInMinutes = types.NumberNull()
        }
    } else {
        data.PollIntervalInMinutes = types.NumberNull()
    }
    if val, ok := item["alertingOnly"].(bool); ok {
        data.AlertingOnly = types.BoolValue(val)
    } else {
        data.AlertingOnly = types.BoolNull()
    }
    if obj, ok := item["lastSuccessfulPollAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSuccessfulPollAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSuccessfulPollAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSuccessfulPollAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSuccessfulPollAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSuccessfulPollAt = types.StringNull()
        }
    } else if val, ok := item["lastSuccessfulPollAt"].(string); ok {
        data.LastSuccessfulPollAt = types.StringValue(val)
    } else {
        data.LastSuccessfulPollAt = types.StringNull()
    }
    if obj, ok := item["lastEventIngestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastEventIngestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastEventIngestedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastEventIngestedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastEventIngestedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastEventIngestedAt = types.StringNull()
        }
    } else if val, ok := item["lastEventIngestedAt"].(string); ok {
        data.LastEventIngestedAt = types.StringValue(val)
    } else {
        data.LastEventIngestedAt = types.StringNull()
    }
    if obj, ok := item["lastPollResult"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastPollResult = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastPollResult = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastPollResult = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastPollResult = types.StringValue(string(jsonBytes))
        } else {
            data.LastPollResult = types.StringNull()
        }
    } else if val, ok := item["lastPollResult"].(string); ok {
        data.LastPollResult = types.StringValue(val)
    } else {
        data.LastPollResult = types.StringNull()
    }
    if obj, ok := item["lastPolledAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastPolledAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastPolledAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastPolledAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastPolledAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastPolledAt = types.StringNull()
        }
    } else if val, ok := item["lastPolledAt"].(string); ok {
        data.LastPolledAt = types.StringValue(val)
    } else {
        data.LastPolledAt = types.StringNull()
    }
    if obj, ok := item["cursor"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Cursor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Cursor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Cursor = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Cursor = types.StringValue(string(jsonBytes))
        } else {
            data.Cursor = types.StringNull()
        }
    } else if val, ok := item["cursor"].(string); ok {
        data.Cursor = types.StringValue(val)
    } else {
        data.Cursor = types.StringNull()
    }
    if obj, ok := item["lastError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastError = types.StringValue(string(jsonBytes))
        } else {
            data.LastError = types.StringNull()
        }
    } else if val, ok := item["lastError"].(string); ok {
        data.LastError = types.StringValue(val)
    } else {
        data.LastError = types.StringNull()
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
