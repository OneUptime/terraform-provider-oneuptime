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
var _ resource.Resource = &StatusPageScimLogResource{}
var _ resource.ResourceWithImportState = &StatusPageScimLogResource{}

func NewStatusPageScimLogResource() resource.Resource {
    return &StatusPageScimLogResource{}
}

// StatusPageScimLogResource defines the resource implementation.
type StatusPageScimLogResource struct {
    client *Client
}

// StatusPageScimLogResourceModel describes the resource data model.
type StatusPageScimLogResourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    StatusPageScimId types.String `tfsdk:"status_page_scim_id"`
    OperationType types.String `tfsdk:"operation_type"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    LogBody types.String `tfsdk:"log_body"`
    HttpMethod types.String `tfsdk:"http_method"`
    RequestPath types.String `tfsdk:"request_path"`
    HttpStatusCode types.Number `tfsdk:"http_status_code"`
    AffectedUserEmail types.String `tfsdk:"affected_user_email"`
}

func (r *StatusPageScimLogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_scim_log"
}

func (r *StatusPageScimLogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_scim_log resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
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
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_page_scim_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "operation_type": schema.StringAttribute{
                MarkdownDescription: "Type of SCIM operation (e.g., CreateUser, UpdateUser, DeleteUser, ListUsers, GetUser, BulkOperation). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page SCIM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of the SCIM operation. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page SCIM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Short error or status description. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page SCIM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "log_body": schema.StringAttribute{
                MarkdownDescription: "Detailed JSON with request/response data. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page SCIM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "http_method": schema.StringAttribute{
                MarkdownDescription: "HTTP method used (GET, POST, PUT, PATCH, DELETE). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page SCIM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "request_path": schema.StringAttribute{
                MarkdownDescription: "The SCIM endpoint path. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page SCIM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "http_status_code": schema.NumberAttribute{
                MarkdownDescription: "Response HTTP status code. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page SCIM Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "affected_user_email": schema.StringAttribute{
                MarkdownDescription: "Email object",
                Computed: true,
            },
        },
    }
}

func (r *StatusPageScimLogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *StatusPageScimLogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data StatusPageScimLogResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    statusPageScimLogRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Post("/status-page-scim-log/count", statusPageScimLogRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status_page_scim_log, got error: %s", err))
        return
    }

    var statusPageScimLogResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageScimLogResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_scim_log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageScimLogResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageScimLogResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
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
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusPageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.StatusPageId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.StatusPageId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusPageId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := dataMap["statusPageScimId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageScimId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusPageScimId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusPageScimId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.StatusPageScimId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageScimId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.StatusPageScimId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageScimId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusPageScimId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageScimId = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageScimId"].(string); ok && val != "" {
        data.StatusPageScimId = types.StringValue(val)
    } else {
        data.StatusPageScimId = types.StringNull()
    }
    if obj, ok := dataMap["operationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OperationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OperationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OperationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.OperationType = types.StringValue(string(jsonBytes))
            } else {
                data.OperationType = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.OperationType = types.StringValue(string(jsonBytes))
            } else {
                data.OperationType = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OperationType = types.StringValue(string(jsonBytes))
        } else {
            data.OperationType = types.StringNull()
        }
    } else if val, ok := dataMap["operationType"].(string); ok && val != "" {
        data.OperationType = types.StringValue(val)
    } else {
        data.OperationType = types.StringNull()
    }
    if obj, ok := dataMap["status"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := dataMap["statusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := dataMap["logBody"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogBody = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LogBody = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LogBody = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.LogBody = types.StringValue(string(jsonBytes))
            } else {
                data.LogBody = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.LogBody = types.StringValue(string(jsonBytes))
            } else {
                data.LogBody = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LogBody = types.StringValue(string(jsonBytes))
        } else {
            data.LogBody = types.StringNull()
        }
    } else if val, ok := dataMap["logBody"].(string); ok && val != "" {
        data.LogBody = types.StringValue(val)
    } else {
        data.LogBody = types.StringNull()
    }
    if obj, ok := dataMap["httpMethod"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HttpMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HttpMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HttpMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.HttpMethod = types.StringValue(string(jsonBytes))
            } else {
                data.HttpMethod = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.HttpMethod = types.StringValue(string(jsonBytes))
            } else {
                data.HttpMethod = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HttpMethod = types.StringValue(string(jsonBytes))
        } else {
            data.HttpMethod = types.StringNull()
        }
    } else if val, ok := dataMap["httpMethod"].(string); ok && val != "" {
        data.HttpMethod = types.StringValue(val)
    } else {
        data.HttpMethod = types.StringNull()
    }
    if obj, ok := dataMap["requestPath"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequestPath = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RequestPath = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RequestPath = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.RequestPath = types.StringValue(string(jsonBytes))
            } else {
                data.RequestPath = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.RequestPath = types.StringValue(string(jsonBytes))
            } else {
                data.RequestPath = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RequestPath = types.StringValue(string(jsonBytes))
        } else {
            data.RequestPath = types.StringNull()
        }
    } else if val, ok := dataMap["requestPath"].(string); ok && val != "" {
        data.RequestPath = types.StringValue(val)
    } else {
        data.RequestPath = types.StringNull()
    }
    if val, ok := dataMap["httpStatusCode"].(float64); ok {
        data.HttpStatusCode = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["httpStatusCode"].(int); ok {
        data.HttpStatusCode = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["httpStatusCode"].(int64); ok {
        data.HttpStatusCode = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["httpStatusCode"] == nil {
        data.HttpStatusCode = types.NumberNull()
    }
    if obj, ok := dataMap["affectedUserEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AffectedUserEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AffectedUserEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AffectedUserEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.AffectedUserEmail = types.StringValue(string(jsonBytes))
            } else {
                data.AffectedUserEmail = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.AffectedUserEmail = types.StringValue(string(jsonBytes))
            } else {
                data.AffectedUserEmail = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AffectedUserEmail = types.StringValue(string(jsonBytes))
        } else {
            data.AffectedUserEmail = types.StringNull()
        }
    } else if val, ok := dataMap["affectedUserEmail"].(string); ok && val != "" {
        data.AffectedUserEmail = types.StringValue(val)
    } else {
        data.AffectedUserEmail = types.StringNull()
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

func (r *StatusPageScimLogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data StatusPageScimLogResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "statusPageId": true,
        "statusPageScimId": true,
        "operationType": true,
        "status": true,
        "statusMessage": true,
        "logBody": true,
        "httpMethod": true,
        "requestPath": true,
        "httpStatusCode": true,
        "affectedUserEmail": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/status-page-scim-log/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_scim_log, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var statusPageScimLogResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageScimLogResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_scim_log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageScimLogResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageScimLogResponse
    }

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
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
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusPageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.StatusPageId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.StatusPageId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusPageId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := dataMap["statusPageScimId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageScimId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusPageScimId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusPageScimId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.StatusPageScimId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageScimId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.StatusPageScimId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageScimId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusPageScimId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageScimId = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageScimId"].(string); ok && val != "" {
        data.StatusPageScimId = types.StringValue(val)
    } else {
        data.StatusPageScimId = types.StringNull()
    }
    if obj, ok := dataMap["operationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OperationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OperationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OperationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.OperationType = types.StringValue(string(jsonBytes))
            } else {
                data.OperationType = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.OperationType = types.StringValue(string(jsonBytes))
            } else {
                data.OperationType = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OperationType = types.StringValue(string(jsonBytes))
        } else {
            data.OperationType = types.StringNull()
        }
    } else if val, ok := dataMap["operationType"].(string); ok && val != "" {
        data.OperationType = types.StringValue(val)
    } else {
        data.OperationType = types.StringNull()
    }
    if obj, ok := dataMap["status"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Status = types.StringValue(string(jsonBytes))
            } else {
                data.Status = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := dataMap["statusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.StatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.StatusMessage = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := dataMap["logBody"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogBody = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LogBody = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LogBody = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.LogBody = types.StringValue(string(jsonBytes))
            } else {
                data.LogBody = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.LogBody = types.StringValue(string(jsonBytes))
            } else {
                data.LogBody = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LogBody = types.StringValue(string(jsonBytes))
        } else {
            data.LogBody = types.StringNull()
        }
    } else if val, ok := dataMap["logBody"].(string); ok && val != "" {
        data.LogBody = types.StringValue(val)
    } else {
        data.LogBody = types.StringNull()
    }
    if obj, ok := dataMap["httpMethod"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HttpMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HttpMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HttpMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.HttpMethod = types.StringValue(string(jsonBytes))
            } else {
                data.HttpMethod = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.HttpMethod = types.StringValue(string(jsonBytes))
            } else {
                data.HttpMethod = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HttpMethod = types.StringValue(string(jsonBytes))
        } else {
            data.HttpMethod = types.StringNull()
        }
    } else if val, ok := dataMap["httpMethod"].(string); ok && val != "" {
        data.HttpMethod = types.StringValue(val)
    } else {
        data.HttpMethod = types.StringNull()
    }
    if obj, ok := dataMap["requestPath"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequestPath = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RequestPath = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RequestPath = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.RequestPath = types.StringValue(string(jsonBytes))
            } else {
                data.RequestPath = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.RequestPath = types.StringValue(string(jsonBytes))
            } else {
                data.RequestPath = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RequestPath = types.StringValue(string(jsonBytes))
        } else {
            data.RequestPath = types.StringNull()
        }
    } else if val, ok := dataMap["requestPath"].(string); ok && val != "" {
        data.RequestPath = types.StringValue(val)
    } else {
        data.RequestPath = types.StringNull()
    }
    if val, ok := dataMap["httpStatusCode"].(float64); ok {
        data.HttpStatusCode = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["httpStatusCode"].(int); ok {
        data.HttpStatusCode = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["httpStatusCode"].(int64); ok {
        data.HttpStatusCode = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["httpStatusCode"] == nil {
        data.HttpStatusCode = types.NumberNull()
    }
    if obj, ok := dataMap["affectedUserEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AffectedUserEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AffectedUserEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AffectedUserEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.AffectedUserEmail = types.StringValue(string(jsonBytes))
            } else {
                data.AffectedUserEmail = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.AffectedUserEmail = types.StringValue(string(jsonBytes))
            } else {
                data.AffectedUserEmail = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AffectedUserEmail = types.StringValue(string(jsonBytes))
        } else {
            data.AffectedUserEmail = types.StringNull()
        }
    } else if val, ok := dataMap["affectedUserEmail"].(string); ok && val != "" {
        data.AffectedUserEmail = types.StringValue(val)
    } else {
        data.AffectedUserEmail = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageScimLogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    resp.Diagnostics.AddError(
        "Update Not Implemented",
        "This resource does not support update operations",
    )
}

func (r *StatusPageScimLogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    resp.Diagnostics.AddError(
        "Delete Not Implemented",
        "This resource does not support delete operations", 
    )
}


func (r *StatusPageScimLogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *StatusPageScimLogResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *StatusPageScimLogResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *StatusPageScimLogResource) parseJSONField(terraformString types.String) interface{} {
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

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *StatusPageScimLogResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType
// Only these types should be marshalled/unmarshalled as typed wrapper objects
// This list is dynamically generated from Common/Types/JSON.ts ObjectType enum
func (r *StatusPageScimLogResource) isValidOneUptimeObjectType(typeStr string) bool {
    validTypes := map[string]bool{
        "ObjectID": true,
        "Decimal": true,
        "Name": true,
        "EqualTo": true,
        "EqualToOrNull": true,
        "MonitorSteps": true,
        "MonitorStep": true,
        "Recurring": true,
        "RestrictionTimes": true,
        "MonitorCriteria": true,
        "PositiveNumber": true,
        "MonitorCriteriaInstance": true,
        "NotEqual": true,
        "Email": true,
        "Phone": true,
        "Color": true,
        "Domain": true,
        "Version": true,
        "IP": true,
        "Route": true,
        "URL": true,
        "Permission": true,
        "Search": true,
        "GreaterThan": true,
        "GreaterThanOrEqual": true,
        "GreaterThanOrNull": true,
        "LessThanOrNull": true,
        "LessThan": true,
        "LessThanOrEqual": true,
        "Port": true,
        "Hostname": true,
        "HashedString": true,
        "DateTime": true,
        "Buffer": true,
        "InBetween": true,
        "NotNull": true,
        "IsNull": true,
        "Includes": true,
        "DashboardComponent": true,
        "DashboardViewConfig": true,
    }
    return validTypes[typeStr]
}
