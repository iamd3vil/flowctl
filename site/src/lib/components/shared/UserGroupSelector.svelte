<script lang="ts">
  import { apiClient } from '$lib/apiClient';
  import { handleInlineError } from '$lib/utils/errorHandling';
  import type { User, Group } from '$lib/types';
  import { IconUsers, IconUser } from '@tabler/icons-svelte';
  
  let {
    type = $bindable('user'),
    selectedSubject = $bindable(null),
    placeholder = 'Search...',
    disabled = false
  }: {
    type: 'user' | 'group';
    selectedSubject: User | Group | null;
    placeholder?: string;
    disabled?: boolean;
  } = $props();

  let searchQuery = $state('');
  let searchResults = $state<(User | Group)[]>([]);
  let showDropdown = $state(false);
  let loading = $state(false);

  async function loadSubjects() {
    loading = true;
    showDropdown = true;
    try {
      if (type === 'user') {
        const response = await apiClient.users.list({
          filter: searchQuery,
          count_per_page: 20
        });
        searchResults = response.users || [];
      } else {
        const response = await apiClient.groups.list({
          filter: searchQuery,
          count_per_page: 20
        });
        searchResults = response.groups || [];
      }
    } catch (error) {
      handleInlineError(error, type === 'user' ? 'Unable to Load Users' : 'Unable to Load Groups');
      searchResults = [];
      showDropdown = false;
    } finally {
      loading = false;
    }
  }

  async function handleFocus() {
    if (searchResults.length === 0) {
      await loadSubjects();
    } else {
      showDropdown = true;
    }
  }

  function selectSubject(subject: User | Group) {
    selectedSubject = subject;
    searchQuery = '';
    showDropdown = false;
    searchResults = [];
  }

  function clearSelection() {
    selectedSubject = null;
    searchQuery = '';
    searchResults = [];
    showDropdown = false;
  }

  // Reset when type changes
  $effect(() => {
    clearSelection();
  });

  // Close dropdown when clicking outside
  function handleOutsideClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.user-group-selector')) {
      showDropdown = false;
    }
  }
</script>

<svelte:window onclick={handleOutsideClick} />

<div class="user-group-selector">
  <div class="relative">
    <input
      type="text"
      bind:value={searchQuery}
      oninput={loadSubjects}
      onfocus={handleFocus}
      {placeholder}
      class="bg-muted border border-input text-foreground text-sm rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent block w-full p-2.5 pr-10"
      autocomplete="off"
      {disabled}
    />
    
    {#if loading}
      <div class="absolute right-3 top-1/2 transform -translate-y-1/2">
        <svg class="animate-spin h-4 w-4 text-muted-foreground" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
    {:else}
      <svg class="w-5 h-5 absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
      </svg>
    {/if}
    
    <!-- Dropdown -->
    {#if showDropdown}
      <div class="absolute z-10 w-full mt-1 bg-card border border-input rounded-lg shadow-lg max-h-48 overflow-y-auto">
        {#if searchResults.length > 0}
          {#each searchResults as subject}
            <button
              type="button"
              class="w-full px-4 py-2 hover:bg-muted cursor-pointer border-b border-border last:border-b-0 text-left"
              onclick={() => selectSubject(subject)}
            >
              <div class="flex items-center">
                <div class="w-8 h-8 rounded-lg flex items-center justify-center mr-3 bg-primary-50">
                  {#if type === 'user'}
                    <IconUser class="w-4 h-4 text-primary-600" />
                  {:else}
                    <IconUsers class="w-4 h-4 text-primary-600" />
                  {/if}
                </div>
                <div>
                  <div class="text-sm font-medium text-foreground">{'name' in subject ? subject.name : subject.username}</div>
                  <div class="text-xs text-muted-foreground">{subject.id}</div>
                </div>
              </div>
            </button>
          {/each}
        {:else if !loading}
          <div class="px-4 py-3 text-sm text-muted-foreground text-center">
            {type === 'user' ? 'No users found' : 'No groups found'}
          </div>
        {/if}
      </div>
    {/if}
  </div>
  
  <!-- Selected subject display -->
  {#if selectedSubject}
    <div class="mt-2 p-2 bg-muted rounded-lg border">
      <div class="flex items-center justify-between">
        <div class="flex items-center">
          <div class="w-8 h-8 rounded-lg flex items-center justify-center mr-3 bg-primary-50">
            {#if type === 'user'}
              <IconUser class="w-4 h-4 text-primary-600" />
            {:else}
              <IconUsers class="w-4 h-4 text-primary-600" />
            {/if}
          </div>
          <div>
            <div class="text-sm font-medium text-foreground">{'name' in selectedSubject ? selectedSubject.name : selectedSubject.username}</div>
            <div class="text-xs text-muted-foreground">{selectedSubject.id}</div>
          </div>
        </div>
        <button type="button" onclick={clearSelection} class="text-muted-foreground hover:text-foreground" {disabled}>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>
    </div>
  {/if}
</div>