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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IncidentResource{}
var _ resource.ResourceWithImportState = &IncidentResource{}

func NewIncidentResource() resource.Resource {
    return &IncidentResource{}
}

// IncidentResource defines the resource implementation.
type IncidentResource struct {
    client *Client
}

// IncidentResourceModel describes the resource data model.
type IncidentResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    DeclaredAt types.String `tfsdk:"declared_at"`
    Monitors types.List `tfsdk:"monitors"`
    OnCallDutyPolicies types.List `tfsdk:"on_call_duty_policies"`
    Labels types.List `tfsdk:"labels"`
    CurrentIncidentStateId types.String `tfsdk:"current_incident_state_id"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    ChangeMonitorStatusToId types.String `tfsdk:"change_monitor_status_to_id"`
    SubscriberNotificationStatusMessage types.String `tfsdk:"subscriber_notification_status_message"`
    SubscriberNotificationStatusMessageOnPostmortemPublished types.String `tfsdk:"subscriber_notification_status_message_on_postmortem_published"`
    ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_incident_created"`
    CustomFields types.String `tfsdk:"custom_fields"`
    RootCause types.String `tfsdk:"root_cause"`
    PostmortemNote types.String `tfsdk:"postmortem_note"`
    ShowPostmortemOnStatusPage types.Bool `tfsdk:"show_postmortem_on_status_page"`
    NotifySubscribersOnPostmortemPublished types.Bool `tfsdk:"notify_subscribers_on_postmortem_published"`
    PostmortemPostedAt types.String `tfsdk:"postmortem_posted_at"`
    PostmortemAttachments types.List `tfsdk:"postmortem_attachments"`
    RemediationNotes types.String `tfsdk:"remediation_notes"`
    TelemetryQuery types.String `tfsdk:"telemetry_query"`
    IsVisibleOnStatusPage types.Bool `tfsdk:"is_visible_on_status_page"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    SubscriberNotificationStatusOnIncidentCreated types.String `tfsdk:"subscriber_notification_status_on_incident_created"`
    SubscriberNotificationStatusOnPostmortemPublished types.String `tfsdk:"subscriber_notification_status_on_postmortem_published"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    CreatedStateLog types.String `tfsdk:"created_state_log"`
    CreatedCriteriaId types.String `tfsdk:"created_criteria_id"`
    CreatedIncidentTemplateId types.String `tfsdk:"created_incident_template_id"`
    CreatedByProbeId types.String `tfsdk:"created_by_probe_id"`
    IsCreatedAutomatically types.Bool `tfsdk:"is_created_automatically"`
    IncidentNumber types.Number `tfsdk:"incident_number"`
}

func (r *IncidentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident"
}

func (r *IncidentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident resource",

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
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Title of this incident. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Short description of this incident. This is in markdown and will be visible on the status page.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "declared_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "monitors": schema.ListAttribute{
                MarkdownDescription: "List of monitors affected by this incident. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.List{
                    listplanmodifier.UseStateForUnknown(),
                },
            },
            "on_call_duty_policies": schema.ListAttribute{
                MarkdownDescription: "List of on-call duty policy affected by this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.List{
                    listplanmodifier.UseStateForUnknown(),
                },
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.List{
                    listplanmodifier.UseStateForUnknown(),
                },
            },
            "current_incident_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "subscriber_notification_status_message": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "subscriber_notification_status_message_on_postmortem_published": schema.StringAttribute{
                MarkdownDescription: "Status message for subscriber notifications on postmortem published - includes success messages, failure reasons, or skip reasons. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "should_status_page_subscribers_be_notified_on_incident_created": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified about this incident?. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "What is the root cause of this incident?. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "postmortem_note": schema.StringAttribute{
                MarkdownDescription: "Document the postmortem summary for this incident.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "show_postmortem_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should the postmortem note and attachments be visible on the status page once published?. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "notify_subscribers_on_postmortem_published": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified when the postmortem is published?. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "postmortem_posted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "postmortem_attachments": schema.ListAttribute{
                MarkdownDescription: "Files that accompany the postmortem note and can be shared publicly when enabled.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.List{
                    listplanmodifier.UseStateForUnknown(),
                },
            },
            "remediation_notes": schema.StringAttribute{
                MarkdownDescription: "Notes on how to remediate this incident. This is in markdown.. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "telemetry_query": schema.StringAttribute{
                MarkdownDescription: "Telemetry query for this incident. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_visible_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should this incident be visible on the status page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "subscriber_notification_status_on_incident_created": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers about this incident. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
            },
            "subscriber_notification_status_on_postmortem_published": schema.StringAttribute{
                MarkdownDescription: "Status of notification sent to subscribers about this incident postmortem. Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [Project Owner, Project Admin, Project Member, Edit Incident]",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of when this resource is created?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_state_log": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_criteria_id": schema.StringAttribute{
                MarkdownDescription: "If this incident was created by a Probe, this is the ID of the criteria that created it.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_incident_template_id": schema.StringAttribute{
                MarkdownDescription: "If this incident was created by a Probe, this is the ID of the incident template that was used for creation.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_probe_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_created_automatically": schema.BoolAttribute{
                MarkdownDescription: "Is this incident created by OneUptime Probe or Workers automatically (and not created manually by a user)?. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "incident_number": schema.NumberAttribute{
                MarkdownDescription: "Incident Number. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *IncidentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *IncidentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data IncidentResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body
    incidentRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "title": data.Title.ValueString(),
        "description": data.Description.ValueString(),
        "declaredAt": r.parseJSONField(data.DeclaredAt),
        "monitors": r.convertTerraformListToInterface(data.Monitors),
        "onCallDutyPolicies": r.convertTerraformListToInterface(data.OnCallDutyPolicies),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "currentIncidentStateId": data.CurrentIncidentStateId.ValueString(),
        "incidentSeverityId": data.IncidentSeverityId.ValueString(),
        "changeMonitorStatusToId": data.ChangeMonitorStatusToId.ValueString(),
        "subscriberNotificationStatusMessage": data.SubscriberNotificationStatusMessage.ValueString(),
        "subscriberNotificationStatusMessageOnPostmortemPublished": data.SubscriberNotificationStatusMessageOnPostmortemPublished.ValueString(),
        "shouldStatusPageSubscribersBeNotifiedOnIncidentCreated": data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated.ValueBool(),
        "customFields": r.parseJSONField(data.CustomFields),
        "rootCause": data.RootCause.ValueString(),
        "postmortemNote": data.PostmortemNote.ValueString(),
        "showPostmortemOnStatusPage": data.ShowPostmortemOnStatusPage.ValueBool(),
        "notifySubscribersOnPostmortemPublished": data.NotifySubscribersOnPostmortemPublished.ValueBool(),
        "postmortemPostedAt": r.parseJSONField(data.PostmortemPostedAt),
        "postmortemAttachments": r.convertTerraformListToInterface(data.PostmortemAttachments),
        "remediationNotes": data.RemediationNotes.ValueString(),
        "telemetryQuery": r.parseJSONField(data.TelemetryQuery),
        "isVisibleOnStatusPage": data.IsVisibleOnStatusPage.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/incident", incidentRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create incident, got error: %s", err))
        return
    }

    var incidentResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incidentResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incidentResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incidentResponse
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
    if obj, ok := dataMap["title"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["declaredAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeclaredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeclaredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.DeclaredAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.DeclaredAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeclaredAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeclaredAt = types.StringNull()
        }
    } else if val, ok := dataMap["declaredAt"].(string); ok && val != "" {
        data.DeclaredAt = types.StringValue(val)
    } else {
        data.DeclaredAt = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Monitors = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.OnCallDutyPolicies = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.OnCallDutyPolicies = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Labels = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["currentIncidentStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CurrentIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentIncidentStateId = types.StringNull()
        }
    } else if val, ok := dataMap["currentIncidentStateId"].(string); ok && val != "" {
        data.CurrentIncidentStateId = types.StringValue(val)
    } else {
        data.CurrentIncidentStateId = types.StringNull()
    }
    if obj, ok := dataMap["incidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentSeverityId"].(string); ok && val != "" {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["changeMonitorStatusToId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChangeMonitorStatusToId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ChangeMonitorStatusToId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
            } else {
                data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
            } else {
                data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
        } else {
            data.ChangeMonitorStatusToId = types.StringNull()
        }
    } else if val, ok := dataMap["changeMonitorStatusToId"].(string); ok && val != "" {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessage"].(string); ok && val != "" {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusMessageOnPostmortemPublished"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessageOnPostmortemPublished"].(string); ok && val != "" {
        data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolValue(val)
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok && val != "" {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
    }
    if obj, ok := dataMap["rootCause"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RootCause = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.RootCause = types.StringValue(string(jsonBytes))
            } else {
                data.RootCause = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.RootCause = types.StringValue(string(jsonBytes))
            } else {
                data.RootCause = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RootCause = types.StringValue(string(jsonBytes))
        } else {
            data.RootCause = types.StringNull()
        }
    } else if val, ok := dataMap["rootCause"].(string); ok && val != "" {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if obj, ok := dataMap["postmortemNote"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostmortemNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostmortemNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostmortemNote = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostmortemNote = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostmortemNote = types.StringValue(string(jsonBytes))
        } else {
            data.PostmortemNote = types.StringNull()
        }
    } else if val, ok := dataMap["postmortemNote"].(string); ok && val != "" {
        data.PostmortemNote = types.StringValue(val)
    } else {
        data.PostmortemNote = types.StringNull()
    }
    if val, ok := dataMap["showPostmortemOnStatusPage"].(bool); ok {
        data.ShowPostmortemOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["notifySubscribersOnPostmortemPublished"].(bool); ok {
        data.NotifySubscribersOnPostmortemPublished = types.BoolValue(val)
    }
    if obj, ok := dataMap["postmortemPostedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostmortemPostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostmortemPostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostmortemPostedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostmortemPostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemPostedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostmortemPostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemPostedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostmortemPostedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PostmortemPostedAt = types.StringNull()
        }
    } else if val, ok := dataMap["postmortemPostedAt"].(string); ok && val != "" {
        data.PostmortemPostedAt = types.StringValue(val)
    } else {
        data.PostmortemPostedAt = types.StringNull()
    }
    if val, ok := dataMap["postmortemAttachments"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.PostmortemAttachments = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.PostmortemAttachments = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["remediationNotes"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.RemediationNotes = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.RemediationNotes = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RemediationNotes = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationNotes = types.StringNull()
        }
    } else if val, ok := dataMap["remediationNotes"].(string); ok && val != "" {
        data.RemediationNotes = types.StringValue(val)
    } else {
        data.RemediationNotes = types.StringNull()
    }
    if obj, ok := dataMap["telemetryQuery"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryQuery = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.TelemetryQuery = types.StringValue(string(jsonBytes))
            } else {
                data.TelemetryQuery = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.TelemetryQuery = types.StringValue(string(jsonBytes))
            } else {
                data.TelemetryQuery = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryQuery = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryQuery = types.StringNull()
        }
    } else if val, ok := dataMap["telemetryQuery"].(string); ok && val != "" {
        data.TelemetryQuery = types.StringValue(val)
    } else {
        data.TelemetryQuery = types.StringNull()
    }
    if val, ok := dataMap["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
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
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusOnIncidentCreated"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusOnIncidentCreated"].(string); ok && val != "" {
        data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnIncidentCreated = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusOnPostmortemPublished"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusOnPostmortemPublished"].(string); ok && val != "" {
        data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnPostmortemPublished = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if obj, ok := dataMap["createdStateLog"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedStateLog = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedStateLog = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedStateLog = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedStateLog = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedStateLog = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedStateLog = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedStateLog = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedStateLog = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedStateLog = types.StringNull()
        }
    } else if val, ok := dataMap["createdStateLog"].(string); ok && val != "" {
        data.CreatedStateLog = types.StringValue(val)
    } else {
        data.CreatedStateLog = types.StringNull()
    }
    if obj, ok := dataMap["createdCriteriaId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedCriteriaId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedCriteriaId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedCriteriaId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedCriteriaId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedCriteriaId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedCriteriaId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedCriteriaId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedCriteriaId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedCriteriaId = types.StringNull()
        }
    } else if val, ok := dataMap["createdCriteriaId"].(string); ok && val != "" {
        data.CreatedCriteriaId = types.StringValue(val)
    } else {
        data.CreatedCriteriaId = types.StringNull()
    }
    if obj, ok := dataMap["createdIncidentTemplateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedIncidentTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedIncidentTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedIncidentTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedIncidentTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedIncidentTemplateId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedIncidentTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedIncidentTemplateId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedIncidentTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedIncidentTemplateId = types.StringNull()
        }
    } else if val, ok := dataMap["createdIncidentTemplateId"].(string); ok && val != "" {
        data.CreatedIncidentTemplateId = types.StringValue(val)
    } else {
        data.CreatedIncidentTemplateId = types.StringNull()
    }
    if obj, ok := dataMap["createdByProbeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedByProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByProbeId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedByProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByProbeId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByProbeId"].(string); ok && val != "" {
        data.CreatedByProbeId = types.StringValue(val)
    } else {
        data.CreatedByProbeId = types.StringNull()
    }
    if val, ok := dataMap["isCreatedAutomatically"].(bool); ok {
        data.IsCreatedAutomatically = types.BoolValue(val)
    }
    if val, ok := dataMap["incidentNumber"].(float64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentNumber"].(int); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentNumber"].(int64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["incidentNumber"] == nil {
        data.IncidentNumber = types.NumberNull()
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

func (r *IncidentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data IncidentResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "title": true,
        "description": true,
        "declaredAt": true,
        "monitors": true,
        "onCallDutyPolicies": true,
        "labels": true,
        "currentIncidentStateId": true,
        "incidentSeverityId": true,
        "changeMonitorStatusToId": true,
        "subscriberNotificationStatusMessage": true,
        "subscriberNotificationStatusMessageOnPostmortemPublished": true,
        "shouldStatusPageSubscribersBeNotifiedOnIncidentCreated": true,
        "customFields": true,
        "rootCause": true,
        "postmortemNote": true,
        "showPostmortemOnStatusPage": true,
        "notifySubscribersOnPostmortemPublished": true,
        "postmortemPostedAt": true,
        "postmortemAttachments": true,
        "remediationNotes": true,
        "telemetryQuery": true,
        "isVisibleOnStatusPage": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "subscriberNotificationStatusOnIncidentCreated": true,
        "subscriberNotificationStatusOnPostmortemPublished": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "createdStateLog": true,
        "createdCriteriaId": true,
        "createdIncidentTemplateId": true,
        "createdByProbeId": true,
        "isCreatedAutomatically": true,
        "incidentNumber": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/incident/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var incidentResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incidentResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incidentResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incidentResponse
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
    if obj, ok := dataMap["title"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["declaredAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeclaredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeclaredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.DeclaredAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.DeclaredAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeclaredAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeclaredAt = types.StringNull()
        }
    } else if val, ok := dataMap["declaredAt"].(string); ok && val != "" {
        data.DeclaredAt = types.StringValue(val)
    } else {
        data.DeclaredAt = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Monitors = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.OnCallDutyPolicies = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.OnCallDutyPolicies = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Labels = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["currentIncidentStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CurrentIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentIncidentStateId = types.StringNull()
        }
    } else if val, ok := dataMap["currentIncidentStateId"].(string); ok && val != "" {
        data.CurrentIncidentStateId = types.StringValue(val)
    } else {
        data.CurrentIncidentStateId = types.StringNull()
    }
    if obj, ok := dataMap["incidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentSeverityId"].(string); ok && val != "" {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["changeMonitorStatusToId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChangeMonitorStatusToId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ChangeMonitorStatusToId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
            } else {
                data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
            } else {
                data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
        } else {
            data.ChangeMonitorStatusToId = types.StringNull()
        }
    } else if val, ok := dataMap["changeMonitorStatusToId"].(string); ok && val != "" {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessage"].(string); ok && val != "" {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusMessageOnPostmortemPublished"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessageOnPostmortemPublished"].(string); ok && val != "" {
        data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolValue(val)
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok && val != "" {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
    }
    if obj, ok := dataMap["rootCause"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RootCause = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.RootCause = types.StringValue(string(jsonBytes))
            } else {
                data.RootCause = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.RootCause = types.StringValue(string(jsonBytes))
            } else {
                data.RootCause = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RootCause = types.StringValue(string(jsonBytes))
        } else {
            data.RootCause = types.StringNull()
        }
    } else if val, ok := dataMap["rootCause"].(string); ok && val != "" {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if obj, ok := dataMap["postmortemNote"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostmortemNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostmortemNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostmortemNote = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostmortemNote = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostmortemNote = types.StringValue(string(jsonBytes))
        } else {
            data.PostmortemNote = types.StringNull()
        }
    } else if val, ok := dataMap["postmortemNote"].(string); ok && val != "" {
        data.PostmortemNote = types.StringValue(val)
    } else {
        data.PostmortemNote = types.StringNull()
    }
    if val, ok := dataMap["showPostmortemOnStatusPage"].(bool); ok {
        data.ShowPostmortemOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["notifySubscribersOnPostmortemPublished"].(bool); ok {
        data.NotifySubscribersOnPostmortemPublished = types.BoolValue(val)
    }
    if obj, ok := dataMap["postmortemPostedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostmortemPostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostmortemPostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostmortemPostedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostmortemPostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemPostedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostmortemPostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemPostedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostmortemPostedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PostmortemPostedAt = types.StringNull()
        }
    } else if val, ok := dataMap["postmortemPostedAt"].(string); ok && val != "" {
        data.PostmortemPostedAt = types.StringValue(val)
    } else {
        data.PostmortemPostedAt = types.StringNull()
    }
    if val, ok := dataMap["postmortemAttachments"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.PostmortemAttachments = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.PostmortemAttachments = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["remediationNotes"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.RemediationNotes = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.RemediationNotes = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RemediationNotes = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationNotes = types.StringNull()
        }
    } else if val, ok := dataMap["remediationNotes"].(string); ok && val != "" {
        data.RemediationNotes = types.StringValue(val)
    } else {
        data.RemediationNotes = types.StringNull()
    }
    if obj, ok := dataMap["telemetryQuery"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryQuery = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.TelemetryQuery = types.StringValue(string(jsonBytes))
            } else {
                data.TelemetryQuery = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.TelemetryQuery = types.StringValue(string(jsonBytes))
            } else {
                data.TelemetryQuery = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryQuery = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryQuery = types.StringNull()
        }
    } else if val, ok := dataMap["telemetryQuery"].(string); ok && val != "" {
        data.TelemetryQuery = types.StringValue(val)
    } else {
        data.TelemetryQuery = types.StringNull()
    }
    if val, ok := dataMap["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
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
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusOnIncidentCreated"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusOnIncidentCreated"].(string); ok && val != "" {
        data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnIncidentCreated = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusOnPostmortemPublished"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusOnPostmortemPublished"].(string); ok && val != "" {
        data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnPostmortemPublished = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if obj, ok := dataMap["createdStateLog"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedStateLog = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedStateLog = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedStateLog = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedStateLog = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedStateLog = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedStateLog = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedStateLog = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedStateLog = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedStateLog = types.StringNull()
        }
    } else if val, ok := dataMap["createdStateLog"].(string); ok && val != "" {
        data.CreatedStateLog = types.StringValue(val)
    } else {
        data.CreatedStateLog = types.StringNull()
    }
    if obj, ok := dataMap["createdCriteriaId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedCriteriaId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedCriteriaId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedCriteriaId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedCriteriaId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedCriteriaId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedCriteriaId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedCriteriaId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedCriteriaId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedCriteriaId = types.StringNull()
        }
    } else if val, ok := dataMap["createdCriteriaId"].(string); ok && val != "" {
        data.CreatedCriteriaId = types.StringValue(val)
    } else {
        data.CreatedCriteriaId = types.StringNull()
    }
    if obj, ok := dataMap["createdIncidentTemplateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedIncidentTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedIncidentTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedIncidentTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedIncidentTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedIncidentTemplateId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedIncidentTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedIncidentTemplateId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedIncidentTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedIncidentTemplateId = types.StringNull()
        }
    } else if val, ok := dataMap["createdIncidentTemplateId"].(string); ok && val != "" {
        data.CreatedIncidentTemplateId = types.StringValue(val)
    } else {
        data.CreatedIncidentTemplateId = types.StringNull()
    }
    if obj, ok := dataMap["createdByProbeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedByProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByProbeId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedByProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByProbeId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByProbeId"].(string); ok && val != "" {
        data.CreatedByProbeId = types.StringValue(val)
    } else {
        data.CreatedByProbeId = types.StringNull()
    }
    if val, ok := dataMap["isCreatedAutomatically"].(bool); ok {
        data.IsCreatedAutomatically = types.BoolValue(val)
    }
    if val, ok := dataMap["incidentNumber"].(float64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentNumber"].(int); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentNumber"].(int64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["incidentNumber"] == nil {
        data.IncidentNumber = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncidentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data IncidentResourceModel
    var state IncidentResourceModel

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
    incidentRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := incidentRequest["data"].(map[string]interface{})

    if !data.Title.IsUnknown() && !state.Title.IsUnknown() && !data.Title.Equal(state.Title) {
        requestDataMap["title"] = data.Title.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.DeclaredAt.IsUnknown() && !state.DeclaredAt.IsUnknown() && !data.DeclaredAt.Equal(state.DeclaredAt) {
        var declaredatData interface{}
        if err := json.Unmarshal([]byte(data.DeclaredAt.ValueString()), &declaredatData); err == nil {
            requestDataMap["declaredAt"] = declaredatData
        } else {
            requestDataMap["declaredAt"] = data.DeclaredAt.ValueString()
        }
    }
    if !data.Monitors.IsUnknown() && !state.Monitors.IsUnknown() && !data.Monitors.Equal(state.Monitors) {
        requestDataMap["monitors"] = r.convertTerraformListToInterface(data.Monitors)
    }
    if !data.OnCallDutyPolicies.IsUnknown() && !state.OnCallDutyPolicies.IsUnknown() && !data.OnCallDutyPolicies.Equal(state.OnCallDutyPolicies) {
        requestDataMap["onCallDutyPolicies"] = r.convertTerraformListToInterface(data.OnCallDutyPolicies)
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformListToInterface(data.Labels)
    }
    if !data.CurrentIncidentStateId.IsUnknown() && !state.CurrentIncidentStateId.IsUnknown() && !data.CurrentIncidentStateId.Equal(state.CurrentIncidentStateId) {
        requestDataMap["currentIncidentStateId"] = data.CurrentIncidentStateId.ValueString()
    }
    if !data.IncidentSeverityId.IsUnknown() && !state.IncidentSeverityId.IsUnknown() && !data.IncidentSeverityId.Equal(state.IncidentSeverityId) {
        requestDataMap["incidentSeverityId"] = data.IncidentSeverityId.ValueString()
    }
    if !data.ChangeMonitorStatusToId.IsUnknown() && !state.ChangeMonitorStatusToId.IsUnknown() && !data.ChangeMonitorStatusToId.Equal(state.ChangeMonitorStatusToId) {
        requestDataMap["changeMonitorStatusToId"] = data.ChangeMonitorStatusToId.ValueString()
    }
    if !data.SubscriberNotificationStatusMessage.IsUnknown() && !state.SubscriberNotificationStatusMessage.IsUnknown() && !data.SubscriberNotificationStatusMessage.Equal(state.SubscriberNotificationStatusMessage) {
        requestDataMap["subscriberNotificationStatusMessage"] = data.SubscriberNotificationStatusMessage.ValueString()
    }
    if !data.SubscriberNotificationStatusMessageOnPostmortemPublished.IsUnknown() && !state.SubscriberNotificationStatusMessageOnPostmortemPublished.IsUnknown() && !data.SubscriberNotificationStatusMessageOnPostmortemPublished.Equal(state.SubscriberNotificationStatusMessageOnPostmortemPublished) {
        requestDataMap["subscriberNotificationStatusMessageOnPostmortemPublished"] = data.SubscriberNotificationStatusMessageOnPostmortemPublished.ValueString()
    }
    if !data.CustomFields.IsUnknown() && !state.CustomFields.IsUnknown() && !data.CustomFields.Equal(state.CustomFields) {
        var customfieldsData interface{}
        if err := json.Unmarshal([]byte(data.CustomFields.ValueString()), &customfieldsData); err == nil {
            requestDataMap["customFields"] = customfieldsData
        } else {
            requestDataMap["customFields"] = data.CustomFields.ValueString()
        }
    }
    if !data.RootCause.IsUnknown() && !state.RootCause.IsUnknown() && !data.RootCause.Equal(state.RootCause) {
        requestDataMap["rootCause"] = data.RootCause.ValueString()
    }
    if !data.PostmortemNote.IsUnknown() && !state.PostmortemNote.IsUnknown() && !data.PostmortemNote.Equal(state.PostmortemNote) {
        requestDataMap["postmortemNote"] = data.PostmortemNote.ValueString()
    }
    if !data.ShowPostmortemOnStatusPage.IsUnknown() && !state.ShowPostmortemOnStatusPage.IsUnknown() && !data.ShowPostmortemOnStatusPage.Equal(state.ShowPostmortemOnStatusPage) {
        requestDataMap["showPostmortemOnStatusPage"] = data.ShowPostmortemOnStatusPage.ValueBool()
    }
    if !data.NotifySubscribersOnPostmortemPublished.IsUnknown() && !state.NotifySubscribersOnPostmortemPublished.IsUnknown() && !data.NotifySubscribersOnPostmortemPublished.Equal(state.NotifySubscribersOnPostmortemPublished) {
        requestDataMap["notifySubscribersOnPostmortemPublished"] = data.NotifySubscribersOnPostmortemPublished.ValueBool()
    }
    if !data.PostmortemPostedAt.IsUnknown() && !state.PostmortemPostedAt.IsUnknown() && !data.PostmortemPostedAt.Equal(state.PostmortemPostedAt) {
        var postmortempostedatData interface{}
        if err := json.Unmarshal([]byte(data.PostmortemPostedAt.ValueString()), &postmortempostedatData); err == nil {
            requestDataMap["postmortemPostedAt"] = postmortempostedatData
        } else {
            requestDataMap["postmortemPostedAt"] = data.PostmortemPostedAt.ValueString()
        }
    }
    if !data.PostmortemAttachments.IsUnknown() && !state.PostmortemAttachments.IsUnknown() && !data.PostmortemAttachments.Equal(state.PostmortemAttachments) {
        requestDataMap["postmortemAttachments"] = r.convertTerraformListToInterface(data.PostmortemAttachments)
    }
    if !data.RemediationNotes.IsUnknown() && !state.RemediationNotes.IsUnknown() && !data.RemediationNotes.Equal(state.RemediationNotes) {
        requestDataMap["remediationNotes"] = data.RemediationNotes.ValueString()
    }
    if !data.TelemetryQuery.IsUnknown() && !state.TelemetryQuery.IsUnknown() && !data.TelemetryQuery.Equal(state.TelemetryQuery) {
        var telemetryqueryData interface{}
        if err := json.Unmarshal([]byte(data.TelemetryQuery.ValueString()), &telemetryqueryData); err == nil {
            requestDataMap["telemetryQuery"] = telemetryqueryData
        } else {
            requestDataMap["telemetryQuery"] = data.TelemetryQuery.ValueString()
        }
    }
    if !data.IsVisibleOnStatusPage.IsUnknown() && !state.IsVisibleOnStatusPage.IsUnknown() && !data.IsVisibleOnStatusPage.Equal(state.IsVisibleOnStatusPage) {
        requestDataMap["isVisibleOnStatusPage"] = data.IsVisibleOnStatusPage.ValueBool()
    }

    // Make API call
    httpResp, err := r.client.Put("/incident/" + data.Id.ValueString() + "", incidentRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update incident, got error: %s", err))
        return
    }

    // Parse the update response
    var incidentResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incidentResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "title": true,
        "description": true,
        "declaredAt": true,
        "monitors": true,
        "onCallDutyPolicies": true,
        "labels": true,
        "currentIncidentStateId": true,
        "incidentSeverityId": true,
        "changeMonitorStatusToId": true,
        "subscriberNotificationStatusMessage": true,
        "subscriberNotificationStatusMessageOnPostmortemPublished": true,
        "shouldStatusPageSubscribersBeNotifiedOnIncidentCreated": true,
        "customFields": true,
        "rootCause": true,
        "postmortemNote": true,
        "showPostmortemOnStatusPage": true,
        "notifySubscribersOnPostmortemPublished": true,
        "postmortemPostedAt": true,
        "postmortemAttachments": true,
        "remediationNotes": true,
        "telemetryQuery": true,
        "isVisibleOnStatusPage": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "subscriberNotificationStatusOnIncidentCreated": true,
        "subscriberNotificationStatusOnPostmortemPublished": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "createdStateLog": true,
        "createdCriteriaId": true,
        "createdIncidentTemplateId": true,
        "createdByProbeId": true,
        "isCreatedAutomatically": true,
        "incidentNumber": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/incident/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident read response, got error: %s", err))
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
    if obj, ok := dataMap["title"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Title = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Title = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Title = types.StringValue(string(jsonBytes))
            } else {
                data.Title = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Title = types.StringValue(string(jsonBytes))
        } else {
            data.Title = types.StringNull()
        }
    } else if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["declaredAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeclaredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeclaredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.DeclaredAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.DeclaredAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeclaredAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeclaredAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeclaredAt = types.StringNull()
        }
    } else if val, ok := dataMap["declaredAt"].(string); ok && val != "" {
        data.DeclaredAt = types.StringValue(val)
    } else {
        data.DeclaredAt = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Monitors = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Monitors = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.OnCallDutyPolicies = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.OnCallDutyPolicies = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.Labels = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["currentIncidentStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CurrentIncidentStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentIncidentStateId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CurrentIncidentStateId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentIncidentStateId = types.StringNull()
        }
    } else if val, ok := dataMap["currentIncidentStateId"].(string); ok && val != "" {
        data.CurrentIncidentStateId = types.StringValue(val)
    } else {
        data.CurrentIncidentStateId = types.StringNull()
    }
    if obj, ok := dataMap["incidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentSeverityId"].(string); ok && val != "" {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["changeMonitorStatusToId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChangeMonitorStatusToId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ChangeMonitorStatusToId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
            } else {
                data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
            } else {
                data.ChangeMonitorStatusToId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ChangeMonitorStatusToId = types.StringValue(string(jsonBytes))
        } else {
            data.ChangeMonitorStatusToId = types.StringNull()
        }
    } else if val, ok := dataMap["changeMonitorStatusToId"].(string); ok && val != "" {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusMessage"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessage = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessage = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessage"].(string); ok && val != "" {
        data.SubscriberNotificationStatusMessage = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessage = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusMessageOnPostmortemPublished"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusMessageOnPostmortemPublished"].(string); ok && val != "" {
        data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusMessageOnPostmortemPublished = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolValue(val)
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CustomFields = types.StringValue(string(jsonBytes))
            } else {
                data.CustomFields = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok && val != "" {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
    }
    if obj, ok := dataMap["rootCause"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RootCause = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RootCause = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.RootCause = types.StringValue(string(jsonBytes))
            } else {
                data.RootCause = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.RootCause = types.StringValue(string(jsonBytes))
            } else {
                data.RootCause = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RootCause = types.StringValue(string(jsonBytes))
        } else {
            data.RootCause = types.StringNull()
        }
    } else if val, ok := dataMap["rootCause"].(string); ok && val != "" {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if obj, ok := dataMap["postmortemNote"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostmortemNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostmortemNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostmortemNote = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostmortemNote = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemNote = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostmortemNote = types.StringValue(string(jsonBytes))
        } else {
            data.PostmortemNote = types.StringNull()
        }
    } else if val, ok := dataMap["postmortemNote"].(string); ok && val != "" {
        data.PostmortemNote = types.StringValue(val)
    } else {
        data.PostmortemNote = types.StringNull()
    }
    if val, ok := dataMap["showPostmortemOnStatusPage"].(bool); ok {
        data.ShowPostmortemOnStatusPage = types.BoolValue(val)
    }
    if val, ok := dataMap["notifySubscribersOnPostmortemPublished"].(bool); ok {
        data.NotifySubscribersOnPostmortemPublished = types.BoolValue(val)
    }
    if obj, ok := dataMap["postmortemPostedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostmortemPostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PostmortemPostedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PostmortemPostedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.PostmortemPostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemPostedAt = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.PostmortemPostedAt = types.StringValue(string(jsonBytes))
            } else {
                data.PostmortemPostedAt = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PostmortemPostedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PostmortemPostedAt = types.StringNull()
        }
    } else if val, ok := dataMap["postmortemPostedAt"].(string); ok && val != "" {
        data.PostmortemPostedAt = types.StringValue(val)
    } else {
        data.PostmortemPostedAt = types.StringNull()
    }
    if val, ok := dataMap["postmortemAttachments"].([]interface{}); ok {
        // Convert API response list to Terraform list
        var listItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    listItems = append(listItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        listItems = append(listItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                listItems = append(listItems, types.StringValue(str))
            }
        }
        // Sort list items by their string value to ensure consistent ordering
        // This fixes idempotency issues where server returns items in different order
        sort.Slice(listItems, func(i, j int) bool {
            iStr := listItems[i].(types.String).ValueString()
            jStr := listItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.PostmortemAttachments = types.ListValueMust(types.StringType, listItems)
    } else {
        // For lists, always use empty list instead of null to match default values
        data.PostmortemAttachments = types.ListValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["remediationNotes"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RemediationNotes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.RemediationNotes = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.RemediationNotes = types.StringValue(string(jsonBytes))
            } else {
                data.RemediationNotes = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RemediationNotes = types.StringValue(string(jsonBytes))
        } else {
            data.RemediationNotes = types.StringNull()
        }
    } else if val, ok := dataMap["remediationNotes"].(string); ok && val != "" {
        data.RemediationNotes = types.StringValue(val)
    } else {
        data.RemediationNotes = types.StringNull()
    }
    if obj, ok := dataMap["telemetryQuery"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryQuery = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.TelemetryQuery = types.StringValue(string(jsonBytes))
            } else {
                data.TelemetryQuery = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.TelemetryQuery = types.StringValue(string(jsonBytes))
            } else {
                data.TelemetryQuery = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryQuery = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryQuery = types.StringNull()
        }
    } else if val, ok := dataMap["telemetryQuery"].(string); ok && val != "" {
        data.TelemetryQuery = types.StringValue(val)
    } else {
        data.TelemetryQuery = types.StringNull()
    }
    if val, ok := dataMap["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
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
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusOnIncidentCreated"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnIncidentCreated = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusOnIncidentCreated"].(string); ok && val != "" {
        data.SubscriberNotificationStatusOnIncidentCreated = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnIncidentCreated = types.StringNull()
    }
    if obj, ok := dataMap["subscriberNotificationStatusOnPostmortemPublished"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberNotificationStatusOnPostmortemPublished = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberNotificationStatusOnPostmortemPublished"].(string); ok && val != "" {
        data.SubscriberNotificationStatusOnPostmortemPublished = types.StringValue(val)
    } else {
        data.SubscriberNotificationStatusOnPostmortemPublished = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if obj, ok := dataMap["createdStateLog"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedStateLog = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedStateLog = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedStateLog = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedStateLog = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedStateLog = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedStateLog = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedStateLog = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedStateLog = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedStateLog = types.StringNull()
        }
    } else if val, ok := dataMap["createdStateLog"].(string); ok && val != "" {
        data.CreatedStateLog = types.StringValue(val)
    } else {
        data.CreatedStateLog = types.StringNull()
    }
    if obj, ok := dataMap["createdCriteriaId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedCriteriaId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedCriteriaId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedCriteriaId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedCriteriaId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedCriteriaId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedCriteriaId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedCriteriaId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedCriteriaId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedCriteriaId = types.StringNull()
        }
    } else if val, ok := dataMap["createdCriteriaId"].(string); ok && val != "" {
        data.CreatedCriteriaId = types.StringValue(val)
    } else {
        data.CreatedCriteriaId = types.StringNull()
    }
    if obj, ok := dataMap["createdIncidentTemplateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedIncidentTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedIncidentTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedIncidentTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedIncidentTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedIncidentTemplateId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedIncidentTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedIncidentTemplateId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedIncidentTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedIncidentTemplateId = types.StringNull()
        }
    } else if val, ok := dataMap["createdIncidentTemplateId"].(string); ok && val != "" {
        data.CreatedIncidentTemplateId = types.StringValue(val)
    } else {
        data.CreatedIncidentTemplateId = types.StringNull()
    }
    if obj, ok := dataMap["createdByProbeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            if jsonBytes, err := json.Marshal(obj); err == nil {
                data.CreatedByProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByProbeId = types.StringValue(fmt.Sprintf("%v", obj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            if jsonBytes, err := json.Marshal(obj["value"]); err == nil {
                data.CreatedByProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByProbeId = types.StringValue(fmt.Sprintf("%v", obj["value"]))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByProbeId"].(string); ok && val != "" {
        data.CreatedByProbeId = types.StringValue(val)
    } else {
        data.CreatedByProbeId = types.StringNull()
    }
    if val, ok := dataMap["isCreatedAutomatically"].(bool); ok {
        data.IsCreatedAutomatically = types.BoolValue(val)
    }
    if val, ok := dataMap["incidentNumber"].(float64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentNumber"].(int); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentNumber"].(int64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["incidentNumber"] == nil {
        data.IncidentNumber = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncidentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data IncidentResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/incident/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete incident, got error: %s", err))
        return
    }
}


func (r *IncidentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *IncidentResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *IncidentResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *IncidentResource) parseJSONField(terraformString types.String) interface{} {
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
func (r *IncidentResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType
// Only these types should be marshalled/unmarshalled as typed wrapper objects
// This list is dynamically generated from Common/Types/JSON.ts ObjectType enum
func (r *IncidentResource) isValidOneUptimeObjectType(typeStr string) bool {
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
