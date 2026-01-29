<script lang="ts">
  let { value }: { value: string } = $props();
</script>

<div class="flex items-center">
  <div class="w-8 h-8 rounded-lg bg-subtle flex items-center justify-center mr-3 flex-shrink-0">
    <svg class="w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
    </svg>
  </div>
  <div class="min-w-0 flex-1">
    <div class="text-xs text-muted-foreground font-mono truncate">{value}</div>
  </div>
</div>