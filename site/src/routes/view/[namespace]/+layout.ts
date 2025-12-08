import { error, redirect } from '@sveltejs/kit';
import { apiClient } from '$lib/apiClient.js';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ params, parent, url }) => {
  const { userPromise } = await parent();
  const user = await userPromise;

  // Redirect to login if not authenticated
  if (!user) {
    const redirectUrl = url.pathname + url.search;
    throw redirect(302, `/login?redirect_url=${encodeURIComponent(redirectUrl)}`);
  }

  const namespacesResponse = await apiClient.namespaces.list();
  const namespace = params.namespace;
  const namespaceObject = namespacesResponse.namespaces.find(ns => ns.name === namespace);

  if (!namespaceObject) {
    throw error(403, 'Access denied. You do not have permission to access this namespace.');
  }

  return {
    user,
    namespace: namespaceObject.name,
    namespaceId: namespaceObject.id
  };
};
