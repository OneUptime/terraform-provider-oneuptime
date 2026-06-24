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
var _ datasource.DataSource = &IoTFleetDataDataSource{}

func NewIoTFleetDataDataSource() datasource.DataSource {
    return &IoTFleetDataDataSource{}
}

// IoTFleetDataDataSource defines the data source implementation.
type IoTFleetDataDataSource struct {
    client *Client
}

// IoTFleetDataDataSourceModel describes the data source data model.
type IoTFleetDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    DeviceCount types.Number `tfsdk:"device_count"`
    OnlineDeviceCount types.Number `tfsdk:"online_device_count"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
}

func (d *IoTFleetDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_io_t_fleet_data"
}

func (d *IoTFleetDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "io_t_fleet_data data source",

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
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read IoT Fleet], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this IoT fleet. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create IoT Fleet], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read IoT Fleet], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit IoT Fleet]",
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Connection status of the OTel Collector agent (connected or disconnected). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read IoT Fleet], Update: [Project Owner, Project Admin, Edit IoT Fleet]",
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the IoT agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read IoT Fleet], Update: [Project Owner, Project Admin, Edit IoT Fleet]",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "device_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of devices in this fleet. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read IoT Fleet], Update: [Project Owner, Project Admin, Edit IoT Fleet]",
                Computed: true,
            },
            "online_device_count": schema.NumberAttribute{
                MarkdownDescription: "Cached count of devices currently online (iot_device_up == 1) in this fleet. Rendered as 'Devices X/Y online' next to deviceCount.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read IoT Fleet], Update: [Project Owner, Project Admin, Edit IoT Fleet]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this IoT fleet archived? Archived IoT fleets are hidden from lists but keep collecting telemetry.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create IoT Fleet], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read IoT Fleet], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit IoT Fleet]",
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create IoT Fleet], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read IoT Fleet], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit IoT Fleet]",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this IoT fleet. Leave blank to use the project-wide default.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create IoT Fleet], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read IoT Fleet], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit IoT Fleet]",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this IoT fleet (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the IoT fleet default, then the project's retention settings.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create IoT Fleet], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read IoT Fleet], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit IoT Fleet]",
                Computed: true,
            },
        },
    }
}

func (d *IoTFleetDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IoTFleetDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IoTFleetDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "iot-fleet" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read io_t_fleet_data, got error: %s", err))
        return
    }

    var ioTFleetDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &ioTFleetDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse io_t_fleet_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := ioTFleetDataResponse["data"].(map[string]interface{}); ok {
        ioTFleetDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := ioTFleetDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := ioTFleetDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["otel_collector_status"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["agent_version"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["device_count"].(float64); ok {
        data.DeviceCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := ioTFleetDataResponse["online_device_count"].(float64); ok {
        data.OnlineDeviceCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := ioTFleetDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["is_archived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := ioTFleetDataResponse["archived_at"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["archived_by_user_id"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := ioTFleetDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Labels = setValue
    }
    if val, ok := ioTFleetDataResponse["retain_telemetry_data_for_days"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := ioTFleetDataResponse["telemetry_retention_config"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
