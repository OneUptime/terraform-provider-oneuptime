package shard13

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard13", []OrderedDataSource{
        {Order: 26, New: NewHuntRunDataSource},
        {Order: 46, New: NewDataStoreFeedDataSource},
        {Order: 52, New: NewDataSecurityIssueDataSource},
        {Order: 129, New: NewIncidentCustomFieldDataSource},
        {Order: 153, New: NewAlertFeedDataSource},
        {Order: 156, New: NewAlertInternalNoteDataSource},
        {Order: 162, New: NewAlertEpisodeMemberDataSource},
        {Order: 224, New: NewPublicNoteApprovalStepDataSource},
        {Order: 254, New: NewRunbookOwnerRuleDataSource},
        {Order: 311, New: NewScheduledMaintenanceUserOwnerDataSource},
        {Order: 318, New: NewPushNotificationLogDataSource},
        {Order: 352, New: NewProjectUserProfileDataSource},
        {Order: 373, New: NewWebAuthnCredentialDataSource},
        {Order: 377, New: NewExceptionActivityDataSource},
        {Order: 488, New: NewCloudCostRecommendationDataSource},
        {Order: 498, New: NewDatabaseTeamOwnerDataSource},
    }, nil)
}
