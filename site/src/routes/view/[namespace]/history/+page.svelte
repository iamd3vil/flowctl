<script lang="ts">
	import { browser } from '$app/environment';
	import { page } from '$app/state';
	import type { PageData } from './$types';
	import PageHeader from '$lib/components/shared/PageHeader.svelte';
	import SearchInput from '$lib/components/shared/SearchInput.svelte';
	import Table from '$lib/components/shared/Table.svelte';
	import Pagination from '$lib/components/shared/Pagination.svelte';
	import ExecutionIdCell from '$lib/components/shared/ExecutionIdCell.svelte';
	import StatusBadge from '$lib/components/shared/StatusBadge.svelte';
	import { apiClient } from '$lib/apiClient';
	import type { ExecutionSummary, ExecutionsPaginateResponse } from '$lib/types';
	import { DEFAULT_PAGE_SIZE } from '$lib/constants';
	import Header from '$lib/components/shared/Header.svelte';
	import { handleInlineError, showSuccess } from '$lib/utils/errorHandling';
	import { formatDateTime, getStartTime } from '$lib/utils';
	import { IconHistory } from '@tabler/icons-svelte';

	let { data }: { data: PageData } = $props();

	// State
	let executions = $state<ExecutionSummary[]>([]);
	let totalCount = $state(0);
	let pageCount = $state(0);
	let currentPage = $state(data.currentPage);
	let searchQuery = $state(data.searchQuery);
	let loading = $state(true);

	// Handle the async data from load function
	$effect(() => {
		let cancelled = false;

		data.executionsPromise
			.then((result: ExecutionsPaginateResponse) => {
				if (!cancelled) {
					executions = result.executions || [];
					totalCount = result.total_count || 0;
					pageCount = result.page_count || 1;
					loading = false;
				}
			})
			.catch((err: Error) => {
				if (!cancelled) {
					executions = [];
					totalCount = 0;
					pageCount = 0;
					handleInlineError(err, "Unable to Load Execution History");
					loading = false;
				}
			});

		return () => {
			cancelled = true;
		};
	});

	// Table configuration
	let tableColumns = [
		{
			key: 'flow_name',
			header: 'Flow Name',
			sortable: true,
			render: (_value: any, execution: ExecutionSummary) => `
					<a href="/view/${data.namespace}/flows/${execution.flow_id}"
					   class="text-sm hover:underline text-foreground hover:text-primary-600 font-medium"
					>
					  ${execution.flow_name}
					</a>
			`
		},
		{
			key: 'id',
			header: 'Exec ID',
			render: (_value: any, execution: ExecutionSummary) => `
     			<a
      				href="/view/${data.namespace}/results/${execution.flow_id}/${execution.id}"
      				class="font-mono text-link hover:underline text-sm"
     			>
                    ${execution.id.substring(0, 8)}
				</a>
			`
		},
		{
			key: 'status',
			header: 'Status',
			sortable: true,
			component: StatusBadge
		},
		{
			key: 'started_at',
			header: 'Started At',
			sortable: true,
			render: (_value: any, execution: ExecutionSummary) => `
				<div class="text-sm text-muted-foreground">${formatDateTime(getStartTime(execution))}</div>
			`
		},
		{
			key: 'duration',
			header: 'Duration',
			render: (_value: any, execution: ExecutionSummary) => `
				<div class="text-sm text-muted-foreground">${formatDuration(getStartTime(execution), execution.completed_at)}</div>
			`
		},
		{
			key: 'triggered_by',
			header: 'Triggered By',
			sortable: true,
			render: (_value: any, execution: ExecutionSummary) => `
				<div class="flex items-center">
					<div class="w-8 h-8 rounded-full bg-primary-100 flex items-center justify-center mr-3">
						<svg class="w-4 h-4 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
						</svg>
					</div>
					<span class="text-sm text-foreground">${execution.triggered_by || 'System'}</span>
				</div>
			`
		},
		{
			key: 'trigger_type',
			header: 'Trigger Type',
			sortable: true,
			render: (_value: any, execution: ExecutionSummary) => `
				<div class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
					execution.trigger_type === 'manual'
						? 'bg-primary-100 text-primary-900'
						: 'bg-success-100 text-success-900'
				}">
					${execution.trigger_type}
				</div>
			`
		}
	];


	// Functions
	async function fetchExecutions(filter: string = '', pageNumber: number = 1) {
		if (!browser) return;

		loading = true;
		try {
			const response = await apiClient.executions.list(data.namespace, {
				page: pageNumber,
				count_per_page: DEFAULT_PAGE_SIZE,
				filter: filter
			});

			executions = response.executions || [];
			totalCount = response.total_count || 0;
			pageCount = response.page_count || 1;
		} catch (error) {
			handleInlineError(error, 'Unable to Load Execution History');
		} finally {
			loading = false;
		}
	}

	function handleSearch(query: string) {
		searchQuery = query;
		fetchExecutions(query);
	}

	function handlePageChange(event: CustomEvent<{ page: number }>) {
		currentPage = event.detail.page;
		fetchExecutions('', currentPage);
	}

	function formatDuration(startedAt: string, completedAt?: string): string {
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
	}
</script>

<svelte:head>
  <title>Execution History - {page.params.namespace} - Flowctl</title>
</svelte:head>

<Header breadcrumbs={[
  { label: page.params.namespace!, url: `/view/${page.params.namespace}/flows` },
  { label: "History" }
]}>
  {#snippet children()}
    <SearchInput
      bind:value={searchQuery}
      placeholder="Search executions..."
      {loading}
      onSearch={handleSearch}
    />
  {/snippet}
</Header>

<div class="p-12">
	<!-- Page Header -->
	<PageHeader
		title="Execution History"
		subtitle="View all flow execution history across all flows in this namespace"
	/>

	<!-- Executions Table -->
	<div class="pt-6">
		<Table
			data={executions}
			columns={tableColumns}
			{loading}
			emptyMessage="No execution history found. Executions will appear here once flows are triggered."
			EmptyIconComponent={IconHistory}
			emptyIconSize={64}
		/>
	</div>

	<!-- Pagination -->
	{#if pageCount > 1}
		<Pagination
			currentPage={currentPage}
			totalPages={pageCount}
			on:page-change={handlePageChange}
		/>
	{/if}
</div>
