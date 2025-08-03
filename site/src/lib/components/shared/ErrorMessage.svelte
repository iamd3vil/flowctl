<script lang="ts">
  let { 
    title = 'Error',
    message,
    onDismiss,
    dismissible = false
  }: {
    title?: string,
    message: string,
    onDismiss?: () => void,
    dismissible?: boolean
  } = $props();
</script>

<div class="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
  <div class="flex">
    <svg class="w-5 h-5 text-red-400 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
    </svg>
    <div class="ml-3 flex-1">
      <h3 class="text-sm font-medium text-red-800">{title}</h3>
      <p class="mt-1 text-sm text-red-700">{message}</p>
    </div>
    {#if dismissible && onDismiss}
      <button 
        onclick={onDismiss}
        class="ml-auto -mx-1.5 -my-1.5 bg-red-50 text-red-500 rounded-lg focus:ring-2 focus:ring-red-600 p-1.5 hover:bg-red-100 inline-flex h-8 w-8"
      >
        <span class="sr-only">Dismiss</span>
        <svg class="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
          <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
        </svg>
      </button>
    {/if}
  </div>
</div>