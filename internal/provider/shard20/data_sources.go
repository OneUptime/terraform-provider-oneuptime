package shard20

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard20", []OrderedDataSource{
        {Order: 57, New: NewComplianceSettingDataSource},
        {Order: 115, New: NewIncomingCallPolicyDataSource},
        {Order: 142, New: NewIncidentPostmortemTemplateDataSource},
        {Order: 175, New: NewAlertEpisodeOnCallRuleDataSource},
        {Order: 215, New: NewServiceLevelObjectiveUserOwnerDataSource},
        {Order: 220, New: NewSloLabelRuleDataSource},
        {Order: 281, New: NewScheduledMaintenanceFeedDataSource},
        {Order: 290, New: NewRunbookTeamOwnerDataSource},
        {Order: 317, New: NewWebhookLogDataSource},
        {Order: 334, New: NewOnCallScheduleLayerDataSource},
        {Order: 337, New: NewOnCallScheduleUserOwnerDataSource},
        {Order: 341, New: NewOnCallDutyPolicyUserOwnerDataSource},
        {Order: 348, New: NewIncomingCallPolicyLabelRuleDataSource},
        {Order: 362, New: NewAiAgentTaskPullRequestDataSource},
        {Order: 436, New: NewNetworkSiteTeamOwnerDataSource},
        {Order: 440, New: NewNetworkSiteStatusEventDataSource},
        {Order: 514, New: NewSpanDataSource},
    }, nil)
}
