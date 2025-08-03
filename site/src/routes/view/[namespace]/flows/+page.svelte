<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { apiClient } from '$lib/apiClient';

  let { data } = $props();
  let searchValue = $state('');
  let debounceTimer: number;
  let flows = $state(data.flows);
  let pageCount = $state(data.pageCount);
  let totalCount = $state(data.totalCount);
  let currentPage = $state(data.currentPage);
  let error = $state(data.error);
  let loading = $state(false);

  const goToFlow = (flowSlug: string) => {
    goto(`/view/${page.params.namespace}/flows/${flowSlug}`);
  };

  const loadFlows = async (filter: string = '', pageNumber: number = 1) => {
    loading = true;
    error = '';
    
    try {
      const result = await apiClient.flows.list(page.params.namespace!, {
        filter,
        page: pageNumber,
        count_per_page: 10
      });
      
      flows = result.flows;
      pageCount = result.page_count;
      totalCount = result.total_count;
      currentPage = pageNumber;
    } catch (err) {
      error = 'Failed to load flows';
      console.error('Failed to load flows:', err);
    } finally {
      loading = false;
    }
  };

  const handleSearch = (event: Event) => {
    let target = event.target as HTMLInputElement;
    searchValue = target.value;
    // Clear existing timer
    clearTimeout(debounceTimer);
    
    // Set new timer for 300ms debounce
    debounceTimer = setTimeout(() => {
      loadFlows(searchValue.trim(), 1);
    }, 300);
  };

  const goToPage = (pageNum: number) => {
    loadFlows('', pageNum);
  };
</script>

<svelte:head>
  <title>Flows - {page.params.namespace} - Flowctl</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@tabler/icons-webfont@latest/tabler-icons.min.css">
</svelte:head>

<div class="max-w-7xl mx-auto">
  <!-- Header -->
  <div class="flex items-center justify-between mb-6">
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Flows</h1>
      <p class="text-gray-600">Manage and run your workflows</p>
    </div>
    
    <!-- Search -->
    <div class="max-w-md">
      <div class="relative">
        <input
          type="text"
          placeholder="Search flows..."
          value={searchValue}
          oninput={handleSearch}
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        {#if loading}
          <div class="absolute right-3 top-1/2 transform -translate-y-1/2">
            <svg class="animate-spin h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
        {/if}
      </div>
    </div>
  </div>

  <!-- Error Message -->
  {#if error}
    <div class="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex">
        <svg class="w-5 h-5 text-red-400 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error</h3>
          <p class="mt-1 text-sm text-red-700">{error}</p>
        </div>
      </div>
    </div>
  {/if}

  <!-- Flows Table -->
  {#if flows.length > 0}
    <div class="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Flow Name
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Description
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Steps
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          {#each flows as flow (flow.id)}
            <tr 
              class="hover:bg-gray-50 cursor-pointer transition-colors"
              onclick={() => goToFlow(flow.slug)}
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-8 w-8 bg-blue-100 rounded-lg flex items-center justify-center">
                    <svg class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                    </svg>
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">{flow.name}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <div class="text-sm text-gray-600 max-w-xs truncate">{flow.description}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center text-sm text-gray-500">
                  <span>{flow.step_count || 0}</span>
                  <span class="ml-1">steps</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button 
                  class="text-blue-600 hover:text-blue-700 transition-colors"
                  onclick={(e) => { e.stopPropagation(); goToFlow(flow.slug); }}
                >
                  Run Flow
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    {#if pageCount > 1}
      <div class="mt-6 flex items-center justify-between">
        <div class="text-sm text-gray-700">
          Showing {flows.length} of {totalCount} flows
        </div>
        <div class="flex space-x-2">
          {#if currentPage > 1}
            <button 
              onclick={() => goToPage(currentPage - 1)}
              class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
            >
              Previous
            </button>
          {/if}
          
          {#each Array.from({length: Math.min(5, pageCount)}, (_, i) => i + Math.max(1, currentPage - 2)) as pageNum}
            {#if pageNum <= pageCount}
              <button 
                onclick={() => goToPage(pageNum)}
                class="px-3 py-2 text-sm font-medium rounded-md"
                class:bg-blue-600={pageNum === currentPage}
                class:text-white={pageNum === currentPage}
                class:text-gray-500={pageNum !== currentPage}
                class:bg-white={pageNum !== currentPage}
                class:border={pageNum !== currentPage}
                class:border-gray-300={pageNum !== currentPage}
                class:hover:bg-gray-50={pageNum !== currentPage}
              >
                {pageNum}
              </button>
            {/if}
          {/each}
          
          {#if currentPage < pageCount}
            <button 
              onclick={() => goToPage(currentPage + 1)}
              class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
            >
              Next
            </button>
          {/if}
        </div>
      </div>
    {/if}
  {:else}
    <!-- Empty State -->
    <div class="flex items-center justify-center h-96 bg-white rounded-lg border border-gray-200">
      <div class="text-center">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No flows found</h3>
        <p class="mt-1 text-sm text-gray-500">
          {searchValue ? 'Try adjusting your search' : 'No flows are available in this namespace'}
        </p>
      </div>
    </div>
  {/if}
</div>