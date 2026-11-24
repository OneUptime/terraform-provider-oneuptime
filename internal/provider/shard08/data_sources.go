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
        {Order: 101, New: NewTraceRecordingRuleDataSource},
        {Order: 110, New: NewEscalationRuleTeamDataSource},
        {Order: 118, New: NewIncomingCallLogDataSource},
        {Order: 137, New: NewIncidentUserOwnerDataSource},
        {Order: 165, New: NewAlertEpisodeTeamOwnerDataSource},
        {Order: 179, New: NewAlertEpisodeLabelRuleDataSource},
        {Order: 217, New: NewSloMonitorRuleDataSource},
        {Order: 262, New: NewDomainDataSource},
        {Order: 280, New: NewScheduledMaintenanceCustomFieldDataSource},
        {Order: 314, New: NewSmsLogDataSource},
        {Order: 326, New: NewScheduledMaintenanceTemplateUserOwnerDataSource},
        {Order: 427, New: NewOidCollectionTemplateDataSource},
        {Order: 460, New: NewOpenstackCloudUserOwnerDataSource},
        {Order: 461, New: NewDockerSwarmClusterUserOwnerDataSource},
        {Order: 501, New: NewDatabaseOwnerRuleDataSource},
        {Order: 528, New: NewEscalationRuleTeamLegacyDataSource},
        {Order: 546, New: NewOpenstackCloudUserOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 110, Name: "escalation_rule_team", LegacyName: "team_on_call_duty_escalation_rule", New: NewEscalationRuleTeamDataSource, NewLegacy: NewEscalationRuleTeamLegacyDataSource},
        {Order: 460, Name: "openstack_cloud_user_owner", LegacyName: "open_stack_cloud_user_owner", New: NewOpenstackCloudUserOwnerDataSource, NewLegacy: NewOpenstackCloudUserOwnerLegacyDataSource},
    })
}
