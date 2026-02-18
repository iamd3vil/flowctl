import { apiClient } from '$lib/apiClient.js';
import { permissionChecker } from '$lib/utils/permissions';
import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ parent }) => {
  const { user, namespaceId } = await parent();
  
  // Check create permissions
  try {
    const permissions = await permissionChecker(user!, 'flow', namespaceId, ['create']);
    if (!permissions.canCreate) {
      error(403, {
        message: 'You do not have permission to create flows in this namespace',
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
  
  try {
    const [executorData, messengerSchemas] = await Promise.all([
      apiClient.executors.list(),
      apiClient.messengers.list(),
    ]);

    const availableExecutors = executorData.executors.map(name => ({
      name: name,
      display_name: name.charAt(0).toUpperCase() + name.slice(1)
    }));

    // Resolve $defs/$ref in each messenger schema
    const messengerConfigs: Record<string, any> = {};
    for (const [name, schema] of Object.entries(messengerSchemas)) {
      if (schema.$defs && schema.$ref) {
        const refPath = schema.$ref.replace('#/$defs/', '');
        messengerConfigs[name] = schema.$defs[refPath] || schema;
      } else {
        messengerConfigs[name] = schema;
      }
    }

    return {
      availableExecutors,
      availableMessengers: Object.keys(messengerSchemas),
      messengerConfigs,
    };
  } catch (loadError) {
    console.error('Error loading executors/messengers:', loadError);
    return {
      availableExecutors: [],
      availableMessengers: [],
      messengerConfigs: {},
    };
  }
};