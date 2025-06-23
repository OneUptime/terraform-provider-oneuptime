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
var _ resource.Resource = &ApiKeyPermissionResource{}
var _ resource.ResourceWithImportState = &ApiKeyPermissionResource{}

func NewApiKeyPermissionResource() resource.Resource {
    return &ApiKeyPermissionResource{}
}

// ApiKeyPermissionResource defines the resource implementation.
type ApiKeyPermissionResource struct {
    client *Client
}

// ApiKeyPermissionResourceModel describes the resource data model.
type ApiKeyPermissionResourceModel struct {
    Id types.String `tfsdk:"id"`
    ApiKeyId types.String `tfsdk:"api_key_id"`
    ProjectId types.String `tfsdk:"project_id"`
    Permission types.Map `tfsdk:"permission"`
    Labels types.List `tfsdk:"labels"`
    IsBlockPermission types.Bool `tfsdk:"is_block_permission"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *ApiKeyPermissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_api_key_permission"
}

func (r *ApiKeyPermissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "api_key_permission resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "api_key_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "permission": schema.MapAttribute{
                MarkdownDescription: "Permission",
                Optional: true,
                ElementType: types.StringType,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Labels",
                Optional: true,
                ElementType: types.StringType,
            },
            "is_block_permission": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create API Key, Edit API Key Permissions], Read: [Project Owner, Project Admin, Read API Key], Update: [Project Owner, Project Admin, Edit API Key Permissions, Edit API Key]",
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
        },
    }
}

func (r *ApiKeyPermissionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *ApiKeyPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ApiKeyPermissionResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    apiKeyPermissionRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "apiKeyId": data.ApiKeyId.ValueString(),
        "projectId": data.ProjectId.ValueString(),
        "permission": r.convertTerraformMapToInterface(data.Permission),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "isBlockPermission": data.IsBlockPermission.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/api-key-permission", apiKeyPermissionRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create api_key_permission, got error: %s", err))
        return
    }

    var apiKeyPermissionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &apiKeyPermissionResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse api_key_permission response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := apiKeyPermissionResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = apiKeyPermissionResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["apiKeyId"].(string); ok && val != "" {
        data.ApiKeyId = types.StringValue(val)
    } else {
        data.ApiKeyId = types.StringNull()
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
    if val, ok := dataMap["permission"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Permission = mapValue
    } else if dataMap["permission"] == nil {
        data.Permission = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["isBlockPermission"].(bool); ok {
        data.IsBlockPermission = types.BoolValue(val)
    } else if dataMap["isBlockPermission"] == nil {
        data.IsBlockPermission = types.BoolNull()
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

func (r *ApiKeyPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ApiKeyPermissionResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "apiKeyId": true,
        "projectId": true,
        "permission": true,
        "labels": true,
        "isBlockPermission": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/api-key-permission/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read api_key_permission, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var apiKeyPermissionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &apiKeyPermissionResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse api_key_permission response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := apiKeyPermissionResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = apiKeyPermissionResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["apiKeyId"].(string); ok && val != "" {
        data.ApiKeyId = types.StringValue(val)
    } else {
        data.ApiKeyId = types.StringNull()
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
    if val, ok := dataMap["permission"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Permission = mapValue
    } else if dataMap["permission"] == nil {
        data.Permission = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["isBlockPermission"].(bool); ok {
        data.IsBlockPermission = types.BoolValue(val)
    } else if dataMap["isBlockPermission"] == nil {
        data.IsBlockPermission = types.BoolNull()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ApiKeyPermissionResourceModel
    var state ApiKeyPermissionResourceModel

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
    apiKeyPermissionRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "apiKeyId": data.ApiKeyId.ValueString(),
        "permission": r.convertTerraformMapToInterface(data.Permission),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "isBlockPermission": data.IsBlockPermission.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/api-key-permission/" + data.Id.ValueString() + "", apiKeyPermissionRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update api_key_permission, got error: %s", err))
        return
    }

    // Parse the update response
    var apiKeyPermissionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &apiKeyPermissionResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse api_key_permission response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "apiKeyId": true,
        "projectId": true,
        "permission": true,
        "labels": true,
        "isBlockPermission": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/api-key-permission/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read api_key_permission after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse api_key_permission read response, got error: %s", err))
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
    if val, ok := dataMap["apiKeyId"].(string); ok && val != "" {
        data.ApiKeyId = types.StringValue(val)
    } else {
        data.ApiKeyId = types.StringNull()
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
    if val, ok := dataMap["permission"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Permission = mapValue
    } else if dataMap["permission"] == nil {
        data.Permission = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["isBlockPermission"].(bool); ok {
        data.IsBlockPermission = types.BoolValue(val)
    } else if dataMap["isBlockPermission"] == nil {
        data.IsBlockPermission = types.BoolNull()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ApiKeyPermissionResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/api-key-permission/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete api_key_permission, got error: %s", err))
        return
    }
}


func (r *ApiKeyPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *ApiKeyPermissionResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *ApiKeyPermissionResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
