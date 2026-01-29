<script lang="ts">
  import Sidebar from '$lib/components/shared/Sidebar.svelte';
  import Header from '$lib/components/shared/Header.svelte';

  let { children, data } = $props();
</script>

<svelte:head>
  <meta name="viewport" content="width=1024, initial-scale=1.0, user-scalable=yes" />
</svelte:head>

<div class="flex h-screen bg-muted">
  <Sidebar namespace={data.namespace} />
  
  <div class="flex-1 flex flex-col overflow-hidden">
    <main class="flex-1 overflow-y-auto">
      {@render children()}
    </main>
  </div>
</div>