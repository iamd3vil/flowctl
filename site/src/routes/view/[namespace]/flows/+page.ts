import { apiClient } from '$lib/apiClient';

export const load = async ({ params, url }: { params: any; url: any }) => {
  const page = Number(url.searchParams.get('page')) || 1;
  const filter = url.searchParams.get('filter') || '';
  
  try {
    const data = await apiClient.flows.list(params.namespace, {
      page,
      count_per_page: 10,
      filter
    });
    
    return {
      flows: data.flows,
      pageCount: data.page_count,
      totalCount: data.total_count,
      currentPage: page,
      filter
    };
  } catch (error) {
    console.error('Failed to load flows:', error);
    return {
      flows: [],
      pageCount: 0,
      totalCount: 0,
      currentPage: 1,
      filter,
      error: 'Failed to load flows'
    };
  }
};