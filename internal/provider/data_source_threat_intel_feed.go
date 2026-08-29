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
var _ datasource.DataSource = &ThreatIntelFeedDataSource{}

func NewThreatIntelFeedDataSource() datasource.DataSource {
    return &ThreatIntelFeedDataSource{}
}

// ThreatIntelFeedDataSource defines the data source implementation.
type ThreatIntelFeedDataSource struct {
    client *Client
}

// ThreatIntelFeedDataSourceModel describes the data source data model.
type ThreatIntelFeedDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    ApiRootUrl types.String `tfsdk:"api_root_url"`
    CollectionId types.String `tfsdk:"collection_id"`
    BasicAuthUsername types.String `tfsdk:"basic_auth_username"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    PollIntervalInMinutes types.Number `tfsdk:"poll_interval_in_minutes"`
    MinimumConfidence types.Number `tfsdk:"minimum_confidence"`
    ShouldCreateAlert types.Bool `tfsdk:"should_create_alert"`
    ShouldWriteDetectionFinding types.Bool `tfsdk:"should_write_detection_finding"`
    ShouldCreateIncident types.Bool `tfsdk:"should_create_incident"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    LastPolledAt types.String `tfsdk:"last_polled_at"`
    Cursor types.String `tfsdk:"cursor"`
    NextPageToken types.String `tfsdk:"next_page_token"`
    LastPollSummary types.String `tfsdk:"last_poll_summary"`
    LastError types.String `tfsdk:"last_error"`
    LastEvaluatedAt types.String `tfsdk:"last_evaluated_at"`
    LastMatchAt types.String `tfsdk:"last_match_at"`
    LastMatchError types.String `tfsdk:"last_match_error"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *ThreatIntelFeedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_threat_intel_feed"
}

func (d *ThreatIntelFeedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "STIX/TAXII 2.1 threat-intelligence feeds. Indicators are polled on an interval and matched against incoming security events. Look up an existing threat_intel_feed by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What this feed carries and why it is subscribed..",
                Computed: true,
            },
            "api_root_url": schema.StringAttribute{
                MarkdownDescription: "The TAXII 2.1 API root, e.g. https://taxii.example.com/api1/. Collections are addressed beneath it..",
                Computed: true,
            },
            "collection_id": schema.StringAttribute{
                MarkdownDescription: "ID of the TAXII collection to poll for indicator objects..",
                Computed: true,
            },
            "basic_auth_username": schema.StringAttribute{
                MarkdownDescription: "Username for basic-auth collections. Leave empty for anonymous or token-authenticated collections..",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this feed is polled and matched..",
                Computed: true,
            },
            "poll_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often the collection is polled for new indicators..",
                Computed: true,
            },
            "minimum_confidence": schema.NumberAttribute{
                MarkdownDescription: "Skip indicators whose STIX confidence is below this (0-100). 0 ingests everything; indicators that carry no confidence always pass..",
                Computed: true,
            },
            "should_create_alert": schema.BoolAttribute{
                MarkdownDescription: "Whether indicator matches open OneUptime alerts..",
                Computed: true,
            },
            "should_write_detection_finding": schema.BoolAttribute{
                MarkdownDescription: "Whether matches also write a Detection Finding security event back into the events table..",
                Computed: true,
            },
            "should_create_incident": schema.BoolAttribute{
                MarkdownDescription: "Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per feed..",
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "last_polled_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "cursor": schema.StringAttribute{
                MarkdownDescription: "Poll cursor: the TAXII added_after timestamp already ingested, as an ISO string..",
                Computed: true,
            },
            "next_page_token": schema.StringAttribute{
                MarkdownDescription: "Resume token for a poll that ended mid-pagination on a server that sends no X-TAXII-Date-Added-Last header. Cleared once the collection drains or the cursor advances..",
                Computed: true,
            },
            "last_poll_summary": schema.StringAttribute{
                MarkdownDescription: "What the most recent successful poll did: objects fetched, indicators ingested, unsupported patterns skipped..",
                Computed: true,
            },
            "last_error": schema.StringAttribute{
                MarkdownDescription: "The most recent poll error, if any. Cleared on the next successful poll..",
                Computed: true,
            },
            "last_evaluated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_match_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_match_error": schema.StringAttribute{
                MarkdownDescription: "The most recent matcher error, if any. Cleared on the next successful evaluation..",
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

func (d *ThreatIntelFeedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ThreatIntelFeedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ThreatIntelFeedDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a threat_intel_feed.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "description": true,
        "apiRootUrl": true,
        "collectionId": true,
        "basicAuthUsername": true,
        "isEnabled": true,
        "pollIntervalInMinutes": true,
        "minimumConfidence": true,
        "shouldCreateAlert": true,
        "shouldWriteDetectionFinding": true,
        "shouldCreateIncident": true,
        "alertSeverityId": true,
        "incidentSeverityId": true,
        "lastPolledAt": true,
        "cursor": true,
        "nextPageToken": true,
        "lastPollSummary": true,
        "lastError": true,
        "lastEvaluatedAt": true,
        "lastMatchAt": true,
        "lastMatchError": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/threat-intel-feed/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read threat_intel_feed, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No threat_intel_feed found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read threat_intel_feed: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/threat-intel-feed/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list threat_intel_feed, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list threat_intel_feed: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No threat_intel_feed found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one threat_intel_feed matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for threat_intel_feed.")
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["apiRootUrl"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApiRootUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ApiRootUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ApiRootUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ApiRootUrl = types.StringValue(string(jsonBytes))
        } else {
            data.ApiRootUrl = types.StringNull()
        }
    } else if val, ok := item["apiRootUrl"].(string); ok {
        data.ApiRootUrl = types.StringValue(val)
    } else {
        data.ApiRootUrl = types.StringNull()
    }
    if obj, ok := item["collectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CollectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CollectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CollectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CollectionId = types.StringValue(string(jsonBytes))
        } else {
            data.CollectionId = types.StringNull()
        }
    } else if val, ok := item["collectionId"].(string); ok {
        data.CollectionId = types.StringValue(val)
    } else {
        data.CollectionId = types.StringNull()
    }
    if obj, ok := item["basicAuthUsername"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BasicAuthUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BasicAuthUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BasicAuthUsername = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BasicAuthUsername = types.StringValue(string(jsonBytes))
        } else {
            data.BasicAuthUsername = types.StringNull()
        }
    } else if val, ok := item["basicAuthUsername"].(string); ok {
        data.BasicAuthUsername = types.StringValue(val)
    } else {
        data.BasicAuthUsername = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["pollIntervalInMinutes"].(float64); ok {
        data.PollIntervalInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["pollIntervalInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PollIntervalInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.PollIntervalInMinutes = types.NumberNull()
        }
    } else {
        data.PollIntervalInMinutes = types.NumberNull()
    }
    if val, ok := item["minimumConfidence"].(float64); ok {
        data.MinimumConfidence = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["minimumConfidence"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.MinimumConfidence = types.NumberValue(big.NewFloat(val))
        } else {
            data.MinimumConfidence = types.NumberNull()
        }
    } else {
        data.MinimumConfidence = types.NumberNull()
    }
    if val, ok := item["shouldCreateAlert"].(bool); ok {
        data.ShouldCreateAlert = types.BoolValue(val)
    } else {
        data.ShouldCreateAlert = types.BoolNull()
    }
    if val, ok := item["shouldWriteDetectionFinding"].(bool); ok {
        data.ShouldWriteDetectionFinding = types.BoolValue(val)
    } else {
        data.ShouldWriteDetectionFinding = types.BoolNull()
    }
    if val, ok := item["shouldCreateIncident"].(bool); ok {
        data.ShouldCreateIncident = types.BoolValue(val)
    } else {
        data.ShouldCreateIncident = types.BoolNull()
    }
    if obj, ok := item["alertSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := item["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if obj, ok := item["incidentSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := item["incidentSeverityId"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if obj, ok := item["lastPolledAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastPolledAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastPolledAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastPolledAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastPolledAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastPolledAt = types.StringNull()
        }
    } else if val, ok := item["lastPolledAt"].(string); ok {
        data.LastPolledAt = types.StringValue(val)
    } else {
        data.LastPolledAt = types.StringNull()
    }
    if obj, ok := item["cursor"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Cursor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Cursor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Cursor = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Cursor = types.StringValue(string(jsonBytes))
        } else {
            data.Cursor = types.StringNull()
        }
    } else if val, ok := item["cursor"].(string); ok {
        data.Cursor = types.StringValue(val)
    } else {
        data.Cursor = types.StringNull()
    }
    if obj, ok := item["nextPageToken"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextPageToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextPageToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextPageToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextPageToken = types.StringValue(string(jsonBytes))
        } else {
            data.NextPageToken = types.StringNull()
        }
    } else if val, ok := item["nextPageToken"].(string); ok {
        data.NextPageToken = types.StringValue(val)
    } else {
        data.NextPageToken = types.StringNull()
    }
    if obj, ok := item["lastPollSummary"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastPollSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastPollSummary = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastPollSummary = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastPollSummary = types.StringValue(string(jsonBytes))
        } else {
            data.LastPollSummary = types.StringNull()
        }
    } else if val, ok := item["lastPollSummary"].(string); ok {
        data.LastPollSummary = types.StringValue(val)
    } else {
        data.LastPollSummary = types.StringNull()
    }
    if obj, ok := item["lastError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastError = types.StringValue(string(jsonBytes))
        } else {
            data.LastError = types.StringNull()
        }
    } else if val, ok := item["lastError"].(string); ok {
        data.LastError = types.StringValue(val)
    } else {
        data.LastError = types.StringNull()
    }
    if obj, ok := item["lastEvaluatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastEvaluatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastEvaluatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastEvaluatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastEvaluatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastEvaluatedAt = types.StringNull()
        }
    } else if val, ok := item["lastEvaluatedAt"].(string); ok {
        data.LastEvaluatedAt = types.StringValue(val)
    } else {
        data.LastEvaluatedAt = types.StringNull()
    }
    if obj, ok := item["lastMatchAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastMatchAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastMatchAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastMatchAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastMatchAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastMatchAt = types.StringNull()
        }
    } else if val, ok := item["lastMatchAt"].(string); ok {
        data.LastMatchAt = types.StringValue(val)
    } else {
        data.LastMatchAt = types.StringNull()
    }
    if obj, ok := item["lastMatchError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastMatchError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastMatchError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastMatchError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastMatchError = types.StringValue(string(jsonBytes))
        } else {
            data.LastMatchError = types.StringNull()
        }
    } else if val, ok := item["lastMatchError"].(string); ok {
        data.LastMatchError = types.StringValue(val)
    } else {
        data.LastMatchError = types.StringNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
