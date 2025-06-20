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
var _ resource.Resource = &CopilotEventResource{}
var _ resource.ResourceWithImportState = &CopilotEventResource{}

func NewCopilotEventResource() resource.Resource {
    return &CopilotEventResource{}
}

// CopilotEventResource defines the resource implementation.
type CopilotEventResource struct {
    client *Client
}

// CopilotEventResourceModel describes the resource data model.
type CopilotEventResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    CodeRepositoryId types.String `tfsdk:"code_repository_id"`
    CommitHash types.String `tfsdk:"commit_hash"`
    CopilotActionType types.String `tfsdk:"copilot_action_type"`
    ServiceCatalogId types.String `tfsdk:"service_catalog_id"`
    ServiceRepositoryId types.String `tfsdk:"service_repository_id"`
    CopilotPullRequestId types.String `tfsdk:"copilot_pull_request_id"`
    CopilotActionStatus types.String `tfsdk:"copilot_action_status"`
    CopilotActionProp types.Map `tfsdk:"copilot_action_prop"`
    StatusMessage types.String `tfsdk:"status_message"`
    Logs types.String `tfsdk:"logs"`
    IsPriority types.Bool `tfsdk:"is_priority"`
    StatusChangedAt types.Map `tfsdk:"status_changed_at"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (r *CopilotEventResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_copilot_event"
}

func (r *CopilotEventResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "copilot_event resource",

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
                Required: true,
            },
            "code_repository_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "commit_hash": schema.StringAttribute{
                MarkdownDescription: "Commit Hash",
                Optional: true,
            },
            "copilot_action_type": schema.StringAttribute{
                MarkdownDescription: "Copilot Event Type",
                Optional: true,
            },
            "service_catalog_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "service_repository_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "copilot_pull_request_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "copilot_action_status": schema.StringAttribute{
                MarkdownDescription: "Copilot Event Status",
                Optional: true,
            },
            "copilot_action_prop": schema.MapAttribute{
                MarkdownDescription: "Action Props",
                Optional: true,
                ElementType: types.StringType,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message",
                Optional: true,
            },
            "logs": schema.StringAttribute{
                MarkdownDescription: "Logs",
                Optional: true,
            },
            "is_priority": schema.BoolAttribute{
                MarkdownDescription: "Is Priority",
                Optional: true,
            },
            "status_changed_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
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

func (r *CopilotEventResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *CopilotEventResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data CopilotEventResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    copilotEventRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "codeRepositoryId": data.CodeRepositoryId.ValueString(),
        "commitHash": data.CommitHash.ValueString(),
        "copilotActionType": data.CopilotActionType.ValueString(),
        "serviceCatalogId": data.ServiceCatalogId.ValueString(),
        "serviceRepositoryId": data.ServiceRepositoryId.ValueString(),
        "copilotPullRequestId": data.CopilotPullRequestId.ValueString(),
        "copilotActionStatus": data.CopilotActionStatus.ValueString(),
        "copilotActionProp": r.convertTerraformMapToInterface(data.CopilotActionProp),
        "statusMessage": data.StatusMessage.ValueString(),
        "logs": data.Logs.ValueString(),
        "isPriority": data.IsPriority.ValueBool(),
        "statusChangedAt": r.convertTerraformMapToInterface(data.StatusChangedAt),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/copilot-action", copilotEventRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create copilot_event, got error: %s", err))
        return
    }

    var copilotEventResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &copilotEventResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse copilot_event response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := copilotEventResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = copilotEventResponse
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
    if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if val, ok := dataMap["commitHash"].(string); ok && val != "" {
        data.CommitHash = types.StringValue(val)
    } else {
        data.CommitHash = types.StringNull()
    }
    if val, ok := dataMap["copilotActionType"].(string); ok && val != "" {
        data.CopilotActionType = types.StringValue(val)
    } else {
        data.CopilotActionType = types.StringNull()
    }
    if val, ok := dataMap["serviceCatalogId"].(string); ok && val != "" {
        data.ServiceCatalogId = types.StringValue(val)
    } else {
        data.ServiceCatalogId = types.StringNull()
    }
    if val, ok := dataMap["serviceRepositoryId"].(string); ok && val != "" {
        data.ServiceRepositoryId = types.StringValue(val)
    } else {
        data.ServiceRepositoryId = types.StringNull()
    }
    if val, ok := dataMap["copilotPullRequestId"].(string); ok && val != "" {
        data.CopilotPullRequestId = types.StringValue(val)
    } else {
        data.CopilotPullRequestId = types.StringNull()
    }
    if val, ok := dataMap["copilotActionStatus"].(string); ok && val != "" {
        data.CopilotActionStatus = types.StringValue(val)
    } else {
        data.CopilotActionStatus = types.StringNull()
    }
    if val, ok := dataMap["copilotActionProp"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CopilotActionProp = mapValue
    } else if dataMap["copilotActionProp"] == nil {
        data.CopilotActionProp = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["logs"].(string); ok && val != "" {
        data.Logs = types.StringValue(val)
    } else {
        data.Logs = types.StringNull()
    }
    if val, ok := dataMap["isPriority"].(bool); ok {
        data.IsPriority = types.BoolValue(val)
    } else if dataMap["isPriority"] == nil {
        data.IsPriority = types.BoolNull()
    }
    if val, ok := dataMap["statusChangedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StatusChangedAt = mapValue
    } else if dataMap["statusChangedAt"] == nil {
        data.StatusChangedAt = types.MapNull(types.StringType)
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

func (r *CopilotEventResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data CopilotEventResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "codeRepositoryId": true,
        "commitHash": true,
        "copilotActionType": true,
        "serviceCatalogId": true,
        "serviceRepositoryId": true,
        "copilotPullRequestId": true,
        "copilotActionStatus": true,
        "copilotActionProp": true,
        "statusMessage": true,
        "logs": true,
        "isPriority": true,
        "statusChangedAt": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/copilot-action/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read copilot_event, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var copilotEventResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &copilotEventResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse copilot_event response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := copilotEventResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = copilotEventResponse
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
    if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if val, ok := dataMap["commitHash"].(string); ok && val != "" {
        data.CommitHash = types.StringValue(val)
    } else {
        data.CommitHash = types.StringNull()
    }
    if val, ok := dataMap["copilotActionType"].(string); ok && val != "" {
        data.CopilotActionType = types.StringValue(val)
    } else {
        data.CopilotActionType = types.StringNull()
    }
    if val, ok := dataMap["serviceCatalogId"].(string); ok && val != "" {
        data.ServiceCatalogId = types.StringValue(val)
    } else {
        data.ServiceCatalogId = types.StringNull()
    }
    if val, ok := dataMap["serviceRepositoryId"].(string); ok && val != "" {
        data.ServiceRepositoryId = types.StringValue(val)
    } else {
        data.ServiceRepositoryId = types.StringNull()
    }
    if val, ok := dataMap["copilotPullRequestId"].(string); ok && val != "" {
        data.CopilotPullRequestId = types.StringValue(val)
    } else {
        data.CopilotPullRequestId = types.StringNull()
    }
    if val, ok := dataMap["copilotActionStatus"].(string); ok && val != "" {
        data.CopilotActionStatus = types.StringValue(val)
    } else {
        data.CopilotActionStatus = types.StringNull()
    }
    if val, ok := dataMap["copilotActionProp"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CopilotActionProp = mapValue
    } else if dataMap["copilotActionProp"] == nil {
        data.CopilotActionProp = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["logs"].(string); ok && val != "" {
        data.Logs = types.StringValue(val)
    } else {
        data.Logs = types.StringNull()
    }
    if val, ok := dataMap["isPriority"].(bool); ok {
        data.IsPriority = types.BoolValue(val)
    } else if dataMap["isPriority"] == nil {
        data.IsPriority = types.BoolNull()
    }
    if val, ok := dataMap["statusChangedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StatusChangedAt = mapValue
    } else if dataMap["statusChangedAt"] == nil {
        data.StatusChangedAt = types.MapNull(types.StringType)
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

func (r *CopilotEventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data CopilotEventResourceModel
    var state CopilotEventResourceModel

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
    copilotEventRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "codeRepositoryId": data.CodeRepositoryId.ValueString(),
        "commitHash": data.CommitHash.ValueString(),
        "copilotActionType": data.CopilotActionType.ValueString(),
        "serviceCatalogId": data.ServiceCatalogId.ValueString(),
        "serviceRepositoryId": data.ServiceRepositoryId.ValueString(),
        "copilotPullRequestId": data.CopilotPullRequestId.ValueString(),
        "copilotActionStatus": data.CopilotActionStatus.ValueString(),
        "copilotActionProp": r.convertTerraformMapToInterface(data.CopilotActionProp),
        "statusMessage": data.StatusMessage.ValueString(),
        "logs": data.Logs.ValueString(),
        "isPriority": data.IsPriority.ValueBool(),
        "statusChangedAt": r.convertTerraformMapToInterface(data.StatusChangedAt),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/copilot-action/" + data.Id.ValueString() + "", copilotEventRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update copilot_event, got error: %s", err))
        return
    }

    // Parse the update response
    var copilotEventResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &copilotEventResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse copilot_event response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "codeRepositoryId": true,
        "commitHash": true,
        "copilotActionType": true,
        "serviceCatalogId": true,
        "serviceRepositoryId": true,
        "copilotPullRequestId": true,
        "copilotActionStatus": true,
        "copilotActionProp": true,
        "statusMessage": true,
        "logs": true,
        "isPriority": true,
        "statusChangedAt": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/copilot-action/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read copilot_event after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse copilot_event read response, got error: %s", err))
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
    if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if val, ok := dataMap["commitHash"].(string); ok && val != "" {
        data.CommitHash = types.StringValue(val)
    } else {
        data.CommitHash = types.StringNull()
    }
    if val, ok := dataMap["copilotActionType"].(string); ok && val != "" {
        data.CopilotActionType = types.StringValue(val)
    } else {
        data.CopilotActionType = types.StringNull()
    }
    if val, ok := dataMap["serviceCatalogId"].(string); ok && val != "" {
        data.ServiceCatalogId = types.StringValue(val)
    } else {
        data.ServiceCatalogId = types.StringNull()
    }
    if val, ok := dataMap["serviceRepositoryId"].(string); ok && val != "" {
        data.ServiceRepositoryId = types.StringValue(val)
    } else {
        data.ServiceRepositoryId = types.StringNull()
    }
    if val, ok := dataMap["copilotPullRequestId"].(string); ok && val != "" {
        data.CopilotPullRequestId = types.StringValue(val)
    } else {
        data.CopilotPullRequestId = types.StringNull()
    }
    if val, ok := dataMap["copilotActionStatus"].(string); ok && val != "" {
        data.CopilotActionStatus = types.StringValue(val)
    } else {
        data.CopilotActionStatus = types.StringNull()
    }
    if val, ok := dataMap["copilotActionProp"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CopilotActionProp = mapValue
    } else if dataMap["copilotActionProp"] == nil {
        data.CopilotActionProp = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["logs"].(string); ok && val != "" {
        data.Logs = types.StringValue(val)
    } else {
        data.Logs = types.StringNull()
    }
    if val, ok := dataMap["isPriority"].(bool); ok {
        data.IsPriority = types.BoolValue(val)
    } else if dataMap["isPriority"] == nil {
        data.IsPriority = types.BoolNull()
    }
    if val, ok := dataMap["statusChangedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StatusChangedAt = mapValue
    } else if dataMap["statusChangedAt"] == nil {
        data.StatusChangedAt = types.MapNull(types.StringType)
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

func (r *CopilotEventResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data CopilotEventResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/copilot-action/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete copilot_event, got error: %s", err))
        return
    }
}


func (r *CopilotEventResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *CopilotEventResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *CopilotEventResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
