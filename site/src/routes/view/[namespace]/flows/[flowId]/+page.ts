import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/apiClient';
import { permissionChecker } from '$lib/utils/permissions';

export const load: PageLoad = async ({ params, parent, url }) => {
  const { user, namespaceId } = await parent();

  try {
    const permissions = await permissionChecker(user!, 'flow', namespaceId, ['view']);
    if (!permissions.canRead) {
      error(403, {
        message: 'You do not have permission to view flows in this namespace',
        code: 'INSUFFICIENT_PERMISSIONS'
      });
    }
  } catch (err) {
    if (err && typeof err === 'object' && 'status' in err) {
      throw err;
    }
    error(500, {
      message: 'Failed to check permissions',
      code: 'PERMISSION_CHECK_FAILED'
    });
  }

  const rerunFromExecId = url.searchParams.get('rerun_from');

  try {
    const [flowInputs, flowMeta, executionData, schedules] = await Promise.all([
      apiClient.flows.getInputs(params.namespace, params.flowId),
      apiClient.flows.getMeta(params.namespace, params.flowId),
      rerunFromExecId
        ? apiClient.executions.getById(params.namespace, rerunFromExecId).catch(() => null)
        : Promise.resolve(null),
      apiClient.flows.schedules.list(params.namespace, params.flowId)
    ]);

    return {
      flowInputs: flowInputs.inputs,
      flowMeta,
      namespaceId,
      rerunFromExecId,
      executionInput: executionData?.input || null,
      userSchedules: schedules.schedules || [],
    };
  } catch (err) {
    error(500, 'Failed to load flow data');
  }
};