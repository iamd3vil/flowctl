<script lang="ts">
  let {
    results
  }: {
    results: Record<string, any>
  } = $props();
</script>

{#if Object.keys(results).length > 0}
  <div class="overflow-hidden rounded-lg border border-border">
    <table class="min-w-full divide-y divide-border">
      <thead class="bg-muted">
        <tr>
          <th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Variable</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Value</th>
        </tr>
      </thead>
      <tbody class="bg-card divide-y divide-border">
        {#each Object.entries(results) as [key, value]}
          <tr class="hover:bg-muted transition-colors">
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-foreground font-mono">{key}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-foreground font-mono">{value}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}
