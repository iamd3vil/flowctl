<script lang="ts">
	import { browser } from '$app/environment';
	import { page } from '$app/state';
	import type { PageData } from './$types';
	import PageHeader from '$lib/components/shared/PageHeader.svelte';
	import Table from '$lib/components/shared/Table.svelte';
	import MemberCell from '$lib/components/members/MemberCell.svelte';
	import MemberTypeBadge from '$lib/components/members/MemberTypeBadge.svelte';
	import MemberRoleBadge from '$lib/components/members/MemberRoleBadge.svelte';
	import AddMemberModal from '$lib/components/members/AddMemberModal.svelte';
	import DeleteModal from '$lib/components/shared/DeleteModal.svelte';
	import { apiClient } from '$lib/apiClient';
	import type { NamespaceMemberResp, NamespaceMemberReq } from '$lib/types';
	import Header from '$lib/components/shared/Header.svelte';

	let { data }: { data: PageData } = $props();

	// State
	let members = $state(data.members);
	let loading = $state(false);
	let showAddModal = $state(false);
	let showDeleteModal = $state(false);
	let deleteMemberId = $state<string | null>(null);
	let deleteMemberName = $state('');

	// Table configuration
	let tableColumns = [
		{
			key: 'subject_name',
			header: 'Member',
			component: MemberCell
		},
		{
			key: 'subject_type',
			header: 'Type',
			component: MemberTypeBadge
		},
		{
			key: 'role',
			header: 'Role',
			component: MemberRoleBadge
		},
		{
			key: 'created_at',
			header: 'Added',
			render: (_value: any, member: NamespaceMemberResp) => formatDate(member.created_at)
		}
	];

	let tableActions = [
		{
			label: 'Remove',
			onClick: (member: NamespaceMemberResp) => handleDelete(member.id, member.subject_name),
			className: 'text-red-600 hover:text-red-800'
		}
	];

	// Functions
	async function fetchMembers() {
		if (!browser) return;
		
		loading = true;
		try {
			const response = await apiClient.namespaces.members.list(data.namespace);
			members = response.members || [];
		} catch (error) {
			console.error('Failed to fetch members:', error);
			notifyError('Failed to fetch members');
		} finally {
			loading = false;
		}
	}

	async function handleMemberSave(memberData: NamespaceMemberReq) {
		try {
			await apiClient.namespaces.members.add(data.namespace, memberData);
			showAddModal = false;
			await fetchMembers();
			notifySuccess('Member added successfully');
		} catch (error) {
			console.error('Failed to add member:', error);
			notifyError('Failed to add member');
			throw error; // Re-throw so modal can handle it
		}
	}

	function handleAdd() {
		showAddModal = true;
	}

	function handleDelete(memberId: string, memberName: string) {
		deleteMemberId = memberId;
		deleteMemberName = memberName;
		showDeleteModal = true;
	}

	async function confirmDelete() {
		if (!deleteMemberId) return;

		try {
			await apiClient.namespaces.members.remove(data.namespace, deleteMemberId);
			closeDeleteModal(); // Close modal after successful deletion
			await fetchMembers();
			notifySuccess('Member removed successfully');
		} catch (error) {
			console.error('Failed to remove member:', error);
			notifyError('Failed to remove member');
			throw error;
		}
	}

	function closeAddModal() {
		showAddModal = false;
	}

	function closeDeleteModal() {
		showDeleteModal = false;
		deleteMemberId = null;
		deleteMemberName = '';
	}

	function formatDate(dateString: string): string {
		if (!dateString) return 'Unknown';
		const date = new Date(dateString);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

		if (diffDays === 0) return 'Today';
		if (diffDays === 1) return 'Yesterday';
		if (diffDays < 7) return `${diffDays} days ago`;
		return date.toLocaleDateString();
	}

	function notifySuccess(message: string) {
		window.dispatchEvent(
			new CustomEvent("notify", {
				detail: { message, type: "success" },
			})
		);
	}

	function notifyError(message: string) {
		window.dispatchEvent(
			new CustomEvent("notify", {
				detail: { message, type: "error" },
			})
		);
	}
</script>

<svelte:head>
  <title>Members - {page.params.namespace} - Flowctl</title>
</svelte:head>

<Header breadcrumbs={[`${page.params.namespace}`, "Members"]}>
  {#snippet children()}
    <!-- Empty slot for now -->
  {/snippet}
</Header>

<div class="p-12">
	<!-- Page Header -->
	<PageHeader 
		title="Members"
		subtitle="Manage user and group access to this namespace"
		actions={[
			{
				label: 'Add Member',
				onClick: handleAdd,
				variant: 'primary',
				icon: '<svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>'
			}
		]}
	/>

	<!-- Members Table -->
	<div class="pt-6">
		<Table
			data={members}
			columns={tableColumns}
			actions={tableActions}
			{loading}
			emptyMessage="No members found. Get started by adding users or groups to this namespace."
		/>
	</div>
</div>

<!-- Add Member Modal -->
{#if showAddModal}
	<AddMemberModal
		show={showAddModal}
		onSave={handleMemberSave}
		onClose={closeAddModal}
	/>
{/if}

<!-- Delete Modal -->
{#if showDeleteModal}
	<DeleteModal
		title="Remove Member"
		itemName={deleteMemberName}
		onConfirm={confirmDelete}
		onClose={closeDeleteModal}
	/>
{/if}