import type { PageLoad } from './$types';
import { apiClient } from '$lib/apiClient';

export const load: PageLoad = async ({ params }) => {
  const { namespace, flowId, logId } = params;

  try {
    // Load flow metadata and execution summary in parallel
    const [flowMeta, executionSummary] = await Promise.all([
      apiClient.flows.getMeta(namespace, flowId),
      apiClient.executions.getById(namespace, logId)
    ]);

    return {
      namespace,
      flowId,
      logId,
      flowMeta,
      executionSummary
    };
  } catch (error) {
    console.error('Failed to load flow status data:', error);
    return {
      namespace,
      flowId,
      logId,
      flowMeta: {
        meta: {
          id: '',
          name: '',
          description: '',
          namespace: ''
        },
        actions: []
      },
      error: 'Failed to load flow status data'
    };
  }
};