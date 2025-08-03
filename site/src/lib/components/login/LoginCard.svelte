<script lang="ts">
  let { 
    onSubmit, 
    error, 
    loading,
    username = $bindable(''),
    password = $bindable('')
  }: { 
    onSubmit: (event: SubmitEvent) => void, 
    error: string, 
    loading: boolean,
    username: string,
    password: string
  } = $props();
</script>

<!-- Login Card -->
<div class="p-8 rounded-lg border bg-white border-slate-200 shadow-sm">
  <form onsubmit={onSubmit} class="space-y-6">
    <!-- Error Message -->
    {#if error}
      <div class="p-3 rounded-md bg-red-50 border border-red-200">
        <div class="text-sm text-red-700">{error}</div>
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
        class="w-full px-3 py-2 text-sm rounded-md border bg-white border-slate-200 text-slate-900 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-colors"
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
        class="w-full px-3 py-2 text-sm rounded-md border bg-white border-slate-200 text-slate-900 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-colors"
        placeholder="Enter your password"
      />
    </div>

    <!-- Sign In Button -->
    <button
      type="submit"
      disabled={loading}
      class="w-full px-4 py-2 text-sm font-medium rounded-md transition-all duration-300 disabled:opacity-50 bg-blue-500 text-white hover:bg-blue-600 disabled:cursor-not-allowed"
    >
      {#if loading}
        <span class="flex items-center justify-center">
          <svg class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Signing In...
        </span>
      {:else}
        Sign In
      {/if}
    </button>

    <!-- OIDC Button -->
    <a
      href="/login/oidc"
      class="block w-full px-4 py-2 text-sm font-medium text-center rounded-md border transition-all duration-300 bg-white border-slate-200 text-slate-600 hover:bg-slate-50 hover:text-slate-900 hover:border-slate-300"
    >
      <div class="flex items-center justify-center">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1721 9z"></path>
        </svg>
        Sign in with OIDC
      </div>
    </a>
  </form>
</div>