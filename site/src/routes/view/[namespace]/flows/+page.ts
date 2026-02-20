import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/apiClient';
import { FLOWS_PER_PAGE } from '$lib/constants';
import { permissionChecker } from '$lib/utils/permissions';

export const load: PageLoad = async ({ params, url, parent }) => {
  const { user, namespaceId } = await parent();

  // Check permissions
  try {
    const permissions = await permissionChecker(user!, 'flow', namespaceId, ['view'], '_');
    if (!permissions.canRead) {
      error(403, {
        message: 'You do not have permission to view flows in this namespace',
        code: 'INSUFFICIENT_PERMISSIONS'
      });
    }
  } catch (err) {
    if (err && typeof err === 'object' && 'status' in err) {
      throw err; // Re-throw SvelteKit errors (like the 403 above)
    }
    error(500, {
      message: 'Failed to check permissions',
      code: 'PERMISSION_CHECK_FAILED'
    });
  }

  const page = Number(url.searchParams.get('page')) || 1;
  const filter = url.searchParams.get('filter') || '';
  const group = url.searchParams.get('group') || '';

  if (group) {
    // Inside a group — load group flows
    const groupFlowsPromise = apiClient.flows.groups.get(params.namespace, group);

    return {
      groupFlowsPromise,
      flowsPromise: Promise.resolve({ flows: [], page_count: 0, total_count: 0 }),
      groupsPromise: Promise.resolve({ groups: [] }),
      group,
      currentPage: 1,
      filter: '',
      namespaceId
    };
  }

  // Root view — load groups and paginated flows
  const groupsPromise = apiClient.flows.groups.me(params.namespace);
  const flowsPromise = apiClient.flows.list(params.namespace, {
    page,
    count_per_page: FLOWS_PER_PAGE,
    filter
  });

  return {
    flowsPromise,
    groupsPromise,
    groupFlowsPromise: null,
    group: '',
    currentPage: page,
    filter,
    namespaceId
  };
};
