<script lang="ts">
    import Sidebar from '$lib/components/shared/Sidebar.svelte';
    import { selectedNamespace } from '$lib/stores/namespace';

    let { children, data } = $props();
  </script>

  <div class="flex h-screen bg-gray-50">
    <Sidebar namespace={$selectedNamespace} />

    <div class="flex-1 flex flex-col overflow-hidden">
      <main class="flex-1 overflow-y-auto">
        {@render children()}
      </main>
    </div>
  </div>
