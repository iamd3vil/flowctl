<script lang="ts">
  import type { Schedule } from '$lib/types';
  import { IconClock } from '@tabler/icons-svelte';

  let {
    schedules,
    title = 'Active Schedules'
  }: {
    schedules: Schedule[],
    title?: string
  } = $props();
</script>

{#if schedules && schedules.length > 0}
  <div class="bg-white rounded-lg border border-gray-200 mt-6">
    <div class="px-4 py-3 border-b border-gray-200">
      <h3 class="text-sm font-semibold text-gray-900">
        {title} ({schedules.length})
      </h3>
    </div>
    <div class="divide-y divide-gray-200">
      {#each schedules as schedule, index}
        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 transition-colors">
          <div class="flex items-center gap-3">
            <div class="w-6 h-6 bg-primary-100 text-primary-600 rounded flex items-center justify-center text-xs font-medium">
              <IconClock class="h-4 w-4" stroke={2} />
            </div>
            <div>
              <div class="text-sm text-gray-900">
                <code class="bg-gray-100 px-2 py-1 rounded text-xs font-mono">{schedule.cron}</code>
              </div>
            </div>
          </div>
          <span class="inline-flex px-2 py-1 text-xs font-medium rounded bg-gray-100 text-gray-700">
            {schedule.timezone}
          </span>
        </div>
      {/each}
    </div>
  </div>
{/if}
