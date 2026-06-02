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
var _ datasource.DataSource = &ExceptionDataDataSource{}

func NewExceptionDataDataSource() datasource.DataSource {
    return &ExceptionDataDataSource{}
}

// ExceptionDataDataSource defines the data source implementation.
type ExceptionDataDataSource struct {
    client *Client
}

// ExceptionDataDataSourceModel describes the data source data model.
type ExceptionDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
    ServiceType types.String `tfsdk:"service_type"`
    Message types.String `tfsdk:"message"`
    StackTrace types.String `tfsdk:"stack_trace"`
    ExceptionType types.String `tfsdk:"exception_type"`
    Fingerprint types.String `tfsdk:"fingerprint"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    MarkedAsResolvedAt types.String `tfsdk:"marked_as_resolved_at"`
    MarkedAsArchivedAt types.String `tfsdk:"marked_as_archived_at"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    AssignToUserId types.String `tfsdk:"assign_to_user_id"`
    AssignToTeamId types.String `tfsdk:"assign_to_team_id"`
    MarkedAsResolvedByUserId types.String `tfsdk:"marked_as_resolved_by_user_id"`
    MarkedAsArchivedByUserId types.String `tfsdk:"marked_as_archived_by_user_id"`
    IsResolved types.Bool `tfsdk:"is_resolved"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    OccuranceCount types.Number `tfsdk:"occurance_count"`
    FirstSeenInRelease types.String `tfsdk:"first_seen_in_release"`
    LastSeenInRelease types.String `tfsdk:"last_seen_in_release"`
    Environment types.String `tfsdk:"environment"`
}

func (d *ExceptionDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_exception_data"
}

func (d *ExceptionDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "exception_data data source",

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
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "service_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "service_type": schema.StringAttribute{
                MarkdownDescription: "Resource type that produced this exception (e.g. OpenTelemetry service, Host, DockerHost, KubernetesCluster, or Unknown for unattributed telemetry).. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "message": schema.StringAttribute{
                MarkdownDescription: "Exception message that was thrown by the telemetry service. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]",
                Computed: true,
            },
            "stack_trace": schema.StringAttribute{
                MarkdownDescription: "Stack trace of the exception that was thrown by the telemetry service. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]",
                Computed: true,
            },
            "exception_type": schema.StringAttribute{
                MarkdownDescription: "Type of the exception that was thrown by the telemetry service. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]",
                Computed: true,
            },
            "fingerprint": schema.StringAttribute{
                MarkdownDescription: "Finger print of the exception that was thrown by the telemetry service. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]",
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
            "marked_as_resolved_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "marked_as_archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "assign_to_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "assign_to_team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "marked_as_resolved_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "marked_as_archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_resolved": schema.BoolAttribute{
                MarkdownDescription: "Is this exception resolved?. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this exception archived?. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]",
                Computed: true,
            },
            "occurance_count": schema.NumberAttribute{
                MarkdownDescription: "Number of times this exception has occurred. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]",
                Computed: true,
            },
            "first_seen_in_release": schema.StringAttribute{
                MarkdownDescription: "The service version / release in which this exception was first observed. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]",
                Computed: true,
            },
            "last_seen_in_release": schema.StringAttribute{
                MarkdownDescription: "The most recent service version / release in which this exception was observed. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]",
                Computed: true,
            },
            "environment": schema.StringAttribute{
                MarkdownDescription: "Deployment environment from deployment.environment resource attribute. Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Viewer, Telemetry Admin, Telemetry Member, Telemetry Viewer, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]",
                Computed: true,
            },
        },
    }
}

func (d *ExceptionDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExceptionDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ExceptionDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "telemetry-exception" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read exception_data, got error: %s", err))
        return
    }

    var exceptionDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &exceptionDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := exceptionDataResponse["data"].(map[string]interface{}); ok {
        exceptionDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := exceptionDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := exceptionDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["service_id"].(string); ok {
        data.ServiceId = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["service_type"].(string); ok {
        data.ServiceType = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["message"].(string); ok {
        data.Message = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["stack_trace"].(string); ok {
        data.StackTrace = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["exception_type"].(string); ok {
        data.ExceptionType = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["fingerprint"].(string); ok {
        data.Fingerprint = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["marked_as_resolved_at"].(string); ok {
        data.MarkedAsResolvedAt = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["marked_as_archived_at"].(string); ok {
        data.MarkedAsArchivedAt = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["first_seen_at"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["last_seen_at"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["assign_to_user_id"].(string); ok {
        data.AssignToUserId = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["assign_to_team_id"].(string); ok {
        data.AssignToTeamId = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["marked_as_resolved_by_user_id"].(string); ok {
        data.MarkedAsResolvedByUserId = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["marked_as_archived_by_user_id"].(string); ok {
        data.MarkedAsArchivedByUserId = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["is_resolved"].(bool); ok {
        data.IsResolved = types.BoolValue(val)
    }
    if val, ok := exceptionDataResponse["is_archived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    if val, ok := exceptionDataResponse["occurance_count"].(float64); ok {
        data.OccuranceCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := exceptionDataResponse["first_seen_in_release"].(string); ok {
        data.FirstSeenInRelease = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["last_seen_in_release"].(string); ok {
        data.LastSeenInRelease = types.StringValue(val)
    }
    if val, ok := exceptionDataResponse["environment"].(string); ok {
        data.Environment = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
