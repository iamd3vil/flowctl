import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/apiClient';
import { permissionChecker } from '$lib/utils/permissions';

export const ssr = false;

export const load: PageLoad = async ({ params, parent }) => {
	const { user, namespaceId } = await parent();

	// Check permissions for namespace secrets
	let permissions;
	try {
		permissions = await permissionChecker(user!, 'namespace_secret', namespaceId, ['view', 'create', 'update', 'delete']);
		if (!permissions.canRead) {
			error(403, {
				message: 'You do not have permission to view secrets in this namespace',
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

	const secretsPromise = apiClient.namespaceSecrets.list(namespace);

	return {
		secretsPromise,
		namespace,
		permissions,
		user,
		namespaceId
	};
};
