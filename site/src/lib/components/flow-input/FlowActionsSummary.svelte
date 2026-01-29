<script lang="ts">
  import type { FlowAction } from '$lib/types';
  
  let { 
    actions,
    title = 'Flow Actions'
  }: {
    actions: FlowAction[],
    title?: string
  } = $props();
</script>

{#if actions && actions.length > 0}
  <div class="bg-card rounded-lg border border-border mt-6">
    <div class="px-4 py-3 border-b border-border">
      <h3 class="text-sm font-semibold text-foreground">
        {title} ({actions.length})
      </h3>
    </div>
    <div class="divide-y divide-border">
      {#each actions as action, index}
        <div class="px-4 py-3 flex items-center justify-between hover:bg-muted transition-colors">
          <div class="flex items-center gap-3">
            <div class="w-6 h-6 bg-primary-100 text-primary-600 rounded flex items-center justify-center text-xs font-medium">
              {index + 1}
            </div>
            <div class="text-sm text-foreground">{action.name}</div>
          </div>
          <span class="inline-flex px-2 py-1 text-xs font-medium rounded bg-subtle text-foreground">
            {action.executor}
          </span>
        </div>
      {/each}
    </div>
  </div>
{/if}