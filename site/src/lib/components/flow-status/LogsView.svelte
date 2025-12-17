<script lang="ts">
    import { onMount, tick } from "svelte";
    import { formatTime } from "$lib/utils";

    type LogMessage = {
        action_id: string;
        message_type: string;
        node_id: string;
        value: string;
        timestamp: string;
    };

    type FormattedLog = {
        timestamp: string | null;
        nodeId: string;
        nodeColor: string;
        value: string;
    };

    type Props = {
        logs: string;
        logMessages?: LogMessage[];
        isRunning?: boolean;
        height?: string;
        showCursor?: boolean;
        autoScroll?: boolean;
        showLineNumbers?: boolean;
        theme?: "dark" | "light";
        fontSize?: "xs" | "sm" | "base";
        filterByActionId?: string;
    };

    let {
        logs = $bindable(),
        logMessages = [],
        isRunning = false,
        height = "h-96",
        showCursor = true,
        autoScroll = true,
        showLineNumbers = false,
        theme = "dark",
        fontSize = "sm",
        filterByActionId,
    }: Props = $props();

    let showTimestamp = $state(false);
    let scrollContainer: HTMLDivElement | undefined;

    const ITEM_HEIGHT = 20;
    const BUFFER_SIZE = 30;
    let scrollTop = $state(0);
    let viewportHeight = $state(500);

    let nodeColorMap = new Map<string, string>();
    const nodeColors = [
        "text-blue-400",
        "text-green-400",
        "text-yellow-400",
        "text-purple-400",
        "text-pink-400",
        "text-cyan-400",
        "text-orange-400",
        "text-teal-400",
        "text-indigo-400",
        "text-rose-400",
    ];

    const getNodeColor = (nodeId: string): string => {
        if (!nodeColorMap.has(nodeId)) {
            const colorIndex = nodeColorMap.size % nodeColors.length;
            nodeColorMap.set(nodeId, nodeColors[colorIndex]);
        }
        return nodeColorMap.get(nodeId)!;
    };

    const getContainerClasses = () => {
        const baseClasses = "rounded-lg p-4 font-mono min-h-128 max-h-128 overflow-y-auto";
        const themeClasses =
            theme === "dark"
                ? "bg-gray-900 text-gray-300"
                : "bg-gray-50 text-gray-900 border border-gray-200";
        const fontClasses = {
            xs: "text-xs",
            sm: "text-sm",
            base: "text-base",
        };
        return `${baseClasses} ${themeClasses} ${fontClasses[fontSize]}`;
    };

    const getCursorClasses = () => {
        const cursorColor = theme === "dark" ? "text-primary-400" : "text-primary-600";
        const blinkColor = theme === "dark" ? "text-gray-500" : "text-gray-400";
        return { cursor: cursorColor, blink: blinkColor };
    };

    const formatLogsWithLineNumbers = (logText: string) => {
        if (!showLineNumbers) return logText;
        return logText
            .split("\n")
            .map((line, index) => `${(index + 1).toString().padStart(4, " ")} | ${line}`)
            .join("\n");
    };

    const hasStructuredLogs = $derived(logMessages && logMessages.length > 0);

    const filteredMessages = $derived(
        filterByActionId
            ? logMessages.filter((msg) => msg.action_id === filterByActionId)
            : logMessages
    );

    const processedLogs = $derived.by(() => {
        const messages = filteredMessages;
        if (!messages || messages.length === 0) return [];

        nodeColorMap = new Map<string, string>();

        const result: FormattedLog[] = [];
        for (const msg of messages) {
            const lines = msg.value.split("\n").filter((line) => line.trim() !== "");
            for (const line of lines) {
                result.push({
                    timestamp: msg.timestamp ? formatTime(msg.timestamp) : null,
                    nodeId: msg.node_id,
                    nodeColor: getNodeColor(msg.node_id),
                    value: line,
                });
            }
        }
        return result;
    });

    const processedRawLogs = $derived(formatLogsWithLineNumbers(logs));
    const cursorClasses = $derived(getCursorClasses());

    const totalHeight = $derived(processedLogs.length * ITEM_HEIGHT);
    const startIndex = $derived(Math.max(0, Math.floor(scrollTop / ITEM_HEIGHT) - BUFFER_SIZE));
    const endIndex = $derived(Math.min(processedLogs.length, Math.ceil((scrollTop + viewportHeight) / ITEM_HEIGHT) + BUFFER_SIZE));
    const visibleLogs = $derived(processedLogs.slice(startIndex, endIndex));
    const offsetY = $derived(startIndex * ITEM_HEIGHT);

    const handleScroll = (e: Event) => {
        const target = e.target as HTMLDivElement;
        scrollTop = target.scrollTop;
    };

    let lastLogCount = 0;
    $effect(() => {
        const currentCount = processedLogs.length;
        if (autoScroll && currentCount > lastLogCount && scrollContainer) {
            lastLogCount = currentCount;
            tick().then(() => {
                if (scrollContainer) {
                    scrollContainer.scrollTop = scrollContainer.scrollHeight;
                }
            });
        }
    });

    onMount(() => {
        if (scrollContainer) {
            viewportHeight = scrollContainer.clientHeight;

            const observer = new ResizeObserver((entries) => {
                for (const entry of entries) {
                    viewportHeight = entry.contentRect.height;
                }
            });
            observer.observe(scrollContainer);

            return () => observer.disconnect();
        }
    });
</script>

<div class="flex flex-col space-y-3">
    {#if logMessages && logMessages.length > 0}
        <div class="flex gap-4 text-sm flex-shrink-0">
            <label class="flex items-center gap-2 cursor-pointer">
                <input
                    type="checkbox"
                    bind:checked={showTimestamp}
                    class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                />
                <span class="text-gray-900">Show Timestamp</span>
            </label>
        </div>
    {/if}

    <div class={getContainerClasses()} bind:this={scrollContainer} onscroll={handleScroll}>
        {#if filterByActionId && processedLogs.length === 0 && !isRunning}
            <div class="flex items-center justify-center h-full text-gray-500 text-sm">
                No logs available for this action
            </div>
        {:else if processedLogs.length > 0}
            <div style="height: {totalHeight}px; width: 100%; position: relative;">
                <div style="position: absolute; top: 0; left: 0; width: 100%; transform: translateY({offsetY}px);">
                    {#each visibleLogs as logMsg, i (startIndex + i)}
                        <div class="truncate" style="height: {ITEM_HEIGHT}px; line-height: {ITEM_HEIGHT}px;">
                            {#if showTimestamp && logMsg.timestamp}<span class="text-gray-500">[{logMsg.timestamp}]</span>{/if}{#if logMsg.nodeId}<span class="font-semibold {logMsg.nodeColor}">[{logMsg.nodeId}]</span>{/if}{logMsg.value}
                        </div>
                    {/each}
                </div>
            </div>
            {#if isRunning && showCursor}
                <div class="sticky bottom-0 bg-gray-900">
                    <span class={cursorClasses.cursor}>█</span>
                    <span class="animate-pulse {cursorClasses.blink}">_</span>
                </div>
            {/if}
        {:else if logs.length > 0}
            <div class="whitespace-pre-wrap break-words">
                {processedRawLogs}
                {#if isRunning && showCursor}
                    <div class="inline-block">
                        <span class={cursorClasses.cursor}>█</span>
                        <span class="animate-pulse {cursorClasses.blink}">_</span>
                    </div>
                {/if}
            </div>
        {:else}
            <div class="flex items-center justify-center h-full text-gray-500 text-sm">
                {#if isRunning}
                    Waiting for logs...
                {:else}
                    No logs available
                {/if}
            </div>
        {/if}
    </div>
</div>
