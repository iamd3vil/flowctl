<script lang="ts">
  import Sidebar from '$lib/components/Sidebar.svelte';
  import Header from '$lib/components/Header.svelte';
  
  let { children, data } = $props();
</script>

<div class="flex h-screen bg-gray-50">
  <Sidebar namespace={data.namespace} />
  
  <div class="flex-1 flex flex-col overflow-hidden">
    <Header />
    <main class="flex-1 overflow-y-auto p-6">
      {@render children()}
    </main>
  </div>
</div>