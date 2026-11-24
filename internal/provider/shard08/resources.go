package shard08

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard08", []OrderedResource{
        {Order: 0, New: NewProbeResource},
        {Order: 96, New: NewMonitorResource},
        {Order: 112, New: NewIncidentTeamOwnerResource},
        {Order: 153, New: NewAlertEpisodePrivacyRuleResource},
        {Order: 171, New: NewIncidentLabelRuleResource},
        {Order: 172, New: NewScheduledMaintenanceOwnerRuleResource},
        {Order: 206, New: NewDockerHostOwnerRuleResource},
        {Order: 216, New: NewOpenstackCloudLabelRuleResource},
        {Order: 232, New: NewMonitorStatusEventResource},
        {Order: 239, New: NewStatusPageAnnouncementResource},
        {Order: 255, New: NewWorkflowVariableResource},
        {Order: 267, New: NewRunnerTeamOwnerResource},
        {Order: 324, New: NewTelemetryIngestionKeyResource},
        {Order: 403, New: NewOpenstackCloudTeamOwnerResource},
        {Order: 440, New: NewDatabaseFeedResource},
        {Order: 449, New: NewQueueOwnerRuleResource},
        {Order: 460, New: NewOpenstackCloudLabelRuleLegacyResource},
        {Order: 473, New: NewOpenstackCloudTeamOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 216, Name: "openstack_cloud_label_rule", LegacyName: "open_stack_cloud_label_rule", New: NewOpenstackCloudLabelRuleResource, NewLegacy: NewOpenstackCloudLabelRuleLegacyResource},
        {Order: 403, Name: "openstack_cloud_team_owner", LegacyName: "open_stack_cloud_team_owner", New: NewOpenstackCloudTeamOwnerResource, NewLegacy: NewOpenstackCloudTeamOwnerLegacyResource},
    })
}
