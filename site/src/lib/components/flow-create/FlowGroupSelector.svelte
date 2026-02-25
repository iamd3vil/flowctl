<script lang="ts">
    import { apiClient } from "$lib/apiClient";
    import { handleInlineError } from "$lib/utils/errorHandling";
    import type { FlowGroupResp } from "$lib/types";
    import { IconFolder, IconX, IconPlus } from "@tabler/icons-svelte";

    let {
        namespace,
        value = $bindable(""),
        allowCreate = true,
    }: {
        namespace: string;
        value: string;
        allowCreate?: boolean;
    } = $props();

    let searchQuery = $state("");
    let groups = $state<FlowGroupResp[]>([]);
    let showDropdown = $state(false);
    let loading = $state(false);

    let filteredGroups = $derived(
        groups.filter((g) =>
            g.prefix.toLowerCase().includes(searchQuery.toLowerCase()),
        ),
    );

    let showCreateOption = $derived(
        allowCreate &&
            searchQuery.trim() !== "" &&
            !groups.some(
                (g) => g.prefix.toLowerCase() === searchQuery.toLowerCase(),
            ),
    );

    async function loadGroups() {
        loading = true;
        try {
            const result = await apiClient.flows.groups.list(namespace);
            groups = result.groups || [];
        } catch (error) {
            handleInlineError(error, "Unable to load flow groups");
            groups = [];
        } finally {
            loading = false;
        }
    }

    async function handleFocus() {
        if (groups.length === 0) {
            await loadGroups();
        }
        showDropdown = true;
    }

    function selectGroup(name: string) {
        value = name;
        searchQuery = "";
        showDropdown = false;
    }

    function clear() {
        value = "";
        searchQuery = "";
    }

    function handleOutsideClick(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest(".flow-group-selector")) {
            showDropdown = false;
        }
    }
</script>

<svelte:window onclick={handleOutsideClick} />

<div class="flow-group-selector">
    <label
        for="flow-group"
        class="block text-sm font-medium text-foreground mb-2"
    >
        Group
        <span class="text-sm text-muted-foreground font-normal">(optional)</span
        >
    </label>

    {#if value}
        <div
            class="flex items-center gap-2 px-3 py-2 bg-card border border-input rounded-md"
        >
            <IconFolder class="w-4 h-4 text-muted-foreground" />
            <span class="text-sm text-foreground flex-1">{value}</span>
            <button
                type="button"
                onclick={clear}
                class="text-muted-foreground hover:text-foreground cursor-pointer"
            >
                <IconX class="w-4 h-4" />
            </button>
        </div>
    {:else}
        <div class="relative">
            <input
                type="text"
                id="flow-group"
                bind:value={searchQuery}
                oninput={loadGroups}
                onfocus={handleFocus}
                placeholder={allowCreate ? "Search or create a group..." : "Search groups..."}
                class="w-full px-3 py-2 text-sm text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                autocomplete="off"
            />

            {#if loading}
                <div
                    class="absolute right-3 top-1/2 transform -translate-y-1/2"
                >
                    <svg
                        class="animate-spin h-4 w-4 text-muted-foreground"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                    >
                        <circle
                            class="opacity-25"
                            cx="12"
                            cy="12"
                            r="10"
                            stroke="currentColor"
                            stroke-width="4"
                        ></circle>
                        <path
                            class="opacity-75"
                            fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        ></path>
                    </svg>
                </div>
            {/if}

            {#if showDropdown}
                <div
                    class="absolute z-10 w-full mt-1 bg-card border border-input rounded-lg shadow-lg max-h-48 overflow-y-auto"
                >
                    {#if showCreateOption}
                        <button
                            type="button"
                            class="w-full px-4 py-2 hover:bg-muted cursor-pointer border-b border-border text-left"
                            onclick={() => selectGroup(searchQuery.trim())}
                        >
                            <div class="flex items-center gap-2">
                                <IconPlus
                                    class="w-4 h-4 text-primary-600"
                                />
                                <span class="text-sm text-foreground"
                                    >Create "{searchQuery.trim()}"</span
                                >
                            </div>
                        </button>
                    {/if}
                    {#if filteredGroups.length > 0}
                        {#each filteredGroups as group}
                            <button
                                type="button"
                                class="w-full px-4 py-2 hover:bg-muted cursor-pointer border-b border-border last:border-b-0 text-left"
                                onclick={() => selectGroup(group.prefix)}
                            >
                                <div class="flex items-center gap-2">
                                    <IconFolder
                                        class="w-4 h-4 text-muted-foreground"
                                    />
                                    <div>
                                        <div
                                            class="text-sm font-medium text-foreground"
                                        >
                                            {group.prefix}
                                        </div>
                                        {#if group.flow_count > 0}
                                            <div
                                                class="text-xs text-muted-foreground"
                                            >
                                                {group.flow_count} flow{group.flow_count !== 1 ? "s" : ""}
                                            </div>
                                        {/if}
                                    </div>
                                </div>
                            </button>
                        {/each}
                    {:else if !loading && !showCreateOption}
                        <div
                            class="px-4 py-3 text-sm text-muted-foreground text-center"
                        >
                            {allowCreate ? "No groups found. Type to create one." : "No groups found."}
                        </div>
                    {/if}
                </div>
            {/if}
        </div>
    {/if}

    <p class="text-xs text-muted-foreground mt-1">
        Assign this flow to a group for organization. Groups are created
        automatically.
    </p>
</div>
