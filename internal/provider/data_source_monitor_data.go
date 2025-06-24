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
var _ datasource.DataSource = &MonitorDataDataSource{}

func NewMonitorDataDataSource() datasource.DataSource {
    return &MonitorDataDataSource{}
}

// MonitorDataDataSource defines the data source implementation.
type MonitorDataDataSource struct {
    client *Client
}

// MonitorDataDataSourceModel describes the data source data model.
type MonitorDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Labels types.List `tfsdk:"labels"`
    MonitorType types.String `tfsdk:"monitor_type"`
    CurrentMonitorStatusId types.String `tfsdk:"current_monitor_status_id"`
    MonitorSteps types.String `tfsdk:"monitor_steps"`
    MonitoringInterval types.String `tfsdk:"monitoring_interval"`
    CustomFields types.String `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    DisableActiveMonitoring types.Bool `tfsdk:"disable_active_monitoring"`
    IncomingRequestMonitorHeartbeatCheckedAt types.String `tfsdk:"incoming_request_monitor_heartbeat_checked_at"`
    TelemetryMonitorNextMonitorAt types.String `tfsdk:"telemetry_monitor_next_monitor_at"`
    TelemetryMonitorLastMonitorAt types.String `tfsdk:"telemetry_monitor_last_monitor_at"`
    DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent types.Bool `tfsdk:"disable_active_monitoring_because_of_scheduled_maintenance_event"`
    DisableActiveMonitoringBecauseOfManualIncident types.Bool `tfsdk:"disable_active_monitoring_because_of_manual_incident"`
    ServerMonitorRequestReceivedAt types.String `tfsdk:"server_monitor_request_received_at"`
    ServerMonitorSecretKey types.String `tfsdk:"server_monitor_secret_key"`
    IncomingRequestSecretKey types.String `tfsdk:"incoming_request_secret_key"`
    IncomingMonitorRequest types.String `tfsdk:"incoming_monitor_request"`
    ServerMonitorResponse types.String `tfsdk:"server_monitor_response"`
    IsAllProbesDisconnectedFromThisMonitor types.Bool `tfsdk:"is_all_probes_disconnected_from_this_monitor"`
    IsNoProbeEnabledOnThisMonitor types.Bool `tfsdk:"is_no_probe_enabled_on_this_monitor"`
}

func (d *MonitorDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_monitor_data"
}

func (d *MonitorDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "monitor_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [Project Owner, Project Admin, Project Member, Edit Monitor]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [Project Owner, Project Admin, Project Member, Edit Monitor]",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "current_monitor_status_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitor_steps": schema.StringAttribute{
                MarkdownDescription: "MonitorSteps object",
                Computed: true,
            },
            "monitoring_interval": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [Project Owner, Project Admin, Project Member, Edit Monitor]",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [Project Owner, Project Admin, Project Member, Edit Monitor]",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "disable_active_monitoring": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [Project Owner, Project Admin, Project Member, Create Monitor]",
                Computed: true,
            },
            "incoming_request_monitor_heartbeat_checked_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "telemetry_monitor_next_monitor_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "telemetry_monitor_last_monitor_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "disable_active_monitoring_because_of_scheduled_maintenance_event": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "disable_active_monitoring_because_of_manual_incident": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "server_monitor_request_received_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "server_monitor_secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incoming_request_secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incoming_monitor_request": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "server_monitor_response": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_all_probes_disconnected_from_this_monitor": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_no_probe_enabled_on_this_monitor": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *MonitorDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MonitorDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "monitor" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor_data, got error: %s", err))
        return
    }

    var monitorDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &monitorDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := monitorDataResponse["data"].(map[string]interface{}); ok {
        monitorDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := monitorDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := monitorDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := monitorDataResponse["monitor_type"].(string); ok {
        data.MonitorType = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["current_monitor_status_id"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["monitor_steps"].(string); ok {
        data.MonitorSteps = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["monitoring_interval"].(string); ok {
        data.MonitoringInterval = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["custom_fields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["is_owner_notified_of_resource_creation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := monitorDataResponse["disable_active_monitoring"].(bool); ok {
        data.DisableActiveMonitoring = types.BoolValue(val)
    }
    if val, ok := monitorDataResponse["incoming_request_monitor_heartbeat_checked_at"].(string); ok {
        data.IncomingRequestMonitorHeartbeatCheckedAt = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["telemetry_monitor_next_monitor_at"].(string); ok {
        data.TelemetryMonitorNextMonitorAt = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["telemetry_monitor_last_monitor_at"].(string); ok {
        data.TelemetryMonitorLastMonitorAt = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["disable_active_monitoring_because_of_scheduled_maintenance_event"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolValue(val)
    }
    if val, ok := monitorDataResponse["disable_active_monitoring_because_of_manual_incident"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolValue(val)
    }
    if val, ok := monitorDataResponse["server_monitor_request_received_at"].(string); ok {
        data.ServerMonitorRequestReceivedAt = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["server_monitor_secret_key"].(string); ok {
        data.ServerMonitorSecretKey = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["incoming_request_secret_key"].(string); ok {
        data.IncomingRequestSecretKey = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["incoming_monitor_request"].(string); ok {
        data.IncomingMonitorRequest = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["server_monitor_response"].(string); ok {
        data.ServerMonitorResponse = types.StringValue(val)
    }
    if val, ok := monitorDataResponse["is_all_probes_disconnected_from_this_monitor"].(bool); ok {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolValue(val)
    }
    if val, ok := monitorDataResponse["is_no_probe_enabled_on_this_monitor"].(bool); ok {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
