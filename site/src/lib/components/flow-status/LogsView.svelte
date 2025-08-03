<script lang="ts">
  import { onMount, onDestroy } from 'svelte';

  type Props = {
    logs: string;
    isRunning?: boolean;
    height?: string;
    showCursor?: boolean;
    autoScroll?: boolean;
    showLineNumbers?: boolean;
    theme?: 'dark' | 'light';
    fontSize?: 'xs' | 'sm' | 'base';
  };

  let {
    logs = $bindable(),
    isRunning = false,
    height = 'h-96',
    showCursor = true,
    autoScroll = true,
    showLineNumbers = false,
    theme = 'dark',
    fontSize = 'sm'
  }: Props = $props();

  let logContainer: HTMLDivElement;
  let previousLogsLength = 0;

  const getContainerClasses = () => {
    const baseClasses = 'rounded-lg p-4 overflow-y-auto font-mono';
    const themeClasses = theme === 'dark' 
      ? 'bg-gray-900 text-gray-300' 
      : 'bg-gray-50 text-gray-900 border border-gray-200';
    const fontClasses = {
      xs: 'text-xs',
      sm: 'text-sm',
      base: 'text-base'
    };

    return `${baseClasses} ${themeClasses} ${fontClasses[fontSize]} ${height}`;
  };

  const getCursorClasses = () => {
    const baseClasses = 'inline-block';
    const cursorColor = theme === 'dark' ? 'text-blue-400' : 'text-blue-600';
    const blinkColor = theme === 'dark' ? 'text-gray-500' : 'text-gray-400';
    
    return { cursor: cursorColor, blink: blinkColor };
  };

  const scrollToBottom = () => {
    if (logContainer && autoScroll) {
      setTimeout(() => {
        logContainer.scrollTop = logContainer.scrollHeight;
      }, 0);
    }
  };

  const formatLogsWithLineNumbers = (logText: string) => {
    if (!showLineNumbers) return logText;
    
    return logText
      .split('\n')
      .map((line, index) => `${(index + 1).toString().padStart(4, ' ')} | ${line}`)
      .join('\n');
  };

  const processedLogs = $derived(formatLogsWithLineNumbers(logs));
  const cursorClasses = $derived(getCursorClasses());

  // Auto-scroll when logs change
  $effect(() => {
    if (logs.length > previousLogsLength) {
      scrollToBottom();
      previousLogsLength = logs.length;
    }
  });

  // Scroll to bottom on mount
  onMount(() => {
    scrollToBottom();
  });
</script>

<div class={getContainerClasses()} bind:this={logContainer}>
  <div class="whitespace-pre-wrap break-words">
    {processedLogs}
    {#if isRunning && showCursor}
      <div class="inline-block">
        <span class={cursorClasses.cursor}>█</span>
        <span class="animate-pulse {cursorClasses.blink}">_</span>
      </div>
    {/if}
  </div>
</div>