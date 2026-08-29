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
var _ datasource.DataSource = &IncidentSlaRuleDataSource{}

func NewIncidentSlaRuleDataSource() datasource.DataSource {
    return &IncidentSlaRuleDataSource{}
}

// IncidentSlaRuleDataSource defines the data source implementation.
type IncidentSlaRuleDataSource struct {
    client *Client
}

// IncidentSlaRuleDataSourceModel describes the data source data model.
type IncidentSlaRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Order types.Number `tfsdk:"order"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    ResponseTimeInMinutes types.Number `tfsdk:"response_time_in_minutes"`
    ResolutionTimeInMinutes types.Number `tfsdk:"resolution_time_in_minutes"`
    AtRiskThresholdInPercentage types.Number `tfsdk:"at_risk_threshold_in_percentage"`
    InternalNoteReminderIntervalInMinutes types.Number `tfsdk:"internal_note_reminder_interval_in_minutes"`
    PublicNoteReminderIntervalInMinutes types.Number `tfsdk:"public_note_reminder_interval_in_minutes"`
    InternalNoteReminderTemplate types.String `tfsdk:"internal_note_reminder_template"`
    PublicNoteReminderTemplate types.String `tfsdk:"public_note_reminder_template"`
    Monitors types.Set `tfsdk:"monitors"`
    IncidentSeverities types.Set `tfsdk:"incident_severities"`
    IncidentLabels types.Set `tfsdk:"incident_labels"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    IncidentTitlePattern types.String `tfsdk:"incident_title_pattern"`
    IncidentDescriptionPattern types.String `tfsdk:"incident_description_pattern"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentSlaRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_sla_rule"
}

func (d *IncidentSlaRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure SLA rules to define response and resolution time targets for incidents Look up an existing incident_sla_rule by `id` or by `name`.",

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
                MarkdownDescription: "Description of this SLA rule.",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Order/priority of this rule. Rules are evaluated in order (lowest first). First matching rule wins..",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this SLA rule is enabled.",
                Computed: true,
            },
            "response_time_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Target response time in minutes. This is the maximum time allowed before the incident must be acknowledged..",
                Computed: true,
            },
            "resolution_time_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Target resolution time in minutes. This is the maximum time allowed before the incident must be resolved..",
                Computed: true,
            },
            "at_risk_threshold_in_percentage": schema.NumberAttribute{
                MarkdownDescription: "Percentage of the deadline at which the SLA status changes to At Risk. For example, 80 means the status becomes At Risk when 80% of the time has elapsed..",
                Computed: true,
            },
            "internal_note_reminder_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often (in minutes) to automatically post internal notes to unresolved incidents. Internal notes are only visible to your team. For example, set to 30 to remind your team every 30 minutes to provide an update. Leave empty to disable..",
                Computed: true,
            },
            "public_note_reminder_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often (in minutes) to automatically post public notes to unresolved incidents. Public notes are visible to external stakeholders on your status page. For example, set to 60 to post a status update every hour. Leave empty to disable..",
                Computed: true,
            },
            "internal_note_reminder_template": schema.StringAttribute{
                MarkdownDescription: "The content of the automatic internal note posted to your team. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used..",
                Computed: true,
            },
            "public_note_reminder_template": schema.StringAttribute{
                MarkdownDescription: "The content of the automatic public note shown on your status page. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used..",
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only apply this SLA rule to incidents affecting these monitors. Leave empty to match incidents from any monitor..",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_severities": schema.SetAttribute{
                MarkdownDescription: "Only apply this SLA rule to incidents with these severities. Leave empty to match incidents of any severity..",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_labels": schema.SetAttribute{
                MarkdownDescription: "Only apply this SLA rule to incidents that have at least one of these labels. Leave empty to match incidents regardless of labels..",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only apply this SLA rule to incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels..",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'..",
                Computed: true,
            },
            "incident_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regular expression pattern to match incident descriptions. Leave empty to match any description..",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IncidentSlaRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentSlaRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentSlaRuleDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a incident_sla_rule.",
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
        "order": true,
        "isEnabled": true,
        "responseTimeInMinutes": true,
        "resolutionTimeInMinutes": true,
        "atRiskThresholdInPercentage": true,
        "internalNoteReminderIntervalInMinutes": true,
        "publicNoteReminderIntervalInMinutes": true,
        "internalNoteReminderTemplate": true,
        "publicNoteReminderTemplate": true,
        "monitors": true,
        "incidentSeverities": true,
        "incidentLabels": true,
        "monitorLabels": true,
        "incidentTitlePattern": true,
        "incidentDescriptionPattern": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-sla-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_sla_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident_sla_rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_sla_rule: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-sla-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_sla_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_sla_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident_sla_rule found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident_sla_rule matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_sla_rule.")
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
    if val, ok := item["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["order"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Order = types.NumberValue(big.NewFloat(val))
        } else {
            data.Order = types.NumberNull()
        }
    } else {
        data.Order = types.NumberNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["responseTimeInMinutes"].(float64); ok {
        data.ResponseTimeInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["responseTimeInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ResponseTimeInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ResponseTimeInMinutes = types.NumberNull()
        }
    } else {
        data.ResponseTimeInMinutes = types.NumberNull()
    }
    if val, ok := item["resolutionTimeInMinutes"].(float64); ok {
        data.ResolutionTimeInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["resolutionTimeInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ResolutionTimeInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.ResolutionTimeInMinutes = types.NumberNull()
        }
    } else {
        data.ResolutionTimeInMinutes = types.NumberNull()
    }
    if val, ok := item["atRiskThresholdInPercentage"].(float64); ok {
        data.AtRiskThresholdInPercentage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["atRiskThresholdInPercentage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AtRiskThresholdInPercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.AtRiskThresholdInPercentage = types.NumberNull()
        }
    } else {
        data.AtRiskThresholdInPercentage = types.NumberNull()
    }
    if val, ok := item["internalNoteReminderIntervalInMinutes"].(float64); ok {
        data.InternalNoteReminderIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["internalNoteReminderIntervalInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.InternalNoteReminderIntervalInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.InternalNoteReminderIntervalInMinutes = types.NumberNull()
        }
    } else {
        data.InternalNoteReminderIntervalInMinutes = types.NumberNull()
    }
    if val, ok := item["publicNoteReminderIntervalInMinutes"].(float64); ok {
        data.PublicNoteReminderIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["publicNoteReminderIntervalInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PublicNoteReminderIntervalInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.PublicNoteReminderIntervalInMinutes = types.NumberNull()
        }
    } else {
        data.PublicNoteReminderIntervalInMinutes = types.NumberNull()
    }
    if obj, ok := item["internalNoteReminderTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InternalNoteReminderTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.InternalNoteReminderTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.InternalNoteReminderTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.InternalNoteReminderTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.InternalNoteReminderTemplate = types.StringNull()
        }
    } else if val, ok := item["internalNoteReminderTemplate"].(string); ok {
        data.InternalNoteReminderTemplate = types.StringValue(val)
    } else {
        data.InternalNoteReminderTemplate = types.StringNull()
    }
    if obj, ok := item["publicNoteReminderTemplate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PublicNoteReminderTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PublicNoteReminderTemplate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PublicNoteReminderTemplate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PublicNoteReminderTemplate = types.StringValue(string(jsonBytes))
        } else {
            data.PublicNoteReminderTemplate = types.StringNull()
        }
    } else if val, ok := item["publicNoteReminderTemplate"].(string); ok {
        data.PublicNoteReminderTemplate = types.StringValue(val)
    } else {
        data.PublicNoteReminderTemplate = types.StringNull()
    }
    if val, ok := item["monitors"].([]interface{}); ok {
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
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Monitors = types.SetNull(types.StringType)
    }
    if val, ok := item["incidentSeverities"].([]interface{}); ok {
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
        data.IncidentSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentSeverities = types.SetNull(types.StringType)
    }
    if val, ok := item["incidentLabels"].([]interface{}); ok {
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
        data.IncidentLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentLabels = types.SetNull(types.StringType)
    }
    if val, ok := item["monitorLabels"].([]interface{}); ok {
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
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.MonitorLabels = types.SetNull(types.StringType)
    }
    if obj, ok := item["incidentTitlePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentTitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentTitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentTitlePattern = types.StringNull()
        }
    } else if val, ok := item["incidentTitlePattern"].(string); ok {
        data.IncidentTitlePattern = types.StringValue(val)
    } else {
        data.IncidentTitlePattern = types.StringNull()
    }
    if obj, ok := item["incidentDescriptionPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentDescriptionPattern = types.StringNull()
        }
    } else if val, ok := item["incidentDescriptionPattern"].(string); ok {
        data.IncidentDescriptionPattern = types.StringValue(val)
    } else {
        data.IncidentDescriptionPattern = types.StringNull()
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
