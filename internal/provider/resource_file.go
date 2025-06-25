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
var _ resource.Resource = &FileResource{}
var _ resource.ResourceWithImportState = &FileResource{}

func NewFileResource() resource.Resource {
    return &FileResource{}
}

// FileResource defines the resource implementation.
type FileResource struct {
    client *Client
}

// FileResourceModel describes the resource data model.
type FileResourceModel struct {
    Id types.String `tfsdk:"id"`
    File types.String `tfsdk:"file"`
    Name types.String `tfsdk:"name"`
    FileType types.String `tfsdk:"file_type"`
    IsPublic types.String `tfsdk:"is_public"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
}

func (r *FileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *FileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "file resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "file": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Logged in User], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object. Permissions - Create: [Logged in User], Read: [Logged in User], Update: [No access - you don't have permission for this operation]",
                Required: true,
            },
            "file_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Logged in User], Update: [No access - you don't have permission for this operation]",
                Required: true,
            },
            "is_public": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Logged in User], Update: [No access - you don't have permission for this operation]",
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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Logged in User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *FileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *FileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data FileResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    fileRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "file": data.File.ValueString(),
        "name": data.Name.ValueString(),
        "fileType": data.FileType.ValueString(),
        "isPublic": data.IsPublic.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/file", fileRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create file, got error: %s", err))
        return
    }

    var fileResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &fileResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse file response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := fileResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = fileResponse
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
    if obj, ok := dataMap["file"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.File = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.File = types.StringValue(val)
        } else {
            data.File = types.StringNull()
        }
    } else if val, ok := dataMap["file"].(string); ok && val != "" {
        data.File = types.StringValue(val)
    } else {
        data.File = types.StringNull()
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
    if obj, ok := dataMap["fileType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FileType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.FileType = types.StringValue(val)
        } else {
            data.FileType = types.StringNull()
        }
    } else if val, ok := dataMap["fileType"].(string); ok && val != "" {
        data.FileType = types.StringValue(val)
    } else {
        data.FileType = types.StringNull()
    }
    if obj, ok := dataMap["isPublic"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IsPublic = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IsPublic = types.StringValue(val)
        } else {
            data.IsPublic = types.StringNull()
        }
    } else if val, ok := dataMap["isPublic"].(string); ok && val != "" {
        data.IsPublic = types.StringValue(val)
    } else {
        data.IsPublic = types.StringNull()
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
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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

func (r *FileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data FileResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "file": true,
        "name": true,
        "fileType": true,
        "isPublic": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/file/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read file, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var fileResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &fileResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse file response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := fileResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = fileResponse
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
    if obj, ok := dataMap["file"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.File = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.File = types.StringValue(val)
        } else {
            data.File = types.StringNull()
        }
    } else if val, ok := dataMap["file"].(string); ok && val != "" {
        data.File = types.StringValue(val)
    } else {
        data.File = types.StringNull()
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
    if obj, ok := dataMap["fileType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FileType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.FileType = types.StringValue(val)
        } else {
            data.FileType = types.StringNull()
        }
    } else if val, ok := dataMap["fileType"].(string); ok && val != "" {
        data.FileType = types.StringValue(val)
    } else {
        data.FileType = types.StringNull()
    }
    if obj, ok := dataMap["isPublic"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IsPublic = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.IsPublic = types.StringValue(val)
        } else {
            data.IsPublic = types.StringNull()
        }
    } else if val, ok := dataMap["isPublic"].(string); ok && val != "" {
        data.IsPublic = types.StringValue(val)
    } else {
        data.IsPublic = types.StringNull()
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
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    resp.Diagnostics.AddError(
        "Update Not Implemented",
        "This resource does not support update operations",
    )
}

func (r *FileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    resp.Diagnostics.AddError(
        "Delete Not Implemented",
        "This resource does not support delete operations", 
    )
}


func (r *FileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *FileResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *FileResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *FileResource) parseJSONField(terraformString types.String) interface{} {
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
