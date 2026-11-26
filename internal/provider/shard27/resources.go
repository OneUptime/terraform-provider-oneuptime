package shard27

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard27", []OrderedResource{
        {Order: 22, New: NewWatchlistResource},
        {Order: 81, New: NewTracePipelineProcessorResource},
        {Order: 87, New: NewApiKeyCustomRoleResource},
        {Order: 101, New: NewMonitorSecretResource},
        {Order: 118, New: NewIncidentRoleResource},
        {Order: 133, New: NewAlertCustomFieldResource},
        {Order: 155, New: NewAlertEpisodeOnCallRuleOwningTeamResource},
        {Order: 161, New: NewIncidentAlertResource},
        {Order: 230, New: NewRunbookOwnerRuleResource},
        {Order: 248, New: NewStatusPageFooterLinkResource},
        {Order: 266, New: NewRunbookRuleResource},
        {Order: 279, New: NewMonitorProbeResource},
        {Order: 359, New: NewProxmoxClusterFeedResource},
        {Order: 369, New: NewKubernetesClusterTeamOwnerResource},
        {Order: 424, New: NewVcenterUserOwnerResource},
        {Order: 485, New: NewVcenterUserOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 424, Name: "vcenter_user_owner", LegacyName: "v_center_user_owner", New: NewVcenterUserOwnerResource, NewLegacy: NewVcenterUserOwnerLegacyResource},
    })
}
