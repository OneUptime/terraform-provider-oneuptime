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
        {Order: 79, New: NewMetricRecordingRuleResource},
        {Order: 126, New: NewFormDomainResource},
        {Order: 135, New: NewAlertInternalNoteResource},
        {Order: 141, New: NewAlertEpisodeMemberResource},
        {Order: 195, New: NewSloOwnerRuleResource},
        {Order: 215, New: NewKubernetesClusterLabelRuleResource},
        {Order: 246, New: NewSubscriberNotificationTemplateResource},
        {Order: 314, New: NewProjectUserProfileResource},
        {Order: 355, New: NewDockerSwarmClusterFeedResource},
        {Order: 365, New: NewServiceAiNoteResource},
        {Order: 422, New: NewVcenterResource},
        {Order: 445, New: NewDatabaseTeamOwnerResource},
        {Order: 483, New: NewVcenterLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 422, Name: "vcenter", LegacyName: "v_center", New: NewVcenterResource, NewLegacy: NewVcenterLegacyResource},
    })
}
