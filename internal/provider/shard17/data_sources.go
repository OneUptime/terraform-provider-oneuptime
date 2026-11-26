package shard17

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard17", []OrderedDataSource{
        {Order: 17, New: NewExceptionSavedViewDataSource},
        {Order: 65, New: NewDataAssetDataSource},
        {Order: 69, New: NewDataPipelineDataSource},
        {Order: 108, New: NewTracePipelineProcessorDataSource},
        {Order: 126, New: NewIncomingCallPolicyEscalationRuleDataSource},
        {Order: 131, New: NewMonitorSecretDataSource},
        {Order: 132, New: NewIncomingRequestSecretDataSource},
        {Order: 187, New: NewAlertEpisodeOwnerRuleDataSource},
        {Order: 333, New: NewUserOnCallLogTimelineDataSource},
        {Order: 342, New: NewMonitorGroupResourceDataSource},
        {Order: 363, New: NewServiceDataSource},
        {Order: 419, New: NewHostFeedDataSource},
        {Order: 428, New: NewKubernetesClusterTeamOwnerDataSource},
        {Order: 430, New: NewDockerHostDataSource},
        {Order: 474, New: NewIotFleetUserOwnerDataSource},
        {Order: 477, New: NewCephClusterTeamOwnerDataSource},
        {Order: 484, New: NewVcenterUserOwnerDataSource},
        {Order: 506, New: NewDatabaseEndpointDataSource},
        {Order: 530, New: NewKubernetesCostAllocationDataSource},
        {Order: 559, New: NewIotFleetUserOwnerLegacyDataSource},
        {Order: 563, New: NewVcenterUserOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 474, Name: "iot_fleet_user_owner", LegacyName: "io_t_fleet_user_owner", New: NewIotFleetUserOwnerDataSource, NewLegacy: NewIotFleetUserOwnerLegacyDataSource},
        {Order: 484, Name: "vcenter_user_owner", LegacyName: "v_center_user_owner", New: NewVcenterUserOwnerDataSource, NewLegacy: NewVcenterUserOwnerLegacyDataSource},
    })
}
