<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import FlowInputForm from '$lib/components/flow-input/FlowInputForm.svelte';
  import Table from '$lib/components/shared/Table.svelte';
  import Header from '$lib/components/shared/Header.svelte';
  import StatusBadge from '$lib/components/shared/StatusBadge.svelte';
  import Tabs from '$lib/components/shared/Tabs.svelte';
  import Pagination from '$lib/components/shared/Pagination.svelte';
  import FlowHero from '$lib/components/flows/FlowHero.svelte';
  import FlowActionsSummary from '$lib/components/flow-input/FlowActionsSummary.svelte';
  import FlowSchedulesList from '$lib/components/flows/FlowSchedulesList.svelte';
  import { handleInlineError } from '$lib/utils/errorHandling';
  import type { PageData } from './$types';
  import type { TableColumn } from '$lib/types';
  import { DEFAULT_PAGE_SIZE } from '$lib/constants';
  import { permissionChecker } from '$lib/utils/permissions';
  import { formatDateTime } from '$lib/utils';
    import { IconPencil } from '@tabler/icons-svelte';

  let { data }: { data: PageData } = $props();

  let activeTab = $state<'run' | 'history'>('run');
  let historyLoading = $state(false);
  let flowExecutions = $state<any[]>([]);
  let historyCurrentPage = $state(1);
  let historyItemsPerPage = $state(DEFAULT_PAGE_SIZE);
  let historyTotalCount = $state(0);
  let historyPageCount = $state(0);
  let canUpdateFlow = $state(false);

  let namespace = $derived(page.params.namespace);
  let flowId = $derived(page.params.flowId);
  let rerunFromExecId = $derived(data.rerunFromExecId);
  let showRerunBanner = $state(!!rerunFromExecId);

  // Check update permission on mount
  permissionChecker(data.user!, 'flow', data.namespaceId, ['update']).then(permissions => {
    canUpdateFlow = permissions.canUpdate;
  });


  const loadFlowHistory = async () => {
    historyLoading = true;

    try {
      const response = await fetch(`/api/v1/${namespace}/flows/${flowId}/executions?page=${historyCurrentPage}&count_per_page=${historyItemsPerPage}`, {
        credentials: 'include',
      });
      const result = await response.json();

      if (!response.ok) {
        const apiError = new Error(result.error || 'Failed to fetch execution history');
        handleInlineError(apiError, 'Unable to Load Flow History');
        flowExecutions = [];
        return;
      }

      flowExecutions = result.executions || [];
      historyTotalCount = result.total_count || 0;
      historyPageCount = result.page_count || 1;
    } catch (error) {
      handleInlineError(error, 'Unable to Load Flow History');
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

  const handleHistoryPageChange = (event: CustomEvent<{ page: number }>) => {
    goToHistoryPage(event.detail.page);
  };

  const formatDuration = (startedAt: string, completedAt?: string) => {
    if (!startedAt) return 'Unknown';

    const start = new Date(startedAt);
    const end = completedAt ? new Date(completedAt) : new Date();
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
      header: 'Exec ID',
      render: (value) => `
        <a
          href="/view/${namespace}/results/${flowId}/${value}"
          class="text-sm text-primary-600 hover:underline hover:text-primary-800 font-mono block"
        >
          ${value.substring(0, 8)}
        </a>
      `
    },
    {
      key: 'started_at',
      header: 'Started At',
      width: 'w-40',
      render: (value) => `<div class="text-sm text-gray-600">${formatDateTime(value)}</div>`
    },
    {
      key: 'duration',
      header: 'Duration',
      render: (_value, row) => `<div class="text-sm text-gray-600">${formatDuration(row.started_at, row.completed_at)}</div>`
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
      render: (value) => `<div class="text-sm text-gray-900">${value || 'System'}</div>`
    },
    {
      key: 'trigger_type',
      header: 'Trigger Type',
      render: (value, row) => `<div class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
        row.trigger_type === 'manual'
          ? 'bg-primary-100 text-primary-900'
          : 'bg-success-100 text-success-900'
      }">${row.trigger_type}</div>`
    }
  ];

  // Tab configuration
  const tabs = [
    { id: 'run', label: 'Run' },
    { id: 'history', label: 'History' }
  ];
</script>

<svelte:head>
  <title>Run Flow - {data.flowMeta?.meta?.name || 'Loading...'}</title>
</svelte:head>

<Header
  breadcrumbs={[
    { label: namespace!, url: `/view/${namespace}/flows` },
    { label: 'Flows', url: `/view/${namespace}/flows` },
    { label: data.flowMeta?.meta?.name || 'Loading...' }
  ]}
  actions={[
    ...(canUpdateFlow ? [{ icon: IconPencil, label: 'Edit', onClick: () => goto(`/view/${namespace}/flows/${flowId}/edit`), variant: 'primary' as const }] : [])
  ]}
/>

<FlowHero
  name={data.flowMeta?.meta?.name || 'Loading...'}
  description={data.flowMeta?.meta?.description || ''}
/>

<div class="bg-white border-b border-gray-200 px-6">
  <div class="max-w-4xl mx-auto">
    <Tabs {tabs} bind:activeTab />
  </div>
</div>

<!-- Tab Content -->
<div class="px-6 py-8 bg-gray-50">
  {#if activeTab === 'run'}
    <div class="max-w-2xl mx-auto">
      {#if showRerunBanner}
        <div class="mb-6">
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 flex items-start justify-between">
            <div class="flex items-start gap-3">
              <svg class="w-5 h-5 text-blue-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <div class="flex-1">
                <h3 class="text-sm font-medium text-blue-900">Rerunning execution</h3>
                <p class="text-sm text-blue-700 mt-1">
                  Inputs have been prepopulated from execution
                  <a
                    href="/view/{namespace}/results/{flowId}/{rerunFromExecId}"
                    class="font-mono underline hover:text-blue-900"
                  >
                    {rerunFromExecId.substring(0, 8)}
                  </a>
                </p>
              </div>
            </div>
            <button
              onclick={() => showRerunBanner = false}
              class="text-blue-400 hover:text-blue-600"
              aria-label="Dismiss"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
        </div>
      {/if}
      <FlowInputForm
        inputs={data.flowInputs || []}
        namespace={namespace!}
        flowId={flowId!}
        executionInput={data.executionInput}
      />
      <div class="mt-6">
        <FlowSchedulesList schedules={data.flowMeta?.meta?.schedules || []} />
      </div>
      <FlowActionsSummary actions={data.flowMeta?.actions || []} />
    </div>
  {/if}

  <!-- History Tab -->
  {#if activeTab === 'history'}
    <div class="max-w-6xl mx-auto">

      <Table
        columns={tableColumns}
        data={flowExecutions}
        loading={historyLoading}
        title="Execution History for {data.flowMeta?.meta?.name || 'Flow'}"
        subtitle="Past executions of this flow"
        emptyMessage="No execution history"
      />

      {#if historyPageCount > 1}
        <div class="flex items-center justify-between mt-6">
          <div class="text-sm text-gray-700">
            Showing {((historyCurrentPage - 1) * historyItemsPerPage) + 1} to {Math.min(historyCurrentPage * historyItemsPerPage, historyTotalCount)} of {historyTotalCount} results
          </div>
          <Pagination
            currentPage={historyCurrentPage}
            totalPages={historyPageCount}
            loading={historyLoading}
            on:page-change={handleHistoryPageChange}
          />
        </div>
      {/if}
    </div>
  {/if}
</div>
