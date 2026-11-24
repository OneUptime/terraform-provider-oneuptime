package shard21

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard21", []OrderedDataSource{
        {Order: 30, New: NewSecurityCaseDataSource},
        {Order: 69, New: NewWarehouseJobDataSource},
        {Order: 147, New: NewFormDomainDataSource},
        {Order: 181, New: NewIncidentEpisodeMemberDataSource},
        {Order: 211, New: NewIncidentSlaDataSource},
        {Order: 236, New: NewPodmanHostOwnerRuleDataSource},
        {Order: 270, New: NewSubscriberNotificationTemplateDataSource},
        {Order: 272, New: NewStatusPageFooterLinkDataSource},
        {Order: 286, New: NewWorkflowTeamOwnerDataSource},
        {Order: 391, New: NewDashboardPlaylistDataSource},
        {Order: 399, New: NewMonitorFeedDataSource},
        {Order: 406, New: NewProxmoxClusterFeedDataSource},
        {Order: 412, New: NewServiceAiNoteDataSource},
        {Order: 502, New: NewQueueDataSource},
        {Order: 516, New: NewSloHistoryDataSource},
        {Order: 526, New: NewChangeEventDataSource},
    }, nil)
}
