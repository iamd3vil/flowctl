<script lang="ts">
	import { browser } from '$app/environment';
	import { page } from '$app/state';
	import type { PageData } from './$types';
	import PageHeader from '$lib/components/shared/PageHeader.svelte';
	import Header from '$lib/components/shared/Header.svelte';
	import NamespaceSecretsModal from '$lib/components/namespace-secrets/NamespaceSecretsModal.svelte';
	import DeleteModal from '$lib/components/shared/DeleteModal.svelte';
	import { apiClient } from '$lib/apiClient';
	import type { NamespaceSecretReq, NamespaceSecretUpdateReq, NamespaceSecretResp } from '$lib/types';
	import { handleInlineError, showSuccess } from '$lib/utils/errorHandling';
	import { formatDateTime } from '$lib/utils';
	import { IconPlus, IconLock } from '@tabler/icons-svelte';

	let { data }: { data: PageData } = $props();

	// State
	let secrets = $state<NamespaceSecretResp[]>([]);
	let loading = $state(true);
	let showModal = $state(false);
	let showDeleteModal = $state(false);
	let selectedSecret = $state<NamespaceSecretResp | null>(null);
	let isEditMode = $state(false);

	// Handle the async data from load function
	$effect(() => {
		let cancelled = false;

		// Load secrets
		data.secretsPromise
			.then((result: NamespaceSecretResp[]) => {
				if (!cancelled) {
					secrets = result || [];
					loading = false;
				}
			})
			.catch((err: Error) => {
				if (!cancelled) {
					secrets = [];
					handleInlineError(err, "Unable to Load Secrets");
					loading = false;
				}
			});

		return () => {
			cancelled = true;
		};
	});

	async function fetchSecrets() {
		if (!browser) return;

		loading = true;
		try {
			const response = await apiClient.namespaceSecrets.list(data.namespace);
			secrets = response;
		} catch (error) {
			handleInlineError(error, 'Unable to Load Secrets');
		} finally {
			loading = false;
		}
	}

	function handleAdd() {
		selectedSecret = null;
		isEditMode = false;
		showModal = true;
	}

	function handleEdit(secret: NamespaceSecretResp) {
		selectedSecret = secret;
		isEditMode = true;
		showModal = true;
	}

	function handleDelete(secret: NamespaceSecretResp) {
		selectedSecret = secret;
		showDeleteModal = true;
	}

	async function handleSave(secretData: NamespaceSecretReq | NamespaceSecretUpdateReq) {
		try {
			if (isEditMode && selectedSecret) {
				await apiClient.namespaceSecrets.update(data.namespace, selectedSecret.id, secretData as NamespaceSecretUpdateReq);
				showSuccess('Secret Updated', 'Secret updated successfully');
			} else {
				await apiClient.namespaceSecrets.create(data.namespace, secretData as NamespaceSecretReq);
				showSuccess('Secret Created', 'Secret created successfully');
			}

			showModal = false;
			await fetchSecrets();
		} catch (error) {
			handleInlineError(error, isEditMode ? 'Unable to Update Secret' : 'Unable to Create Secret');
		}
	}

	async function confirmDelete() {
		if (!selectedSecret) return;

		try {
			await apiClient.namespaceSecrets.delete(data.namespace, selectedSecret.id);
			showSuccess('Secret Deleted', 'Secret deleted successfully');
			closeDeleteModal();
			await fetchSecrets();
		} catch (error) {
			handleInlineError(error, 'Unable to Delete Secret');
		}
	}

	function closeDeleteModal() {
		showDeleteModal = false;
		selectedSecret = null;
	}

	function handleModalClose() {
		showModal = false;
		isEditMode = false;
		selectedSecret = null;
	}
</script>

<svelte:head>
  <title>Secrets - {page.params.namespace} - Flowctl</title>
</svelte:head>

<Header breadcrumbs={[
  { label: page.params.namespace!, url: `/view/${page.params.namespace}/flows` },
  { label: "Secrets" }
]}>
  {#snippet children()}
    <div class="mb-10"></div>
  {/snippet}
</Header>

<div class="p-12">
	<!-- Page Header -->
	<PageHeader
		title="Namespace Secrets"
		subtitle="Manage encrypted secrets available to all flows in this namespace. Flow-level secrets with the same key will override these."
		actions={data.permissions?.canCreate ? [
			{
				label: 'Add Secret',
				onClick: handleAdd,
				variant: 'primary',
				IconComponent: IconPlus,
				iconSize: 16
			}
		] : []}
	/>

	<!-- Secrets List -->
	{#if loading}
		<div class="flex items-center justify-center py-8">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
		</div>
	{:else if secrets.length === 0}
		<div class="text-center py-8">
			<div class="text-gray-500">
				<IconLock size={48} class="mx-auto h-12 w-12 text-gray-400" />
				<h3 class="mt-2 text-sm font-medium text-gray-900">No namespace secrets yet</h3>
				<p class="mt-1 text-sm text-gray-500">Add secrets that will be available to all flows in this namespace.</p>
			</div>
		</div>
	{:else}
		<div class="bg-white shadow overflow-hidden sm:rounded-md">
			<ul role="list" class="divide-y divide-gray-200">
				{#each secrets as secret}
					<li class="px-4 py-4 flex items-center justify-between">
						<div class="flex-1 min-w-0">
							<div class="flex items-center space-x-3">
								<div class="flex-shrink-0">
									<svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
									</svg>
								</div>
								<div class="min-w-0 flex-1">
									<p class="text-sm font-medium text-gray-900 truncate">
										{secret.key}
									</p>
									{#if secret.description}
										<p class="text-sm text-gray-500 truncate">
											{secret.description}
										</p>
									{/if}
									<p class="text-xs text-gray-400">
										Created: {formatDateTime(secret.created_at)}
									</p>
								</div>
							</div>
						</div>

						<div class="flex items-center space-x-2">
							<button
								onclick={() => handleEdit(secret)}
								disabled={!data.permissions?.canUpdate}
								class="text-primary-600 hover:text-primary-500 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
								title="Edit secret"
							>
								<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
								</svg>
							</button>
							<button
								onclick={() => handleDelete(secret)}
								disabled={!data.permissions?.canDelete}
								class="text-danger-600 hover:text-danger-500 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
								title="Delete secret"
							>
								<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
								</svg>
							</button>
						</div>
					</li>
				{/each}
			</ul>
		</div>
	{/if}
</div>

<!-- Secret Modal -->
{#if showModal}
	<NamespaceSecretsModal
		{isEditMode}
		secretData={selectedSecret}
		onSave={handleSave}
		onClose={handleModalClose}
	/>
{/if}

<!-- Delete Modal -->
{#if showDeleteModal}
	<DeleteModal
		title="Delete Secret"
		itemName={selectedSecret?.key ?? ''}
		description="This may break flow executions that depend on this secret."
		onConfirm={confirmDelete}
		onClose={closeDeleteModal}
	/>
{/if}
