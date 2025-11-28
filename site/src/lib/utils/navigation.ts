import { apiClient } from '$lib/apiClient';

export async function getDefaultNamespace(): Promise<string> {
  const storedNamespace = localStorage.getItem('selectedNamespace');

  try {
    const response = await apiClient.namespaces.list({ count_per_page: 100 });
    const namespaces = response.namespaces || [];

    if (namespaces.length === 0) {
      return 'default';
    }

    if (storedNamespace && namespaces.some(ns => ns.name === storedNamespace)) {
      return storedNamespace;
    }

    return namespaces[0].name;
  } catch {
    return storedNamespace || 'default';
  }
}
