package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SpanDataSource{}

func NewSpanDataSource() datasource.DataSource {
    return &SpanDataSource{}
}

// SpanDataSource defines the data source implementation.
type SpanDataSource struct {
    client *Client
}

// SpanDataSourceModel describes the data source data model.
type SpanDataSourceModel struct {
    Id types.String `tfsdk:"id"`
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
    SessionId types.String `tfsdk:"session_id"`
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
    Name types.String `tfsdk:"name"`
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
    LlmConversationId types.String `tfsdk:"llm_conversation_id"`
    LlmUserId types.String `tfsdk:"llm_user_id"`
    LlmUserEmail types.String `tfsdk:"llm_user_email"`
    LlmTeam types.String `tfsdk:"llm_team"`
    LlmCallKind types.String `tfsdk:"llm_call_kind"`
    LlmIssues types.Set `tfsdk:"llm_issues"`
    LlmUserMessagePreview types.String `tfsdk:"llm_user_message_preview"`
}

func (d *SpanDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_span"
}

func (d *SpanDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for Span Look up an existing span by `id`, or by any of its other arguments (`name`, `attributes`, `container_entity_key`, ...): each one set must match, and exactly one span may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID.",
                Computed: true,
            },
            "primary_entity_id": schema.StringAttribute{
                MarkdownDescription: "Service ID.",
                Optional: true,
                Computed: true,
            },
            "primary_entity_type": schema.StringAttribute{
                MarkdownDescription: "Service Type.",
                Optional: true,
                Computed: true,
            },
            "start_time": schema.StringAttribute{
                MarkdownDescription: "Start Time.",
                Optional: true,
                Computed: true,
            },
            "end_time": schema.StringAttribute{
                MarkdownDescription: "End Time.",
                Optional: true,
                Computed: true,
            },
            "start_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "Start Time in Unix Nano.",
                Optional: true,
                Computed: true,
            },
            "duration_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Duration in Unix Nano.",
                Optional: true,
                Computed: true,
            },
            "end_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "End Time.",
                Optional: true,
                Computed: true,
            },
            "trace_id": schema.StringAttribute{
                MarkdownDescription: "Trace ID.",
                Optional: true,
                Computed: true,
            },
            "span_id": schema.StringAttribute{
                MarkdownDescription: "Span ID.",
                Optional: true,
                Computed: true,
            },
            "session_id": schema.StringAttribute{
                MarkdownDescription: "Session ID.",
                Optional: true,
                Computed: true,
            },
            "parent_span_id": schema.StringAttribute{
                MarkdownDescription: "Parent Span ID.",
                Optional: true,
                Computed: true,
            },
            "trace_state": schema.StringAttribute{
                MarkdownDescription: "Trace State.",
                Optional: true,
                Computed: true,
            },
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes.",
                Optional: true,
                Computed: true,
            },
            "attribute_keys": schema.SetAttribute{
                MarkdownDescription: "Attribute Keys.",
                Computed: true,
                ElementType: types.StringType,
            },
            "entity_keys": schema.SetAttribute{
                MarkdownDescription: "Entity Keys.",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_entity_key": schema.StringAttribute{
                MarkdownDescription: "Service Entity Key.",
                Optional: true,
                Computed: true,
            },
            "host_entity_key": schema.StringAttribute{
                MarkdownDescription: "Host Entity Key.",
                Optional: true,
                Computed: true,
            },
            "k8s_pod_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Pod Entity Key.",
                Optional: true,
                Computed: true,
            },
            "k8s_node_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Node Entity Key.",
                Optional: true,
                Computed: true,
            },
            "k8s_cluster_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Cluster Entity Key.",
                Optional: true,
                Computed: true,
            },
            "container_entity_key": schema.StringAttribute{
                MarkdownDescription: "Container Entity Key.",
                Optional: true,
                Computed: true,
            },
            "events": schema.SetAttribute{
                MarkdownDescription: "Events.",
                Computed: true,
                ElementType: types.StringType,
            },
            "links": schema.StringAttribute{
                MarkdownDescription: "Links. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "status_code": schema.NumberAttribute{
                MarkdownDescription: "Status Code.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name.",
                Optional: true,
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "Kind.",
                Optional: true,
                Computed: true,
            },
            "has_exception": schema.BoolAttribute{
                MarkdownDescription: "Has Exception.",
                Optional: true,
                Computed: true,
            },
            "is_root_span": schema.BoolAttribute{
                MarkdownDescription: "Is Root Span.",
                Optional: true,
                Computed: true,
            },
            "is_llm_span": schema.BoolAttribute{
                MarkdownDescription: "Is LLM Span.",
                Optional: true,
                Computed: true,
            },
            "llm_system": schema.StringAttribute{
                MarkdownDescription: "LLM System.",
                Optional: true,
                Computed: true,
            },
            "llm_operation": schema.StringAttribute{
                MarkdownDescription: "LLM Operation.",
                Optional: true,
                Computed: true,
            },
            "llm_request_model": schema.StringAttribute{
                MarkdownDescription: "LLM Request Model.",
                Optional: true,
                Computed: true,
            },
            "llm_response_model": schema.StringAttribute{
                MarkdownDescription: "LLM Response Model.",
                Optional: true,
                Computed: true,
            },
            "llm_agent_name": schema.StringAttribute{
                MarkdownDescription: "LLM Agent Name.",
                Optional: true,
                Computed: true,
            },
            "llm_tool_name": schema.StringAttribute{
                MarkdownDescription: "LLM Tool Name.",
                Optional: true,
                Computed: true,
            },
            "llm_input_tokens": schema.NumberAttribute{
                MarkdownDescription: "LLM Input Tokens.",
                Optional: true,
                Computed: true,
            },
            "llm_output_tokens": schema.NumberAttribute{
                MarkdownDescription: "LLM Output Tokens.",
                Optional: true,
                Computed: true,
            },
            "llm_total_tokens": schema.NumberAttribute{
                MarkdownDescription: "LLM Total Tokens.",
                Optional: true,
                Computed: true,
            },
            "llm_cost": schema.NumberAttribute{
                MarkdownDescription: "LLM Cost (USD).",
                Optional: true,
                Computed: true,
            },
            "llm_conversation_id": schema.StringAttribute{
                MarkdownDescription: "LLM Conversation ID.",
                Optional: true,
                Computed: true,
            },
            "llm_user_id": schema.StringAttribute{
                MarkdownDescription: "LLM User ID.",
                Optional: true,
                Computed: true,
            },
            "llm_user_email": schema.StringAttribute{
                MarkdownDescription: "LLM User Email.",
                Optional: true,
                Computed: true,
            },
            "llm_team": schema.StringAttribute{
                MarkdownDescription: "LLM Team.",
                Optional: true,
                Computed: true,
            },
            "llm_call_kind": schema.StringAttribute{
                MarkdownDescription: "LLM Call Kind.",
                Optional: true,
                Computed: true,
            },
            "llm_issues": schema.SetAttribute{
                MarkdownDescription: "LLM Answer Issues.",
                Computed: true,
                ElementType: types.StringType,
            },
            "llm_user_message_preview": schema.StringAttribute{
                MarkdownDescription: "LLM User Message Preview.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *SpanDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SpanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SpanDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.PrimaryEntityId.IsNull() && !data.PrimaryEntityId.IsUnknown() {
        filters["primaryEntityId"] = data.PrimaryEntityId.ValueString()
        filterNames = append(filterNames, "primary_entity_id = "+fmt.Sprintf("%q", data.PrimaryEntityId.ValueString()))
    }
    if !data.PrimaryEntityType.IsNull() && !data.PrimaryEntityType.IsUnknown() {
        filters["primaryEntityType"] = data.PrimaryEntityType.ValueString()
        filterNames = append(filterNames, "primary_entity_type = "+fmt.Sprintf("%q", data.PrimaryEntityType.ValueString()))
    }
    if !data.StartTime.IsNull() && !data.StartTime.IsUnknown() {
        filters["startTime"] = data.StartTime.ValueString()
        filterNames = append(filterNames, "start_time = "+fmt.Sprintf("%q", data.StartTime.ValueString()))
    }
    if !data.EndTime.IsNull() && !data.EndTime.IsUnknown() {
        filters["endTime"] = data.EndTime.ValueString()
        filterNames = append(filterNames, "end_time = "+fmt.Sprintf("%q", data.EndTime.ValueString()))
    }
    if !data.StartTimeUnixNano.IsNull() && !data.StartTimeUnixNano.IsUnknown() {
        filters["startTimeUnixNano"] = data.StartTimeUnixNano.ValueString()
        filterNames = append(filterNames, "start_time_unix_nano = "+fmt.Sprintf("%q", data.StartTimeUnixNano.ValueString()))
    }
    if !data.DurationUnixNano.IsNull() && !data.DurationUnixNano.IsUnknown() {
        filters["durationUnixNano"] = lookupNumber(data.DurationUnixNano)
        filterNames = append(filterNames, "duration_unix_nano = "+data.DurationUnixNano.ValueBigFloat().String())
    }
    if !data.EndTimeUnixNano.IsNull() && !data.EndTimeUnixNano.IsUnknown() {
        filters["endTimeUnixNano"] = data.EndTimeUnixNano.ValueString()
        filterNames = append(filterNames, "end_time_unix_nano = "+fmt.Sprintf("%q", data.EndTimeUnixNano.ValueString()))
    }
    if !data.TraceId.IsNull() && !data.TraceId.IsUnknown() {
        filters["traceId"] = data.TraceId.ValueString()
        filterNames = append(filterNames, "trace_id = "+fmt.Sprintf("%q", data.TraceId.ValueString()))
    }
    if !data.SpanId.IsNull() && !data.SpanId.IsUnknown() {
        filters["spanId"] = data.SpanId.ValueString()
        filterNames = append(filterNames, "span_id = "+fmt.Sprintf("%q", data.SpanId.ValueString()))
    }
    if !data.SessionId.IsNull() && !data.SessionId.IsUnknown() {
        filters["sessionId"] = data.SessionId.ValueString()
        filterNames = append(filterNames, "session_id = "+fmt.Sprintf("%q", data.SessionId.ValueString()))
    }
    if !data.ParentSpanId.IsNull() && !data.ParentSpanId.IsUnknown() {
        filters["parentSpanId"] = data.ParentSpanId.ValueString()
        filterNames = append(filterNames, "parent_span_id = "+fmt.Sprintf("%q", data.ParentSpanId.ValueString()))
    }
    if !data.TraceState.IsNull() && !data.TraceState.IsUnknown() {
        filters["traceState"] = data.TraceState.ValueString()
        filterNames = append(filterNames, "trace_state = "+fmt.Sprintf("%q", data.TraceState.ValueString()))
    }
    if !data.Attributes.IsNull() && !data.Attributes.IsUnknown() {
        filters["attributes"] = data.Attributes.ValueString()
        filterNames = append(filterNames, "attributes = "+fmt.Sprintf("%q", data.Attributes.ValueString()))
    }
    if !data.ServiceEntityKey.IsNull() && !data.ServiceEntityKey.IsUnknown() {
        filters["serviceEntityKey"] = data.ServiceEntityKey.ValueString()
        filterNames = append(filterNames, "service_entity_key = "+fmt.Sprintf("%q", data.ServiceEntityKey.ValueString()))
    }
    if !data.HostEntityKey.IsNull() && !data.HostEntityKey.IsUnknown() {
        filters["hostEntityKey"] = data.HostEntityKey.ValueString()
        filterNames = append(filterNames, "host_entity_key = "+fmt.Sprintf("%q", data.HostEntityKey.ValueString()))
    }
    if !data.K8sPodEntityKey.IsNull() && !data.K8sPodEntityKey.IsUnknown() {
        filters["k8sPodEntityKey"] = data.K8sPodEntityKey.ValueString()
        filterNames = append(filterNames, "k8s_pod_entity_key = "+fmt.Sprintf("%q", data.K8sPodEntityKey.ValueString()))
    }
    if !data.K8sNodeEntityKey.IsNull() && !data.K8sNodeEntityKey.IsUnknown() {
        filters["k8sNodeEntityKey"] = data.K8sNodeEntityKey.ValueString()
        filterNames = append(filterNames, "k8s_node_entity_key = "+fmt.Sprintf("%q", data.K8sNodeEntityKey.ValueString()))
    }
    if !data.K8sClusterEntityKey.IsNull() && !data.K8sClusterEntityKey.IsUnknown() {
        filters["k8sClusterEntityKey"] = data.K8sClusterEntityKey.ValueString()
        filterNames = append(filterNames, "k8s_cluster_entity_key = "+fmt.Sprintf("%q", data.K8sClusterEntityKey.ValueString()))
    }
    if !data.ContainerEntityKey.IsNull() && !data.ContainerEntityKey.IsUnknown() {
        filters["containerEntityKey"] = data.ContainerEntityKey.ValueString()
        filterNames = append(filterNames, "container_entity_key = "+fmt.Sprintf("%q", data.ContainerEntityKey.ValueString()))
    }
    if !data.StatusCode.IsNull() && !data.StatusCode.IsUnknown() {
        filters["statusCode"] = lookupNumber(data.StatusCode)
        filterNames = append(filterNames, "status_code = "+data.StatusCode.ValueBigFloat().String())
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        filters["kind"] = data.Kind.ValueString()
        filterNames = append(filterNames, "kind = "+fmt.Sprintf("%q", data.Kind.ValueString()))
    }
    if !data.HasException.IsNull() && !data.HasException.IsUnknown() {
        filters["hasException"] = data.HasException.ValueBool()
        filterNames = append(filterNames, "has_exception = "+fmt.Sprintf("%t", data.HasException.ValueBool()))
    }
    if !data.IsRootSpan.IsNull() && !data.IsRootSpan.IsUnknown() {
        filters["isRootSpan"] = data.IsRootSpan.ValueBool()
        filterNames = append(filterNames, "is_root_span = "+fmt.Sprintf("%t", data.IsRootSpan.ValueBool()))
    }
    if !data.IsLlmSpan.IsNull() && !data.IsLlmSpan.IsUnknown() {
        filters["isLlmSpan"] = data.IsLlmSpan.ValueBool()
        filterNames = append(filterNames, "is_llm_span = "+fmt.Sprintf("%t", data.IsLlmSpan.ValueBool()))
    }
    if !data.LlmSystem.IsNull() && !data.LlmSystem.IsUnknown() {
        filters["llmSystem"] = data.LlmSystem.ValueString()
        filterNames = append(filterNames, "llm_system = "+fmt.Sprintf("%q", data.LlmSystem.ValueString()))
    }
    if !data.LlmOperation.IsNull() && !data.LlmOperation.IsUnknown() {
        filters["llmOperation"] = data.LlmOperation.ValueString()
        filterNames = append(filterNames, "llm_operation = "+fmt.Sprintf("%q", data.LlmOperation.ValueString()))
    }
    if !data.LlmRequestModel.IsNull() && !data.LlmRequestModel.IsUnknown() {
        filters["llmRequestModel"] = data.LlmRequestModel.ValueString()
        filterNames = append(filterNames, "llm_request_model = "+fmt.Sprintf("%q", data.LlmRequestModel.ValueString()))
    }
    if !data.LlmResponseModel.IsNull() && !data.LlmResponseModel.IsUnknown() {
        filters["llmResponseModel"] = data.LlmResponseModel.ValueString()
        filterNames = append(filterNames, "llm_response_model = "+fmt.Sprintf("%q", data.LlmResponseModel.ValueString()))
    }
    if !data.LlmAgentName.IsNull() && !data.LlmAgentName.IsUnknown() {
        filters["llmAgentName"] = data.LlmAgentName.ValueString()
        filterNames = append(filterNames, "llm_agent_name = "+fmt.Sprintf("%q", data.LlmAgentName.ValueString()))
    }
    if !data.LlmToolName.IsNull() && !data.LlmToolName.IsUnknown() {
        filters["llmToolName"] = data.LlmToolName.ValueString()
        filterNames = append(filterNames, "llm_tool_name = "+fmt.Sprintf("%q", data.LlmToolName.ValueString()))
    }
    if !data.LlmInputTokens.IsNull() && !data.LlmInputTokens.IsUnknown() {
        filters["llmInputTokens"] = lookupNumber(data.LlmInputTokens)
        filterNames = append(filterNames, "llm_input_tokens = "+data.LlmInputTokens.ValueBigFloat().String())
    }
    if !data.LlmOutputTokens.IsNull() && !data.LlmOutputTokens.IsUnknown() {
        filters["llmOutputTokens"] = lookupNumber(data.LlmOutputTokens)
        filterNames = append(filterNames, "llm_output_tokens = "+data.LlmOutputTokens.ValueBigFloat().String())
    }
    if !data.LlmTotalTokens.IsNull() && !data.LlmTotalTokens.IsUnknown() {
        filters["llmTotalTokens"] = lookupNumber(data.LlmTotalTokens)
        filterNames = append(filterNames, "llm_total_tokens = "+data.LlmTotalTokens.ValueBigFloat().String())
    }
    if !data.LlmCost.IsNull() && !data.LlmCost.IsUnknown() {
        filters["llmCost"] = lookupNumber(data.LlmCost)
        filterNames = append(filterNames, "llm_cost = "+data.LlmCost.ValueBigFloat().String())
    }
    if !data.LlmConversationId.IsNull() && !data.LlmConversationId.IsUnknown() {
        filters["llmConversationId"] = data.LlmConversationId.ValueString()
        filterNames = append(filterNames, "llm_conversation_id = "+fmt.Sprintf("%q", data.LlmConversationId.ValueString()))
    }
    if !data.LlmUserId.IsNull() && !data.LlmUserId.IsUnknown() {
        filters["llmUserId"] = data.LlmUserId.ValueString()
        filterNames = append(filterNames, "llm_user_id = "+fmt.Sprintf("%q", data.LlmUserId.ValueString()))
    }
    if !data.LlmUserEmail.IsNull() && !data.LlmUserEmail.IsUnknown() {
        filters["llmUserEmail"] = data.LlmUserEmail.ValueString()
        filterNames = append(filterNames, "llm_user_email = "+fmt.Sprintf("%q", data.LlmUserEmail.ValueString()))
    }
    if !data.LlmTeam.IsNull() && !data.LlmTeam.IsUnknown() {
        filters["llmTeam"] = data.LlmTeam.ValueString()
        filterNames = append(filterNames, "llm_team = "+fmt.Sprintf("%q", data.LlmTeam.ValueString()))
    }
    if !data.LlmCallKind.IsNull() && !data.LlmCallKind.IsUnknown() {
        filters["llmCallKind"] = data.LlmCallKind.ValueString()
        filterNames = append(filterNames, "llm_call_kind = "+fmt.Sprintf("%q", data.LlmCallKind.ValueString()))
    }
    if !data.LlmUserMessagePreview.IsNull() && !data.LlmUserMessagePreview.IsUnknown() {
        filters["llmUserMessagePreview"] = data.LlmUserMessagePreview.ValueString()
        filterNames = append(filterNames, "llm_user_message_preview = "+fmt.Sprintf("%q", data.LlmUserMessagePreview.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the span up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the span up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "projectId": true,
        "primaryEntityId": true,
        "primaryEntityType": true,
        "startTime": true,
        "endTime": true,
        "startTimeUnixNano": true,
        "durationUnixNano": true,
        "endTimeUnixNano": true,
        "traceId": true,
        "spanId": true,
        "sessionId": true,
        "parentSpanId": true,
        "traceState": true,
        "attributes": true,
        "attributeKeys": true,
        "entityKeys": true,
        "serviceEntityKey": true,
        "hostEntityKey": true,
        "k8sPodEntityKey": true,
        "k8sNodeEntityKey": true,
        "k8sClusterEntityKey": true,
        "containerEntityKey": true,
        "events": true,
        "links": true,
        "statusCode": true,
        "statusMessage": true,
        "name": true,
        "kind": true,
        "hasException": true,
        "isRootSpan": true,
        "isLlmSpan": true,
        "llmSystem": true,
        "llmOperation": true,
        "llmRequestModel": true,
        "llmResponseModel": true,
        "llmAgentName": true,
        "llmToolName": true,
        "llmInputTokens": true,
        "llmOutputTokens": true,
        "llmTotalTokens": true,
        "llmCost": true,
        "llmConversationId": true,
        "llmUserId": true,
        "llmUserEmail": true,
        "llmTeam": true,
        "llmCallKind": true,
        "llmIssues": true,
        "llmUserMessagePreview": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/span/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read span, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No span found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read span: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/span/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list span, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list span: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No span matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one span matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for span.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["primaryEntityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if val, ok := item["primaryEntityId"].(string); ok {
        data.PrimaryEntityId = types.StringValue(val)
    } else {
        data.PrimaryEntityId = types.StringNull()
    }
    if obj, ok := item["primaryEntityType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if val, ok := item["primaryEntityType"].(string); ok {
        data.PrimaryEntityType = types.StringValue(val)
    } else {
        data.PrimaryEntityType = types.StringNull()
    }
    if obj, ok := item["startTime"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartTime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartTime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartTime = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartTime = types.StringValue(string(jsonBytes))
        } else {
            data.StartTime = types.StringNull()
        }
    } else if val, ok := item["startTime"].(string); ok {
        data.StartTime = types.StringValue(val)
    } else {
        data.StartTime = types.StringNull()
    }
    if obj, ok := item["endTime"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndTime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndTime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndTime = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndTime = types.StringValue(string(jsonBytes))
        } else {
            data.EndTime = types.StringNull()
        }
    } else if val, ok := item["endTime"].(string); ok {
        data.EndTime = types.StringValue(val)
    } else {
        data.EndTime = types.StringNull()
    }
    if obj, ok := item["startTimeUnixNano"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartTimeUnixNano = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartTimeUnixNano = types.StringValue(string(jsonBytes))
        } else {
            data.StartTimeUnixNano = types.StringNull()
        }
    } else if val, ok := item["startTimeUnixNano"].(string); ok {
        data.StartTimeUnixNano = types.StringValue(val)
    } else {
        data.StartTimeUnixNano = types.StringNull()
    }
    if val, ok := item["durationUnixNano"].(float64); ok {
        data.DurationUnixNano = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["durationUnixNano"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DurationUnixNano = types.NumberValue(big.NewFloat(val))
        } else {
            data.DurationUnixNano = types.NumberNull()
        }
    } else {
        data.DurationUnixNano = types.NumberNull()
    }
    if obj, ok := item["endTimeUnixNano"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndTimeUnixNano = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndTimeUnixNano = types.StringValue(string(jsonBytes))
        } else {
            data.EndTimeUnixNano = types.StringNull()
        }
    } else if val, ok := item["endTimeUnixNano"].(string); ok {
        data.EndTimeUnixNano = types.StringValue(val)
    } else {
        data.EndTimeUnixNano = types.StringNull()
    }
    if obj, ok := item["traceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TraceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TraceId = types.StringValue(string(jsonBytes))
        } else {
            data.TraceId = types.StringNull()
        }
    } else if val, ok := item["traceId"].(string); ok {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
    }
    if obj, ok := item["spanId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SpanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SpanId = types.StringValue(string(jsonBytes))
        } else {
            data.SpanId = types.StringNull()
        }
    } else if val, ok := item["spanId"].(string); ok {
        data.SpanId = types.StringValue(val)
    } else {
        data.SpanId = types.StringNull()
    }
    if obj, ok := item["sessionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionId = types.StringValue(string(jsonBytes))
        } else {
            data.SessionId = types.StringNull()
        }
    } else if val, ok := item["sessionId"].(string); ok {
        data.SessionId = types.StringValue(val)
    } else {
        data.SessionId = types.StringNull()
    }
    if obj, ok := item["parentSpanId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ParentSpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ParentSpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ParentSpanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ParentSpanId = types.StringValue(string(jsonBytes))
        } else {
            data.ParentSpanId = types.StringNull()
        }
    } else if val, ok := item["parentSpanId"].(string); ok {
        data.ParentSpanId = types.StringValue(val)
    } else {
        data.ParentSpanId = types.StringNull()
    }
    if obj, ok := item["traceState"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TraceState = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TraceState = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TraceState = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TraceState = types.StringValue(string(jsonBytes))
        } else {
            data.TraceState = types.StringNull()
        }
    } else if val, ok := item["traceState"].(string); ok {
        data.TraceState = types.StringValue(val)
    } else {
        data.TraceState = types.StringNull()
    }
    if obj, ok := item["attributes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Attributes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Attributes = types.StringValue(string(jsonBytes))
        } else {
            data.Attributes = types.StringNull()
        }
    } else if val, ok := item["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    } else {
        data.Attributes = types.StringNull()
    }
    if val, ok := item["attributeKeys"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.AttributeKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AttributeKeys = types.SetNull(types.StringType)
    }
    if val, ok := item["entityKeys"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.EntityKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        data.EntityKeys = types.SetNull(types.StringType)
    }
    if obj, ok := item["serviceEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceEntityKey = types.StringNull()
        }
    } else if val, ok := item["serviceEntityKey"].(string); ok {
        data.ServiceEntityKey = types.StringValue(val)
    } else {
        data.ServiceEntityKey = types.StringNull()
    }
    if obj, ok := item["hostEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HostEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HostEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HostEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HostEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.HostEntityKey = types.StringNull()
        }
    } else if val, ok := item["hostEntityKey"].(string); ok {
        data.HostEntityKey = types.StringValue(val)
    } else {
        data.HostEntityKey = types.StringNull()
    }
    if obj, ok := item["k8sPodEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sPodEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.K8sPodEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.K8sPodEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.K8sPodEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sPodEntityKey = types.StringNull()
        }
    } else if val, ok := item["k8sPodEntityKey"].(string); ok {
        data.K8sPodEntityKey = types.StringValue(val)
    } else {
        data.K8sPodEntityKey = types.StringNull()
    }
    if obj, ok := item["k8sNodeEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sNodeEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.K8sNodeEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.K8sNodeEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.K8sNodeEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sNodeEntityKey = types.StringNull()
        }
    } else if val, ok := item["k8sNodeEntityKey"].(string); ok {
        data.K8sNodeEntityKey = types.StringValue(val)
    } else {
        data.K8sNodeEntityKey = types.StringNull()
    }
    if obj, ok := item["k8sClusterEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sClusterEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.K8sClusterEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.K8sClusterEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.K8sClusterEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sClusterEntityKey = types.StringNull()
        }
    } else if val, ok := item["k8sClusterEntityKey"].(string); ok {
        data.K8sClusterEntityKey = types.StringValue(val)
    } else {
        data.K8sClusterEntityKey = types.StringNull()
    }
    if obj, ok := item["containerEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContainerEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ContainerEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ContainerEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ContainerEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.ContainerEntityKey = types.StringNull()
        }
    } else if val, ok := item["containerEntityKey"].(string); ok {
        data.ContainerEntityKey = types.StringValue(val)
    } else {
        data.ContainerEntityKey = types.StringNull()
    }
    if val, ok := item["events"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Events = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Events = types.SetNull(types.StringType)
    }
    if obj, ok := item["links"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Links = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Links = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Links = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Links = types.StringValue(string(jsonBytes))
        } else {
            data.Links = types.StringNull()
        }
    } else if val, ok := item["links"].(string); ok {
        data.Links = types.StringValue(val)
    } else {
        data.Links = types.StringNull()
    }
    if val, ok := item["statusCode"].(float64); ok {
        data.StatusCode = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["statusCode"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.StatusCode = types.NumberValue(big.NewFloat(val))
        } else {
            data.StatusCode = types.NumberNull()
        }
    } else {
        data.StatusCode = types.NumberNull()
    }
    if obj, ok := item["statusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := item["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["kind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := item["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if val, ok := item["hasException"].(bool); ok {
        data.HasException = types.BoolValue(val)
    } else {
        data.HasException = types.BoolNull()
    }
    if val, ok := item["isRootSpan"].(bool); ok {
        data.IsRootSpan = types.BoolValue(val)
    } else {
        data.IsRootSpan = types.BoolNull()
    }
    if val, ok := item["isLlmSpan"].(bool); ok {
        data.IsLlmSpan = types.BoolValue(val)
    } else {
        data.IsLlmSpan = types.BoolNull()
    }
    if obj, ok := item["llmSystem"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmSystem = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmSystem = types.StringValue(string(jsonBytes))
        } else {
            data.LlmSystem = types.StringNull()
        }
    } else if val, ok := item["llmSystem"].(string); ok {
        data.LlmSystem = types.StringValue(val)
    } else {
        data.LlmSystem = types.StringNull()
    }
    if obj, ok := item["llmOperation"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmOperation = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmOperation = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmOperation = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmOperation = types.StringValue(string(jsonBytes))
        } else {
            data.LlmOperation = types.StringNull()
        }
    } else if val, ok := item["llmOperation"].(string); ok {
        data.LlmOperation = types.StringValue(val)
    } else {
        data.LlmOperation = types.StringNull()
    }
    if obj, ok := item["llmRequestModel"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmRequestModel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmRequestModel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmRequestModel = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmRequestModel = types.StringValue(string(jsonBytes))
        } else {
            data.LlmRequestModel = types.StringNull()
        }
    } else if val, ok := item["llmRequestModel"].(string); ok {
        data.LlmRequestModel = types.StringValue(val)
    } else {
        data.LlmRequestModel = types.StringNull()
    }
    if obj, ok := item["llmResponseModel"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmResponseModel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmResponseModel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmResponseModel = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmResponseModel = types.StringValue(string(jsonBytes))
        } else {
            data.LlmResponseModel = types.StringNull()
        }
    } else if val, ok := item["llmResponseModel"].(string); ok {
        data.LlmResponseModel = types.StringValue(val)
    } else {
        data.LlmResponseModel = types.StringNull()
    }
    if obj, ok := item["llmAgentName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmAgentName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmAgentName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmAgentName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmAgentName = types.StringValue(string(jsonBytes))
        } else {
            data.LlmAgentName = types.StringNull()
        }
    } else if val, ok := item["llmAgentName"].(string); ok {
        data.LlmAgentName = types.StringValue(val)
    } else {
        data.LlmAgentName = types.StringNull()
    }
    if obj, ok := item["llmToolName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmToolName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmToolName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmToolName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmToolName = types.StringValue(string(jsonBytes))
        } else {
            data.LlmToolName = types.StringNull()
        }
    } else if val, ok := item["llmToolName"].(string); ok {
        data.LlmToolName = types.StringValue(val)
    } else {
        data.LlmToolName = types.StringNull()
    }
    if val, ok := item["llmInputTokens"].(float64); ok {
        data.LlmInputTokens = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["llmInputTokens"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.LlmInputTokens = types.NumberValue(big.NewFloat(val))
        } else {
            data.LlmInputTokens = types.NumberNull()
        }
    } else {
        data.LlmInputTokens = types.NumberNull()
    }
    if val, ok := item["llmOutputTokens"].(float64); ok {
        data.LlmOutputTokens = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["llmOutputTokens"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.LlmOutputTokens = types.NumberValue(big.NewFloat(val))
        } else {
            data.LlmOutputTokens = types.NumberNull()
        }
    } else {
        data.LlmOutputTokens = types.NumberNull()
    }
    if val, ok := item["llmTotalTokens"].(float64); ok {
        data.LlmTotalTokens = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["llmTotalTokens"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.LlmTotalTokens = types.NumberValue(big.NewFloat(val))
        } else {
            data.LlmTotalTokens = types.NumberNull()
        }
    } else {
        data.LlmTotalTokens = types.NumberNull()
    }
    if val, ok := item["llmCost"].(float64); ok {
        data.LlmCost = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["llmCost"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.LlmCost = types.NumberValue(big.NewFloat(val))
        } else {
            data.LlmCost = types.NumberNull()
        }
    } else {
        data.LlmCost = types.NumberNull()
    }
    if obj, ok := item["llmConversationId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmConversationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmConversationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmConversationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmConversationId = types.StringValue(string(jsonBytes))
        } else {
            data.LlmConversationId = types.StringNull()
        }
    } else if val, ok := item["llmConversationId"].(string); ok {
        data.LlmConversationId = types.StringValue(val)
    } else {
        data.LlmConversationId = types.StringNull()
    }
    if obj, ok := item["llmUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LlmUserId = types.StringNull()
        }
    } else if val, ok := item["llmUserId"].(string); ok {
        data.LlmUserId = types.StringValue(val)
    } else {
        data.LlmUserId = types.StringNull()
    }
    if obj, ok := item["llmUserEmail"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmUserEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmUserEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmUserEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmUserEmail = types.StringValue(string(jsonBytes))
        } else {
            data.LlmUserEmail = types.StringNull()
        }
    } else if val, ok := item["llmUserEmail"].(string); ok {
        data.LlmUserEmail = types.StringValue(val)
    } else {
        data.LlmUserEmail = types.StringNull()
    }
    if obj, ok := item["llmTeam"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmTeam = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmTeam = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmTeam = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmTeam = types.StringValue(string(jsonBytes))
        } else {
            data.LlmTeam = types.StringNull()
        }
    } else if val, ok := item["llmTeam"].(string); ok {
        data.LlmTeam = types.StringValue(val)
    } else {
        data.LlmTeam = types.StringNull()
    }
    if obj, ok := item["llmCallKind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmCallKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmCallKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmCallKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmCallKind = types.StringValue(string(jsonBytes))
        } else {
            data.LlmCallKind = types.StringNull()
        }
    } else if val, ok := item["llmCallKind"].(string); ok {
        data.LlmCallKind = types.StringValue(val)
    } else {
        data.LlmCallKind = types.StringNull()
    }
    if val, ok := item["llmIssues"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.LlmIssues = types.SetValueMust(types.StringType, setItems)
    } else {
        data.LlmIssues = types.SetNull(types.StringType)
    }
    if obj, ok := item["llmUserMessagePreview"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LlmUserMessagePreview = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LlmUserMessagePreview = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LlmUserMessagePreview = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LlmUserMessagePreview = types.StringValue(string(jsonBytes))
        } else {
            data.LlmUserMessagePreview = types.StringNull()
        }
    } else if val, ok := item["llmUserMessagePreview"].(string); ok {
        data.LlmUserMessagePreview = types.StringValue(val)
    } else {
        data.LlmUserMessagePreview = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
