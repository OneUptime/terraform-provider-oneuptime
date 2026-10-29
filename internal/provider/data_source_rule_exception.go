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
var _ datasource.DataSource = &RuleExceptionDataSource{}

func NewRuleExceptionDataSource() datasource.DataSource {
    return &RuleExceptionDataSource{}
}

// RuleExceptionDataSource defines the data source implementation.
type RuleExceptionDataSource struct {
    client *Client
}

// RuleExceptionDataSourceModel describes the data source data model.
type RuleExceptionDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    DetectionRuleId types.String `tfsdk:"detection_rule_id"`
    Kind types.String `tfsdk:"kind"`
    EntityType types.String `tfsdk:"entity_type"`
    FieldName types.String `tfsdk:"field_name"`
    Value types.String `tfsdk:"value"`
    StartsAt types.String `tfsdk:"starts_at"`
    EndsAt types.String `tfsdk:"ends_at"`
    Reason types.String `tfsdk:"reason"`
    ExpiresAt types.String `tfsdk:"expires_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *RuleExceptionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_rule_exception"
}

func (d *RuleExceptionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "What a detection rule skips: an entity, a value or a time window, with who added it, why, and when it expires. Look up an existing rule exception by `id`, or by any of its other arguments (`created_by_user_id`, `detection_rule_id`, `entity_type`, ...): each one set must match, and exactly one rule exception may match them all.",

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
                MarkdownDescription: "ID of the project this exception belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "detection_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the detection rule that skips what this names. The ID of a `oneuptime_detection_rule`.",
                Optional: true,
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "What the rule skips: Entity (a user, host or IP address), Value (one event field's value) or Time Window (every event between two times).",
                Optional: true,
                Computed: true,
            },
            "entity_type": schema.StringAttribute{
                MarkdownDescription: "For an Entity: User, Host or IP Address. The rule skips every event that names it, as who acted, as who was acted on, or among its observables.",
                Optional: true,
                Computed: true,
            },
            "field_name": schema.StringAttribute{
                MarkdownDescription: "For a Value: the event field to match, named as in a Sigma rule - principalUser, processName, an attribute key such as http.user_agent.",
                Optional: true,
                Computed: true,
            },
            "value": schema.StringAttribute{
                MarkdownDescription: "For an Entity or a Value: what the rule skips, matched whatever the case. * stands for any run of characters (svc-* skips every service account), and an IP address may be a range such as 10.20.0.0/16.",
                Optional: true,
                Computed: true,
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "For a Time Window: the first moment the rule skips.",
                Computed: true,
            },
            "ends_at": schema.StringAttribute{
                MarkdownDescription: "For a Time Window: when the rule sees events again. A window is at most 30 days.",
                Computed: true,
            },
            "reason": schema.StringAttribute{
                MarkdownDescription: "Why it is safe for the rule to skip this, for whoever reads the list later.",
                Optional: true,
                Computed: true,
            },
            "expires_at": schema.StringAttribute{
                MarkdownDescription: "For an Entity or a Value: when the rule sees it again. Empty: until someone deletes the exception.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who added the exception. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *RuleExceptionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RuleExceptionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RuleExceptionDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.DetectionRuleId.IsNull() && !data.DetectionRuleId.IsUnknown() {
        filters["detectionRuleId"] = data.DetectionRuleId.ValueString()
        filterNames = append(filterNames, "detection_rule_id = "+fmt.Sprintf("%q", data.DetectionRuleId.ValueString()))
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        filters["kind"] = data.Kind.ValueString()
        filterNames = append(filterNames, "kind = "+fmt.Sprintf("%q", data.Kind.ValueString()))
    }
    if !data.EntityType.IsNull() && !data.EntityType.IsUnknown() {
        filters["entityType"] = data.EntityType.ValueString()
        filterNames = append(filterNames, "entity_type = "+fmt.Sprintf("%q", data.EntityType.ValueString()))
    }
    if !data.FieldName.IsNull() && !data.FieldName.IsUnknown() {
        filters["fieldName"] = data.FieldName.ValueString()
        filterNames = append(filterNames, "field_name = "+fmt.Sprintf("%q", data.FieldName.ValueString()))
    }
    if !data.Value.IsNull() && !data.Value.IsUnknown() {
        filters["value"] = data.Value.ValueString()
        filterNames = append(filterNames, "value = "+fmt.Sprintf("%q", data.Value.ValueString()))
    }
    if !data.Reason.IsNull() && !data.Reason.IsUnknown() {
        filters["reason"] = data.Reason.ValueString()
        filterNames = append(filterNames, "reason = "+fmt.Sprintf("%q", data.Reason.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the rule exception up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the rule exception up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "detectionRuleId": true,
        "kind": true,
        "entityType": true,
        "fieldName": true,
        "value": true,
        "startsAt": true,
        "endsAt": true,
        "reason": true,
        "expiresAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/detection-rule-exception/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rule_exception, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No rule exception found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read rule_exception: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/detection-rule-exception/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list rule_exception, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list rule_exception: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No rule exception matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one rule exception matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for rule_exception.")
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
    if obj, ok := item["entityType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EntityType = types.StringValue(string(jsonBytes))
        } else {
            data.EntityType = types.StringNull()
        }
    } else if val, ok := item["entityType"].(string); ok {
        data.EntityType = types.StringValue(val)
    } else {
        data.EntityType = types.StringNull()
    }
    if obj, ok := item["fieldName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FieldName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FieldName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FieldName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FieldName = types.StringValue(string(jsonBytes))
        } else {
            data.FieldName = types.StringNull()
        }
    } else if val, ok := item["fieldName"].(string); ok {
        data.FieldName = types.StringValue(val)
    } else {
        data.FieldName = types.StringNull()
    }
    if obj, ok := item["value"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Value = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if val, ok := item["value"].(string); ok {
        data.Value = types.StringValue(val)
    } else {
        data.Value = types.StringNull()
    }
    if obj, ok := item["startsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartsAt = types.StringValue(string(jsonBytes))
        } else {
            data.StartsAt = types.StringNull()
        }
    } else if val, ok := item["startsAt"].(string); ok {
        data.StartsAt = types.StringValue(val)
    } else {
        data.StartsAt = types.StringNull()
    }
    if obj, ok := item["endsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndsAt = types.StringValue(string(jsonBytes))
        } else {
            data.EndsAt = types.StringNull()
        }
    } else if val, ok := item["endsAt"].(string); ok {
        data.EndsAt = types.StringValue(val)
    } else {
        data.EndsAt = types.StringNull()
    }
    if obj, ok := item["reason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Reason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if val, ok := item["reason"].(string); ok {
        data.Reason = types.StringValue(val)
    } else {
        data.Reason = types.StringNull()
    }
    if obj, ok := item["expiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExpiresAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExpiresAt = types.StringValue(string(jsonBytes))
        } else {
            data.ExpiresAt = types.StringNull()
        }
    } else if val, ok := item["expiresAt"].(string); ok {
        data.ExpiresAt = types.StringValue(val)
    } else {
        data.ExpiresAt = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
