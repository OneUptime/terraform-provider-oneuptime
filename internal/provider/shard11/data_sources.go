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
        {Order: 49, New: NewDataStoreOwnerRuleDataSource},
        {Order: 77, New: NewPipelineOwnerRuleDataSource},
        {Order: 168, New: NewAlertGroupingRuleDataSource},
        {Order: 172, New: NewAlertOwnerRuleDataSource},
        {Order: 218, New: NewSloFeedDataSource},
        {Order: 289, New: NewRunbookExecutionDataSource},
        {Order: 293, New: NewAiRuleDataSource},
        {Order: 303, New: NewRunbookVariableDataSource},
        {Order: 306, New: NewStatusPageScimDataSource},
        {Order: 328, New: NewRecommendationDismissalDataSource},
        {Order: 372, New: NewTwoFactorBackupCodeDataSource},
        {Order: 438, New: NewNetworkSiteTypeDataSource},
        {Order: 463, New: NewIotFleetTeamOwnerDataSource},
        {Order: 473, New: NewVcenterTeamOwnerDataSource},
        {Order: 495, New: NewDatabaseDataSource},
        {Order: 517, New: NewExceptionInstanceDataSource},
        {Order: 536, New: NewAiRuleLegacyDataSource},
        {Order: 548, New: NewIotFleetTeamOwnerLegacyDataSource},
        {Order: 552, New: NewVcenterTeamOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 293, Name: "ai_rule", LegacyName: "auto_remediation_rule", New: NewAiRuleDataSource, NewLegacy: NewAiRuleLegacyDataSource},
        {Order: 463, Name: "iot_fleet_team_owner", LegacyName: "io_t_fleet_team_owner", New: NewIotFleetTeamOwnerDataSource, NewLegacy: NewIotFleetTeamOwnerLegacyDataSource},
        {Order: 473, Name: "vcenter_team_owner", LegacyName: "v_center_team_owner", New: NewVcenterTeamOwnerDataSource, NewLegacy: NewVcenterTeamOwnerLegacyDataSource},
    })
}
