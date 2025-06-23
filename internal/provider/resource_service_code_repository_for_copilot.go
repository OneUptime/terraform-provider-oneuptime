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
var _ resource.Resource = &ServiceCodeRepositoryForCopilotResource{}
var _ resource.ResourceWithImportState = &ServiceCodeRepositoryForCopilotResource{}

func NewServiceCodeRepositoryForCopilotResource() resource.Resource {
    return &ServiceCodeRepositoryForCopilotResource{}
}

// ServiceCodeRepositoryForCopilotResource defines the resource implementation.
type ServiceCodeRepositoryForCopilotResource struct {
    client *Client
}

// ServiceCodeRepositoryForCopilotResourceModel describes the resource data model.
type ServiceCodeRepositoryForCopilotResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    ServicePathInRepository types.String `tfsdk:"service_path_in_repository"`
    LimitNumberOfOpenPullRequestsCount types.Number `tfsdk:"limit_number_of_open_pull_requests_count"`
    EnablePullRequests types.Bool `tfsdk:"enable_pull_requests"`
    CodeRepositoryId types.String `tfsdk:"code_repository_id"`
    ServiceCatalogId types.String `tfsdk:"service_catalog_id"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (r *ServiceCodeRepositoryForCopilotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service_code_repository_for_copilot"
}

func (r *ServiceCodeRepositoryForCopilotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "service_code_repository_for_copilot resource",

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
            "service_path_in_repository": schema.StringAttribute{
                MarkdownDescription: "Path in Repository",
                Optional: true,
            },
            "limit_number_of_open_pull_requests_count": schema.NumberAttribute{
                MarkdownDescription: "Limit Number of Open Pull Requests Count",
                Optional: true,
            },
            "enable_pull_requests": schema.BoolAttribute{
                MarkdownDescription: "Enable Pull Requests",
                Optional: true,
            },
            "code_repository_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "service_catalog_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
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

func (r *ServiceCodeRepositoryForCopilotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *ServiceCodeRepositoryForCopilotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ServiceCodeRepositoryForCopilotResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    serviceCodeRepositoryForCopilotRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "servicePathInRepository": data.ServicePathInRepository.ValueString(),
        "limitNumberOfOpenPullRequestsCount": data.LimitNumberOfOpenPullRequestsCount.ValueBigFloat(),
        "enablePullRequests": data.EnablePullRequests.ValueBool(),
        "codeRepositoryId": data.CodeRepositoryId.ValueString(),
        "serviceCatalogId": data.ServiceCatalogId.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/service-copilot-code-repository", serviceCodeRepositoryForCopilotRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service_code_repository_for_copilot, got error: %s", err))
        return
    }

    var serviceCodeRepositoryForCopilotResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &serviceCodeRepositoryForCopilotResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_code_repository_for_copilot response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := serviceCodeRepositoryForCopilotResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = serviceCodeRepositoryForCopilotResponse
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
    if val, ok := dataMap["servicePathInRepository"].(string); ok && val != "" {
        data.ServicePathInRepository = types.StringValue(val)
    } else {
        data.ServicePathInRepository = types.StringNull()
    }
    if val, ok := dataMap["limitNumberOfOpenPullRequestsCount"].(float64); ok {
        data.LimitNumberOfOpenPullRequestsCount = types.NumberValue(big.NewFloat(val))
    } else if dataMap["limitNumberOfOpenPullRequestsCount"] == nil {
        data.LimitNumberOfOpenPullRequestsCount = types.NumberNull()
    }
    if val, ok := dataMap["enablePullRequests"].(bool); ok {
        data.EnablePullRequests = types.BoolValue(val)
    } else if dataMap["enablePullRequests"] == nil {
        data.EnablePullRequests = types.BoolNull()
    }
    if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if val, ok := dataMap["serviceCatalogId"].(string); ok && val != "" {
        data.ServiceCatalogId = types.StringValue(val)
    } else {
        data.ServiceCatalogId = types.StringNull()
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

func (r *ServiceCodeRepositoryForCopilotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ServiceCodeRepositoryForCopilotResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "servicePathInRepository": true,
        "limitNumberOfOpenPullRequestsCount": true,
        "enablePullRequests": true,
        "codeRepositoryId": true,
        "serviceCatalogId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/service-copilot-code-repository/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_code_repository_for_copilot, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var serviceCodeRepositoryForCopilotResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &serviceCodeRepositoryForCopilotResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_code_repository_for_copilot response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := serviceCodeRepositoryForCopilotResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = serviceCodeRepositoryForCopilotResponse
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
    if val, ok := dataMap["servicePathInRepository"].(string); ok && val != "" {
        data.ServicePathInRepository = types.StringValue(val)
    } else {
        data.ServicePathInRepository = types.StringNull()
    }
    if val, ok := dataMap["limitNumberOfOpenPullRequestsCount"].(float64); ok {
        data.LimitNumberOfOpenPullRequestsCount = types.NumberValue(big.NewFloat(val))
    } else if dataMap["limitNumberOfOpenPullRequestsCount"] == nil {
        data.LimitNumberOfOpenPullRequestsCount = types.NumberNull()
    }
    if val, ok := dataMap["enablePullRequests"].(bool); ok {
        data.EnablePullRequests = types.BoolValue(val)
    } else if dataMap["enablePullRequests"] == nil {
        data.EnablePullRequests = types.BoolNull()
    }
    if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if val, ok := dataMap["serviceCatalogId"].(string); ok && val != "" {
        data.ServiceCatalogId = types.StringValue(val)
    } else {
        data.ServiceCatalogId = types.StringNull()
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

func (r *ServiceCodeRepositoryForCopilotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ServiceCodeRepositoryForCopilotResourceModel
    var state ServiceCodeRepositoryForCopilotResourceModel

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
    serviceCodeRepositoryForCopilotRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "servicePathInRepository": data.ServicePathInRepository.ValueString(),
        "limitNumberOfOpenPullRequestsCount": data.LimitNumberOfOpenPullRequestsCount.ValueBigFloat(),
        "enablePullRequests": data.EnablePullRequests.ValueBool(),
        "codeRepositoryId": data.CodeRepositoryId.ValueString(),
        "serviceCatalogId": data.ServiceCatalogId.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/service-copilot-code-repository/" + data.Id.ValueString() + "", serviceCodeRepositoryForCopilotRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update service_code_repository_for_copilot, got error: %s", err))
        return
    }

    // Parse the update response
    var serviceCodeRepositoryForCopilotResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &serviceCodeRepositoryForCopilotResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_code_repository_for_copilot response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "servicePathInRepository": true,
        "limitNumberOfOpenPullRequestsCount": true,
        "enablePullRequests": true,
        "codeRepositoryId": true,
        "serviceCatalogId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/service-copilot-code-repository/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_code_repository_for_copilot after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_code_repository_for_copilot read response, got error: %s", err))
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
    if val, ok := dataMap["servicePathInRepository"].(string); ok && val != "" {
        data.ServicePathInRepository = types.StringValue(val)
    } else {
        data.ServicePathInRepository = types.StringNull()
    }
    if val, ok := dataMap["limitNumberOfOpenPullRequestsCount"].(float64); ok {
        data.LimitNumberOfOpenPullRequestsCount = types.NumberValue(big.NewFloat(val))
    } else if dataMap["limitNumberOfOpenPullRequestsCount"] == nil {
        data.LimitNumberOfOpenPullRequestsCount = types.NumberNull()
    }
    if val, ok := dataMap["enablePullRequests"].(bool); ok {
        data.EnablePullRequests = types.BoolValue(val)
    } else if dataMap["enablePullRequests"] == nil {
        data.EnablePullRequests = types.BoolNull()
    }
    if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if val, ok := dataMap["serviceCatalogId"].(string); ok && val != "" {
        data.ServiceCatalogId = types.StringValue(val)
    } else {
        data.ServiceCatalogId = types.StringNull()
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

func (r *ServiceCodeRepositoryForCopilotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ServiceCodeRepositoryForCopilotResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/service-copilot-code-repository/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service_code_repository_for_copilot, got error: %s", err))
        return
    }
}


func (r *ServiceCodeRepositoryForCopilotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *ServiceCodeRepositoryForCopilotResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *ServiceCodeRepositoryForCopilotResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
