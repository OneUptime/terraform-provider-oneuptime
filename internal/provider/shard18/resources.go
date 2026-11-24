package shard18

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard18", []OrderedResource{
        {Order: 52, New: NewBiToolResource},
        {Order: 76, New: NewTracePipelineResource},
        {Order: 113, New: NewIncidentUserOwnerResource},
        {Order: 116, New: NewIncidentSeverityResource},
        {Order: 118, New: NewIncidentPostmortemTemplateResource},
        {Order: 138, New: NewAlertEpisodeStateTimelineResource},
        {Order: 148, New: NewAlertPrivacyRuleResource},
        {Order: 155, New: NewIncidentEpisodeResource},
        {Order: 194, New: NewAlertReminderRuleResource},
        {Order: 200, New: NewStatusPageLabelRuleResource},
        {Order: 210, New: NewKubernetesClusterOwnerRuleResource},
        {Order: 231, New: NewDashboardLabelRuleResource},
        {Order: 268, New: NewRunnerUserOwnerResource},
        {Order: 272, New: NewStatusPageSsoResource},
        {Order: 292, New: NewOnCallScheduleLayerResource},
        {Order: 306, New: NewIncomingCallPolicyLabelRuleResource},
        {Order: 332, New: NewTableViewResource},
        {Order: 336, New: NewDashboardSsoResource},
        {Order: 353, New: NewStorageArrayFeedResource},
        {Order: 397, New: NewPodmanHostTeamOwnerResource},
        {Order: 406, New: NewOpenstackCloudUserOwnerResource},
        {Order: 474, New: NewOpenstackCloudUserOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 406, Name: "openstack_cloud_user_owner", LegacyName: "open_stack_cloud_user_owner", New: NewOpenstackCloudUserOwnerResource, NewLegacy: NewOpenstackCloudUserOwnerLegacyResource},
    })
}
