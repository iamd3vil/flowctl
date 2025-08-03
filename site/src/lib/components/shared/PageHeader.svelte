<script lang="ts">
  let { 
    title,
    subtitle,
    actions = []
  }: {
    title: string,
    subtitle?: string,
    actions?: Array<{ 
      label: string, 
      onClick: () => void, 
      variant?: 'primary' | 'secondary',
      icon?: string 
    }>
  } = $props();
</script>

<div class="flex items-center justify-between mb-6">
  <div>
    <h1 class="text-2xl font-bold text-gray-900">{title}</h1>
    {#if subtitle}
      <p class="text-gray-600">{subtitle}</p>
    {/if}
  </div>
  
  {#if actions.length > 0}
    <div class="flex items-center gap-3">
      {#each actions as action}
        <button 
          onclick={action.onClick}
          class="inline-flex items-center gap-2 px-4 py-2 rounded-md transition-colors {action.variant === 'primary' ? 'bg-blue-600 text-white hover:bg-blue-700' : 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'}"
        >
          {#if action.icon}
            {@html action.icon}
          {/if}
          {action.label}
        </button>
      {/each}
    </div>
  {/if}
</div>