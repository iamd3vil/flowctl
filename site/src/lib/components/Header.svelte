<script lang="ts">
  import UserDropdown from './UserDropdown.svelte';
</script>

<!-- Header -->
<header class="bg-white shadow-sm border-b border-gray-200 px-6 py-4">
  <div class="flex items-center justify-between">
    <div class="flex items-center space-x-4">
      <!-- Breadcrumb or page title can go here -->
      <h1 class="text-xl font-semibold text-gray-900">
        <slot name="title">Flowctl</slot>
      </h1>
    </div>
    
    <div class="flex items-center space-x-4">
      <!-- Notifications or other actions can go here -->
      <UserDropdown />
    </div>
  </div>
</header>