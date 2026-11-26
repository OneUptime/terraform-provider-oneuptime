package shard08

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard08", []OrderedDataSource{
        {Order: 8, New: NewCustomRoleDataSource},
        {Order: 12, New: NewApiKeyDataSource},
        {Order: 29, New: NewThreatIntelFeedDataSource},
        {Order: 90, New: NewCloudWafDataSource},
        {Order: 111, New: NewTraceRecordingRuleDataSource},
        {Order: 120, New: NewEscalationRuleTeamDataSource},
        {Order: 128, New: NewIncomingCallLogDataSource},
        {Order: 147, New: NewIncidentUserOwnerDataSource},
        {Order: 175, New: NewAlertEpisodeTeamOwnerDataSource},
        {Order: 189, New: NewAlertEpisodeLabelRuleDataSource},
        {Order: 227, New: NewSloMonitorRuleDataSource},
        {Order: 272, New: NewDomainDataSource},
        {Order: 290, New: NewScheduledMaintenanceCustomFieldDataSource},
        {Order: 324, New: NewSmsLogDataSource},
        {Order: 336, New: NewScheduledMaintenanceTemplateUserOwnerDataSource},
        {Order: 437, New: NewOidCollectionTemplateDataSource},
        {Order: 470, New: NewOpenstackCloudUserOwnerDataSource},
        {Order: 471, New: NewDockerSwarmClusterUserOwnerDataSource},
        {Order: 511, New: NewDatabaseOwnerRuleDataSource},
        {Order: 538, New: NewEscalationRuleTeamLegacyDataSource},
        {Order: 556, New: NewOpenstackCloudUserOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 120, Name: "escalation_rule_team", LegacyName: "team_on_call_duty_escalation_rule", New: NewEscalationRuleTeamDataSource, NewLegacy: NewEscalationRuleTeamLegacyDataSource},
        {Order: 470, Name: "openstack_cloud_user_owner", LegacyName: "open_stack_cloud_user_owner", New: NewOpenstackCloudUserOwnerDataSource, NewLegacy: NewOpenstackCloudUserOwnerLegacyDataSource},
    })
}
