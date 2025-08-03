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
  <div class="bg-gray-50 rounded-lg p-4 mt-6">
    <h3 class="text-sm font-medium text-gray-900 mb-3 flex items-center">
      <svg class="w-4 h-4 text-gray-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
      </svg>
      <span>{title} ({actions.length} step{actions.length !== 1 ? 's' : ''})</span>
    </h3>
    <div class="space-y-2">
      {#each actions as action, index}
        <div>
          <div class="flex items-center justify-between p-3 bg-white border border-gray-200 rounded-md">
            <div class="flex items-center">
              <div class="w-6 h-6 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center text-xs font-medium mr-3">
                {index + 1}
              </div>
              <div class="text-sm font-medium text-gray-900">{action.name}</div>
            </div>
            <span class="inline-flex px-2 py-1 text-xs font-medium rounded-md bg-blue-100 text-blue-800">
              {action.executor}
            </span>
          </div>
          
          <!-- Arrow connecting actions -->
          {#if index < actions.length - 1}
            <div class="flex justify-center">
              <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"/>
              </svg>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  </div>
{/if}