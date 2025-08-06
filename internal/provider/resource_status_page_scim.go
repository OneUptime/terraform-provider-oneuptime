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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StatusPageScimResource{}
var _ resource.ResourceWithImportState = &StatusPageScimResource{}

func NewStatusPageScimResource() resource.Resource {
    return &StatusPageScimResource{}
}

// StatusPageScimResource defines the resource implementation.
type StatusPageScimResource struct {
    client *Client
}

// StatusPageScimResourceModel describes the resource data model.
type StatusPageScimResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    BearerToken types.String `tfsdk:"bearer_token"`
    AutoProvisionUsers types.Bool `tfsdk:"auto_provision_users"`
    AutoDeprovisionUsers types.Bool `tfsdk:"auto_deprovision_users"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (r *StatusPageScimResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_scim"
}

func (r *StatusPageScimResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_scim resource",

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
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name for this SCIM configuration. Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description to help you remember. Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Optional: true,
                Computed: true,
            },
            "bearer_token": schema.StringAttribute{
                MarkdownDescription: "Bearer token for SCIM authentication. Keep this secure.. Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Required: true,
            },
            "auto_provision_users": schema.BoolAttribute{
                MarkdownDescription: "Automatically create status page users when they are added via SCIM. Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
            },
            "auto_deprovision_users": schema.BoolAttribute{
                MarkdownDescription: "Automatically remove status page users when they are removed via SCIM. Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
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

func (r *StatusPageScimResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *StatusPageScimResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data StatusPageScimResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    statusPageScimRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "statusPageId": data.StatusPageId.ValueString(),
        "name": data.Name.ValueString(),
        "description": data.Description.ValueString(),
        "bearerToken": data.BearerToken.ValueString(),
        "autoProvisionUsers": data.AutoProvisionUsers.ValueBool(),
        "autoDeprovisionUsers": data.AutoDeprovisionUsers.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/status-page-scim", statusPageScimRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status_page_scim, got error: %s", err))
        return
    }

    var statusPageScimResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageScimResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_scim response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageScimResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageScimResponse
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
    if obj, ok := dataMap["statusPageId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["bearerToken"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BearerToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.BearerToken = types.StringValue(val)
        } else {
            data.BearerToken = types.StringNull()
        }
    } else if val, ok := dataMap["bearerToken"].(string); ok && val != "" {
        data.BearerToken = types.StringValue(val)
    } else {
        data.BearerToken = types.StringNull()
    }
    if val, ok := dataMap["autoProvisionUsers"].(bool); ok {
        data.AutoProvisionUsers = types.BoolValue(val)
    }
    if val, ok := dataMap["autoDeprovisionUsers"].(bool); ok {
        data.AutoDeprovisionUsers = types.BoolValue(val)
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

func (r *StatusPageScimResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data StatusPageScimResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "statusPageId": true,
        "name": true,
        "description": true,
        "bearerToken": true,
        "autoProvisionUsers": true,
        "autoDeprovisionUsers": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/status-page-scim/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_scim, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var statusPageScimResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageScimResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_scim response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageScimResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageScimResponse
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
    if obj, ok := dataMap["statusPageId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["bearerToken"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BearerToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.BearerToken = types.StringValue(val)
        } else {
            data.BearerToken = types.StringNull()
        }
    } else if val, ok := dataMap["bearerToken"].(string); ok && val != "" {
        data.BearerToken = types.StringValue(val)
    } else {
        data.BearerToken = types.StringNull()
    }
    if val, ok := dataMap["autoProvisionUsers"].(bool); ok {
        data.AutoProvisionUsers = types.BoolValue(val)
    }
    if val, ok := dataMap["autoDeprovisionUsers"].(bool); ok {
        data.AutoDeprovisionUsers = types.BoolValue(val)
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

func (r *StatusPageScimResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data StatusPageScimResourceModel
    var state StatusPageScimResourceModel

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
    statusPageScimRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := statusPageScimRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && data.Name.ValueString() != "" {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && data.Description.ValueString() != "" {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.BearerToken.IsNull() && data.BearerToken.ValueString() != "" {
        requestDataMap["bearerToken"] = data.BearerToken.ValueString()
    }
    requestDataMap["autoProvisionUsers"] = data.AutoProvisionUsers.ValueBool()
    requestDataMap["autoDeprovisionUsers"] = data.AutoDeprovisionUsers.ValueBool()

    // Make API call
    httpResp, err := r.client.Put("/status-page-scim/" + data.Id.ValueString() + "", statusPageScimRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status_page_scim, got error: %s", err))
        return
    }

    // Parse the update response
    var statusPageScimResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageScimResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_scim response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "statusPageId": true,
        "name": true,
        "description": true,
        "bearerToken": true,
        "autoProvisionUsers": true,
        "autoDeprovisionUsers": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/status-page-scim/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_scim after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_scim read response, got error: %s", err))
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
    if obj, ok := dataMap["statusPageId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["bearerToken"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BearerToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.BearerToken = types.StringValue(val)
        } else {
            data.BearerToken = types.StringNull()
        }
    } else if val, ok := dataMap["bearerToken"].(string); ok && val != "" {
        data.BearerToken = types.StringValue(val)
    } else {
        data.BearerToken = types.StringNull()
    }
    if val, ok := dataMap["autoProvisionUsers"].(bool); ok {
        data.AutoProvisionUsers = types.BoolValue(val)
    }
    if val, ok := dataMap["autoDeprovisionUsers"].(bool); ok {
        data.AutoDeprovisionUsers = types.BoolValue(val)
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

func (r *StatusPageScimResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data StatusPageScimResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/status-page-scim/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete status_page_scim, got error: %s", err))
        return
    }
}


func (r *StatusPageScimResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *StatusPageScimResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *StatusPageScimResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *StatusPageScimResource) parseJSONField(terraformString types.String) interface{} {
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
