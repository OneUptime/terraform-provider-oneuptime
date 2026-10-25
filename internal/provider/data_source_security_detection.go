package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SecurityDetectionDataSource{}

func NewSecurityDetectionDataSource() datasource.DataSource {
    return &SecurityDetectionDataSource{}
}

// SecurityDetectionDataSource defines the data source implementation.
type SecurityDetectionDataSource struct {
    client *Client
}

// SecurityDetectionDataSourceModel describes the data source data model.
type SecurityDetectionDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    SecurityCaseId types.String `tfsdk:"security_case_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    Severity types.String `tfsdk:"severity"`
    Source types.String `tfsdk:"source"`
    Fingerprint types.String `tfsdk:"fingerprint"`
    RuleName types.String `tfsdk:"rule_name"`
    DetectionRuleId types.String `tfsdk:"detection_rule_id"`
    ThreatIntelFeedId types.String `tfsdk:"threat_intel_feed_id"`
    SecurityEventConnectionId types.String `tfsdk:"security_event_connection_id"`
    GroupValue types.String `tfsdk:"group_value"`
    Entities types.String `tfsdk:"entities"`
    MitreTactics types.String `tfsdk:"mitre_tactics"`
    MitreTechniques types.String `tfsdk:"mitre_techniques"`
    EventCount types.Number `tfsdk:"event_count"`
    FireCount types.Number `tfsdk:"fire_count"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    SampleMessage types.String `tfsdk:"sample_message"`
    EventUid types.String `tfsdk:"event_uid"`
}

func (d *SecurityDetectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_detection"
}

func (d *SecurityDetectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A signal inside a security case: a detection rule's match, a threat intelligence hit, or an alert from a connected security product, with how often it fired. Look up an existing security detection by `id`, or by any of its other arguments (`description`, `detection_rule_id`, `event_count`, ...): each one set must match, and exactly one security detection may match them all.",

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
                MarkdownDescription: "ID of the project this detection belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "security_case_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case this detection belongs to. The ID of a `oneuptime_security_case`.",
                Optional: true,
                Computed: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "What fired, about what: the rule's name and the value it grouped on.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the detection saw, in Markdown.",
                Optional: true,
                Computed: true,
            },
            "severity": schema.StringAttribute{
                MarkdownDescription: "OCSF severity of the detection: Informational, Low, Medium, High, Critical or Fatal.",
                Optional: true,
                Computed: true,
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "What raised it: a Detection Rule, Threat Intel, a Connected Product, Entity Risk, or a person (Manual).",
                Optional: true,
                Computed: true,
            },
            "fingerprint": schema.StringAttribute{
                MarkdownDescription: "Identifies the same signal across repeats: the source and what it fired for.",
                Optional: true,
                Computed: true,
            },
            "rule_name": schema.StringAttribute{
                MarkdownDescription: "The name of what fired: the detection rule, the threat intel feed, or the connected product's own rule.",
                Optional: true,
                Computed: true,
            },
            "detection_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the detection rule that fired. The ID of a `oneuptime_detection_rule`.",
                Optional: true,
                Computed: true,
            },
            "threat_intel_feed_id": schema.StringAttribute{
                MarkdownDescription: "ID of the feed whose indicator matched. The ID of a `oneuptime_threat_intel_feed`.",
                Optional: true,
                Computed: true,
            },
            "security_event_connection_id": schema.StringAttribute{
                MarkdownDescription: "ID of the connection the alert was imported through. The ID of a `oneuptime_security_event_connection`.",
                Optional: true,
                Computed: true,
            },
            "group_value": schema.StringAttribute{
                MarkdownDescription: "The value the rule grouped its matches by, such as the host or IP address they share.",
                Optional: true,
                Computed: true,
            },
            "entities": schema.StringAttribute{
                MarkdownDescription: "The users, hosts, IP addresses, domains and files the detection names, each as { type, value }. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "mitre_tactics": schema.StringAttribute{
                MarkdownDescription: "ATT&CK tactic ids the detection is tagged with. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "mitre_techniques": schema.StringAttribute{
                MarkdownDescription: "ATT&CK technique ids the detection is tagged with. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "event_count": schema.NumberAttribute{
                MarkdownDescription: "How many security events are behind the detection, over every time it fired.",
                Optional: true,
                Computed: true,
            },
            "fire_count": schema.NumberAttribute{
                MarkdownDescription: "How many times the detection fired while its case was open.",
                Optional: true,
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the detection first fired.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the detection last fired.",
                Computed: true,
            },
            "sample_message": schema.StringAttribute{
                MarkdownDescription: "The message of one of the events behind the detection.",
                Optional: true,
                Computed: true,
            },
            "event_uid": schema.StringAttribute{
                MarkdownDescription: "eventUid of the security event that records the detection: the Detection Finding a rule or feed wrote, or the imported alert.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *SecurityDetectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurityDetectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SecurityDetectionDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.SecurityCaseId.IsNull() && !data.SecurityCaseId.IsUnknown() {
        filters["securityCaseId"] = data.SecurityCaseId.ValueString()
        filterNames = append(filterNames, "security_case_id = "+fmt.Sprintf("%q", data.SecurityCaseId.ValueString()))
    }
    if !data.Title.IsNull() && !data.Title.IsUnknown() {
        filters["title"] = data.Title.ValueString()
        filterNames = append(filterNames, "title = "+fmt.Sprintf("%q", data.Title.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.Severity.IsNull() && !data.Severity.IsUnknown() {
        filters["severity"] = data.Severity.ValueString()
        filterNames = append(filterNames, "severity = "+fmt.Sprintf("%q", data.Severity.ValueString()))
    }
    if !data.Source.IsNull() && !data.Source.IsUnknown() {
        filters["source"] = data.Source.ValueString()
        filterNames = append(filterNames, "source = "+fmt.Sprintf("%q", data.Source.ValueString()))
    }
    if !data.Fingerprint.IsNull() && !data.Fingerprint.IsUnknown() {
        filters["fingerprint"] = data.Fingerprint.ValueString()
        filterNames = append(filterNames, "fingerprint = "+fmt.Sprintf("%q", data.Fingerprint.ValueString()))
    }
    if !data.RuleName.IsNull() && !data.RuleName.IsUnknown() {
        filters["ruleName"] = data.RuleName.ValueString()
        filterNames = append(filterNames, "rule_name = "+fmt.Sprintf("%q", data.RuleName.ValueString()))
    }
    if !data.DetectionRuleId.IsNull() && !data.DetectionRuleId.IsUnknown() {
        filters["detectionRuleId"] = data.DetectionRuleId.ValueString()
        filterNames = append(filterNames, "detection_rule_id = "+fmt.Sprintf("%q", data.DetectionRuleId.ValueString()))
    }
    if !data.ThreatIntelFeedId.IsNull() && !data.ThreatIntelFeedId.IsUnknown() {
        filters["threatIntelFeedId"] = data.ThreatIntelFeedId.ValueString()
        filterNames = append(filterNames, "threat_intel_feed_id = "+fmt.Sprintf("%q", data.ThreatIntelFeedId.ValueString()))
    }
    if !data.SecurityEventConnectionId.IsNull() && !data.SecurityEventConnectionId.IsUnknown() {
        filters["securityEventConnectionId"] = data.SecurityEventConnectionId.ValueString()
        filterNames = append(filterNames, "security_event_connection_id = "+fmt.Sprintf("%q", data.SecurityEventConnectionId.ValueString()))
    }
    if !data.GroupValue.IsNull() && !data.GroupValue.IsUnknown() {
        filters["groupValue"] = data.GroupValue.ValueString()
        filterNames = append(filterNames, "group_value = "+fmt.Sprintf("%q", data.GroupValue.ValueString()))
    }
    if !data.EventCount.IsNull() && !data.EventCount.IsUnknown() {
        filters["eventCount"] = lookupNumber(data.EventCount)
        filterNames = append(filterNames, "event_count = "+data.EventCount.ValueBigFloat().String())
    }
    if !data.FireCount.IsNull() && !data.FireCount.IsUnknown() {
        filters["fireCount"] = lookupNumber(data.FireCount)
        filterNames = append(filterNames, "fire_count = "+data.FireCount.ValueBigFloat().String())
    }
    if !data.SampleMessage.IsNull() && !data.SampleMessage.IsUnknown() {
        filters["sampleMessage"] = data.SampleMessage.ValueString()
        filterNames = append(filterNames, "sample_message = "+fmt.Sprintf("%q", data.SampleMessage.ValueString()))
    }
    if !data.EventUid.IsNull() && !data.EventUid.IsUnknown() {
        filters["eventUid"] = data.EventUid.ValueString()
        filterNames = append(filterNames, "event_uid = "+fmt.Sprintf("%q", data.EventUid.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the security detection up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the security detection up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "securityCaseId": true,
        "title": true,
        "description": true,
        "severity": true,
        "source": true,
        "fingerprint": true,
        "ruleName": true,
        "detectionRuleId": true,
        "threatIntelFeedId": true,
        "securityEventConnectionId": true,
        "groupValue": true,
        "entities": true,
        "mitreTactics": true,
        "mitreTechniques": true,
        "eventCount": true,
        "fireCount": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "sampleMessage": true,
        "eventUid": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-detection/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_detection, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security detection found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_detection: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-detection/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list security_detection, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list security_detection: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security detection matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security detection matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for security_detection.")
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
    if obj, ok := item["securityCaseId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if val, ok := item["securityCaseId"].(string); ok {
        data.SecurityCaseId = types.StringValue(val)
    } else {
        data.SecurityCaseId = types.StringNull()
    }
    if obj, ok := item["title"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := item["title"].(string); ok {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
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
    if obj, ok := item["severity"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Severity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Severity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Severity = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Severity = types.StringValue(string(jsonBytes))
        } else {
            data.Severity = types.StringNull()
        }
    } else if val, ok := item["severity"].(string); ok {
        data.Severity = types.StringValue(val)
    } else {
        data.Severity = types.StringNull()
    }
    if obj, ok := item["source"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := item["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
    }
    if obj, ok := item["fingerprint"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Fingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Fingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.Fingerprint = types.StringNull()
        }
    } else if val, ok := item["fingerprint"].(string); ok {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
    }
    if obj, ok := item["ruleName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuleName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuleName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuleName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuleName = types.StringValue(string(jsonBytes))
        } else {
            data.RuleName = types.StringNull()
        }
    } else if val, ok := item["ruleName"].(string); ok {
        data.RuleName = types.StringValue(val)
    } else {
        data.RuleName = types.StringNull()
    }
    if obj, ok := item["detectionRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DetectionRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DetectionRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DetectionRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DetectionRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.DetectionRuleId = types.StringNull()
        }
    } else if val, ok := item["detectionRuleId"].(string); ok {
        data.DetectionRuleId = types.StringValue(val)
    } else {
        data.DetectionRuleId = types.StringNull()
    }
    if obj, ok := item["threatIntelFeedId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ThreatIntelFeedId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ThreatIntelFeedId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ThreatIntelFeedId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ThreatIntelFeedId = types.StringValue(string(jsonBytes))
        } else {
            data.ThreatIntelFeedId = types.StringNull()
        }
    } else if val, ok := item["threatIntelFeedId"].(string); ok {
        data.ThreatIntelFeedId = types.StringValue(val)
    } else {
        data.ThreatIntelFeedId = types.StringNull()
    }
    if obj, ok := item["securityEventConnectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityEventConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityEventConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityEventConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityEventConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityEventConnectionId = types.StringNull()
        }
    } else if val, ok := item["securityEventConnectionId"].(string); ok {
        data.SecurityEventConnectionId = types.StringValue(val)
    } else {
        data.SecurityEventConnectionId = types.StringNull()
    }
    if obj, ok := item["groupValue"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GroupValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GroupValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GroupValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GroupValue = types.StringValue(string(jsonBytes))
        } else {
            data.GroupValue = types.StringNull()
        }
    } else if val, ok := item["groupValue"].(string); ok {
        data.GroupValue = types.StringValue(val)
    } else {
        data.GroupValue = types.StringNull()
    }
    if obj, ok := item["entities"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Entities = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Entities = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Entities = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Entities = types.StringValue(string(jsonBytes))
        } else {
            data.Entities = types.StringNull()
        }
    } else if val, ok := item["entities"].(string); ok {
        data.Entities = types.StringValue(val)
    } else {
        data.Entities = types.StringNull()
    }
    if obj, ok := item["mitreTactics"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MitreTactics = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MitreTactics = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MitreTactics = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MitreTactics = types.StringValue(string(jsonBytes))
        } else {
            data.MitreTactics = types.StringNull()
        }
    } else if val, ok := item["mitreTactics"].(string); ok {
        data.MitreTactics = types.StringValue(val)
    } else {
        data.MitreTactics = types.StringNull()
    }
    if obj, ok := item["mitreTechniques"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MitreTechniques = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MitreTechniques = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MitreTechniques = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MitreTechniques = types.StringValue(string(jsonBytes))
        } else {
            data.MitreTechniques = types.StringNull()
        }
    } else if val, ok := item["mitreTechniques"].(string); ok {
        data.MitreTechniques = types.StringValue(val)
    } else {
        data.MitreTechniques = types.StringNull()
    }
    data.EventCount = numberFromAPI(data.EventCount, item["eventCount"])
    data.FireCount = numberFromAPI(data.FireCount, item["fireCount"])
    if obj, ok := item["firstSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenAt = types.StringNull()
        }
    } else if val, ok := item["firstSeenAt"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    } else {
        data.FirstSeenAt = types.StringNull()
    }
    if obj, ok := item["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenAt = types.StringNull()
        }
    } else if val, ok := item["lastSeenAt"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    } else {
        data.LastSeenAt = types.StringNull()
    }
    if obj, ok := item["sampleMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SampleMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SampleMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SampleMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SampleMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SampleMessage = types.StringNull()
        }
    } else if val, ok := item["sampleMessage"].(string); ok {
        data.SampleMessage = types.StringValue(val)
    } else {
        data.SampleMessage = types.StringNull()
    }
    if obj, ok := item["eventUid"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EventUid = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EventUid = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EventUid = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EventUid = types.StringValue(string(jsonBytes))
        } else {
            data.EventUid = types.StringNull()
        }
    } else if val, ok := item["eventUid"].(string); ok {
        data.EventUid = types.StringValue(val)
    } else {
        data.EventUid = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
