package shard20

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard20", []OrderedDataSource{
        {Order: 62, New: NewComplianceSettingDataSource},
        {Order: 125, New: NewIncomingCallPolicyDataSource},
        {Order: 152, New: NewIncidentPostmortemTemplateDataSource},
        {Order: 185, New: NewAlertEpisodeOnCallRuleDataSource},
        {Order: 225, New: NewServiceLevelObjectiveUserOwnerDataSource},
        {Order: 230, New: NewSloLabelRuleDataSource},
        {Order: 291, New: NewScheduledMaintenanceFeedDataSource},
        {Order: 300, New: NewRunbookTeamOwnerDataSource},
        {Order: 327, New: NewWebhookLogDataSource},
        {Order: 344, New: NewOnCallScheduleLayerDataSource},
        {Order: 347, New: NewOnCallScheduleUserOwnerDataSource},
        {Order: 351, New: NewOnCallDutyPolicyUserOwnerDataSource},
        {Order: 358, New: NewIncomingCallPolicyLabelRuleDataSource},
        {Order: 372, New: NewAiAgentTaskPullRequestDataSource},
        {Order: 446, New: NewNetworkSiteTeamOwnerDataSource},
        {Order: 450, New: NewNetworkSiteStatusEventDataSource},
        {Order: 524, New: NewSpanDataSource},
    }, nil)
}
