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
<div class="bg-card rounded-lg border border-input p-6 mb-6">
  <div class="flex justify-between items-start">
    <div>
      <div class="flex items-center gap-3">
        <h1 class="text-2xl font-semibold text-foreground">{flowName}</h1>
        {#if triggerType}
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {triggerType === 'manual' ? 'bg-primary-100 text-primary-900' : 'bg-success-100 text-success-900'}">
            {triggerType}
          </span>
        {/if}
      </div>
      <p class="text-muted-foreground mt-1">Started at {startTime}</p>
      {#if triggeredBy}
        <p class="text-sm text-muted-foreground mt-3">Triggered By</p>
        <p class="text-sm text-foreground">{extractName(triggeredBy)}</p>
      {/if}
    </div>
    <div class="text-right">
      <p class="text-sm text-muted-foreground">Execution ID</p>
      <p class="font-mono text-sm text-foreground">{executionId}</p>
      {#if scheduledAt}
        <p class="text-sm text-muted-foreground mt-3">Scheduled At</p>
        <p class="text-sm text-foreground">{scheduledAt}</p>
      {/if}
    </div>
  </div>
</div>