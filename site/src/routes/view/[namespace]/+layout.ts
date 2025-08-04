import { error } from '@sveltejs/kit';
import { apiClient } from '$lib/apiClient.js';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ params, parent }) => {
  const { user } = await parent();
  const namespace = params.namespace;
  
  // Check if user has access to this namespace
  try {
    const namespacesResponse = await apiClient.namespaces.list();
    const accessibleNamespaces = namespacesResponse.namespaces.map(ns => ns.name);
    
    if (!accessibleNamespaces.includes(namespace)) {
      error(403, 'Access denied. You do not have permission to access this namespace.');
    }
  } catch (err) {
    error(500, 'Could not retrieve the namespace');
  }
  
  return {
    namespace
  };
};