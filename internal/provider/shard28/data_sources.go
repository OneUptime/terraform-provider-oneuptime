package shard28

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard28", []OrderedDataSource{
        {Order: 38, New: NewSecurityCaseCustomFieldDataSource},
        {Order: 59, New: NewDataSecuritySettingDataSource},
        {Order: 61, New: NewDataAssetFeedDataSource},
        {Order: 67, New: NewBiToolDataSource},
        {Order: 109, New: NewEscalationRuleDataSource},
        {Order: 148, New: NewFormPrivateUserDataSource},
        {Order: 183, New: NewIncidentEpisodeRoleMemberDataSource},
        {Order: 195, New: NewIncidentPrivacyRuleDataSource},
        {Order: 229, New: NewStatusPageMonitorRuleDataSource},
        {Order: 238, New: NewKubernetesClusterOwnerRuleDataSource},
        {Order: 248, New: NewCephClusterOwnerRuleDataSource},
        {Order: 255, New: NewRunbookLabelRuleDataSource},
        {Order: 288, New: NewRunbookDataSource},
        {Order: 331, New: NewMonitorGroupUserOwnerDataSource},
        {Order: 358, New: NewProbeUserOwnerDataSource},
        {Order: 360, New: NewAiAgentTeamOwnerDataSource},
        {Order: 364, New: NewAiConversationMessageDataSource},
        {Order: 370, New: NewDataSourceDataSource},
        {Order: 392, New: NewWorkspaceNotificationRuleDataSource},
        {Order: 401, New: NewDockerHostFeedDataSource},
        {Order: 411, New: NewServiceFeedDataSource},
        {Order: 435, New: NewNetworkSiteDataSource},
        {Order: 451, New: NewPodmanHostTeamOwnerDataSource},
        {Order: 540, New: NewAiAgentTeamOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 360, Name: "ai_agent_team_owner", LegacyName: "ai_agent_owner_team", New: NewAiAgentTeamOwnerDataSource, NewLegacy: NewAiAgentTeamOwnerLegacyDataSource},
    })
}
