package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AiRunEventDataSource{}

func NewAiRunEventDataSource() datasource.DataSource {
    return &AiRunEventDataSource{}
}

// AiRunEventDataSource defines the data source implementation.
type AiRunEventDataSource struct {
    client *Client
}

// AiRunEventDataSourceModel describes the data source data model.
type AiRunEventDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    AiRunId types.String `tfsdk:"ai_run_id"`
    UserId types.String `tfsdk:"user_id"`
    Sequence types.Number `tfsdk:"sequence"`
    EventType types.String `tfsdk:"event_type"`
    ToolName types.String `tfsdk:"tool_name"`
    ToolArguments types.String `tfsdk:"tool_arguments"`
    ResultSummary types.String `tfsdk:"result_summary"`
    CitationId types.String `tfsdk:"citation_id"`
}

func (d *AiRunEventDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_ai_run_event"
}

func (d *AiRunEventDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "An event in an AI run: LLM calls, tool calls with validated arguments, and lifecycle transitions. Look up an existing ai run event by `id`, or by any of its other arguments (`ai_run_id`, `citation_id`, `event_type`, ...): each one set must match, and exactly one ai run event may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this event belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "ai_run_id": schema.StringAttribute{
                MarkdownDescription: "ID of the run this event belongs to. The ID of a `oneuptime_ai_run` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user whose run this event belongs to. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "sequence": schema.NumberAttribute{
                MarkdownDescription: "Order of this event within the run.",
                Optional: true,
                Computed: true,
            },
            "event_type": schema.StringAttribute{
                MarkdownDescription: "Type of event.",
                Optional: true,
                Computed: true,
            },
            "tool_name": schema.StringAttribute{
                MarkdownDescription: "Name of the tool for tool-call events.",
                Optional: true,
                Computed: true,
            },
            "tool_arguments": schema.StringAttribute{
                MarkdownDescription: "Validated tool arguments as executed. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "result_summary": schema.StringAttribute{
                MarkdownDescription: "Summary of the result: row count, duration, truncation and bytes sent to the LLM. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "citation_id": schema.StringAttribute{
                MarkdownDescription: "ID of the citation this event minted (e.g. C1), if it produced one.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *AiRunEventDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AiRunEventDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AiRunEventDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.AiRunId.IsNull() && !data.AiRunId.IsUnknown() {
        filters["aiRunId"] = data.AiRunId.ValueString()
        filterNames = append(filterNames, "ai_run_id = "+fmt.Sprintf("%q", data.AiRunId.ValueString()))
    }
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }
    if !data.Sequence.IsNull() && !data.Sequence.IsUnknown() {
        filters["sequence"] = lookupNumber(data.Sequence)
        filterNames = append(filterNames, "sequence = "+data.Sequence.ValueBigFloat().String())
    }
    if !data.EventType.IsNull() && !data.EventType.IsUnknown() {
        filters["eventType"] = data.EventType.ValueString()
        filterNames = append(filterNames, "event_type = "+fmt.Sprintf("%q", data.EventType.ValueString()))
    }
    if !data.ToolName.IsNull() && !data.ToolName.IsUnknown() {
        filters["toolName"] = data.ToolName.ValueString()
        filterNames = append(filterNames, "tool_name = "+fmt.Sprintf("%q", data.ToolName.ValueString()))
    }
    if !data.CitationId.IsNull() && !data.CitationId.IsUnknown() {
        filters["citationId"] = data.CitationId.ValueString()
        filterNames = append(filterNames, "citation_id = "+fmt.Sprintf("%q", data.CitationId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the ai run event up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the ai run event up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "aiRunId": true,
        "userId": true,
        "sequence": true,
        "eventType": true,
        "toolName": true,
        "toolArguments": true,
        "resultSummary": true,
        "citationId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/ai-run-event/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ai_run_event, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai run event found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read ai_run_event: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/ai-run-event/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list ai_run_event, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list ai_run_event: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No ai run event matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one ai run event matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for ai_run_event.")
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
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
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
    if obj, ok := item["aiRunId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiRunId = types.StringValue(string(jsonBytes))
        } else {
            data.AiRunId = types.StringNull()
        }
    } else if val, ok := item["aiRunId"].(string); ok {
        data.AiRunId = types.StringValue(val)
    } else {
        data.AiRunId = types.StringNull()
    }
    if obj, ok := item["userId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserId = types.StringValue(string(jsonBytes))
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := item["userId"].(string); ok {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if val, ok := item["sequence"].(float64); ok {
        data.Sequence = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sequence"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Sequence = types.NumberValue(big.NewFloat(val))
        } else {
            data.Sequence = types.NumberNull()
        }
    } else {
        data.Sequence = types.NumberNull()
    }
    if obj, ok := item["eventType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EventType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EventType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EventType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EventType = types.StringValue(string(jsonBytes))
        } else {
            data.EventType = types.StringNull()
        }
    } else if val, ok := item["eventType"].(string); ok {
        data.EventType = types.StringValue(val)
    } else {
        data.EventType = types.StringNull()
    }
    if obj, ok := item["toolName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ToolName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ToolName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ToolName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ToolName = types.StringValue(string(jsonBytes))
        } else {
            data.ToolName = types.StringNull()
        }
    } else if val, ok := item["toolName"].(string); ok {
        data.ToolName = types.StringValue(val)
    } else {
        data.ToolName = types.StringNull()
    }
    if obj, ok := item["toolArguments"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ToolArguments = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ToolArguments = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ToolArguments = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ToolArguments = types.StringValue(string(jsonBytes))
        } else {
            data.ToolArguments = types.StringNull()
        }
    } else if val, ok := item["toolArguments"].(string); ok {
        data.ToolArguments = types.StringValue(val)
    } else {
        data.ToolArguments = types.StringNull()
    }
    if obj, ok := item["resultSummary"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResultSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResultSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResultSummary = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResultSummary = types.StringValue(string(jsonBytes))
        } else {
            data.ResultSummary = types.StringNull()
        }
    } else if val, ok := item["resultSummary"].(string); ok {
        data.ResultSummary = types.StringValue(val)
    } else {
        data.ResultSummary = types.StringNull()
    }
    if obj, ok := item["citationId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CitationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CitationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CitationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CitationId = types.StringValue(string(jsonBytes))
        } else {
            data.CitationId = types.StringNull()
        }
    } else if val, ok := item["citationId"].(string); ok {
        data.CitationId = types.StringValue(val)
    } else {
        data.CitationId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
