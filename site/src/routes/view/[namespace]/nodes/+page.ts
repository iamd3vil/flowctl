import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/apiClient';
import { permissionChecker } from '$lib/utils/permissions';
import { DEFAULT_PAGE_SIZE } from '$lib/constants';


export const load: PageLoad = async ({ params, url, parent }) => {
	const { user, namespaceId } = await parent();

	// Check permissions
	try {
		const permissions = await permissionChecker(user!, 'node', namespaceId, ['view']);
		if (!permissions.canRead) {
			error(403, {
				message: 'You do not have permission to view nodes in this namespace',
				code: 'INSUFFICIENT_PERMISSIONS'
			});
		}
	} catch (err) {
		if (err && typeof err === 'object' && 'status' in err) {
			throw err; // Re-throw SvelteKit errors
		}
		error(500, {
			message: 'Failed to check permissions',
			code: 'PERMISSION_CHECK_FAILED'
		});
	}

	const { namespace } = params;
	const page = Number(url.searchParams.get('page') || '1');
	const search = url.searchParams.get('search') || '';

	const nodesPromise = apiClient.nodes.list(namespace, {
		page,
		count_per_page: DEFAULT_PAGE_SIZE,
		filter: search || undefined
	});
	const statsPromise = apiClient.nodes.getStats(namespace);
	const credentialsPromise = apiClient.credentials.list(namespace);

	return {
		nodesPromise,
		statsPromise,
		credentialsPromise,
		currentPage: page,
		searchQuery: search,
		namespace
	};
};
