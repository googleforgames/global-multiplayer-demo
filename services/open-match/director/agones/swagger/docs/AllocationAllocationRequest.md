# AllocationAllocationRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Namespace** | **string** |  | [optional] [default to null]
**MultiClusterSetting** | [***AllocationMultiClusterSetting**](allocationMultiClusterSetting.md) | If specified, multi-cluster policies are applied. Otherwise, allocation will happen locally. | [optional] [default to null]
**RequiredGameServerSelector** | [***AllocationGameServerSelector**](allocationGameServerSelector.md) | Deprecated: Please use gameServerSelectors instead. This field is ignored if the gameServerSelectors field is set The required allocation. Defaults to all GameServers. | [optional] [default to null]
**PreferredGameServerSelectors** | [**[]AllocationGameServerSelector**](allocationGameServerSelector.md) | Deprecated: Please use gameServerSelectors instead. This field is ignored if the gameServerSelectors field is set The ordered list of preferred allocations out of the &#x60;required&#x60; set. If the first selector is not matched, the selection attempts the second selector, and so on. | [optional] [default to null]
**Scheduling** | [***AllocationRequestSchedulingStrategy**](AllocationRequestSchedulingStrategy.md) | Scheduling strategy. Defaults to \&quot;Packed\&quot;. | [optional] [default to null]
**MetaPatch** | [***AllocationMetaPatch**](allocationMetaPatch.md) |  | [optional] [default to null]
**Metadata** | [***AllocationMetaPatch**](allocationMetaPatch.md) |  | [optional] [default to null]
**GameServerSelectors** | [**[]AllocationGameServerSelector**](allocationGameServerSelector.md) | Ordered list of GameServer label selectors. If the first selector is not matched, the selection attempts the second selector, and so on. This is useful for things like smoke testing of new game servers. Note: This field can only be set if neither Required or Preferred is set. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


