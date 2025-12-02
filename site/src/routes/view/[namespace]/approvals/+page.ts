import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/apiClient';
import { DEFAULT_PAGE_SIZE } from '$lib/constants';
import { permissionChecker } from '$lib/utils/permissions';

export const ssr = false;

export const load: PageLoad = async ({ params, url, parent }) => {
	const { user, namespaceId } = await parent();

	try {
		const permissions = await permissionChecker(user!, 'approval', namespaceId, ['view']);
		if (!permissions.canRead) {
			error(403, {
				message: 'You do not have permission to view approvals in this namespace',
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

	const { namespace } = params;
	const page = Number(url.searchParams.get('page') || '1');
	const search = url.searchParams.get('search') || '';
	const status = url.searchParams.get('status') || '';

	// Return promise without awaiting - page renders immediately while data loads
	type ApprovalStatus = "pending" | "approved" | "rejected" | "";
	const approvalsPromise = apiClient.approvals.list(namespace, {
		page,
		count_per_page: DEFAULT_PAGE_SIZE,
		filter: search || undefined,
		status: (status as ApprovalStatus) || undefined
	});

	return {
		approvalsPromise,
		currentPage: page,
		searchQuery: search,
		statusFilter: status,
		namespace
	};
};