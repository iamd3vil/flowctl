<script lang="ts">
  let {
    flowName,
    startTime,
    executionId,
    status,
    scheduledAt,
    triggerType,
    triggeredBy
  }: {
    flowName: string,
    startTime: string,
    executionId: string,
    status?: string,
    scheduledAt?: string,
    triggerType?: string,
    triggeredBy?: string
  } = $props();

  // Extract just the name from "Name <username>" format
  function extractName(triggeredBy: string): string {
    const match = triggeredBy.match(/^(.+?)\s*</);
    return match ? match[1].trim() : triggeredBy;
  }
</script>

<!-- Flow Info Card -->
<div class="bg-white rounded-lg border border-gray-300 p-6 mb-6">
  <div class="flex justify-between items-start">
    <div>
      <div class="flex items-center gap-3">
        <h1 class="text-2xl font-semibold text-gray-900">{flowName}</h1>
        {#if triggerType}
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {triggerType === 'manual' ? 'bg-primary-100 text-primary-900' : 'bg-success-100 text-success-900'}">
            {triggerType}
          </span>
        {/if}
      </div>
      <p class="text-gray-600 mt-1">Started at {startTime}</p>
      {#if triggeredBy}
        <p class="text-sm text-gray-500 mt-3">Triggered By</p>
        <p class="text-sm text-gray-900">{extractName(triggeredBy)}</p>
      {/if}
    </div>
    <div class="text-right">
      <p class="text-sm text-gray-500">Execution ID</p>
      <p class="font-mono text-sm text-gray-900">{executionId}</p>
      {#if scheduledAt}
        <p class="text-sm text-gray-500 mt-3">Scheduled At</p>
        <p class="text-sm text-gray-900">{scheduledAt}</p>
      {/if}
    </div>
  </div>
</div>