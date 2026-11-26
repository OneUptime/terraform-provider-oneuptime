package shard11

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard11", []OrderedDataSource{
        {Order: 25, New: NewHuntDataSource},
        {Order: 52, New: NewDataStoreOwnerRuleDataSource},
        {Order: 82, New: NewPipelineOwnerRuleDataSource},
        {Order: 92, New: NewBlocklistFeedFetchDataSource},
        {Order: 178, New: NewAlertGroupingRuleDataSource},
        {Order: 182, New: NewAlertOwnerRuleDataSource},
        {Order: 228, New: NewSloFeedDataSource},
        {Order: 299, New: NewRunbookExecutionDataSource},
        {Order: 303, New: NewAiRuleDataSource},
        {Order: 313, New: NewRunbookVariableDataSource},
        {Order: 316, New: NewStatusPageScimDataSource},
        {Order: 338, New: NewRecommendationDismissalDataSource},
        {Order: 382, New: NewTwoFactorBackupCodeDataSource},
        {Order: 448, New: NewNetworkSiteTypeDataSource},
        {Order: 473, New: NewIotFleetTeamOwnerDataSource},
        {Order: 483, New: NewVcenterTeamOwnerDataSource},
        {Order: 505, New: NewDatabaseDataSource},
        {Order: 527, New: NewExceptionInstanceDataSource},
        {Order: 546, New: NewAiRuleLegacyDataSource},
        {Order: 558, New: NewIotFleetTeamOwnerLegacyDataSource},
        {Order: 562, New: NewVcenterTeamOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 303, Name: "ai_rule", LegacyName: "auto_remediation_rule", New: NewAiRuleDataSource, NewLegacy: NewAiRuleLegacyDataSource},
        {Order: 473, Name: "iot_fleet_team_owner", LegacyName: "io_t_fleet_team_owner", New: NewIotFleetTeamOwnerDataSource, NewLegacy: NewIotFleetTeamOwnerLegacyDataSource},
        {Order: 483, Name: "vcenter_team_owner", LegacyName: "v_center_team_owner", New: NewVcenterTeamOwnerDataSource, NewLegacy: NewVcenterTeamOwnerLegacyDataSource},
    })
}
