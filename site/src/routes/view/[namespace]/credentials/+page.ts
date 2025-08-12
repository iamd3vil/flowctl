import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/apiClient';
import { DEFAULT_PAGE_SIZE } from '$lib/constants';
import { permissionChecker, type ResourcePermissions } from '$lib/utils/permissions';

export const load: PageLoad = async ({ params, url, parent }) => {
	const { user, namespaceId } = await parent();
	// Check permissions
	let permissions: ResourcePermissions
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

	try {
		const { namespace } = params;
		const page = Number(url.searchParams.get('page') || '1');
		const search = url.searchParams.get('search') || '';

		// Fetch credentials data
		const credentialsResponse = await apiClient.credentials.list(namespace, {
			page,
			count_per_page: DEFAULT_PAGE_SIZE,
			filter: search || undefined
		});

		return {
			credentials: credentialsResponse.credentials || [],
			totalCount: credentialsResponse.total_count || 0,
			pageCount: credentialsResponse.page_count || 1,
			currentPage: page,
			searchQuery: search,
			namespace,
			permissions
		};
	} catch (err) {
		console.error('Failed to load credentials data:', err);
		throw error(500, 'Failed to load credentials data');
	}
};