package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AlertDataDataSource{}

func NewAlertDataDataSource() datasource.DataSource {
    return &AlertDataDataSource{}
}

// AlertDataDataSource defines the data source implementation.
type AlertDataDataSource struct {
    client *Client
}

// AlertDataDataSourceModel describes the data source data model.
type AlertDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    OnCallDutyPolicies types.List `tfsdk:"on_call_duty_policies"`
    Labels types.List `tfsdk:"labels"`
    CurrentAlertStateId types.String `tfsdk:"current_alert_state_id"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    MonitorStatusWhenThisAlertWasCreatedId types.String `tfsdk:"monitor_status_when_this_alert_was_created_id"`
    CustomFields types.String `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfAlertCreation types.Bool `tfsdk:"is_owner_notified_of_alert_creation"`
    RootCause types.String `tfsdk:"root_cause"`
    CreatedStateLog types.String `tfsdk:"created_state_log"`
    CreatedCriteriaId types.String `tfsdk:"created_criteria_id"`
    CreatedByProbeId types.String `tfsdk:"created_by_probe_id"`
    IsCreatedAutomatically types.Bool `tfsdk:"is_created_automatically"`
    RemediationNotes types.String `tfsdk:"remediation_notes"`
    TelemetryQuery types.String `tfsdk:"telemetry_query"`
    AlertNumber types.Number `tfsdk:"alert_number"`
}

func (d *AlertDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_data"
}

func (d *AlertDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "alert_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [Project Owner, Project Admin, Project Member, Edit Alert]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [Project Owner, Project Admin, Project Member, Edit Alert]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policies": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [Project Owner, Project Admin, Project Member, Edit Alert]",
                Computed: true,
                ElementType: types.StringType,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [Project Owner, Project Admin, Project Member, Edit Alert]",
                Computed: true,
                ElementType: types.StringType,
            },
            "current_alert_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitor_status_when_this_alert_was_created_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [Project Owner, Project Admin, Project Member, Edit Alert]",
                Computed: true,
            },
            "is_owner_notified_of_alert_creation": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_state_log": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_criteria_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_probe_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_created_automatically": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "remediation_notes": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [Project Owner, Project Admin, Project Member, Edit Alert]",
                Computed: true,
            },
            "telemetry_query": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [Project Owner, Project Admin, Project Member, Edit Alert]",
                Computed: true,
            },
            "alert_number": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert], Read: [Project Owner, Project Admin, Project Member, Read Alert], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *AlertDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AlertDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "alert" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_data, got error: %s", err))
        return
    }

    var alertDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &alertDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := alertDataResponse["data"].(map[string]interface{}); ok {
        alertDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := alertDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := alertDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := alertDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := alertDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := alertDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := alertDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := alertDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := alertDataResponse["title"].(string); ok {
        data.Title = types.StringValue(val)
    }
    if val, ok := alertDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := alertDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := alertDataResponse["monitor_id"].(string); ok {
        data.MonitorId = types.StringValue(val)
    }
    if val, ok := alertDataResponse["on_call_duty_policies"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.OnCallDutyPolicies = listValue
    }
    if val, ok := alertDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.Labels = listValue
    }
    if val, ok := alertDataResponse["current_alert_state_id"].(string); ok {
        data.CurrentAlertStateId = types.StringValue(val)
    }
    if val, ok := alertDataResponse["alert_severity_id"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    }
    if val, ok := alertDataResponse["monitor_status_when_this_alert_was_created_id"].(string); ok {
        data.MonitorStatusWhenThisAlertWasCreatedId = types.StringValue(val)
    }
    if val, ok := alertDataResponse["custom_fields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    }
    if val, ok := alertDataResponse["is_owner_notified_of_alert_creation"].(bool); ok {
        data.IsOwnerNotifiedOfAlertCreation = types.BoolValue(val)
    }
    if val, ok := alertDataResponse["root_cause"].(string); ok {
        data.RootCause = types.StringValue(val)
    }
    if val, ok := alertDataResponse["created_state_log"].(string); ok {
        data.CreatedStateLog = types.StringValue(val)
    }
    if val, ok := alertDataResponse["created_criteria_id"].(string); ok {
        data.CreatedCriteriaId = types.StringValue(val)
    }
    if val, ok := alertDataResponse["created_by_probe_id"].(string); ok {
        data.CreatedByProbeId = types.StringValue(val)
    }
    if val, ok := alertDataResponse["is_created_automatically"].(bool); ok {
        data.IsCreatedAutomatically = types.BoolValue(val)
    }
    if val, ok := alertDataResponse["remediation_notes"].(string); ok {
        data.RemediationNotes = types.StringValue(val)
    }
    if val, ok := alertDataResponse["telemetry_query"].(string); ok {
        data.TelemetryQuery = types.StringValue(val)
    }
    if val, ok := alertDataResponse["alert_number"].(float64); ok {
        data.AlertNumber = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
