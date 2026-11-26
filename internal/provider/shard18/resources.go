package shard18

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard18", []OrderedResource{
        {Order: 53, New: NewBiToolResource},
        {Order: 66, New: NewCloudWafResource},
        {Order: 80, New: NewTracePipelineResource},
        {Order: 117, New: NewIncidentUserOwnerResource},
        {Order: 120, New: NewIncidentSeverityResource},
        {Order: 122, New: NewIncidentPostmortemTemplateResource},
        {Order: 142, New: NewAlertEpisodeStateTimelineResource},
        {Order: 152, New: NewAlertPrivacyRuleResource},
        {Order: 159, New: NewIncidentEpisodeResource},
        {Order: 198, New: NewAlertReminderRuleResource},
        {Order: 204, New: NewStatusPageLabelRuleResource},
        {Order: 214, New: NewKubernetesClusterOwnerRuleResource},
        {Order: 235, New: NewDashboardLabelRuleResource},
        {Order: 272, New: NewRunnerUserOwnerResource},
        {Order: 276, New: NewStatusPageSsoResource},
        {Order: 296, New: NewOnCallScheduleLayerResource},
        {Order: 310, New: NewIncomingCallPolicyLabelRuleResource},
        {Order: 336, New: NewTableViewResource},
        {Order: 340, New: NewDashboardSsoResource},
        {Order: 357, New: NewStorageArrayFeedResource},
        {Order: 401, New: NewPodmanHostTeamOwnerResource},
        {Order: 410, New: NewOpenstackCloudUserOwnerResource},
        {Order: 478, New: NewOpenstackCloudUserOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 410, Name: "openstack_cloud_user_owner", LegacyName: "open_stack_cloud_user_owner", New: NewOpenstackCloudUserOwnerResource, NewLegacy: NewOpenstackCloudUserOwnerLegacyResource},
    })
}
