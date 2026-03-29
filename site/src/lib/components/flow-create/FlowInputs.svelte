<script lang="ts">
    let {
        inputs = $bindable(),
        addInput,
        disabled = false,
    }: {
        inputs: any[];
        addInput: () => void;
        disabled?: boolean;
    } = $props();

    // Keep max_file_size in sync with maxFileSizeMB for file inputs
    $effect(() => {
        for (const input of inputs) {
            if (input.type === 'file') {
                const mb = input.maxFileSizeMB;
                input.max_file_size = mb ? mb * 1024 * 1024 : undefined;
            }
        }
    });

    function removeInput(index: number) {
        inputs.splice(index, 1);
    }

    function sanitizeName(value: string) {
        return value.replace(/[^a-zA-Z0-9_]/g, "");
    }

    function onInputTypeChange(input: any) {
        if (input.type !== "select") {
            input.options = [];
            input.optionsText = "";
            input.useRemoteOptions = false;
            input.remote_options = undefined;
        }
        if (input.type === "file") {
            input.default = "";
        }
        if (input.type !== "file") {
            input.max_file_size = undefined;
            input.maxFileSizeMB = undefined;
        }
    }

    function updateOptions(input: any) {
        input.options = input.optionsText
            .split("\n")
            .filter((opt: string) => opt.trim());
    }

    function toggleRemoteOptions(input: any) {
        input.useRemoteOptions = !input.useRemoteOptions;
        if (input.useRemoteOptions) {
            input.remote_options = input.remote_options ?? { url: "", method: "GET", headers: {} };
            input.options = [];
            input.optionsText = "";
        } else {
            input.remote_options = undefined;
        }
    }

    function addRemoteHeader(input: any) {
        if (!input.remote_options.headers) {
            input.remote_options.headers = {};
        }
        // Use a temporary array to track header key/value pairs for the UI
        if (!input.remoteHeaders) {
            input.remoteHeaders = [];
        }
        input.remoteHeaders = [...input.remoteHeaders, { key: "", value: "" }];
    }

    function removeRemoteHeader(input: any, index: number) {
        input.remoteHeaders.splice(index, 1);
        input.remoteHeaders = [...input.remoteHeaders];
        syncHeaders(input);
    }

    function syncHeaders(input: any) {
        const headers: Record<string, string> = {};
        for (const h of (input.remoteHeaders ?? [])) {
            if (h.key.trim()) {
                headers[h.key.trim()] = h.value;
            }
        }
        input.remote_options.headers = headers;
    }
</script>

<!-- Flow Inputs Section -->
<div>
    <div class="flex items-center justify-between mb-6">
        <h3 class="text-base font-medium text-foreground">Flow Inputs</h3>
        {#if !disabled}
            <button
                onclick={addInput}
                class="px-4 py-2 text-sm font-medium bg-primary-500 text-white rounded-md hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 cursor-pointer"
            >
                + Add Input
            </button>
        {/if}
    </div>

    <div class="space-y-4">
        {#each inputs as input, index (index)}
            <div class="border border-border rounded-lg p-4 relative">
                {#if !disabled}
                <button
                    onclick={() => removeInput(index)}
                    class="absolute top-4 right-4 text-muted-foreground hover:text-danger-600 cursor-pointer"
                >
                    <svg
                        class="w-5 h-5"
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
                {/if}

                <div class="grid grid-cols-3 gap-4">
                    <div>
                        <label
                            class="block text-sm font-medium text-foreground mb-1"
                            >Input Name *</label
                        >
                        <input
                            type="text"
                            bind:value={input.name}
                            oninput={(e) =>
                                (input.name = sanitizeName(
                                    e.currentTarget.value,
                                ))}
                            class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                            placeholder="input_name"
                            required
                        />
                    </div>
                    <div>
                        <label
                            class="block text-sm font-medium text-foreground mb-1"
                            >Type *</label
                        >
                        <select
                            bind:value={input.type}
                            onchange={() => onInputTypeChange(input)}
                            class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                            required
                        >
                            <option value="string">String</option>
                            <option value="number">Number</option>
                            <option value="checkbox">Checkbox</option>
                            <option value="password">Password</option>
                            <option value="file">File</option>
                            <option value="datetime">DateTime</option>
                            <option value="select">Select</option>
                        </select>
                    </div>
                    <div>
                        <label
                            class="block text-sm font-medium text-foreground mb-1"
                            >Label</label
                        >
                        <input
                            type="text"
                            bind:value={input.label}
                            class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                            placeholder="Display Label"
                        />
                    </div>
                    <div class="col-span-2">
                        <label
                            class="block text-sm font-medium text-foreground mb-1"
                            >Description</label
                        >
                        <input
                            type="text"
                            bind:value={input.description}
                            class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                            placeholder="Help text for this input"
                        />
                    </div>
                    <div>
                        <label
                            class="block text-sm font-medium mb-1"
                            class:text-foreground={input.type !== "file"}
                            class:text-muted-foreground={input.type === "file"}
                            >Default Value</label
                        >
                        <input
                            type="text"
                            bind:value={input.default}
                            disabled={input.type === "file"}
                            class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm disabled:bg-subtle disabled:text-muted-foreground disabled:cursor-not-allowed"
                            placeholder={input.type === "file" ? "Not available for file inputs" : "Default value"}
                        />
                    </div>
                    <div class="col-span-2">
                        <label
                            class="block text-sm font-medium text-foreground mb-1"
                            >Validation</label
                        >
                        <input
                            type="text"
                            bind:value={input.validation}
                            class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm font-mono"
                            placeholder="len(input_name) > 3"
                        />
                    </div>
                    <div class="flex items-center">
                        <input
                            type="checkbox"
                            bind:checked={input.required}
                            class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-input rounded"
                        />
                        <label class="ml-2 block text-sm text-foreground"
                            >Required</label
                        >
                    </div>
                </div>

                {#if input.type === "select"}
                    <div class="mt-4 p-3 bg-muted rounded-md space-y-3">
                        <div class="flex items-center justify-between">
                            <span class="text-sm font-medium text-foreground">Options Source</span>
                            <div class="flex items-center gap-2">
                                <span class="text-sm text-muted-foreground">Static</span>
                                <button
                                    type="button"
                                    onclick={() => toggleRemoteOptions(input)}
                                    class="relative inline-flex h-5 w-9 cursor-pointer rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-1 {input.useRemoteOptions ? 'bg-primary-500' : 'bg-muted-foreground/30'}"
                                    role="switch"
                                    aria-checked={input.useRemoteOptions ?? false}
                                >
                                    <span
                                        class="inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform mt-0.5 {input.useRemoteOptions ? 'translate-x-4' : 'translate-x-0.5'}"
                                    ></span>
                                </button>
                                <span class="text-sm text-muted-foreground">Remote API</span>
                            </div>
                        </div>

                        {#if input.useRemoteOptions}
                            <div class="space-y-3">
                                <div class="grid grid-cols-4 gap-2">
                                    <div class="col-span-1">
                                        <label class="block text-xs font-medium text-foreground mb-1">Method</label>
                                        <select
                                            bind:value={input.remote_options.method}
                                            class="w-full px-2 py-1.5 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                                        >
                                            <option value="GET">GET</option>
                                            <option value="POST">POST</option>
                                        </select>
                                    </div>
                                    <div class="col-span-3">
                                        <label class="block text-xs font-medium text-foreground mb-1">URL *</label>
                                        <input
                                            type="url"
                                            bind:value={input.remote_options.url}
                                            class="w-full px-2 py-1.5 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm font-mono"
                                            placeholder="https://api.example.com/options"
                                            required
                                        />
                                    </div>
                                </div>

                                <div>
                                    <div class="flex items-center justify-between mb-1">
                                        <label class="block text-xs font-medium text-foreground">Headers</label>
                                        <button
                                            type="button"
                                            onclick={() => addRemoteHeader(input)}
                                            class="text-xs text-primary-600 hover:text-primary-700 font-medium cursor-pointer"
                                        >
                                            + Add Header
                                        </button>
                                    </div>
                                    {#if input.remoteHeaders && input.remoteHeaders.length > 0}
                                        <div class="space-y-1.5">
                                            {#each input.remoteHeaders as header, hi (hi)}
                                                <div class="flex gap-2 items-center">
                                                    <input
                                                        type="text"
                                                        bind:value={header.key}
                                                        oninput={() => syncHeaders(input)}
                                                        class="flex-1 px-2 py-1.5 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-xs font-mono"
                                                        placeholder="Header name"
                                                    />
                                                    <input
                                                        type="text"
                                                        bind:value={header.value}
                                                        oninput={() => syncHeaders(input)}
                                                        class="flex-1 px-2 py-1.5 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-xs font-mono"
                                                        placeholder="Value or {'{{ secrets.KEY }}'}"
                                                    />
                                                    <button
                                                        type="button"
                                                        onclick={() => removeRemoteHeader(input, hi)}
                                                        class="text-muted-foreground hover:text-danger-600 cursor-pointer flex-shrink-0"
                                                    >
                                                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                                        </svg>
                                                    </button>
                                                </div>
                                            {/each}
                                        </div>
                                    {:else}
                                        <p class="text-xs text-muted-foreground">No headers. Use <code class="font-mono bg-card px-1 rounded">{'{{ secrets.KEY }}'}</code> for interpolation.</p>
                                    {/if}
                                </div>
                            </div>
                        {:else}
                            <div>
                                <label
                                    class="block text-sm font-medium text-foreground mb-2"
                                    >Options (one per line)</label
                                >
                                <textarea
                                    bind:value={input.optionsText}
                                    oninput={() => updateOptions(input)}
                                    class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm font-mono h-20"
                                    placeholder="option1&#10;option2&#10;option3"
                                ></textarea>
                            </div>
                        {/if}
                    </div>
                {/if}

                {#if input.type === "file"}
                    <div class="mt-4 p-3 bg-muted rounded-md">
                        <label
                            class="block text-sm font-medium text-foreground mb-2"
                            >Max File Size (MB)</label
                        >
                        <input
                            type="number"
                            bind:value={input.maxFileSizeMB}
                            class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                            placeholder="Leave empty for default (100MB)"
                            min="1"
                        />
                        <p class="text-xs text-muted-foreground mt-1">Optional. Leave empty to use server default.</p>
                    </div>
                {/if}
            </div>
        {/each}

        {#if inputs.length === 0}
            <div class="text-center py-8 text-muted-foreground">
                <svg
                    class="mx-auto h-12 w-12 text-muted-foreground mb-3"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    />
                </svg>
                <p>No inputs defined yet</p>
                {#if !disabled}
                    <button
                        onclick={addInput}
                        class="mt-2 text-sm text-primary-600 hover:text-primary-700 font-medium cursor-pointer"
                    >
                        Add your first input
                    </button>
                {/if}
            </div>
        {/if}
    </div>
</div>
