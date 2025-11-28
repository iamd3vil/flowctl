<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { isAuthenticated, isLoading } from '$lib/stores/auth';
  import { getDefaultNamespace } from '$lib/utils/navigation';

  onMount(async () => {
    if ($isLoading) {
      await new Promise(resolve => setTimeout(resolve, 100));
    }

    if (!$isAuthenticated) {
      goto('/login');
      return;
    }

    try {
      const namespace = await getDefaultNamespace();
      goto(`/view/${namespace}/flows`);
    } catch {
      goto('/login');
    }
  });
</script>

<svelte:head>
  <title>Flowctl</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
  <div class="text-center">
    <div class="w-16 h-16 mx-auto mb-4 flex items-center justify-center rounded-lg bg-primary-500">
      <span class="text-2xl font-bold text-white">F</span>
    </div>
    <h1 class="text-2xl font-bold text-gray-900 mb-2">Flowctl</h1>
    <p class="text-gray-600">Loading...</p>
  </div>
</div>