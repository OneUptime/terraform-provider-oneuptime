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
        MarkdownDescription: "The risk of each user, host and IP address: what the detection rules, threat intel and connected products found about it in the last 7 days adds up to. Look up an existing entity risk by `id`, or by any of its other arguments (`contributor_count`, `criticality`, `entity_key`, ...): each one set must match, and exactly one entity risk may match them all.",

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
                MarkdownDescription: "ID of the project this entity's risk belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "entity_type": schema.StringAttribute{
                MarkdownDescription: "What the entity is: a User, a Host or an IP Address.",
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
                MarkdownDescription: "What adds to its risk: each source with its severity, its points, how often it fired and when. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "mitre_tactics": schema.StringAttribute{
                MarkdownDescription: "ATT&CK tactic ids (TA0006) its sources are tagged with. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "mitre_techniques": schema.StringAttribute{
                MarkdownDescription: "ATT&CK technique ids (T1110) its sources are tagged with. A JSON value: write it with `jsonencode()`.",
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
                MarkdownDescription: "Its risk score at the end of each of the last 30 days. A JSON value: write it with `jsonencode()`.",
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
        filters["riskScore"] = lookupNumber(data.RiskScore)
        filterNames = append(filterNames, "risk_score = "+data.RiskScore.ValueBigFloat().String())
    }
    if !data.Criticality.IsNull() && !data.Criticality.IsUnknown() {
        filters["criticality"] = data.Criticality.ValueString()
        filterNames = append(filterNames, "criticality = "+fmt.Sprintf("%q", data.Criticality.ValueString()))
    }
    if !data.ContributorCount.IsNull() && !data.ContributorCount.IsUnknown() {
        filters["contributorCount"] = lookupNumber(data.ContributorCount)
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No entity risk matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one entity risk matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
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
    } else if val, ok := item["entityKey"].(string); ok {
        data.EntityKey = types.StringValue(val)
    } else {
        data.EntityKey = types.StringNull()
    }
    data.RiskScore = numberFromAPI(data.RiskScore, item["riskScore"])
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
    } else if val, ok := item["criticality"].(string); ok {
        data.Criticality = types.StringValue(val)
    } else {
        data.Criticality = types.StringNull()
    }
    data.ContributorCount = numberFromAPI(data.ContributorCount, item["contributorCount"])
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
