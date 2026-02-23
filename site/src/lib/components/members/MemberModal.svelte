<script lang="ts">
  import { handleInlineError } from '$lib/utils/errorHandling';
  import { autofocus } from '$lib/utils/autofocus';
  import UserGroupSelector from '$lib/components/shared/UserGroupSelector.svelte';
  import type { NamespaceMemberReq, NamespaceMemberResp, FlowGroupResp, User, Group } from '$lib/types';
  import { apiClient } from '$lib/apiClient';
  import { IconUsers, IconUser, IconX, IconPlus, IconFolder } from '@tabler/icons-svelte';

  interface Props {
    isEditMode?: boolean;
    memberData?: NamespaceMemberResp | null;
    namespace: string;
    onSave: (memberData: NamespaceMemberReq) => void;
    onClose: () => void;
  }

  let {
    isEditMode = false,
    memberData = null,
    namespace,
    onSave,
    onClose
  }: Props = $props();


  // Form state
  let memberForm = $state<NamespaceMemberReq>({
    subject_type: 'user',
    subject_id: '',
    role: 'user'
  });

  let selectedSubject = $state<User | Group | null>(null);
  let loading = $state(false);

  // Group access state
  let memberPrefixes = $state<FlowGroupResp[]>([]);
  let availablePrefixes = $state<FlowGroupResp[]>([]);
  let prefixLoading = $state(false);
  let newPrefix = $state('');

  // Initialize form data when memberData changes
  $effect(() => {
    if (isEditMode && memberData) {
      // Pre-populate form with existing member data
      memberForm.subject_type = memberData.subject_type as 'user' | 'group';
      memberForm.subject_id = memberData.subject_id;
      memberForm.role = memberData.role as 'user' | 'reviewer' | 'admin';

      // Create a selectedSubject object for the UserGroupSelector
      selectedSubject = {
        id: memberData.subject_id,
        name: memberData.subject_name,
        username: memberData.subject_name // For users, this would be the username
      } as User | Group;

      // Load prefixes for this member
      loadMemberPrefixes();
    } else {
      resetForm();
    }
  });

  // Load available prefixes when role is user (for the add dropdown)
  $effect(() => {
    if (isEditMode && memberForm.role === 'user') {
      loadAvailablePrefixes();
    }
  });

  // Update subject_id when selectedSubject changes
  $effect(() => {
    memberForm.subject_id = selectedSubject?.id || '';
  });

  async function loadMemberPrefixes() {
    if (!isEditMode || !memberData) return;
    try {
      const result = await apiClient.namespaces.members.groups.list(namespace, memberData.id);
      memberPrefixes = result.groups || [];
    } catch {
      memberPrefixes = [];
    }
  }

  async function loadAvailablePrefixes() {
    try {
      const result = await apiClient.flows.groups.list(namespace);
      availablePrefixes = result.groups || [];
    } catch {
      availablePrefixes = [];
    }
  }

  async function addPrefix(prefix: string) {
    if (!memberData || !prefix.trim()) return;
    prefixLoading = true;
    try {
      await apiClient.namespaces.members.groups.add(namespace, memberData.id, { prefix: prefix.trim() });
      await loadMemberPrefixes();
      newPrefix = '';
    } catch (err) {
      handleInlineError(err, 'Unable to Grant Group Access');
    } finally {
      prefixLoading = false;
    }
  }

  async function removePrefix(prefix: string) {
    if (!memberData) return;
    prefixLoading = true;
    try {
      await apiClient.namespaces.members.groups.remove(namespace, memberData.id, prefix);
      await loadMemberPrefixes();
    } catch (err) {
      handleInlineError(err, 'Unable to Revoke Group Access');
    } finally {
      prefixLoading = false;
    }
  }

  // Get groups that haven't been assigned yet
  const unassignedPrefixes = $derived(() => {
    const assignedPrefixes = new Set(memberPrefixes.map(p => p.prefix));
    return availablePrefixes.filter(p => !assignedPrefixes.has(p.prefix));
  });

  function onSubjectTypeChange() {
    selectedSubject = null;
    memberForm.subject_id = '';
  }

  async function handleSubmit() {
    try {
      loading = true;

      // Basic client-side validation
      if (!selectedSubject) {
        handleInlineError(new Error('Please select a member'), 'Validation Error');
        return;
      }

      if (!memberForm.role) {
        handleInlineError(new Error('Please select a role'), 'Validation Error');
        return;
      }

      onSave(memberForm);
    } catch (err) {
      handleInlineError(err, isEditMode ? 'Unable to Update Member Role' : 'Unable to Add Member to Namespace');
    } finally {
      loading = false;
    }
  }

  function handleClose() {
    onClose();
  }

  function resetForm() {
    memberForm = {
      subject_type: 'user',
      subject_id: '',
      role: 'user'
    };
    selectedSubject = null;
    memberPrefixes = [];
    newPrefix = '';
  }

  // Close on Escape key
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      handleClose();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- Modal Backdrop -->
<div class="fixed inset-0 z-50 flex items-center justify-center bg-overlay p-4" onclick={handleClose} onkeydown={(e) => e.key === 'Escape' && handleClose()} role="dialog" aria-modal="true" tabindex="-1">
  <!-- Modal Content -->
  <div class="bg-card rounded-lg shadow-lg w-full max-w-lg max-h-[90vh] overflow-y-auto" onclick={(e) => e.stopPropagation()}>
      <div class="p-6">
        <h3 class="font-bold text-lg mb-4 text-foreground">
          {isEditMode ? 'Edit Member' : 'Add Member'}
        </h3>


        <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
          <!-- Subject Type Selection -->
          <div class="mb-4">
            <label class="block mb-1 font-medium text-foreground">Member Type *</label>
            <select
              bind:value={memberForm.subject_type}
              onchange={onSubjectTypeChange}
              class="bg-muted border border-input text-foreground text-sm rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent block w-full p-2.5"
              required
              disabled={loading || isEditMode}
              use:autofocus
            >
              <option value="user">User</option>
              <option value="group">Group</option>
            </select>
            {#if isEditMode}
              <p class="text-xs text-muted-foreground mt-1">Member type cannot be changed when editing.</p>
            {/if}
          </div>

          <!-- User/Group Selection -->
          <div class="mb-4">
            <label class="block mb-1 font-medium text-foreground">
              {memberForm.subject_type === 'user' ? 'User' : 'Group'} *
            </label>
            {#if isEditMode && memberData}
              <!-- Show selected member in edit mode -->
              <div class="p-3 bg-muted rounded-lg border">
                <div class="flex items-center">
                  <div class="w-8 h-8 rounded-lg flex items-center justify-center mr-3 bg-primary-100">
                    {#if memberData.subject_type === 'user'}
                      <IconUser class="w-4 h-4 text-primary-600" />
                    {:else}
                      <IconUsers class="w-4 h-4 text-primary-600" />
                    {/if}
                  </div>
                  <div>
                    <div class="text-sm font-medium text-foreground">{memberData.subject_name}</div>
                    <div class="text-xs text-muted-foreground">{memberData.subject_id}</div>
                  </div>
                </div>
              </div>
              <p class="text-xs text-muted-foreground mt-1">Member cannot be changed when editing.</p>
            {:else}
              <UserGroupSelector
                bind:type={memberForm.subject_type}
                bind:selectedSubject={selectedSubject}
                placeholder="Search {memberForm.subject_type}s..."
                disabled={loading}
              />
            {/if}
          </div>

          <!-- Role Selection -->
          <div class="mb-4">
            <label class="block mb-1 font-medium text-foreground">Role *</label>
            <select
              bind:value={memberForm.role}
              class="bg-muted border border-input text-foreground text-sm rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent block w-full p-2.5"
              required
              disabled={loading}
            >
              <option value="user">User - Can view and trigger flows</option>
              <option value="reviewer">Reviewer - Can approve flows and view all content</option>
              <option value="admin">Admin - Full access to namespace management</option>
            </select>
          </div>

          <!-- Prefix Access (edit mode, user role only) -->
          {#if isEditMode && memberData && memberData.role === 'user' && memberForm.role === 'user'}
            <div class="mb-4">
              <label class="block mb-1 font-medium text-foreground">Prefix Access</label>
              <p class="text-xs text-muted-foreground mb-2">
                Users can only see ungrouped flows by default. Grant access to specific prefixes below.
              </p>

              <!-- Current prefixes -->
              {#if memberPrefixes.length > 0}
                <div class="flex flex-wrap gap-2 mb-3">
                  {#each memberPrefixes as prefix}
                    <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-md text-sm bg-primary-100 text-primary-700">
                      <IconFolder class="w-3.5 h-3.5" />
                      {prefix.prefix}
                      <button
                        type="button"
                        onclick={() => removePrefix(prefix.prefix)}
                        disabled={prefixLoading}
                        class="ml-0.5 hover:text-danger-600 cursor-pointer disabled:opacity-50"
                        title="Remove access"
                      >
                        <IconX class="w-3.5 h-3.5" />
                      </button>
                    </span>
                  {/each}
                </div>
              {/if}

              <!-- Add from available prefixes -->
              {#if unassignedPrefixes().length > 0}
                <div class="flex gap-2">
                  <select
                    bind:value={newPrefix}
                    class="bg-muted border border-input text-foreground text-sm rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent block flex-1 p-2"
                    disabled={prefixLoading}
                  >
                    <option value="">Select a group...</option>
                    {#each unassignedPrefixes() as prefix}
                      <option value={prefix.prefix}>{prefix.prefix}{prefix.description ? ` — ${prefix.description}` : ''}</option>
                    {/each}
                  </select>
                  <button
                    type="button"
                    onclick={() => newPrefix && addPrefix(newPrefix)}
                    disabled={prefixLoading || !newPrefix}
                    class="inline-flex items-center px-3 py-2 text-sm font-medium text-white bg-primary-500 rounded-lg hover:bg-primary-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                  >
                    <IconPlus class="w-4 h-4" />
                  </button>
                </div>
              {:else if memberPrefixes.length === 0}
                <p class="text-xs text-muted-foreground italic">No flow prefixes available in this namespace.</p>
              {/if}
            </div>
          {/if}

          <!-- Actions -->
          <div class="flex justify-end gap-2 mt-6">
            <button
              type="button"
              onclick={handleClose}
              class="inline-flex items-center px-5 py-2.5 text-sm font-medium text-foreground bg-subtle rounded-lg hover:bg-subtle-hover disabled:opacity-50 cursor-pointer"
              disabled={loading}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              class="inline-flex items-center px-5 py-2.5 text-sm font-medium text-white bg-primary-500 rounded-lg hover:bg-primary-600 focus:ring-4 focus:outline-none focus:ring-primary-300 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
            >
              {#if loading}
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              {/if}
              {isEditMode ? 'Update' : 'Add'}
            </button>
          </div>
        </form>
      </div>
    </div>
</div>
