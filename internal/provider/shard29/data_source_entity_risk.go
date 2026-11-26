package shard29

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
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
var _ datasource.DataSource = &EntityRiskDataSource{}

func NewEntityRiskDataSource() datasource.DataSource {
    return &EntityRiskDataSource{}
}

// EntityRiskDataSource defines the data source implementation.
type EntityRiskDataSource struct {
    client *Client
}

// EntityRiskDataSourceModel describes the data source data model.
type EntityRiskDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    ProjectId types.String `tfsdk:"project_id"`
    EntityType types.String `tfsdk:"entity_type"`
    EntityValue types.String `tfsdk:"entity_value"`
    EntityKey types.String `tfsdk:"entity_key"`
    RiskScore types.Number `tfsdk:"risk_score"`
    Criticality types.String `tfsdk:"criticality"`
    ContributorCount types.Number `tfsdk:"contributor_count"`
    Contributions types.String `tfsdk:"contributions"`
    MitreTactics types.String `tfsdk:"mitre_tactics"`
    MitreTechniques types.String `tfsdk:"mitre_techniques"`
    FirstRiskAt types.String `tfsdk:"first_risk_at"`
    LastRiskAt types.String `tfsdk:"last_risk_at"`
    RiskExpiresAt types.String `tfsdk:"risk_expires_at"`
    RiskCaseAt types.String `tfsdk:"risk_case_at"`
    RiskHistory types.String `tfsdk:"risk_history"`
}

func (d *EntityRiskDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_entity_risk"
}

func (d *EntityRiskDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The risk of each user, host, IP address and data store: what the detection rules, threat intel, connected products and Data Security found about it in the last 7 days adds up to. Look up an existing entity risk by `id`, or by any of its other arguments (`contributor_count`, `created_by_actor_id`, `created_by_actor_type`, ...): each one set must match, and exactly one entity risk may match them all.",

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
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this entity's risk belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "entity_type": schema.StringAttribute{
                MarkdownDescription: "What the entity is: a User, a Host, an IP Address or a Data Store.",
                Optional: true,
                Computed: true,
            },
            "entity_value": schema.StringAttribute{
                MarkdownDescription: "The user's sign-in name, the host's name or the address, as the first finding about it named it.",
                Optional: true,
                Computed: true,
            },
            "entity_key": schema.StringAttribute{
                MarkdownDescription: "The entity's value, trimmed and in lower case: one entity's risk, whatever case its findings write it in.",
                Optional: true,
                Computed: true,
            },
            "risk_score": schema.NumberAttribute{
                MarkdownDescription: "What its sources of the last 7 days add up to, weighted by how much it matters. A case opens at the project's threshold.",
                Optional: true,
                Computed: true,
            },
            "criticality": schema.StringAttribute{
                MarkdownDescription: "How much it matters, from the project's Entity Criticality list: Low (x0.5), High (x1.5) or Extreme (x2). Empty counts as it is.",
                Optional: true,
                Computed: true,
            },
            "contributor_count": schema.NumberAttribute{
                MarkdownDescription: "How many different detection rules, threat intel feeds and connected products' rules add to its risk.",
                Optional: true,
                Computed: true,
            },
            "contributions": schema.StringAttribute{
                MarkdownDescription: "What adds to its risk: each source with its severity, its points, how often it fired and when. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "mitre_tactics": schema.StringAttribute{
                MarkdownDescription: "ATT&CK tactic ids (TA0006) its sources are tagged with. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "mitre_techniques": schema.StringAttribute{
                MarkdownDescription: "ATT&CK technique ids (T1110) its sources are tagged with. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "first_risk_at": schema.StringAttribute{
                MarkdownDescription: "When the sources that count now began counting.",
                Computed: true,
            },
            "last_risk_at": schema.StringAttribute{
                MarkdownDescription: "When a finding last added to its risk.",
                Computed: true,
            },
            "risk_expires_at": schema.StringAttribute{
                MarkdownDescription: "When its score next falls, as the first of its sources stops counting.",
                Computed: true,
            },
            "risk_case_at": schema.StringAttribute{
                MarkdownDescription: "When its risk last reached the threshold and opened a case, or joined the open case about it.",
                Computed: true,
            },
            "risk_history": schema.StringAttribute{
                MarkdownDescription: "When its risk score changed over the last 30 days, and to what: each day's score is read from it in the reader's time zone. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
        },
    }
}

func (d *EntityRiskDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EntityRiskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data EntityRiskDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.EntityType.IsNull() && !data.EntityType.IsUnknown() {
        filters["entityType"] = data.EntityType.ValueString()
        filterNames = append(filterNames, "entity_type = "+fmt.Sprintf("%q", data.EntityType.ValueString()))
    }
    if !data.EntityValue.IsNull() && !data.EntityValue.IsUnknown() {
        filters["entityValue"] = data.EntityValue.ValueString()
        filterNames = append(filterNames, "entity_value = "+fmt.Sprintf("%q", data.EntityValue.ValueString()))
    }
    if !data.EntityKey.IsNull() && !data.EntityKey.IsUnknown() {
        filters["entityKey"] = data.EntityKey.ValueString()
        filterNames = append(filterNames, "entity_key = "+fmt.Sprintf("%q", data.EntityKey.ValueString()))
    }
    if !data.RiskScore.IsNull() && !data.RiskScore.IsUnknown() {
        filters["riskScore"] = LookupNumber(data.RiskScore)
        filterNames = append(filterNames, "risk_score = "+data.RiskScore.ValueBigFloat().String())
    }
    if !data.Criticality.IsNull() && !data.Criticality.IsUnknown() {
        filters["criticality"] = data.Criticality.ValueString()
        filterNames = append(filterNames, "criticality = "+fmt.Sprintf("%q", data.Criticality.ValueString()))
    }
    if !data.ContributorCount.IsNull() && !data.ContributorCount.IsUnknown() {
        filters["contributorCount"] = LookupNumber(data.ContributorCount)
        filterNames = append(filterNames, "contributor_count = "+data.ContributorCount.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the entity risk up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the entity risk up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "projectId": true,
        "entityType": true,
        "entityValue": true,
        "entityKey": true,
        "riskScore": true,
        "criticality": true,
        "contributorCount": true,
        "contributions": true,
        "mitreTactics": true,
        "mitreTechniques": true,
        "firstRiskAt": true,
        "lastRiskAt": true,
        "riskExpiresAt": true,
        "riskCaseAt": true,
        "riskHistory": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-entity-risk/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read entity_risk, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No entity risk found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read entity_risk: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-entity-risk/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list entity_risk, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list entity_risk: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No entity risk matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one entity risk matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for entity_risk.")
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
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
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
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
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
    } else if arr, ok := item["entityType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EntityType = types.StringValue(string(jsonBytes))
        } else {
            data.EntityType = types.StringNull()
        }
    } else if val, ok := item["entityType"].(string); ok {
        data.EntityType = types.StringValue(val)
    } else {
        data.EntityType = types.StringNull()
    }
    if obj, ok := item["entityValue"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EntityValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EntityValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EntityValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EntityValue = types.StringValue(string(jsonBytes))
        } else {
            data.EntityValue = types.StringNull()
        }
    } else if arr, ok := item["entityValue"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EntityValue = types.StringValue(string(jsonBytes))
        } else {
            data.EntityValue = types.StringNull()
        }
    } else if val, ok := item["entityValue"].(string); ok {
        data.EntityValue = types.StringValue(val)
    } else {
        data.EntityValue = types.StringNull()
    }
    if obj, ok := item["entityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.EntityKey = types.StringNull()
        }
    } else if arr, ok := item["entityKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.EntityKey = types.StringNull()
        }
    } else if val, ok := item["entityKey"].(string); ok {
        data.EntityKey = types.StringValue(val)
    } else {
        data.EntityKey = types.StringNull()
    }
    data.RiskScore = NumberFromAPI(data.RiskScore, item["riskScore"])
    if obj, ok := item["criticality"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Criticality = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Criticality = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Criticality = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Criticality = types.StringValue(string(jsonBytes))
        } else {
            data.Criticality = types.StringNull()
        }
    } else if arr, ok := item["criticality"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Criticality = types.StringValue(string(jsonBytes))
        } else {
            data.Criticality = types.StringNull()
        }
    } else if val, ok := item["criticality"].(string); ok {
        data.Criticality = types.StringValue(val)
    } else {
        data.Criticality = types.StringNull()
    }
    data.ContributorCount = NumberFromAPI(data.ContributorCount, item["contributorCount"])
    if obj, ok := item["contributions"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Contributions = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Contributions = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Contributions = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Contributions = types.StringValue(string(jsonBytes))
        } else {
            data.Contributions = types.StringNull()
        }
    } else if arr, ok := item["contributions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Contributions = types.StringValue(string(jsonBytes))
        } else {
            data.Contributions = types.StringNull()
        }
    } else if val, ok := item["contributions"].(string); ok {
        data.Contributions = types.StringValue(val)
    } else {
        data.Contributions = types.StringNull()
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
    } else if arr, ok := item["mitreTactics"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["mitreTechniques"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MitreTechniques = types.StringValue(string(jsonBytes))
        } else {
            data.MitreTechniques = types.StringNull()
        }
    } else if val, ok := item["mitreTechniques"].(string); ok {
        data.MitreTechniques = types.StringValue(val)
    } else {
        data.MitreTechniques = types.StringNull()
    }
    if obj, ok := item["firstRiskAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstRiskAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstRiskAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstRiskAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstRiskAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstRiskAt = types.StringNull()
        }
    } else if arr, ok := item["firstRiskAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FirstRiskAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstRiskAt = types.StringNull()
        }
    } else if val, ok := item["firstRiskAt"].(string); ok {
        data.FirstRiskAt = types.StringValue(val)
    } else {
        data.FirstRiskAt = types.StringNull()
    }
    if obj, ok := item["lastRiskAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRiskAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRiskAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRiskAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRiskAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastRiskAt = types.StringNull()
        }
    } else if arr, ok := item["lastRiskAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastRiskAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastRiskAt = types.StringNull()
        }
    } else if val, ok := item["lastRiskAt"].(string); ok {
        data.LastRiskAt = types.StringValue(val)
    } else {
        data.LastRiskAt = types.StringNull()
    }
    if obj, ok := item["riskExpiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RiskExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RiskExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RiskExpiresAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RiskExpiresAt = types.StringValue(string(jsonBytes))
        } else {
            data.RiskExpiresAt = types.StringNull()
        }
    } else if arr, ok := item["riskExpiresAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RiskExpiresAt = types.StringValue(string(jsonBytes))
        } else {
            data.RiskExpiresAt = types.StringNull()
        }
    } else if val, ok := item["riskExpiresAt"].(string); ok {
        data.RiskExpiresAt = types.StringValue(val)
    } else {
        data.RiskExpiresAt = types.StringNull()
    }
    if obj, ok := item["riskCaseAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RiskCaseAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RiskCaseAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RiskCaseAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RiskCaseAt = types.StringValue(string(jsonBytes))
        } else {
            data.RiskCaseAt = types.StringNull()
        }
    } else if arr, ok := item["riskCaseAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RiskCaseAt = types.StringValue(string(jsonBytes))
        } else {
            data.RiskCaseAt = types.StringNull()
        }
    } else if val, ok := item["riskCaseAt"].(string); ok {
        data.RiskCaseAt = types.StringValue(val)
    } else {
        data.RiskCaseAt = types.StringNull()
    }
    if obj, ok := item["riskHistory"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RiskHistory = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RiskHistory = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RiskHistory = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RiskHistory = types.StringValue(string(jsonBytes))
        } else {
            data.RiskHistory = types.StringNull()
        }
    } else if arr, ok := item["riskHistory"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RiskHistory = types.StringValue(string(jsonBytes))
        } else {
            data.RiskHistory = types.StringNull()
        }
    } else if val, ok := item["riskHistory"].(string); ok {
        data.RiskHistory = types.StringValue(val)
    } else {
        data.RiskHistory = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
