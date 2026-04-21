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
var _ datasource.DataSource = &TraceScrubRuleDataDataSource{}

func NewTraceScrubRuleDataDataSource() datasource.DataSource {
    return &TraceScrubRuleDataDataSource{}
}

// TraceScrubRuleDataDataSource defines the data source implementation.
type TraceScrubRuleDataDataSource struct {
    client *Client
}

// TraceScrubRuleDataDataSourceModel describes the data source data model.
type TraceScrubRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    PatternType types.String `tfsdk:"pattern_type"`
    CustomRegex types.String `tfsdk:"custom_regex"`
    ScrubAction types.String `tfsdk:"scrub_action"`
    FieldsToScrub types.String `tfsdk:"fields_to_scrub"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    SortOrder types.Number `tfsdk:"sort_order"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *TraceScrubRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_trace_scrub_rule_data"
}

func (d *TraceScrubRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "trace_scrub_rule_data data source",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this scrub rule does.. Permissions - Create: [Project Owner, Project Admin, Create Trace Scrub Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Scrub Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Scrub Rule]",
                Computed: true,
            },
            "pattern_type": schema.StringAttribute{
                MarkdownDescription: "The type of sensitive data pattern to detect: email, creditCard, ssn, phoneNumber, ipAddress, or custom.. Permissions - Create: [Project Owner, Project Admin, Create Trace Scrub Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Scrub Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Scrub Rule]",
                Computed: true,
            },
            "custom_regex": schema.StringAttribute{
                MarkdownDescription: "A custom regular expression pattern to match. Only used when patternType is 'custom'.. Permissions - Create: [Project Owner, Project Admin, Create Trace Scrub Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Scrub Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Scrub Rule]",
                Computed: true,
            },
            "scrub_action": schema.StringAttribute{
                MarkdownDescription: "How to scrub matched data: 'mask' partially hides it, 'hash' replaces with a hash, 'redact' removes entirely.. Permissions - Create: [Project Owner, Project Admin, Create Trace Scrub Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Scrub Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Scrub Rule]",
                Computed: true,
            },
            "fields_to_scrub": schema.StringAttribute{
                MarkdownDescription: "Which span fields to scrub: 'name' (span name), 'attributes' (attribute values), 'events' (span event attributes), or 'all'.. Permissions - Create: [Project Owner, Project Admin, Create Trace Scrub Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Scrub Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Scrub Rule]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this scrub rule is active.. Permissions - Create: [Project Owner, Project Admin, Create Trace Scrub Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Scrub Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Scrub Rule]",
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Determines the evaluation order of this rule relative to others.. Permissions - Create: [Project Owner, Project Admin, Create Trace Scrub Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Manager, Read Trace Scrub Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Trace Scrub Rule]",
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
        },
    }
}

func (d *TraceScrubRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TraceScrubRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TraceScrubRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "trace-scrub-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read trace_scrub_rule_data, got error: %s", err))
        return
    }

    var traceScrubRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &traceScrubRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse trace_scrub_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := traceScrubRuleDataResponse["data"].(map[string]interface{}); ok {
        traceScrubRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := traceScrubRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := traceScrubRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["pattern_type"].(string); ok {
        data.PatternType = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["custom_regex"].(string); ok {
        data.CustomRegex = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["scrub_action"].(string); ok {
        data.ScrubAction = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["fields_to_scrub"].(string); ok {
        data.FieldsToScrub = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["sort_order"].(float64); ok {
        data.SortOrder = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := traceScrubRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := traceScrubRuleDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
