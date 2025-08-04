<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { isAuthenticated } from '$lib/stores/auth';

  const errorMessage = page.error?.message || 'An unexpected error occurred';

  const handleGoHome = () => {
    if ($isAuthenticated) {
      goto('/view/default/flows');
    } else {
      goto('/');
    }
  };

</script>

<svelte:head>
  <title>Error - Flowctl</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@tabler/icons-webfont@latest/tabler-icons.min.css">
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50 px-4">
  <div class="max-w-lg w-full text-center">
    <!-- Error Icon -->
    <div class="mb-8">
      <div class="mx-auto w-24 h-24 bg-red-100 rounded-full flex items-center justify-center">
        <i class="ti ti-alert-triangle text-4xl text-red-600"></i>
      </div>
    </div>

    <!-- Error Content -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-4">
        Something went wrong
      </h1>
      <p class="text-lg text-gray-600 mb-2">
        {errorMessage || 'An unexpected error occurred'}
      </p>
    </div>

    <!-- Action Buttons -->
    <div class="flex flex-col sm:flex-row gap-3 justify-center">
      <button 
        onclick={handleGoHome}
        class="inline-flex items-center gap-2 px-6 py-3 bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition-colors font-medium"
      >
        <i class="ti ti-home text-lg"></i>
        Home
      </button>
    </div>

  </div>
</div>