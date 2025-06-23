package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ExceptionResource{}
var _ resource.ResourceWithImportState = &ExceptionResource{}

func NewExceptionResource() resource.Resource {
    return &ExceptionResource{}
}

// ExceptionResource defines the resource implementation.
type ExceptionResource struct {
    client *Client
}

// ExceptionResourceModel describes the resource data model.
type ExceptionResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    TelemetryServiceId types.String `tfsdk:"telemetry_service_id"`
    Message types.String `tfsdk:"message"`
    StackTrace types.String `tfsdk:"stack_trace"`
    ExceptionType types.String `tfsdk:"exception_type"`
    Fingerprint types.String `tfsdk:"fingerprint"`
    MarkedAsResolvedAt types.Map `tfsdk:"marked_as_resolved_at"`
    MarkedAsArchivedAt types.Map `tfsdk:"marked_as_archived_at"`
    FirstSeenAt types.Map `tfsdk:"first_seen_at"`
    LastSeenAt types.Map `tfsdk:"last_seen_at"`
    AssignToUserId types.String `tfsdk:"assign_to_user_id"`
    AssignToTeamId types.String `tfsdk:"assign_to_team_id"`
    MarkedAsResolvedByUserId types.String `tfsdk:"marked_as_resolved_by_user_id"`
    MarkedAsArchivedByUserId types.String `tfsdk:"marked_as_archived_by_user_id"`
    IsResolved types.Bool `tfsdk:"is_resolved"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    OccuranceCount types.Number `tfsdk:"occurance_count"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (r *ExceptionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_exception"
}

func (r *ExceptionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "exception resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "telemetry_service_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "message": schema.StringAttribute{
                MarkdownDescription: "Exception Message",
                Optional: true,
            },
            "stack_trace": schema.StringAttribute{
                MarkdownDescription: "Stack Trace",
                Optional: true,
            },
            "exception_type": schema.StringAttribute{
                MarkdownDescription: "Exception Type",
                Optional: true,
            },
            "fingerprint": schema.StringAttribute{
                MarkdownDescription: "Finger Print",
                Optional: true,
            },
            "marked_as_resolved_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "marked_as_archived_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "first_seen_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "last_seen_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "assign_to_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "assign_to_team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "marked_as_resolved_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "marked_as_archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "is_resolved": schema.BoolAttribute{
                MarkdownDescription: "Is Resolved",
                Optional: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is Archived",
                Optional: true,
            },
            "occurance_count": schema.NumberAttribute{
                MarkdownDescription: "Occurances",
                Optional: true,
            },
            "created_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "updated_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "deleted_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Version",
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

func (r *ExceptionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *ExceptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ExceptionResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    exceptionRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "telemetryServiceId": data.TelemetryServiceId.ValueString(),
        "message": data.Message.ValueString(),
        "stackTrace": data.StackTrace.ValueString(),
        "exceptionType": data.ExceptionType.ValueString(),
        "fingerprint": data.Fingerprint.ValueString(),
        "markedAsResolvedAt": r.convertTerraformMapToInterface(data.MarkedAsResolvedAt),
        "markedAsArchivedAt": r.convertTerraformMapToInterface(data.MarkedAsArchivedAt),
        "firstSeenAt": r.convertTerraformMapToInterface(data.FirstSeenAt),
        "lastSeenAt": r.convertTerraformMapToInterface(data.LastSeenAt),
        "assignToUserId": data.AssignToUserId.ValueString(),
        "assignToTeamId": data.AssignToTeamId.ValueString(),
        "markedAsResolvedByUserId": data.MarkedAsResolvedByUserId.ValueString(),
        "markedAsArchivedByUserId": data.MarkedAsArchivedByUserId.ValueString(),
        "isResolved": data.IsResolved.ValueBool(),
        "isArchived": data.IsArchived.ValueBool(),
        "occuranceCount": data.OccuranceCount.ValueBigFloat(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/telemetry-exception-status", exceptionRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create exception, got error: %s", err))
        return
    }

    var exceptionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &exceptionResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := exceptionResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = exceptionResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["telemetryServiceId"].(string); ok && val != "" {
        data.TelemetryServiceId = types.StringValue(val)
    } else {
        data.TelemetryServiceId = types.StringNull()
    }
    if val, ok := dataMap["message"].(string); ok && val != "" {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
    }
    if val, ok := dataMap["stackTrace"].(string); ok && val != "" {
        data.StackTrace = types.StringValue(val)
    } else {
        data.StackTrace = types.StringNull()
    }
    if val, ok := dataMap["exceptionType"].(string); ok && val != "" {
        data.ExceptionType = types.StringValue(val)
    } else {
        data.ExceptionType = types.StringNull()
    }
    if val, ok := dataMap["fingerprint"].(string); ok && val != "" {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
    }
    if val, ok := dataMap["markedAsResolvedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.MarkedAsResolvedAt = mapValue
    } else if dataMap["markedAsResolvedAt"] == nil {
        data.MarkedAsResolvedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["markedAsArchivedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.MarkedAsArchivedAt = mapValue
    } else if dataMap["markedAsArchivedAt"] == nil {
        data.MarkedAsArchivedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["firstSeenAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.FirstSeenAt = mapValue
    } else if dataMap["firstSeenAt"] == nil {
        data.FirstSeenAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.LastSeenAt = mapValue
    } else if dataMap["lastSeenAt"] == nil {
        data.LastSeenAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["assignToUserId"].(string); ok && val != "" {
        data.AssignToUserId = types.StringValue(val)
    } else {
        data.AssignToUserId = types.StringNull()
    }
    if val, ok := dataMap["assignToTeamId"].(string); ok && val != "" {
        data.AssignToTeamId = types.StringValue(val)
    } else {
        data.AssignToTeamId = types.StringNull()
    }
    if val, ok := dataMap["markedAsResolvedByUserId"].(string); ok && val != "" {
        data.MarkedAsResolvedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsResolvedByUserId = types.StringNull()
    }
    if val, ok := dataMap["markedAsArchivedByUserId"].(string); ok && val != "" {
        data.MarkedAsArchivedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isResolved"].(bool); ok {
        data.IsResolved = types.BoolValue(val)
    } else if dataMap["isResolved"] == nil {
        data.IsResolved = types.BoolNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else if dataMap["isArchived"] == nil {
        data.IsArchived = types.BoolNull()
    }
    if val, ok := dataMap["occuranceCount"].(float64); ok {
        data.OccuranceCount = types.NumberValue(big.NewFloat(val))
    } else if dataMap["occuranceCount"] == nil {
        data.OccuranceCount = types.NumberNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ExceptionResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "telemetryServiceId": true,
        "message": true,
        "stackTrace": true,
        "exceptionType": true,
        "fingerprint": true,
        "markedAsResolvedAt": true,
        "markedAsArchivedAt": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "assignToUserId": true,
        "assignToTeamId": true,
        "markedAsResolvedByUserId": true,
        "markedAsArchivedByUserId": true,
        "isResolved": true,
        "isArchived": true,
        "occuranceCount": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/telemetry-exception-status/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read exception, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var exceptionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &exceptionResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := exceptionResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = exceptionResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["telemetryServiceId"].(string); ok && val != "" {
        data.TelemetryServiceId = types.StringValue(val)
    } else {
        data.TelemetryServiceId = types.StringNull()
    }
    if val, ok := dataMap["message"].(string); ok && val != "" {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
    }
    if val, ok := dataMap["stackTrace"].(string); ok && val != "" {
        data.StackTrace = types.StringValue(val)
    } else {
        data.StackTrace = types.StringNull()
    }
    if val, ok := dataMap["exceptionType"].(string); ok && val != "" {
        data.ExceptionType = types.StringValue(val)
    } else {
        data.ExceptionType = types.StringNull()
    }
    if val, ok := dataMap["fingerprint"].(string); ok && val != "" {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
    }
    if val, ok := dataMap["markedAsResolvedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.MarkedAsResolvedAt = mapValue
    } else if dataMap["markedAsResolvedAt"] == nil {
        data.MarkedAsResolvedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["markedAsArchivedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.MarkedAsArchivedAt = mapValue
    } else if dataMap["markedAsArchivedAt"] == nil {
        data.MarkedAsArchivedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["firstSeenAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.FirstSeenAt = mapValue
    } else if dataMap["firstSeenAt"] == nil {
        data.FirstSeenAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.LastSeenAt = mapValue
    } else if dataMap["lastSeenAt"] == nil {
        data.LastSeenAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["assignToUserId"].(string); ok && val != "" {
        data.AssignToUserId = types.StringValue(val)
    } else {
        data.AssignToUserId = types.StringNull()
    }
    if val, ok := dataMap["assignToTeamId"].(string); ok && val != "" {
        data.AssignToTeamId = types.StringValue(val)
    } else {
        data.AssignToTeamId = types.StringNull()
    }
    if val, ok := dataMap["markedAsResolvedByUserId"].(string); ok && val != "" {
        data.MarkedAsResolvedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsResolvedByUserId = types.StringNull()
    }
    if val, ok := dataMap["markedAsArchivedByUserId"].(string); ok && val != "" {
        data.MarkedAsArchivedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isResolved"].(bool); ok {
        data.IsResolved = types.BoolValue(val)
    } else if dataMap["isResolved"] == nil {
        data.IsResolved = types.BoolNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else if dataMap["isArchived"] == nil {
        data.IsArchived = types.BoolNull()
    }
    if val, ok := dataMap["occuranceCount"].(float64); ok {
        data.OccuranceCount = types.NumberValue(big.NewFloat(val))
    } else if dataMap["occuranceCount"] == nil {
        data.OccuranceCount = types.NumberNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ExceptionResourceModel
    var state ExceptionResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // Create API request body
    exceptionRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "telemetryServiceId": data.TelemetryServiceId.ValueString(),
        "message": data.Message.ValueString(),
        "stackTrace": data.StackTrace.ValueString(),
        "exceptionType": data.ExceptionType.ValueString(),
        "fingerprint": data.Fingerprint.ValueString(),
        "markedAsResolvedAt": r.convertTerraformMapToInterface(data.MarkedAsResolvedAt),
        "markedAsArchivedAt": r.convertTerraformMapToInterface(data.MarkedAsArchivedAt),
        "firstSeenAt": r.convertTerraformMapToInterface(data.FirstSeenAt),
        "lastSeenAt": r.convertTerraformMapToInterface(data.LastSeenAt),
        "assignToUserId": data.AssignToUserId.ValueString(),
        "assignToTeamId": data.AssignToTeamId.ValueString(),
        "markedAsResolvedByUserId": data.MarkedAsResolvedByUserId.ValueString(),
        "markedAsArchivedByUserId": data.MarkedAsArchivedByUserId.ValueString(),
        "isResolved": data.IsResolved.ValueBool(),
        "isArchived": data.IsArchived.ValueBool(),
        "occuranceCount": data.OccuranceCount.ValueBigFloat(),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/telemetry-exception-status/" + data.Id.ValueString() + "", exceptionRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update exception, got error: %s", err))
        return
    }

    // Parse the update response
    var exceptionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &exceptionResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "telemetryServiceId": true,
        "message": true,
        "stackTrace": true,
        "exceptionType": true,
        "fingerprint": true,
        "markedAsResolvedAt": true,
        "markedAsArchivedAt": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "assignToUserId": true,
        "assignToTeamId": true,
        "markedAsResolvedByUserId": true,
        "markedAsArchivedByUserId": true,
        "isResolved": true,
        "isArchived": true,
        "occuranceCount": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/telemetry-exception-status/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read exception after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception read response, got error: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["telemetryServiceId"].(string); ok && val != "" {
        data.TelemetryServiceId = types.StringValue(val)
    } else {
        data.TelemetryServiceId = types.StringNull()
    }
    if val, ok := dataMap["message"].(string); ok && val != "" {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
    }
    if val, ok := dataMap["stackTrace"].(string); ok && val != "" {
        data.StackTrace = types.StringValue(val)
    } else {
        data.StackTrace = types.StringNull()
    }
    if val, ok := dataMap["exceptionType"].(string); ok && val != "" {
        data.ExceptionType = types.StringValue(val)
    } else {
        data.ExceptionType = types.StringNull()
    }
    if val, ok := dataMap["fingerprint"].(string); ok && val != "" {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
    }
    if val, ok := dataMap["markedAsResolvedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.MarkedAsResolvedAt = mapValue
    } else if dataMap["markedAsResolvedAt"] == nil {
        data.MarkedAsResolvedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["markedAsArchivedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.MarkedAsArchivedAt = mapValue
    } else if dataMap["markedAsArchivedAt"] == nil {
        data.MarkedAsArchivedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["firstSeenAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.FirstSeenAt = mapValue
    } else if dataMap["firstSeenAt"] == nil {
        data.FirstSeenAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.LastSeenAt = mapValue
    } else if dataMap["lastSeenAt"] == nil {
        data.LastSeenAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["assignToUserId"].(string); ok && val != "" {
        data.AssignToUserId = types.StringValue(val)
    } else {
        data.AssignToUserId = types.StringNull()
    }
    if val, ok := dataMap["assignToTeamId"].(string); ok && val != "" {
        data.AssignToTeamId = types.StringValue(val)
    } else {
        data.AssignToTeamId = types.StringNull()
    }
    if val, ok := dataMap["markedAsResolvedByUserId"].(string); ok && val != "" {
        data.MarkedAsResolvedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsResolvedByUserId = types.StringNull()
    }
    if val, ok := dataMap["markedAsArchivedByUserId"].(string); ok && val != "" {
        data.MarkedAsArchivedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isResolved"].(bool); ok {
        data.IsResolved = types.BoolValue(val)
    } else if dataMap["isResolved"] == nil {
        data.IsResolved = types.BoolNull()
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else if dataMap["isArchived"] == nil {
        data.IsArchived = types.BoolNull()
    }
    if val, ok := dataMap["occuranceCount"].(float64); ok {
        data.OccuranceCount = types.NumberValue(big.NewFloat(val))
    } else if dataMap["occuranceCount"] == nil {
        data.OccuranceCount = types.NumberNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ExceptionResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/telemetry-exception-status/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete exception, got error: %s", err))
        return
    }
}


func (r *ExceptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *ExceptionResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *ExceptionResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
