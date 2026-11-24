package shard19

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard19", []OrderedResource{
        {Order: 15, New: NewTraceSavedViewResource},
        {Order: 75, New: NewMetricRecordingRuleResource},
        {Order: 122, New: NewFormDomainResource},
        {Order: 131, New: NewAlertInternalNoteResource},
        {Order: 137, New: NewAlertEpisodeMemberResource},
        {Order: 191, New: NewSloOwnerRuleResource},
        {Order: 211, New: NewKubernetesClusterLabelRuleResource},
        {Order: 242, New: NewSubscriberNotificationTemplateResource},
        {Order: 310, New: NewProjectUserProfileResource},
        {Order: 351, New: NewDockerSwarmClusterFeedResource},
        {Order: 361, New: NewServiceAiNoteResource},
        {Order: 418, New: NewVcenterResource},
        {Order: 441, New: NewDatabaseTeamOwnerResource},
        {Order: 479, New: NewVcenterLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 418, Name: "vcenter", LegacyName: "v_center", New: NewVcenterResource, NewLegacy: NewVcenterLegacyResource},
    })
}
