import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/apiClient';
import { permissionChecker } from '$lib/utils/permissions';

export const ssr = false;

export const load: PageLoad = async ({ params, parent }) => {
	const { user, namespaceId } = await parent();

	// Check permissions
	let permissions;
	try {
		permissions = await permissionChecker(user!, 'member', namespaceId, ['view', 'create', 'update', 'delete']);
		if (!permissions.canRead) {
			error(403, {
				message: 'You do not have permission to view members in this namespace',
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

	const membersPromise = apiClient.namespaces.members.list(namespace);

	return {
		membersPromise,
		namespace,
		permissions,
		user,
		namespaceId
	};
};
