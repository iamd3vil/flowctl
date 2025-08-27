<script lang="ts">
  import UserDropdown from './UserDropdown.svelte';
  import { goto } from '$app/navigation';
  
  type BreadcrumbItem = {
    label: string;
    url?: string;
  };
  
  let { 
    breadcrumbs = [], 
    actions = [],
    showUserDropdown = true,
    children
  }: { 
    breadcrumbs?: (string | BreadcrumbItem)[],
    actions?: Array<{ label: string, onClick: () => void, variant?: 'primary' | 'secondary' | 'danger' }>,
    showUserDropdown?: boolean,
    children?: any
  } = $props();

  // Convert breadcrumbs to uniform format
  const normalizedBreadcrumbs = $derived(
    breadcrumbs.map(crumb => 
      typeof crumb === 'string' ? { label: crumb } : crumb
    )
  );

  const handleBreadcrumbClick = (crumb: BreadcrumbItem) => {
    if (crumb.url) {
      goto(crumb.url);
    }
  };
</script>

<!-- Header -->
<header class="bg-white shadow-sm border-b border-gray-200 px-6 py-4">
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-5">
      {#if normalizedBreadcrumbs.length > 0}
        <div class="flex items-center text-sm text-gray-500">
          {#each normalizedBreadcrumbs as crumb, index}
            {#if crumb.url && index < normalizedBreadcrumbs.length - 1}
              <button 
                onclick={() => handleBreadcrumbClick(crumb)}
                class="hover:text-blue-600 hover:underline transition-colors cursor-pointer"
              >
                {crumb.label}
              </button>
            {:else}
              <span class={index === normalizedBreadcrumbs.length - 1 ? 'text-gray-900' : ''}>{crumb.label}</span>
            {/if}
            {#if index < normalizedBreadcrumbs.length - 1}
              <span class="mx-2">/</span>
            {/if}
          {/each}
        </div>
      {/if}
    </div>
    
    <div class="flex items-center space-x-4">
      {#if children}
        {@render children()}
      {/if}
      
      {#each actions as action}
        <button 
          onclick={action.onClick}
          class="inline-flex items-center gap-2 px-4 py-2 rounded-md transition-colors {action.variant === 'primary' ? 'bg-blue-600 text-white hover:bg-blue-700' : action.variant === 'danger' ? 'bg-red-600 text-white hover:bg-red-700' : 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'}"
        >
          {action.label}
        </button>
      {/each}
      
      {#if showUserDropdown}
        <UserDropdown />
      {/if}
    </div>
  </div>
</header>