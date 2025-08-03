import { apiClient } from '$lib/apiClient';

export const load = async ({ params }: { params: any }) => {
  try {
    const [flowInputs, flowMeta] = await Promise.all([
      apiClient.flows.getInputs(params.namespace, params.flowId),
      apiClient.flows.getMeta(params.namespace, params.flowId)
    ]);
    
    return {
      flowInputs: flowInputs.inputs,
      flowMeta,
      flowId: params.flowId
    };
  } catch (error) {
    console.error('Failed to load flow data:', error);
    return {
      flowInputs: [],
      flowMeta: { meta: { name: '', description: '' }, actions: [] },
      flowId: params.flowId,
      error: 'Failed to load flow data'
    };
  }
};