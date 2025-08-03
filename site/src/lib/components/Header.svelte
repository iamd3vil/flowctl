<script lang="ts">
  import UserDropdown from './UserDropdown.svelte';
  let {title}: {title: string} = $props();
</script>

<!-- Header -->
<header class="bg-white shadow-sm border-b border-gray-200 px-6 py-4">
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-5">
      <div class="flex items-center text-sm text-gray-500">
        <span>Flows</span>
        <span class="mx-2">/</span>
        <span>{title}</span>
        <span class="mx-2">/</span>
        <span class="text-gray-900">Run Flow</span>
      </div>
    </div>
    
    <div class="flex items-center space-x-4">
      <!-- Notifications or other actions can go here -->
      <UserDropdown />
    </div>
  </div>
</header>