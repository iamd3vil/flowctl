<script lang="ts">
	import { browser } from '$app/environment';
	import { page } from '$app/state';
	import type { PageData } from './$types';
	import PageHeader from '$lib/components/shared/PageHeader.svelte';
	import SearchInput from '$lib/components/shared/SearchInput.svelte';
	import Table from '$lib/components/shared/Table.svelte';
	import Pagination from '$lib/components/shared/Pagination.svelte';
	import StatCard from '$lib/components/shared/StatCard.svelte';
	import NodeModal from '$lib/components/nodes/NodeModal.svelte';
	import DeleteModal from '$lib/components/shared/DeleteModal.svelte';
	import { apiClient } from '$lib/apiClient';
	import type { NodeResp, NodeReq, NodeStatsResp, NodesPaginateResponse, CredentialResp, CredentialsPaginateResponse } from '$lib/types';
	import { DEFAULT_PAGE_SIZE } from '$lib/constants';
	import Header from '$lib/components/shared/Header.svelte';
	import { handleInlineError, showSuccess } from '$lib/utils/errorHandling';
	import { IconPlus, IconServer } from '@tabler/icons-svelte';

	let { data }: { data: PageData } = $props();

	// State
	let nodes = $state<NodeResp[]>([]);
	let totalCount = $state(0);
	let pageCount = $state(0);
	let currentPage = $state(data.currentPage);
	let searchQuery = $state(data.searchQuery);
	let stats = $state<NodeStatsResp>({ total_hosts: 0, qssh_hosts: 0, ssh_hosts: 0 });
	let credentials = $state<CredentialResp[]>([]);
	let loading = $state(true);
	let showModal = $state(false);

	// Handle the async data from load function
	$effect(() => {
		let cancelled = false;

		// Load nodes
		data.nodesPromise
			.then((result: NodesPaginateResponse) => {
				if (!cancelled) {
					nodes = result.nodes || [];
					totalCount = result.total_count || 0;
					pageCount = result.page_count || 1;
					loading = false;
				}
			})
			.catch((err: Error) => {
				if (!cancelled) {
					nodes = [];
					totalCount = 0;
					pageCount = 0;
					handleInlineError(err, "Unable to Load Nodes");
					loading = false;
				}
			});

		// Load stats
		data.statsPromise
			.then((result: NodeStatsResp) => {
				if (!cancelled) {
					stats = result;
				}
			})
			.catch((err: Error) => {
				if (!cancelled) {
					stats = { total_hosts: 0, qssh_hosts: 0, ssh_hosts: 0 };
					handleInlineError(err, "Unable to Load Node Statistics");
				}
			});

		// Load credentials
		data.credentialsPromise
			.then((result: CredentialsPaginateResponse) => {
				if (!cancelled) {
					credentials = result.credentials || [];
				}
			})
			.catch((err: Error) => {
				if (!cancelled) {
					credentials = [];
					handleInlineError(err, "Unable to Load Credentials");
				}
			});

		return () => {
			cancelled = true;
		};
	});
	let isEditMode = $state(false);
	let editingNodeId = $state<string | null>(null);
	let editingNodeData = $state<NodeResp | null>(null);
	let showDeleteModal = $state(false);
	let deleteNodeId = $state<string | null>(null);
	let deleteNodeName = $state('');


	// Table configuration
	let tableColumns = [
		{
			key: 'name',
			header: 'Node',
			sortable: true,
			render: (_value: any, node: NodeResp) => `
				<div class="flex items-center">
					<div class="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center mr-3">
						<svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2 4h.01M17 16h.01"></path>
						</svg>
					</div>
					<div>
						<a href="#" class="text-sm hover:underline font-medium text-gray-900 cursor-pointer hover:text-primary-600 transition-colors" onclick="event.preventDefault(); document.dispatchEvent(new CustomEvent('editNode', {detail: {id: '${node.id}'}}))">${node.name}</a>
						<div class="text-sm text-gray-500">${node.id}</div>
					</div>
				</div>
			`
		},
		{ key: 'hostname', header: 'Hostname', sortable: true },
		{ key: 'port', header: 'Port', sortable: true },
		{ key: 'username', header: 'Username', sortable: true },
		{ key: 'os_family', header: 'OS Family', sortable: true },
		{
			key: 'connection_type',
			header: 'Connection Type',
			sortable: true,
			render: (_value: any, node: NodeResp) => `
				<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
					node.connection_type === 'qssh'
						? 'bg-success-100 text-success-800'
						: 'bg-blue-100 text-blue-800'
				}">${node.connection_type?.toUpperCase() || 'N/A'}</span>
			`
		},
		{
			key: 'tags',
			header: 'Tags',
			render: (_value: any, node: NodeResp) => node.tags && node.tags.length > 0
				? `<div class="flex flex-wrap gap-1">
					${node.tags.map(tag =>
						`<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-primary-100 text-primary-800">${tag}</span>`
					).join('')}
				</div>`
				: '<span class="text-xs text-gray-400">No tags</span>'
		}
	];

	let tableActions = [
		{
			label: 'Edit',
			onClick: (node: NodeResp) => handleEdit(node.id),
			className: 'text-primary-600 border-primary-600 hover:bg-primary-100'
		},
		{
			label: 'Delete',
			onClick: (node: NodeResp) => handleDelete(node.id),
			className: 'text-danger-600 border-danger-600 hover:bg-danger-100'
		}
	];

	async function fetchStats() {
		if (!browser) return;

		try {
			stats = await apiClient.nodes.getStats(data.namespace);
		} catch (error) {
			handleInlineError(error, 'Unable to Load Node Statistics');
		}
	}

	async function fetchNodes(filter: string = '', pageNumber: number = 1) {
		if (!browser) return;

		loading = true;
		try {
			const response = await apiClient.nodes.list(data.namespace, {
				page: pageNumber,
				count_per_page: DEFAULT_PAGE_SIZE,
				filter: filter
			});

			nodes = response.nodes || [];
			totalCount = response.total_count || 0;
			pageCount = response.page_count || 1;
		} catch (error) {
			handleInlineError(error, 'Unable to Load Nodes List');
		} finally {
			loading = false;
		}
	}

	function handleSearch(query: string) {
		searchQuery = query;
		fetchNodes(query);
	}

	function handlePageChange(event: CustomEvent<{ page: number }>) {
		currentPage = event.detail.page;
		fetchNodes('', currentPage);
	}

	function handleAdd() {
		isEditMode = false;
		editingNodeId = null;
		editingNodeData = null;
		showModal = true;
	}

	async function handleEdit(nodeId: string) {
		try {
			const node = await apiClient.nodes.getById(data.namespace, nodeId);

			isEditMode = true;
			editingNodeId = nodeId;
			editingNodeData = node;
			showModal = true;
		} catch (error) {
			handleInlineError(error, 'Unable to Load Node Details');
		}
	}

	function handleDelete(nodeId: string) {
		const node = nodes.find(n => n.id === nodeId);
		if (node) {
			deleteNodeId = nodeId;
			deleteNodeName = node.name;
			showDeleteModal = true;
		}
	}

	async function confirmDelete() {
		if (!deleteNodeId) return;

		try {
			await apiClient.nodes.delete(data.namespace, deleteNodeId);
			closeDeleteModal(); // Close modal after successful deletion
			showSuccess('Node deleted', `Node ${deleteNodeName} has been successfully deleted.`);
			await Promise.all([fetchNodes(), fetchStats()]);
		} catch (error) {
			handleInlineError(error, 'Unable to Delete Node');
		}
	}

	function closeDeleteModal() {
		showDeleteModal = false;
		deleteNodeId = null;
		deleteNodeName = '';
	}

	async function handleNodeSave(nodeData: NodeReq) {
		try {
			if (isEditMode && editingNodeId) {
				await apiClient.nodes.update(data.namespace, editingNodeId, nodeData);
				showSuccess('Node updated', `Node ${nodeData.name} has been successfully updated.`);
			} else {
				await apiClient.nodes.create(data.namespace, nodeData);
				showSuccess('Node created', `Node ${nodeData.name} has been successfully created.`);
			}
			showModal = false;
			await Promise.all([fetchNodes(), fetchStats()]);
		} catch (error) {
			handleInlineError(error, 'Unable to Save Node');
		}
	}

	function handleModalClose() {
		showModal = false;
		isEditMode = false;
		editingNodeId = null;
		editingNodeData = null;
	}

	// Handle node name clicks
	if (browser) {
		document.addEventListener('editNode', ((event: CustomEvent) => {
			handleEdit(event.detail.id);
		}) as EventListener);
	}
</script>

<svelte:head>
  <title>Nodes - {page.params.namespace} - Flowctl</title>
</svelte:head>

<Header breadcrumbs={[
  { label: page.params.namespace!, url: `/view/${page.params.namespace}/flows` },
  { label: "Nodes" }
]}>
  {#snippet children()}
    <SearchInput
      bind:value={searchQuery}
      placeholder="Search Nodes..."
      {loading}
      onSearch={handleSearch}
    />
  {/snippet}
</Header>

<div class="p-12">
	<!-- Page Header -->
	<PageHeader
		title="Nodes"
		subtitle="Manage remote nodes that run flows"
		actions={[
			{
				label: 'Add',
				onClick: handleAdd,
				variant: 'primary',
				IconComponent: IconPlus,
				iconSize: 16
			}
		]}
	/>

	<!-- Statistics Cards -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<StatCard
			title="Total Hosts"
			value={stats.total_hosts}
			IconComponent={IconServer}
		iconSize={24}
			color="blue"
		/>
		<StatCard
			title="QSSH Hosts"
			value={stats.qssh_hosts}
			IconComponent={IconServer}
		iconSize={24}
			color="green"
		/>
		<StatCard
			title="SSH Hosts"
			value={stats.ssh_hosts}
			IconComponent={IconServer}
		iconSize={24}
			color="blue"
		/>
	</div>

	<!-- Nodes Table -->
	<div class="pt-6">
		<Table
			data={nodes}
			columns={tableColumns}
			actions={tableActions}
			{loading}
			emptyMessage="No nodes found. Get started by adding your first node."
			EmptyIconComponent={IconServer}
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

<!-- Node Modal -->
{#if showModal}
	<NodeModal
		{isEditMode}
		nodeData={editingNodeData}
		{credentials}
		onSave={handleNodeSave}
		onClose={handleModalClose}
	/>
{/if}

<!-- Delete Modal -->
{#if showDeleteModal}
	<DeleteModal
		title="Delete Node"
		itemName={deleteNodeName}
		onConfirm={confirmDelete}
		onClose={closeDeleteModal}
	/>
{/if}
