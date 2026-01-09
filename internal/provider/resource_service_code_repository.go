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
    "encoding/json"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ServiceCodeRepositoryResource{}
var _ resource.ResourceWithImportState = &ServiceCodeRepositoryResource{}

func NewServiceCodeRepositoryResource() resource.Resource {
    return &ServiceCodeRepositoryResource{}
}

// ServiceCodeRepositoryResource defines the resource implementation.
type ServiceCodeRepositoryResource struct {
    client *Client
}

// ServiceCodeRepositoryResourceModel describes the resource data model.
type ServiceCodeRepositoryResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
    CodeRepositoryId types.String `tfsdk:"code_repository_id"`
    ServicePathInRepository types.String `tfsdk:"service_path_in_repository"`
    EnableAutomaticImprovements types.Bool `tfsdk:"enable_automatic_improvements"`
    MaxOpenPullRequests types.Number `tfsdk:"max_open_pull_requests"`
    RestrictedImprovementActions types.String `tfsdk:"restricted_improvement_actions"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (r *ServiceCodeRepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service_code_repository"
}

func (r *ServiceCodeRepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "service_code_repository resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
            },
            "service_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "code_repository_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "service_path_in_repository": schema.StringAttribute{
                MarkdownDescription: "The path in the repository where the service code lives (e.g., /services/api or /src/backend). Permissions - Create: [Project Owner, Project Admin, Project Member, Create Service Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Service Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Service Code Repository]",
                Optional: true,
                Computed: true,
            },
            "enable_automatic_improvements": schema.BoolAttribute{
                MarkdownDescription: "Enable OneUptime to automatically create pull requests to improve the code for this service.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Service Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Service Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Service Code Repository]",
                Optional: true,
                Computed: true,
            },
            "max_open_pull_requests": schema.NumberAttribute{
                MarkdownDescription: "Maximum number of open pull requests that OneUptime can create for this service at any given time.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Service Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Service Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Service Code Repository]",
                Optional: true,
                Computed: true,
            },
            "restricted_improvement_actions": schema.StringAttribute{
                MarkdownDescription: "Restrict code improvements to only these actions. If empty, all improvement actions are allowed.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Service Code Repository], Read: [Project Owner, Project Admin, Project Member, Read Service Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Service Code Repository]",
                Optional: true,
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

func (r *ServiceCodeRepositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *ServiceCodeRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ServiceCodeRepositoryResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    serviceCodeRepositoryRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "serviceId": data.ServiceId.ValueString(),
        "codeRepositoryId": data.CodeRepositoryId.ValueString(),
        "servicePathInRepository": data.ServicePathInRepository.ValueString(),
        "enableAutomaticImprovements": data.EnableAutomaticImprovements.ValueBool(),
        "maxOpenPullRequests": data.MaxOpenPullRequests.ValueBigFloat(),
        "restrictedImprovementActions": r.parseJSONField(data.RestrictedImprovementActions),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/service-code-repository", serviceCodeRepositoryRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service_code_repository, got error: %s", err))
        return
    }

    var serviceCodeRepositoryResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &serviceCodeRepositoryResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_code_repository response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := serviceCodeRepositoryResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = serviceCodeRepositoryResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if obj, ok := dataMap["serviceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceId"].(string); ok && val != "" {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := dataMap["codeRepositoryId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else {
            data.CodeRepositoryId = types.StringNull()
        }
    } else if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if obj, ok := dataMap["servicePathInRepository"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServicePathInRepository = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ServicePathInRepository = types.StringValue(val)
        } else {
            data.ServicePathInRepository = types.StringNull()
        }
    } else if val, ok := dataMap["servicePathInRepository"].(string); ok && val != "" {
        data.ServicePathInRepository = types.StringValue(val)
    } else {
        data.ServicePathInRepository = types.StringNull()
    }
    if val, ok := dataMap["enableAutomaticImprovements"].(bool); ok {
        data.EnableAutomaticImprovements = types.BoolValue(val)
    } else if dataMap["enableAutomaticImprovements"] == nil {
        data.EnableAutomaticImprovements = types.BoolNull()
    }
    if val, ok := dataMap["maxOpenPullRequests"].(float64); ok {
        data.MaxOpenPullRequests = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["maxOpenPullRequests"].(int); ok {
        data.MaxOpenPullRequests = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["maxOpenPullRequests"].(int64); ok {
        data.MaxOpenPullRequests = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["maxOpenPullRequests"] == nil {
        data.MaxOpenPullRequests = types.NumberNull()
    }
    if val, ok := dataMap["restrictedImprovementActions"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.RestrictedImprovementActions = types.StringValue(string(jsonBytes))
        } else {
            data.RestrictedImprovementActions = types.StringNull()
        }
    } else if val, ok := dataMap["restrictedImprovementActions"].(string); ok && val != "" {
        data.RestrictedImprovementActions = types.StringValue(val)
    } else {
        data.RestrictedImprovementActions = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
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

func (r *ServiceCodeRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ServiceCodeRepositoryResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceId": true,
        "codeRepositoryId": true,
        "servicePathInRepository": true,
        "enableAutomaticImprovements": true,
        "maxOpenPullRequests": true,
        "restrictedImprovementActions": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/service-code-repository/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_code_repository, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var serviceCodeRepositoryResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &serviceCodeRepositoryResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_code_repository response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := serviceCodeRepositoryResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = serviceCodeRepositoryResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if obj, ok := dataMap["serviceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceId"].(string); ok && val != "" {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := dataMap["codeRepositoryId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else {
            data.CodeRepositoryId = types.StringNull()
        }
    } else if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if obj, ok := dataMap["servicePathInRepository"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServicePathInRepository = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ServicePathInRepository = types.StringValue(val)
        } else {
            data.ServicePathInRepository = types.StringNull()
        }
    } else if val, ok := dataMap["servicePathInRepository"].(string); ok && val != "" {
        data.ServicePathInRepository = types.StringValue(val)
    } else {
        data.ServicePathInRepository = types.StringNull()
    }
    if val, ok := dataMap["enableAutomaticImprovements"].(bool); ok {
        data.EnableAutomaticImprovements = types.BoolValue(val)
    } else if dataMap["enableAutomaticImprovements"] == nil {
        data.EnableAutomaticImprovements = types.BoolNull()
    }
    if val, ok := dataMap["maxOpenPullRequests"].(float64); ok {
        data.MaxOpenPullRequests = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["maxOpenPullRequests"].(int); ok {
        data.MaxOpenPullRequests = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["maxOpenPullRequests"].(int64); ok {
        data.MaxOpenPullRequests = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["maxOpenPullRequests"] == nil {
        data.MaxOpenPullRequests = types.NumberNull()
    }
    if val, ok := dataMap["restrictedImprovementActions"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.RestrictedImprovementActions = types.StringValue(string(jsonBytes))
        } else {
            data.RestrictedImprovementActions = types.StringNull()
        }
    } else if val, ok := dataMap["restrictedImprovementActions"].(string); ok && val != "" {
        data.RestrictedImprovementActions = types.StringValue(val)
    } else {
        data.RestrictedImprovementActions = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
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

func (r *ServiceCodeRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ServiceCodeRepositoryResourceModel
    var state ServiceCodeRepositoryResourceModel

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
    serviceCodeRepositoryRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := serviceCodeRepositoryRequest["data"].(map[string]interface{})

    if !data.ServicePathInRepository.IsUnknown() && !state.ServicePathInRepository.IsUnknown() && !data.ServicePathInRepository.Equal(state.ServicePathInRepository) {
        requestDataMap["servicePathInRepository"] = data.ServicePathInRepository.ValueString()
    }
    if !data.EnableAutomaticImprovements.IsUnknown() && !state.EnableAutomaticImprovements.IsUnknown() && !data.EnableAutomaticImprovements.Equal(state.EnableAutomaticImprovements) {
        requestDataMap["enableAutomaticImprovements"] = data.EnableAutomaticImprovements.ValueBool()
    }
    if !data.MaxOpenPullRequests.IsUnknown() && !state.MaxOpenPullRequests.IsUnknown() && !data.MaxOpenPullRequests.Equal(state.MaxOpenPullRequests) {
        requestDataMap["maxOpenPullRequests"] = data.MaxOpenPullRequests.ValueBigFloat()
    }
    if !data.RestrictedImprovementActions.IsUnknown() && !state.RestrictedImprovementActions.IsUnknown() && !data.RestrictedImprovementActions.Equal(state.RestrictedImprovementActions) {
        var restrictedimprovementactionsData interface{}
        if err := json.Unmarshal([]byte(data.RestrictedImprovementActions.ValueString()), &restrictedimprovementactionsData); err == nil {
            requestDataMap["restrictedImprovementActions"] = restrictedimprovementactionsData
        }
    }

    // Make API call
    httpResp, err := r.client.Put("/service-code-repository/" + data.Id.ValueString() + "", serviceCodeRepositoryRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update service_code_repository, got error: %s", err))
        return
    }

    // Parse the update response
    var serviceCodeRepositoryResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &serviceCodeRepositoryResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_code_repository response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceId": true,
        "codeRepositoryId": true,
        "servicePathInRepository": true,
        "enableAutomaticImprovements": true,
        "maxOpenPullRequests": true,
        "restrictedImprovementActions": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/service-code-repository/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service_code_repository after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service_code_repository read response, got error: %s", err))
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

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
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
    if obj, ok := dataMap["serviceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceId"].(string); ok && val != "" {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := dataMap["codeRepositoryId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CodeRepositoryId = types.StringValue(val)
        } else {
            data.CodeRepositoryId = types.StringNull()
        }
    } else if val, ok := dataMap["codeRepositoryId"].(string); ok && val != "" {
        data.CodeRepositoryId = types.StringValue(val)
    } else {
        data.CodeRepositoryId = types.StringNull()
    }
    if obj, ok := dataMap["servicePathInRepository"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServicePathInRepository = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.ServicePathInRepository = types.StringValue(val)
        } else {
            data.ServicePathInRepository = types.StringNull()
        }
    } else if val, ok := dataMap["servicePathInRepository"].(string); ok && val != "" {
        data.ServicePathInRepository = types.StringValue(val)
    } else {
        data.ServicePathInRepository = types.StringNull()
    }
    if val, ok := dataMap["enableAutomaticImprovements"].(bool); ok {
        data.EnableAutomaticImprovements = types.BoolValue(val)
    } else if dataMap["enableAutomaticImprovements"] == nil {
        data.EnableAutomaticImprovements = types.BoolNull()
    }
    if val, ok := dataMap["maxOpenPullRequests"].(float64); ok {
        data.MaxOpenPullRequests = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["maxOpenPullRequests"].(int); ok {
        data.MaxOpenPullRequests = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["maxOpenPullRequests"].(int64); ok {
        data.MaxOpenPullRequests = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["maxOpenPullRequests"] == nil {
        data.MaxOpenPullRequests = types.NumberNull()
    }
    if val, ok := dataMap["restrictedImprovementActions"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.RestrictedImprovementActions = types.StringValue(string(jsonBytes))
        } else {
            data.RestrictedImprovementActions = types.StringNull()
        }
    } else if val, ok := dataMap["restrictedImprovementActions"].(string); ok && val != "" {
        data.RestrictedImprovementActions = types.StringValue(val)
    } else {
        data.RestrictedImprovementActions = types.StringNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if jsonBytes, err := json.Marshal(val); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
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

func (r *ServiceCodeRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ServiceCodeRepositoryResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/service-code-repository/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service_code_repository, got error: %s", err))
        return
    }
}


func (r *ServiceCodeRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *ServiceCodeRepositoryResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *ServiceCodeRepositoryResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)
    
    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to parse JSON field for complex objects
func (r *ServiceCodeRepositoryResource) parseJSONField(terraformString types.String) interface{} {
    if terraformString.IsNull() || terraformString.IsUnknown() || terraformString.ValueString() == "" {
        return nil
    }
    
    var result interface{}
    if err := json.Unmarshal([]byte(terraformString.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return terraformString.ValueString()
    }
    
    return result
}
