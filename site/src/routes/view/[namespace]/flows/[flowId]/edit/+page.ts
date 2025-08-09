import { apiClient } from '$lib/apiClient.js';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
  try {
    const data = await apiClient.executors.list();
    
    const availableExecutors = data.executors.map(name => ({
      name: name,
      display_name: name.charAt(0).toUpperCase() + name.slice(1)
    }));

    return {
      availableExecutors
    };
  } catch (error) {
    console.error('Error loading executors:', error);
    return {
      availableExecutors: []
    };
  }
};