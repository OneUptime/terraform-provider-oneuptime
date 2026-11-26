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
        {Order: 67, New: NewBlocklistFeedResource},
        {Order: 100, New: NewMonitorResource},
        {Order: 116, New: NewIncidentTeamOwnerResource},
        {Order: 157, New: NewAlertEpisodePrivacyRuleResource},
        {Order: 175, New: NewIncidentLabelRuleResource},
        {Order: 176, New: NewScheduledMaintenanceOwnerRuleResource},
        {Order: 210, New: NewDockerHostOwnerRuleResource},
        {Order: 220, New: NewOpenstackCloudLabelRuleResource},
        {Order: 236, New: NewMonitorStatusEventResource},
        {Order: 243, New: NewStatusPageAnnouncementResource},
        {Order: 259, New: NewWorkflowVariableResource},
        {Order: 271, New: NewRunnerTeamOwnerResource},
        {Order: 328, New: NewTelemetryIngestionKeyResource},
        {Order: 407, New: NewOpenstackCloudTeamOwnerResource},
        {Order: 444, New: NewDatabaseFeedResource},
        {Order: 453, New: NewQueueOwnerRuleResource},
        {Order: 464, New: NewOpenstackCloudLabelRuleLegacyResource},
        {Order: 477, New: NewOpenstackCloudTeamOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 220, Name: "openstack_cloud_label_rule", LegacyName: "open_stack_cloud_label_rule", New: NewOpenstackCloudLabelRuleResource, NewLegacy: NewOpenstackCloudLabelRuleLegacyResource},
        {Order: 407, Name: "openstack_cloud_team_owner", LegacyName: "open_stack_cloud_team_owner", New: NewOpenstackCloudTeamOwnerResource, NewLegacy: NewOpenstackCloudTeamOwnerLegacyResource},
    })
}
