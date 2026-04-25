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
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &LogResource{}
var _ resource.ResourceWithImportState = &LogResource{}

func NewLogResource() resource.Resource {
    return &LogResource{}
}

// LogResource defines the resource implementation.
type LogResource struct {
    client *Client
}

// LogResourceModel describes the resource data model.
type LogResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
    Time types.String `tfsdk:"time"`
    TimeUnixNano types.Number `tfsdk:"time_unix_nano"`
    SeverityText types.String `tfsdk:"severity_text"`
    SeverityNumber types.Number `tfsdk:"severity_number"`
    Attributes types.String `tfsdk:"attributes"`
    AttributeKeys types.Set `tfsdk:"attribute_keys"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
    Body types.String `tfsdk:"body"`
    ObservedTimeUnixNano types.Number `tfsdk:"observed_time_unix_nano"`
    DroppedAttributesCount types.Number `tfsdk:"dropped_attributes_count"`
    Flags types.Number `tfsdk:"flags"`
}

func (r *LogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_log"
}

func (r *LogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "log resource",

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
                MarkdownDescription: "Project ID",
                Computed: true,
            },
            "service_id": schema.StringAttribute{
                MarkdownDescription: "Service ID",
                Computed: true,
            },
            "time": schema.StringAttribute{
                MarkdownDescription: "Time",
                Computed: true,
            },
            "time_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Time (in Unix Nano)",
                Computed: true,
            },
            "severity_text": schema.StringAttribute{
                MarkdownDescription: "Severity Text",
                Computed: true,
            },
            "severity_number": schema.NumberAttribute{
                MarkdownDescription: "Severity Number",
                Computed: true,
            },
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
            },
            "attribute_keys": schema.SetAttribute{
                MarkdownDescription: "Attribute Keys",
                Computed: true,
                ElementType: types.StringType,
            },
            "trace_id": schema.StringAttribute{
                MarkdownDescription: "Trace ID",
                Computed: true,
            },
            "span_id": schema.StringAttribute{
                MarkdownDescription: "Span ID",
                Computed: true,
            },
            "body": schema.StringAttribute{
                MarkdownDescription: "Log Body",
                Computed: true,
            },
            "observed_time_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Observed Time (in Unix Nano)",
                Computed: true,
            },
            "dropped_attributes_count": schema.NumberAttribute{
                MarkdownDescription: "Dropped Attributes Count",
                Computed: true,
            },
            "flags": schema.NumberAttribute{
                MarkdownDescription: "Flags",
                Computed: true,
            },
        },
    }
}

func (r *LogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *LogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data LogResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    logRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Post("/logs", logRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create log, got error: %s", err))
        return
    }

    var logResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &logResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := logResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = logResponse
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
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
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
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceId"].(string); ok && val != "" {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := dataMap["time"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Time = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Time = types.StringValue(string(jsonBytes))
            } else {
                data.Time = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Time = types.StringValue(string(jsonBytes))
            } else {
                data.Time = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Time = types.StringValue(string(jsonBytes))
        } else {
            data.Time = types.StringNull()
        }
    } else if val, ok := dataMap["time"].(string); ok && val != "" {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if val, ok := dataMap["timeUnixNano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["timeUnixNano"].(int); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["timeUnixNano"].(int64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["timeUnixNano"] == nil {
        data.TimeUnixNano = types.NumberNull()
    }
    if obj, ok := dataMap["severityText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SeverityText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SeverityText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SeverityText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SeverityText = types.StringValue(string(jsonBytes))
            } else {
                data.SeverityText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SeverityText = types.StringValue(string(jsonBytes))
            } else {
                data.SeverityText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SeverityText = types.StringValue(string(jsonBytes))
        } else {
            data.SeverityText = types.StringNull()
        }
    } else if val, ok := dataMap["severityText"].(string); ok && val != "" {
        data.SeverityText = types.StringValue(val)
    } else {
        data.SeverityText = types.StringNull()
    }
    if val, ok := dataMap["severityNumber"].(float64); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["severityNumber"].(int); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["severityNumber"].(int64); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["severityNumber"] == nil {
        data.SeverityNumber = types.NumberNull()
    }
    if obj, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Attributes = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Attributes = types.StringValue(string(jsonBytes))
            } else {
                data.Attributes = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Attributes = types.StringValue(string(jsonBytes))
            } else {
                data.Attributes = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Attributes = types.StringValue(string(jsonBytes))
        } else {
            data.Attributes = types.StringNull()
        }
    } else if val, ok := dataMap["attributes"].(string); ok && val != "" {
        data.Attributes = types.StringValue(val)
    } else {
        data.Attributes = types.StringNull()
    }
    if val, ok := dataMap["attributeKeys"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AttributeKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AttributeKeys = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["traceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TraceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TraceId = types.StringValue(string(jsonBytes))
            } else {
                data.TraceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TraceId = types.StringValue(string(jsonBytes))
            } else {
                data.TraceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TraceId = types.StringValue(string(jsonBytes))
        } else {
            data.TraceId = types.StringNull()
        }
    } else if val, ok := dataMap["traceId"].(string); ok && val != "" {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
    }
    if obj, ok := dataMap["spanId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SpanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SpanId = types.StringValue(string(jsonBytes))
            } else {
                data.SpanId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SpanId = types.StringValue(string(jsonBytes))
            } else {
                data.SpanId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SpanId = types.StringValue(string(jsonBytes))
        } else {
            data.SpanId = types.StringNull()
        }
    } else if val, ok := dataMap["spanId"].(string); ok && val != "" {
        data.SpanId = types.StringValue(val)
    } else {
        data.SpanId = types.StringNull()
    }
    if obj, ok := dataMap["body"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Body = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Body = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Body = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Body = types.StringValue(string(jsonBytes))
            } else {
                data.Body = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Body = types.StringValue(string(jsonBytes))
            } else {
                data.Body = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Body = types.StringValue(string(jsonBytes))
        } else {
            data.Body = types.StringNull()
        }
    } else if val, ok := dataMap["body"].(string); ok && val != "" {
        data.Body = types.StringValue(val)
    } else {
        data.Body = types.StringNull()
    }
    if val, ok := dataMap["observedTimeUnixNano"].(float64); ok {
        data.ObservedTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["observedTimeUnixNano"].(int); ok {
        data.ObservedTimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["observedTimeUnixNano"].(int64); ok {
        data.ObservedTimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["observedTimeUnixNano"] == nil {
        data.ObservedTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["droppedAttributesCount"].(float64); ok {
        data.DroppedAttributesCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["droppedAttributesCount"].(int); ok {
        data.DroppedAttributesCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["droppedAttributesCount"].(int64); ok {
        data.DroppedAttributesCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["droppedAttributesCount"] == nil {
        data.DroppedAttributesCount = types.NumberNull()
    }
    if val, ok := dataMap["flags"].(float64); ok {
        data.Flags = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["flags"].(int); ok {
        data.Flags = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["flags"].(int64); ok {
        data.Flags = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["flags"] == nil {
        data.Flags = types.NumberNull()
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

func (r *LogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data LogResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceId": true,
        "time": true,
        "timeUnixNano": true,
        "severityText": true,
        "severityNumber": true,
        "attributes": true,
        "attributeKeys": true,
        "traceId": true,
        "spanId": true,
        "body": true,
        "observedTimeUnixNano": true,
        "droppedAttributesCount": true,
        "flags": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/logs/" + data.Id.ValueString() + "", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read log, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var logResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &logResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := logResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = logResponse
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
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
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
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceId"].(string); ok && val != "" {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := dataMap["time"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Time = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Time = types.StringValue(string(jsonBytes))
            } else {
                data.Time = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Time = types.StringValue(string(jsonBytes))
            } else {
                data.Time = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Time = types.StringValue(string(jsonBytes))
        } else {
            data.Time = types.StringNull()
        }
    } else if val, ok := dataMap["time"].(string); ok && val != "" {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if val, ok := dataMap["timeUnixNano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["timeUnixNano"].(int); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["timeUnixNano"].(int64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["timeUnixNano"] == nil {
        data.TimeUnixNano = types.NumberNull()
    }
    if obj, ok := dataMap["severityText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SeverityText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SeverityText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SeverityText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SeverityText = types.StringValue(string(jsonBytes))
            } else {
                data.SeverityText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SeverityText = types.StringValue(string(jsonBytes))
            } else {
                data.SeverityText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SeverityText = types.StringValue(string(jsonBytes))
        } else {
            data.SeverityText = types.StringNull()
        }
    } else if val, ok := dataMap["severityText"].(string); ok && val != "" {
        data.SeverityText = types.StringValue(val)
    } else {
        data.SeverityText = types.StringNull()
    }
    if val, ok := dataMap["severityNumber"].(float64); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["severityNumber"].(int); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["severityNumber"].(int64); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["severityNumber"] == nil {
        data.SeverityNumber = types.NumberNull()
    }
    if obj, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Attributes = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Attributes = types.StringValue(string(jsonBytes))
            } else {
                data.Attributes = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Attributes = types.StringValue(string(jsonBytes))
            } else {
                data.Attributes = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Attributes = types.StringValue(string(jsonBytes))
        } else {
            data.Attributes = types.StringNull()
        }
    } else if val, ok := dataMap["attributes"].(string); ok && val != "" {
        data.Attributes = types.StringValue(val)
    } else {
        data.Attributes = types.StringNull()
    }
    if val, ok := dataMap["attributeKeys"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AttributeKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AttributeKeys = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["traceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TraceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TraceId = types.StringValue(string(jsonBytes))
            } else {
                data.TraceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TraceId = types.StringValue(string(jsonBytes))
            } else {
                data.TraceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TraceId = types.StringValue(string(jsonBytes))
        } else {
            data.TraceId = types.StringNull()
        }
    } else if val, ok := dataMap["traceId"].(string); ok && val != "" {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
    }
    if obj, ok := dataMap["spanId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SpanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SpanId = types.StringValue(string(jsonBytes))
            } else {
                data.SpanId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SpanId = types.StringValue(string(jsonBytes))
            } else {
                data.SpanId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SpanId = types.StringValue(string(jsonBytes))
        } else {
            data.SpanId = types.StringNull()
        }
    } else if val, ok := dataMap["spanId"].(string); ok && val != "" {
        data.SpanId = types.StringValue(val)
    } else {
        data.SpanId = types.StringNull()
    }
    if obj, ok := dataMap["body"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Body = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Body = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Body = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Body = types.StringValue(string(jsonBytes))
            } else {
                data.Body = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Body = types.StringValue(string(jsonBytes))
            } else {
                data.Body = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Body = types.StringValue(string(jsonBytes))
        } else {
            data.Body = types.StringNull()
        }
    } else if val, ok := dataMap["body"].(string); ok && val != "" {
        data.Body = types.StringValue(val)
    } else {
        data.Body = types.StringNull()
    }
    if val, ok := dataMap["observedTimeUnixNano"].(float64); ok {
        data.ObservedTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["observedTimeUnixNano"].(int); ok {
        data.ObservedTimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["observedTimeUnixNano"].(int64); ok {
        data.ObservedTimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["observedTimeUnixNano"] == nil {
        data.ObservedTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["droppedAttributesCount"].(float64); ok {
        data.DroppedAttributesCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["droppedAttributesCount"].(int); ok {
        data.DroppedAttributesCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["droppedAttributesCount"].(int64); ok {
        data.DroppedAttributesCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["droppedAttributesCount"] == nil {
        data.DroppedAttributesCount = types.NumberNull()
    }
    if val, ok := dataMap["flags"].(float64); ok {
        data.Flags = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["flags"].(int); ok {
        data.Flags = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["flags"].(int64); ok {
        data.Flags = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["flags"] == nil {
        data.Flags = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data LogResourceModel
    var state LogResourceModel

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
    logRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }


    // Make API call
    httpResp, err := r.client.Put("/logs/" + data.Id.ValueString() + "", logRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update log, got error: %s", err))
        return
    }

    // Parse the update response
    var logResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &logResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse log response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "serviceId": true,
        "time": true,
        "timeUnixNano": true,
        "severityText": true,
        "severityNumber": true,
        "attributes": true,
        "attributeKeys": true,
        "traceId": true,
        "spanId": true,
        "body": true,
        "observedTimeUnixNano": true,
        "droppedAttributesCount": true,
        "flags": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/logs/" + data.Id.ValueString() + "", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read log after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse log read response, got error: %s", err))
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
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
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
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceId"].(string); ok && val != "" {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := dataMap["time"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Time = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Time = types.StringValue(string(jsonBytes))
            } else {
                data.Time = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Time = types.StringValue(string(jsonBytes))
            } else {
                data.Time = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Time = types.StringValue(string(jsonBytes))
        } else {
            data.Time = types.StringNull()
        }
    } else if val, ok := dataMap["time"].(string); ok && val != "" {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if val, ok := dataMap["timeUnixNano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["timeUnixNano"].(int); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["timeUnixNano"].(int64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["timeUnixNano"] == nil {
        data.TimeUnixNano = types.NumberNull()
    }
    if obj, ok := dataMap["severityText"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SeverityText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SeverityText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SeverityText = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SeverityText = types.StringValue(string(jsonBytes))
            } else {
                data.SeverityText = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SeverityText = types.StringValue(string(jsonBytes))
            } else {
                data.SeverityText = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SeverityText = types.StringValue(string(jsonBytes))
        } else {
            data.SeverityText = types.StringNull()
        }
    } else if val, ok := dataMap["severityText"].(string); ok && val != "" {
        data.SeverityText = types.StringValue(val)
    } else {
        data.SeverityText = types.StringNull()
    }
    if val, ok := dataMap["severityNumber"].(float64); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["severityNumber"].(int); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["severityNumber"].(int64); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["severityNumber"] == nil {
        data.SeverityNumber = types.NumberNull()
    }
    if obj, ok := dataMap["attributes"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Attributes = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Attributes = types.StringValue(string(jsonBytes))
            } else {
                data.Attributes = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Attributes = types.StringValue(string(jsonBytes))
            } else {
                data.Attributes = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Attributes = types.StringValue(string(jsonBytes))
        } else {
            data.Attributes = types.StringNull()
        }
    } else if val, ok := dataMap["attributes"].(string); ok && val != "" {
        data.Attributes = types.StringValue(val)
    } else {
        data.Attributes = types.StringNull()
    }
    if val, ok := dataMap["attributeKeys"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.AttributeKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.AttributeKeys = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["traceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TraceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TraceId = types.StringValue(string(jsonBytes))
            } else {
                data.TraceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TraceId = types.StringValue(string(jsonBytes))
            } else {
                data.TraceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TraceId = types.StringValue(string(jsonBytes))
        } else {
            data.TraceId = types.StringNull()
        }
    } else if val, ok := dataMap["traceId"].(string); ok && val != "" {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
    }
    if obj, ok := dataMap["spanId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SpanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SpanId = types.StringValue(string(jsonBytes))
            } else {
                data.SpanId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SpanId = types.StringValue(string(jsonBytes))
            } else {
                data.SpanId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SpanId = types.StringValue(string(jsonBytes))
        } else {
            data.SpanId = types.StringNull()
        }
    } else if val, ok := dataMap["spanId"].(string); ok && val != "" {
        data.SpanId = types.StringValue(val)
    } else {
        data.SpanId = types.StringNull()
    }
    if obj, ok := dataMap["body"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Body = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Body = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Body = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Body = types.StringValue(string(jsonBytes))
            } else {
                data.Body = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Body = types.StringValue(string(jsonBytes))
            } else {
                data.Body = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Body = types.StringValue(string(jsonBytes))
        } else {
            data.Body = types.StringNull()
        }
    } else if val, ok := dataMap["body"].(string); ok && val != "" {
        data.Body = types.StringValue(val)
    } else {
        data.Body = types.StringNull()
    }
    if val, ok := dataMap["observedTimeUnixNano"].(float64); ok {
        data.ObservedTimeUnixNano = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["observedTimeUnixNano"].(int); ok {
        data.ObservedTimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["observedTimeUnixNano"].(int64); ok {
        data.ObservedTimeUnixNano = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["observedTimeUnixNano"] == nil {
        data.ObservedTimeUnixNano = types.NumberNull()
    }
    if val, ok := dataMap["droppedAttributesCount"].(float64); ok {
        data.DroppedAttributesCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["droppedAttributesCount"].(int); ok {
        data.DroppedAttributesCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["droppedAttributesCount"].(int64); ok {
        data.DroppedAttributesCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["droppedAttributesCount"] == nil {
        data.DroppedAttributesCount = types.NumberNull()
    }
    if val, ok := dataMap["flags"].(float64); ok {
        data.Flags = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["flags"].(int); ok {
        data.Flags = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["flags"].(int64); ok {
        data.Flags = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["flags"] == nil {
        data.Flags = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data LogResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/logs/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete log, got error: %s", err))
        return
    }
}


func (r *LogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *LogResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *LogResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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

// Helper method to convert Terraform set to Go interface{}
func (r *LogResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }
    
    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)
    
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
func (r *LogResource) parseJSONField(terraformString types.String) interface{} {
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

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *LogResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *LogResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *LogResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType
// Only these types should be marshalled/unmarshalled as typed wrapper objects
// This list is dynamically generated from Common/Types/JSON.ts ObjectType enum
func (r *LogResource) isValidOneUptimeObjectType(typeStr string) bool {
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
        "IncludesAll": true,
        "IncludesNone": true,
        "StartsWith": true,
        "EndsWith": true,
        "NotContains": true,
        "DashboardComponent": true,
        "DashboardViewConfig": true,
    }
    return validTypes[typeStr]
}
