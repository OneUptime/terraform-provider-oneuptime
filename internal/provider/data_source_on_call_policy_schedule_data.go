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
var _ datasource.DataSource = &OnCallPolicyScheduleDataDataSource{}

func NewOnCallPolicyScheduleDataDataSource() datasource.DataSource {
    return &OnCallPolicyScheduleDataDataSource{}
}

// OnCallPolicyScheduleDataDataSource defines the data source implementation.
type OnCallPolicyScheduleDataDataSource struct {
    client *Client
}

// OnCallPolicyScheduleDataDataSourceModel describes the data source data model.
type OnCallPolicyScheduleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Labels types.List `tfsdk:"labels"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    CurrentUserIdOnRoster types.String `tfsdk:"current_user_id_on_roster"`
    NextUserIdOnRoster types.String `tfsdk:"next_user_id_on_roster"`
    RosterHandoffAt types.String `tfsdk:"roster_handoff_at"`
    RosterNextHandoffAt types.String `tfsdk:"roster_next_handoff_at"`
    RosterNextStartAt types.String `tfsdk:"roster_next_start_at"`
    RosterStartAt types.String `tfsdk:"roster_start_at"`
}

func (d *OnCallPolicyScheduleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_on_call_policy_schedule_data"
}

func (d *OnCallPolicyScheduleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "on_call_policy_schedule_data data source",

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
            "labels": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Duty Policy Schedule], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Schedule], Update: [Project Owner, Project Admin, Project Member, Edit On-Call Duty Policy Schedule]",
                Computed: true,
                ElementType: types.StringType,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Duty Policy Schedule], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Schedule], Update: [Project Owner, Project Admin, Project Member, Edit On-Call Duty Policy Schedule]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create On-Call Duty Policy Schedule], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Schedule], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "current_user_id_on_roster": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "next_user_id_on_roster": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "roster_handoff_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "roster_next_handoff_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "roster_next_start_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "roster_start_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *OnCallPolicyScheduleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OnCallPolicyScheduleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data OnCallPolicyScheduleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "on-call-duty-policy-schedule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read on_call_policy_schedule_data, got error: %s", err))
        return
    }

    var onCallPolicyScheduleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &onCallPolicyScheduleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse on_call_policy_schedule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := onCallPolicyScheduleDataResponse["data"].(map[string]interface{}); ok {
        onCallPolicyScheduleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := onCallPolicyScheduleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := onCallPolicyScheduleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := onCallPolicyScheduleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["current_user_id_on_roster"].(string); ok {
        data.CurrentUserIdOnRoster = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["next_user_id_on_roster"].(string); ok {
        data.NextUserIdOnRoster = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["roster_handoff_at"].(string); ok {
        data.RosterHandoffAt = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["roster_next_handoff_at"].(string); ok {
        data.RosterNextHandoffAt = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["roster_next_start_at"].(string); ok {
        data.RosterNextStartAt = types.StringValue(val)
    }
    if val, ok := onCallPolicyScheduleDataResponse["roster_start_at"].(string); ok {
        data.RosterStartAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
