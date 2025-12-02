<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { isAuthenticated, isLoading } from '$lib/stores/auth';
  import { getDefaultNamespace } from '$lib/utils/navigation';
  import logo from '$lib/assets/full-logo.svg';
  import LoadingSpinner from '$lib/components/shared/LoadingSpinner.svelte';

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

<div class="min-h-screen flex items-center justify-center bg-white">
  <div class="flex flex-col items-center gap-6">
    <img src={logo} alt="Flowctl" class="h-16 w-auto" />
    <LoadingSpinner label="Loading..." />
  </div>
</div>