<script lang="ts">
    function parseInitial(json: string | undefined): Array<{ name: string; value: string }> {
        if (!json) return [{ name: "", value: "" }];
        try {
            const obj = JSON.parse(json);
            if (typeof obj === "object" && obj !== null) {
                const entries = Object.entries(obj as Record<string, unknown>);
                if (entries.length > 0) {
                    return entries.map(([k, v]) => ({ name: k, value: String(v) }));
                }
            }
        } catch {
            // ignore
        }
        return [{ name: "", value: "" }];
    }

    let {
        pairs = $bindable(parseInitial(initialValue)),
        initialValue,
        onchange,
        keyPlaceholder = "KEY",
        valuePlaceholder = "value",
    }: {
        pairs?: Array<{ name: string; value: string }>;
        initialValue?: string;
        onchange?: (json: string) => void;
        keyPlaceholder?: string;
        valuePlaceholder?: string;
    } = $props();

    function serialize() {
        const obj: Record<string, string> = {};
        for (const pair of pairs) {
            if (pair.name.trim() !== "") {
                obj[pair.name.trim()] = pair.value;
            }
        }
        return JSON.stringify(obj);
    }

    function notifyChange() {
        onchange?.(serialize());
    }

    function addPair() {
        pairs.push({ name: "", value: "" });
        notifyChange();
    }

    function removePair(i: number) {
        pairs.splice(i, 1);
        if (pairs.length === 0) {
            pairs.push({ name: "", value: "" });
        }
        notifyChange();
    }

    function handleKeyInput(e: Event, i: number) {
        pairs[i].name = (e.currentTarget as HTMLInputElement).value;
        notifyChange();
    }

    function handleValueInput(e: Event, i: number) {
        pairs[i].value = (e.currentTarget as HTMLInputElement).value;
        notifyChange();
    }
</script>

<div class="space-y-2">
    {#each pairs as pair, i}
        <div class="flex items-center gap-2">
            <input
                type="text"
                value={pair.name}
                oninput={(e) => handleKeyInput(e, i)}
                placeholder={keyPlaceholder}
                class="flex-1 px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm font-mono"
            />
            <span class="text-muted-foreground">=</span>
            <input
                type="text"
                value={pair.value}
                oninput={(e) => handleValueInput(e, i)}
                placeholder={valuePlaceholder}
                class="flex-1 px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm font-mono"
            />
            {#if pairs.length > 1 || pair.name !== "" || pair.value !== ""}
                <button
                    onclick={() => removePair(i)}
                    type="button"
                    class="text-muted-foreground hover:text-danger-600 cursor-pointer"
                >
                    <svg
                        class="w-4 h-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M6 18L18 6M6 6l12 12"
                        />
                    </svg>
                </button>
            {:else}
                <div class="w-4"></div>
            {/if}
        </div>
    {/each}
    <button
        onclick={addPair}
        type="button"
        class="text-xs text-primary-600 hover:text-primary-700 cursor-pointer"
    >
        + Add
    </button>
</div>
