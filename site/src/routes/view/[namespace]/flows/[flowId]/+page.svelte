<script lang="ts">
  import { page } from '$app/state';
  import FlowInputForm from '$lib/components/FlowInputForm.svelte';
  import Table from '$lib/components/Table.svelte';
  import Header from '$lib/components/Header.svelte';
  import StatusBadge from '$lib/components/StatusBadge.svelte';
  import ExecutionIdCell from '$lib/components/ExecutionIdCell.svelte';
  import type { PageData } from './$types';
  import type { TableColumn } from '$lib/types';
    import { goto } from '$app/navigation';
    import { DEFAULT_PAGE_SIZE } from '$lib/constants';

  let { data }: { data: PageData } = $props();
  
  let activeTab = $state<'run' | 'history'>('run');
  let historyLoading = $state(false);
  let historyError = $state('');
  let flowExecutions = $state<any[]>([]);
  let historyCurrentPage = $state(1);
  let historyItemsPerPage = $state(DEFAULT_PAGE_SIZE);
  let historyTotalCount = $state(0);
  let historyPageCount = $state(0);

  let namespace = $derived(page.params.namespace);
  let flowId = $derived(page.params.flowId);

  let historyPaginationPages = $derived(() => {
    const pages = [];
    const start = Math.max(1, historyCurrentPage - 2);
    const end = Math.min(historyPageCount, historyCurrentPage + 2);
    
    for (let i = start; i <= end; i++) {
      pages.push(i);
    }
    return pages;
  });

  const loadFlowHistory = async () => {
    historyLoading = true;
    historyError = '';

    try {
      const response = await fetch(`/api/v1/${namespace}/flows/${flowId}/executions?page=${historyCurrentPage}&count_per_page=${historyItemsPerPage}`);
      const result = await response.json();

      if (!response.ok) {
        historyError = result.error || 'Failed to fetch execution history';
        flowExecutions = [];
        return;
      }

      flowExecutions = result.executions || [];
      historyTotalCount = result.total_count || 0;
      historyPageCount = result.page_count || 1;
    } catch (error) {
      console.error('Error loading flow history:', error);
      historyError = 'Failed to load execution history';
      flowExecutions = [];
    } finally {
      historyLoading = false;
    }
  };

  const goToHistoryPage = (pageNum: number) => {
    if (pageNum < 1 || pageNum > historyPageCount) return;
    historyCurrentPage = pageNum;
    loadFlowHistory();
  };

  const viewExecution = (executionId: string) => {
    goto(`/view/${namespace}/results/${flowId}/${executionId}`)
  };

  const formatDateTime = (dateString: string) => {
    if (!dateString) return 'Unknown';
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  const formatDuration = (startedAt: string, completedAt?: string) => {
    if (!startedAt) return 'Unknown';
    if (!completedAt) return 'Running...';

    const start = new Date(startedAt);
    const end = new Date(completedAt);
    const durationMs = end.getTime() - start.getTime();

    if (durationMs < 1000) return '<1s';

    const seconds = Math.floor(durationMs / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);

    if (hours > 0) {
      return `${hours}h ${minutes % 60}m`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds % 60}s`;
    } else {
      return `${seconds}s`;
    }
  };

  // Watch for tab changes and load history when needed
  $effect(() => {
    if (activeTab === 'history') {
      loadFlowHistory();
    }
  });

  // Table configuration
  const tableColumns: TableColumn<any>[] = [
    {
      key: 'id',
      header: 'Execution ID',
      width: 'w-48',
      component: ExecutionIdCell
    },
    {
      key: 'status',
      header: 'Status',
      component: StatusBadge
    },
    {
      key: 'triggered_by',
      header: 'Triggered By',
      width: 'w-32',
      render: (value) => value || 'System'
    },
    {
      key: 'started_at',
      header: 'Started At',
      width: 'w-40',
      render: (value) => formatDateTime(value)
    },
    {
      key: 'duration',
      header: 'Duration',
      render: (value, row) => value || formatDuration(row.started_at, row.completed_at)
    }
  ];
</script>

<svelte:head>
  <title>Run Flow - {data.flowMeta?.meta?.name || 'Loading...'}</title>
</svelte:head>

<style>
  .gradient-bg {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }
</style>

<Header title={data.flowMeta.meta.name}></Header>

<!-- Hero Section -->
<div class="gradient-bg px-6 py-8">
  <div class="max-w-4xl mx-auto text-center">
    <div class="flex items-center justify-center mb-4">
      <div class="w-12 h-12 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center">
        <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/>
        </svg>
      </div>
    </div>
    <h1 class="text-3xl font-bold text-white mb-2">{data.flowMeta?.meta?.name || 'Loading...'}</h1>
    <p class="text-lg text-white/90 max-w-xl mx-auto">{data.flowMeta?.meta?.description || 'Loading flow description...'}</p>
  </div>
</div>

<!-- Tab Navigation -->
<div class="bg-white border-b border-gray-200 px-6">
  <nav class="max-w-4xl mx-auto flex space-x-8" aria-label="Tabs">
    <button 
      onclick={() => activeTab = 'run'} 
      class="whitespace-nowrap border-b-2 py-4 px-1 text-sm font-medium {activeTab === 'run' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
    >
      Run Flow
    </button>
    <button 
      onclick={() => activeTab = 'history'} 
      class="whitespace-nowrap border-b-2 py-4 px-1 text-sm font-medium {activeTab === 'history' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
    >
      History
    </button>
  </nav>
</div>

<!-- Tab Content -->
<div class="px-6 py-8 bg-gray-50">
  {#if activeTab === 'run'}
    <div class="max-w-2xl mx-auto">
      <FlowInputForm inputs={data.flowInputs || []} namespace={namespace!} flowId={flowId!} />

      <!-- Flow Actions Summary -->
      {#if data.flowMeta?.actions && data.flowMeta.actions.length > 0}
        <div class="bg-gray-50 rounded-lg p-4 mt-6">
          <h3 class="text-sm font-medium text-gray-900 mb-3 flex items-center">
            <svg class="w-4 h-4 text-gray-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
            </svg>
            <span>Flow Actions ({data.flowMeta.actions.length} step{data.flowMeta.actions.length !== 1 ? 's' : ''})</span>
          </h3>
          <div class="space-y-2">
            {#each data.flowMeta.actions as action, index}
              <div>
                <div class="flex items-center justify-between p-3 bg-white border border-gray-200 rounded-md">
                  <div class="flex items-center">
                    <div class="w-6 h-6 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center text-xs font-medium mr-3">
                      {index + 1}
                    </div>
                    <div class="text-sm font-medium text-gray-900">{action.name}</div>
                  </div>
                  <span class="inline-flex px-2 py-1 text-xs font-medium rounded-md bg-blue-100 text-blue-800">
                    {action.executor}
                  </span>
                </div>
                
                <!-- Arrow connecting actions -->
                {#if index < data.flowMeta.actions.length - 1}
                  <div class="flex justify-center">
                    <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"/>
                    </svg>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <!-- History Tab -->
  {#if activeTab === 'history'}
    <div class="max-w-6xl mx-auto">
      {#if historyError}
        <!-- Error Message -->
        <div class="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
          <div class="flex">
            <svg class="w-5 h-5 text-red-400 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error</h3>
              <p class="mt-1 text-sm text-red-700">{historyError}</p>
            </div>
          </div>
        </div>
      {/if}

      <Table
        columns={tableColumns}
        data={flowExecutions}
        onRowClick={viewExecution}
        loading={historyLoading}
        title="Execution History for {data.flowMeta?.meta?.name || 'Flow'}"
        subtitle="Past executions of this flow"
        emptyMessage="No execution history"
      />

      <!-- Pagination Controls -->
      {#if historyPageCount > 1}
        <div class="flex items-center justify-between mt-6">
          <!-- Results Info -->
          <div class="flex items-center text-sm text-gray-700">
            <span>Showing </span>
            <span class="font-medium">{((historyCurrentPage - 1) * historyItemsPerPage) + 1}</span>
            <span> to </span>
            <span class="font-medium">{Math.min(historyCurrentPage * historyItemsPerPage, historyTotalCount)}</span>
            <span> of </span>
            <span class="font-medium">{historyTotalCount}</span>
            <span> results</span>
          </div>

          <!-- Pagination Buttons -->
          <div class="flex space-x-1">
            <button 
              onclick={() => goToHistoryPage(historyCurrentPage - 1)}
              disabled={historyCurrentPage === 1}
              class="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg {historyCurrentPage === 1 ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-100'}"
            >
              Previous
            </button>

            {#each historyPaginationPages() as pageNum}
              <button 
                onclick={() => goToHistoryPage(pageNum)}
                class="px-3 py-2 text-sm font-medium border rounded-lg {pageNum === historyCurrentPage ? 'bg-blue-600 text-white border-blue-600' : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-100'}"
              >
                {pageNum}
              </button>
            {/each}

            <button 
              onclick={() => goToHistoryPage(historyCurrentPage + 1)}
              disabled={historyCurrentPage === historyPageCount}
              class="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg {historyCurrentPage === historyPageCount ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-100'}"
            >
              Next
            </button>
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>