<script lang="ts">
  import type { ScheduledExecution, UserSchedule } from '$lib/types';
  import { getNextCronRun } from '$lib/utils/cronParser';

  interface UpcomingRun {
    type: 'cron' | 'scheduled';
    label: string;
    scheduledAt: Date;
    execId?: string;
  }

  let {
    schedules = [],
    cronSchedules = [],
    namespace,
    flowId,
    title = 'Upcoming Scheduled Runs'
  }: {
    schedules: ScheduledExecution[];
    cronSchedules?: UserSchedule[];
    namespace: string;
    flowId: string;
    title?: string;
  } = $props();

  // Compute combined list of upcoming runs
  let upcomingRuns = $derived.by(() => {
    const runs: UpcomingRun[] = [];

    // Add cron-based runs
    for (const cron of cronSchedules) {
      const nextRun = getNextCronRun(cron.cron, cron.timezone);
      if (nextRun) {
        runs.push({
          type: 'cron',
          label: cron.cron,
          scheduledAt: nextRun
        });
      }
    }

    // Add manually scheduled runs
    for (const schedule of schedules) {
      runs.push({
        type: 'scheduled',
        label: 'Scheduled',
        scheduledAt: new Date(schedule.scheduled_at),
        execId: schedule.exec_id
      });
    }

    // Sort by scheduled time ascending
    return runs.sort((a, b) => a.scheduledAt.getTime() - b.scheduledAt.getTime());
  });

  function formatScheduledTime(date: Date): string {
    return date.toLocaleString();
  }
</script>

{#if upcomingRuns.length > 0}
  <div class="bg-white rounded-lg border border-gray-200">
    <div class="px-4 py-4 border-b border-gray-200">
      <h3 class="text-sm font-semibold text-gray-900">{title}</h3>
      <p class="text-xs text-gray-500 mt-0.5">{upcomingRuns.length} {upcomingRuns.length === 1 ? 'run' : 'runs'} scheduled</p>
    </div>
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-4 py-2.5 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Type
            </th>
            <th scope="col" class="px-4 py-2.5 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Scheduled Time
            </th>
            <th scope="col" class="px-4 py-2.5 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Exec ID
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          {#each upcomingRuns as run}
            <tr class="hover:bg-gray-50 transition-colors">
              <td class="px-4 py-3 whitespace-nowrap">
                {#if run.type === 'cron'}
                  <code class="text-xs font-mono bg-gray-100 px-2 py-0.5 rounded text-gray-700">{run.label}</code>
                {:else}
                  <span class="text-sm text-gray-700">{run.label}</span>
                {/if}
              </td>
              <td class="px-4 py-3 whitespace-nowrap text-sm text-gray-700">
                {formatScheduledTime(run.scheduledAt)}
              </td>
              <td class="px-4 py-3 whitespace-nowrap">
                {#if run.execId}
                  <a
                    href="/view/{namespace}/results/{flowId}/{run.execId}"
                    class="text-sm font-mono text-primary-600 hover:text-primary-800 hover:underline"
                  >
                    {run.execId.substring(0, 8)}
                  </a>
                {:else}
                  <span class="text-sm text-gray-400">-</span>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
{/if}
