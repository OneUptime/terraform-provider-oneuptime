package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SpanDataDataSource{}

func NewSpanDataDataSource() datasource.DataSource {
    return &SpanDataDataSource{}
}

// SpanDataDataSource defines the data source implementation.
type SpanDataDataSource struct {
    client *Client
}

// SpanDataDataSourceModel describes the data source data model.
type SpanDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    StartTime types.String `tfsdk:"start_time"`
    EndTime types.String `tfsdk:"end_time"`
    StartTimeUnixNano types.String `tfsdk:"start_time_unix_nano"`
    DurationUnixNano types.Number `tfsdk:"duration_unix_nano"`
    EndTimeUnixNano types.String `tfsdk:"end_time_unix_nano"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
    ParentSpanId types.String `tfsdk:"parent_span_id"`
    TraceState types.String `tfsdk:"trace_state"`
    Attributes types.String `tfsdk:"attributes"`
    AttributeKeys types.Set `tfsdk:"attribute_keys"`
    EntityKeys types.Set `tfsdk:"entity_keys"`
    ServiceEntityKey types.String `tfsdk:"service_entity_key"`
    HostEntityKey types.String `tfsdk:"host_entity_key"`
    K8sPodEntityKey types.String `tfsdk:"k8s_pod_entity_key"`
    K8sNodeEntityKey types.String `tfsdk:"k8s_node_entity_key"`
    K8sClusterEntityKey types.String `tfsdk:"k8s_cluster_entity_key"`
    ContainerEntityKey types.String `tfsdk:"container_entity_key"`
    Events types.Set `tfsdk:"events"`
    Links types.String `tfsdk:"links"`
    StatusCode types.Number `tfsdk:"status_code"`
    StatusMessage types.String `tfsdk:"status_message"`
    Kind types.String `tfsdk:"kind"`
    HasException types.Bool `tfsdk:"has_exception"`
    IsRootSpan types.Bool `tfsdk:"is_root_span"`
    IsLlmSpan types.Bool `tfsdk:"is_llm_span"`
    LlmSystem types.String `tfsdk:"llm_system"`
    LlmOperation types.String `tfsdk:"llm_operation"`
    LlmRequestModel types.String `tfsdk:"llm_request_model"`
    LlmResponseModel types.String `tfsdk:"llm_response_model"`
    LlmAgentName types.String `tfsdk:"llm_agent_name"`
    LlmToolName types.String `tfsdk:"llm_tool_name"`
    LlmInputTokens types.Number `tfsdk:"llm_input_tokens"`
    LlmOutputTokens types.Number `tfsdk:"llm_output_tokens"`
    LlmTotalTokens types.Number `tfsdk:"llm_total_tokens"`
    LlmCost types.Number `tfsdk:"llm_cost"`
}

func (d *SpanDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_span_data"
}

func (d *SpanDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "span_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID",
                Computed: true,
            },
            "primary_entity_id": schema.StringAttribute{
                MarkdownDescription: "Service ID",
                Computed: true,
            },
            "primary_entity_type": schema.StringAttribute{
                MarkdownDescription: "Service Type",
                Computed: true,
            },
            "start_time": schema.StringAttribute{
                MarkdownDescription: "Start Time",
                Computed: true,
            },
            "end_time": schema.StringAttribute{
                MarkdownDescription: "End Time",
                Computed: true,
            },
            "start_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "Start Time in Unix Nano",
                Computed: true,
            },
            "duration_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Duration in Unix Nano",
                Computed: true,
            },
            "end_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "End Time",
                Computed: true,
            },
            "trace_id": schema.StringAttribute{
                MarkdownDescription: "Trace ID",
                Computed: true,
            },
            "span_id": schema.StringAttribute{
                MarkdownDescription: "Span ID",
                Computed: true,
            },
            "parent_span_id": schema.StringAttribute{
                MarkdownDescription: "Parent Span ID",
                Computed: true,
            },
            "trace_state": schema.StringAttribute{
                MarkdownDescription: "Trace State",
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
            "entity_keys": schema.SetAttribute{
                MarkdownDescription: "Entity Keys",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_entity_key": schema.StringAttribute{
                MarkdownDescription: "Service Entity Key",
                Computed: true,
            },
            "host_entity_key": schema.StringAttribute{
                MarkdownDescription: "Host Entity Key",
                Computed: true,
            },
            "k8s_pod_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Pod Entity Key",
                Computed: true,
            },
            "k8s_node_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Node Entity Key",
                Computed: true,
            },
            "k8s_cluster_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Cluster Entity Key",
                Computed: true,
            },
            "container_entity_key": schema.StringAttribute{
                MarkdownDescription: "Container Entity Key",
                Computed: true,
            },
            "events": schema.SetAttribute{
                MarkdownDescription: "Events",
                Computed: true,
                ElementType: types.StringType,
            },
            "links": schema.StringAttribute{
                MarkdownDescription: "Links",
                Computed: true,
            },
            "status_code": schema.NumberAttribute{
                MarkdownDescription: "Status Code",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message",
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "Kind",
                Computed: true,
            },
            "has_exception": schema.BoolAttribute{
                MarkdownDescription: "Has Exception",
                Computed: true,
            },
            "is_root_span": schema.BoolAttribute{
                MarkdownDescription: "Is Root Span",
                Computed: true,
            },
            "is_llm_span": schema.BoolAttribute{
                MarkdownDescription: "Is LLM Span",
                Computed: true,
            },
            "llm_system": schema.StringAttribute{
                MarkdownDescription: "LLM System",
                Computed: true,
            },
            "llm_operation": schema.StringAttribute{
                MarkdownDescription: "LLM Operation",
                Computed: true,
            },
            "llm_request_model": schema.StringAttribute{
                MarkdownDescription: "LLM Request Model",
                Computed: true,
            },
            "llm_response_model": schema.StringAttribute{
                MarkdownDescription: "LLM Response Model",
                Computed: true,
            },
            "llm_agent_name": schema.StringAttribute{
                MarkdownDescription: "LLM Agent Name",
                Computed: true,
            },
            "llm_tool_name": schema.StringAttribute{
                MarkdownDescription: "LLM Tool Name",
                Computed: true,
            },
            "llm_input_tokens": schema.NumberAttribute{
                MarkdownDescription: "LLM Input Tokens",
                Computed: true,
            },
            "llm_output_tokens": schema.NumberAttribute{
                MarkdownDescription: "LLM Output Tokens",
                Computed: true,
            },
            "llm_total_tokens": schema.NumberAttribute{
                MarkdownDescription: "LLM Total Tokens",
                Computed: true,
            },
            "llm_cost": schema.NumberAttribute{
                MarkdownDescription: "LLM Cost (USD)",
                Computed: true,
            },
        },
    }
}

func (d *SpanDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SpanDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SpanDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "span" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read span_data, got error: %s", err))
        return
    }

    var spanDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &spanDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse span_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := spanDataResponse["data"].(map[string]interface{}); ok {
        spanDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := spanDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := spanDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := spanDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := spanDataResponse["primary_entity_id"].(string); ok {
        data.PrimaryEntityId = types.StringValue(val)
    }
    if val, ok := spanDataResponse["primary_entity_type"].(string); ok {
        data.PrimaryEntityType = types.StringValue(val)
    }
    if val, ok := spanDataResponse["start_time"].(string); ok {
        data.StartTime = types.StringValue(val)
    }
    if val, ok := spanDataResponse["end_time"].(string); ok {
        data.EndTime = types.StringValue(val)
    }
    if val, ok := spanDataResponse["start_time_unix_nano"].(string); ok {
        data.StartTimeUnixNano = types.StringValue(val)
    }
    if val, ok := spanDataResponse["duration_unix_nano"].(float64); ok {
        data.DurationUnixNano = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := spanDataResponse["end_time_unix_nano"].(string); ok {
        data.EndTimeUnixNano = types.StringValue(val)
    }
    if val, ok := spanDataResponse["trace_id"].(string); ok {
        data.TraceId = types.StringValue(val)
    }
    if val, ok := spanDataResponse["span_id"].(string); ok {
        data.SpanId = types.StringValue(val)
    }
    if val, ok := spanDataResponse["parent_span_id"].(string); ok {
        data.ParentSpanId = types.StringValue(val)
    }
    if val, ok := spanDataResponse["trace_state"].(string); ok {
        data.TraceState = types.StringValue(val)
    }
    if val, ok := spanDataResponse["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    }
    if val, ok := spanDataResponse["attribute_keys"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.AttributeKeys = setValue
    }
    if val, ok := spanDataResponse["entity_keys"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.EntityKeys = setValue
    }
    if val, ok := spanDataResponse["service_entity_key"].(string); ok {
        data.ServiceEntityKey = types.StringValue(val)
    }
    if val, ok := spanDataResponse["host_entity_key"].(string); ok {
        data.HostEntityKey = types.StringValue(val)
    }
    if val, ok := spanDataResponse["k8s_pod_entity_key"].(string); ok {
        data.K8sPodEntityKey = types.StringValue(val)
    }
    if val, ok := spanDataResponse["k8s_node_entity_key"].(string); ok {
        data.K8sNodeEntityKey = types.StringValue(val)
    }
    if val, ok := spanDataResponse["k8s_cluster_entity_key"].(string); ok {
        data.K8sClusterEntityKey = types.StringValue(val)
    }
    if val, ok := spanDataResponse["container_entity_key"].(string); ok {
        data.ContainerEntityKey = types.StringValue(val)
    }
    if val, ok := spanDataResponse["events"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Events = setValue
    }
    if val, ok := spanDataResponse["links"].(string); ok {
        data.Links = types.StringValue(val)
    }
    if val, ok := spanDataResponse["status_code"].(float64); ok {
        data.StatusCode = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := spanDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := spanDataResponse["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    }
    if val, ok := spanDataResponse["has_exception"].(bool); ok {
        data.HasException = types.BoolValue(val)
    }
    if val, ok := spanDataResponse["is_root_span"].(bool); ok {
        data.IsRootSpan = types.BoolValue(val)
    }
    if val, ok := spanDataResponse["is_llm_span"].(bool); ok {
        data.IsLlmSpan = types.BoolValue(val)
    }
    if val, ok := spanDataResponse["llm_system"].(string); ok {
        data.LlmSystem = types.StringValue(val)
    }
    if val, ok := spanDataResponse["llm_operation"].(string); ok {
        data.LlmOperation = types.StringValue(val)
    }
    if val, ok := spanDataResponse["llm_request_model"].(string); ok {
        data.LlmRequestModel = types.StringValue(val)
    }
    if val, ok := spanDataResponse["llm_response_model"].(string); ok {
        data.LlmResponseModel = types.StringValue(val)
    }
    if val, ok := spanDataResponse["llm_agent_name"].(string); ok {
        data.LlmAgentName = types.StringValue(val)
    }
    if val, ok := spanDataResponse["llm_tool_name"].(string); ok {
        data.LlmToolName = types.StringValue(val)
    }
    if val, ok := spanDataResponse["llm_input_tokens"].(float64); ok {
        data.LlmInputTokens = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := spanDataResponse["llm_output_tokens"].(float64); ok {
        data.LlmOutputTokens = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := spanDataResponse["llm_total_tokens"].(float64); ok {
        data.LlmTotalTokens = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := spanDataResponse["llm_cost"].(float64); ok {
        data.LlmCost = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
