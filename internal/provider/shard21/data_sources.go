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
        {Order: 74, New: NewWarehouseJobDataSource},
        {Order: 157, New: NewFormDomainDataSource},
        {Order: 191, New: NewIncidentEpisodeMemberDataSource},
        {Order: 221, New: NewIncidentSlaDataSource},
        {Order: 246, New: NewPodmanHostOwnerRuleDataSource},
        {Order: 280, New: NewSubscriberNotificationTemplateDataSource},
        {Order: 282, New: NewStatusPageFooterLinkDataSource},
        {Order: 296, New: NewWorkflowTeamOwnerDataSource},
        {Order: 401, New: NewDashboardPlaylistDataSource},
        {Order: 409, New: NewMonitorFeedDataSource},
        {Order: 416, New: NewProxmoxClusterFeedDataSource},
        {Order: 422, New: NewServiceAiNoteDataSource},
        {Order: 512, New: NewQueueDataSource},
        {Order: 526, New: NewSloHistoryDataSource},
        {Order: 536, New: NewChangeEventDataSource},
    }, nil)
}
