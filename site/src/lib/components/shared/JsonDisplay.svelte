<script lang="ts">
  interface Props {
    data: any;
    title?: string;
    expanded?: boolean;
  }

  let { data, title = 'JSON Data', expanded = false }: Props = $props();

  let isExpanded = $state(expanded);

  const toggleExpanded = () => {
    isExpanded = !isExpanded;
  };

  const formatJson = (obj: any): string => {
    if (obj === null || obj === undefined) {
      return '';
    }

    try {
      return JSON.stringify(obj, null, 2);
    } catch (error) {
      return String(obj);
    }
  };

  let jsonString = $derived(formatJson(data));
  let hasData = $derived(data !== null && data !== undefined && jsonString.trim() !== '');
</script>

{#if hasData}
  <div class="bg-card rounded-lg border border-input">
    <button
      class="w-full flex items-center justify-between px-6 py-3 hover:bg-muted transition-colors duration-200 cursor-pointer"
      onclick={toggleExpanded}
      type="button"
    >
      <span class="font-bold text-base text-foreground">{title}</span>
      <div class="flex items-center space-x-2">
        <span class="text-xs text-muted-foreground bg-subtle px-2 py-1 rounded">JSON</span>
        <svg
          class="w-4 h-4 text-muted-foreground transition-transform duration-200 {isExpanded ? 'transform rotate-180' : ''}"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
        </svg>
      </div>
    </button>

    {#if isExpanded}
      <div class="p-4">
        <pre class="bg-gray-900 dark:bg-gray-950 text-gray-100 p-4 rounded-md text-sm overflow-x-auto font-mono leading-relaxed whitespace-pre-wrap">{jsonString}</pre>
      </div>
    {/if}
  </div>
{/if}
