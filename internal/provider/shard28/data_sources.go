package shard28

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard28", []OrderedDataSource{
        {Order: 41, New: NewSecurityCaseCustomFieldDataSource},
        {Order: 64, New: NewDataSecuritySettingDataSource},
        {Order: 66, New: NewDataAssetFeedDataSource},
        {Order: 72, New: NewBiToolDataSource},
        {Order: 119, New: NewEscalationRuleDataSource},
        {Order: 158, New: NewFormPrivateUserDataSource},
        {Order: 193, New: NewIncidentEpisodeRoleMemberDataSource},
        {Order: 205, New: NewIncidentPrivacyRuleDataSource},
        {Order: 239, New: NewStatusPageMonitorRuleDataSource},
        {Order: 248, New: NewKubernetesClusterOwnerRuleDataSource},
        {Order: 258, New: NewCephClusterOwnerRuleDataSource},
        {Order: 265, New: NewRunbookLabelRuleDataSource},
        {Order: 298, New: NewRunbookDataSource},
        {Order: 341, New: NewMonitorGroupUserOwnerDataSource},
        {Order: 368, New: NewProbeUserOwnerDataSource},
        {Order: 370, New: NewAiAgentTeamOwnerDataSource},
        {Order: 374, New: NewAiConversationMessageDataSource},
        {Order: 380, New: NewDataSourceDataSource},
        {Order: 402, New: NewWorkspaceNotificationRuleDataSource},
        {Order: 411, New: NewDockerHostFeedDataSource},
        {Order: 421, New: NewServiceFeedDataSource},
        {Order: 445, New: NewNetworkSiteDataSource},
        {Order: 461, New: NewPodmanHostTeamOwnerDataSource},
        {Order: 550, New: NewAiAgentTeamOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 370, Name: "ai_agent_team_owner", LegacyName: "ai_agent_owner_team", New: NewAiAgentTeamOwnerDataSource, NewLegacy: NewAiAgentTeamOwnerLegacyDataSource},
    })
}
