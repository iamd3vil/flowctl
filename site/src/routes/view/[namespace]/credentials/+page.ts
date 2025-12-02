import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/apiClient';
import { DEFAULT_PAGE_SIZE } from '$lib/constants';
import { permissionChecker, type ResourcePermissions } from '$lib/utils/permissions';

export const load: PageLoad = async ({ params, url, parent }) => {
	const { user, namespaceId } = await parent();
	// Check permissions
	let permissions: ResourcePermissions;
	try {
		permissions = await permissionChecker(user!, 'credential', namespaceId, ['view', 'create', 'update', 'delete']);
		if (!permissions.canRead) {
			error(403, {
				message: 'You do not have permission to view credentials in this namespace',
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

	const { namespace } = params;
	const page = Number(url.searchParams.get('page') || '1');
	const search = url.searchParams.get('search') || '';

	// Return promise without awaiting - page renders immediately while data loads
	const credentialsPromise = apiClient.credentials.list(namespace, {
		page,
		count_per_page: DEFAULT_PAGE_SIZE,
		filter: search || undefined
	});

	return {
		credentialsPromise,
		currentPage: page,
		searchQuery: search,
		namespace,
		permissions
	};
};