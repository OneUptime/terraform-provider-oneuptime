package provider

import (
    "context"
    "fmt"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OnCallScheduleLayerDataDataSource{}

func NewOnCallScheduleLayerDataDataSource() datasource.DataSource {
    return &OnCallScheduleLayerDataDataSource{}
}

// OnCallScheduleLayerDataDataSource defines the data source implementation.
type OnCallScheduleLayerDataDataSource struct {
    client *Client
}

// OnCallScheduleLayerDataDataSourceModel describes the data source data model.
type OnCallScheduleLayerDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    OnCallDutyPolicyScheduleId types.String `tfsdk:"on_call_duty_policy_schedule_id"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Order types.Number `tfsdk:"order"`
    StartsAt types.String `tfsdk:"starts_at"`
    Rotation types.String `tfsdk:"rotation"`
    HandOffTime types.String `tfsdk:"hand_off_time"`
    RestrictionTimes types.String `tfsdk:"restriction_times"`
}

func (d *OnCallScheduleLayerDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_on_call_schedule_layer_data"
}

func (d *OnCallScheduleLayerDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "on_call_schedule_layer_data data source",

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
            "on_call_duty_policy_schedule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Schedule Layer], Read: [Project Owner, Project Admin, Project Member, Read On-Call Schedule Layer], Update: [Project Owner, Project Admin, Edit On-Call Schedule Layer]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Schedule Layer], Read: [Project Owner, Project Admin, Project Member, Read On-Call Schedule Layer], Update: [Project Owner, Project Admin, Project Member, Edit On-Call Schedule Layer]",
                Computed: true,
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "rotation": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Schedule Layer], Read: [Project Owner, Project Admin, Project Member, Read On-Call Schedule Layer], Update: [Project Owner, Project Admin, Project Member, Edit On-Call Schedule Layer]",
                Computed: true,
            },
            "hand_off_time": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "restriction_times": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Schedule Layer], Read: [Project Owner, Project Admin, Project Member, Read On-Call Schedule Layer], Update: [Project Owner, Project Admin, Project Member, Edit On-Call Schedule Layer]",
                Computed: true,
            },
        },
    }
}

func (d *OnCallScheduleLayerDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OnCallScheduleLayerDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data OnCallScheduleLayerDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "on-call-duty-schedule-layer" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read on_call_schedule_layer_data, got error: %s", err))
        return
    }

    var onCallScheduleLayerDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &onCallScheduleLayerDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse on_call_schedule_layer_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := onCallScheduleLayerDataResponse["data"].(map[string]interface{}); ok {
        onCallScheduleLayerDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := onCallScheduleLayerDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := onCallScheduleLayerDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["on_call_duty_policy_schedule_id"].(string); ok {
        data.OnCallDutyPolicyScheduleId = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := onCallScheduleLayerDataResponse["starts_at"].(string); ok {
        data.StartsAt = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["rotation"].(string); ok {
        data.Rotation = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["hand_off_time"].(string); ok {
        data.HandOffTime = types.StringValue(val)
    }
    if val, ok := onCallScheduleLayerDataResponse["restriction_times"].(string); ok {
        data.RestrictionTimes = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
