<script lang="ts">
  import { IconKey } from '@tabler/icons-svelte';
  import { onMount } from 'svelte';
  import { apiClient } from '$lib/apiClient';
  import LoadingSpinner from '$lib/components/shared/LoadingSpinner.svelte';
  import type { SSOProvider } from '$lib/types';

  let {
    onSubmit,
    error,
    loading,
    username = $bindable(''),
    password = $bindable(''),
    redirectUrl = null
  }: {
    onSubmit: (event: SubmitEvent) => void,
    error: string,
    loading: boolean,
    username: string,
    password: string,
    redirectUrl: string | null
  } = $props();

  let ssoProviders: SSOProvider[] = $state([]);
  let oidcLoadingProvider = $state<string | null>(null);

  onMount(async () => {
    try {
      ssoProviders = await apiClient.auth.getSSOProviders();
    } catch (err) {
      console.error('Failed to fetch SSO providers:', err);
    }
  });

  const handleOIDCLogin = (providerId: string) => {
    oidcLoadingProvider = providerId;
    const url = redirectUrl
      ? `/login/oidc/${providerId}?redirect_url=${encodeURIComponent(redirectUrl)}`
      : `/login/oidc/${providerId}`;
    window.location.href = url;
  };
</script>

<!-- Login Card -->
<article class="p-8 rounded-lg border bg-white border-slate-200 shadow-sm">
  <form onsubmit={onSubmit} class="space-y-6" aria-label="Login form">
    <!-- Error Message -->
    {#if error}
      <div class="p-3 rounded-md bg-danger-50 border border-danger-200" role="alert" aria-live="assertive">
        <div class="text-sm text-danger-900">{error}</div>
      </div>
    {/if}

    <!-- Username Field -->
    <div>
      <label for="username" class="block text-sm font-medium mb-2 text-slate-900">
        Username
      </label>
      <input
        type="text"
        bind:value={username}
        id="username"
        name="username"
        required
        class="w-full px-3 py-2 text-sm rounded-md border bg-white border-slate-200 text-slate-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
        placeholder="Enter your username"
      />
    </div>

    <!-- Password Field -->
    <div>
      <label for="password" class="block text-sm font-medium mb-2 text-slate-900">
        Password
      </label>
      <input
        type="password"
        bind:value={password}
        id="password"
        name="password"
        required
        class="w-full px-3 py-2 text-sm rounded-md border bg-white border-slate-200 text-slate-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
        placeholder="Enter your password"
      />
    </div>

    <!-- Sign In Button -->
    <button
      type="submit"
      disabled={loading}
      class="w-full px-4 py-2 text-sm font-medium rounded-md transition-all duration-300 disabled:opacity-50 bg-primary-500 text-white hover:bg-primary-600 disabled:cursor-not-allowed"
    >
      {#if loading}
        <span class="flex items-center justify-center">
          <svg class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" aria-hidden="true">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Signing In...
        </span>
      {:else}
        Sign In
      {/if}
    </button>

    {#each ssoProviders as provider}
      <button
        type="button"
        onclick={() => handleOIDCLogin(provider.id)}
        disabled={oidcLoadingProvider !== null}
        class="w-full px-4 py-2 text-sm font-medium rounded-md border transition-all duration-300 bg-white border-slate-200 text-slate-600 hover:bg-slate-50 hover:text-slate-900 hover:border-slate-300 disabled:opacity-50 disabled:cursor-not-allowed"
        aria-label={provider.label}
      >
        {#if oidcLoadingProvider === provider.id}
          <div class="flex items-center justify-center">
            <LoadingSpinner size="sm" label="Redirecting..." />
          </div>
        {:else}
          <div class="flex items-center justify-center">
            <IconKey class="w-4 h-4 mr-2" aria-hidden="true" />
            {provider.label}
          </div>
        {/if}
      </button>
    {/each}
  </form>
</article>
