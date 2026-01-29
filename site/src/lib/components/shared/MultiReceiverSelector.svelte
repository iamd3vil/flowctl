<script lang="ts">
    import { apiClient } from "$lib/apiClient";
    import { handleInlineError } from "$lib/utils/errorHandling";
    import type { User, Group } from "$lib/types";
    import { IconUsers, IconUser, IconX, IconMail } from "@tabler/icons-svelte";

    let {
        selectedReceivers = $bindable([]),
        disabled = false,
    }: {
        selectedReceivers: string[];
        disabled?: boolean;
    } = $props();

    let searchQuery = $state("");
    let searchType = $state<"user" | "group">("user");
    let searchResults = $state<(User | Group)[]>([]);
    let showDropdown = $state(false);
    let loading = $state(false);

    interface SelectedItem {
        type: "user" | "group";
        id: string;
        name: string;
        value: string; // The formatted string: email or group:name
    }

    function isValidEmail(email: string): boolean {
        if (email.length < 3 || email.length > 254) return false;
        return /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/.test(
            email,
        );
    }

    // Show "Add email" option when user types a valid email that's not already selected
    let showAddEmailOption = $derived(
        searchType === "user" &&
            isValidEmail(searchQuery) &&
            !selectedReceivers.includes(searchQuery),
    );

    // Derived state: parse receivers into items for display
    let selectedItems = $derived.by(() => {
        if (!selectedReceivers || selectedReceivers.length === 0) {
            return [];
        }
        return selectedReceivers.map((r) => {
            // Check if it's a group (has "group:" prefix) or a user (no prefix)
            if (r.startsWith("group:")) {
                const groupName = r.substring(6); // Remove "group:" prefix
                return {
                    type: "group" as const,
                    id: groupName,
                    name: groupName,
                    value: r,
                };
            } else {
                // User email (no prefix needed)
                return {
                    type: "user" as const,
                    id: r,
                    name: r,
                    value: r,
                };
            }
        });
    });

    async function loadSubjects() {
        loading = true;
        showDropdown = true;
        try {
            if (searchType === "user") {
                const response = await apiClient.users.list({
                    filter: searchQuery,
                    count_per_page: 20,
                });
                searchResults = response.users || [];
            } else {
                const response = await apiClient.groups.list({
                    filter: searchQuery,
                    count_per_page: 20,
                });
                searchResults = response.groups || [];
            }
        } catch (error) {
            handleInlineError(
                error,
                searchType === "user"
                    ? "Unable to Load Users"
                    : "Unable to Load Groups",
            );
            searchResults = [];
            showDropdown = false;
        } finally {
            loading = false;
        }
    }

    async function handleFocus() {
        if (searchResults.length === 0) {
            await loadSubjects();
        } else {
            showDropdown = true;
        }
    }

    async function handleTypeChange() {
        searchQuery = "";
        searchResults = [];
        await loadSubjects();
    }

    function selectSubject(subject: User | Group) {
        const isUser = "username" in subject;
        const name = isUser
            ? (subject as User).username
            : (subject as Group).name;
        // Users don't need a prefix, groups use "group:" prefix
        const value = isUser ? name : `group:${name}`;

        // Check if already selected
        if (selectedReceivers.includes(value)) {
            return;
        }

        // Update the bindable array
        selectedReceivers = [...selectedReceivers, value];

        searchQuery = "";
        showDropdown = false;
        searchResults = [];
    }

    function addCustomEmail() {
        if (
            !isValidEmail(searchQuery) ||
            selectedReceivers.includes(searchQuery)
        )
            return;
        selectedReceivers = [...selectedReceivers, searchQuery];
        searchQuery = "";
        showDropdown = false;
    }

    function removeReceiver(index: number) {
        selectedReceivers = selectedReceivers.filter((_, i) => i !== index);
    }

    // Close dropdown when clicking outside
    function handleOutsideClick(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest(".multi-receiver-selector")) {
            showDropdown = false;
        }
    }
</script>

<svelte:window onclick={handleOutsideClick} />

<div class="multi-receiver-selector">
    <!-- Search Input -->
    <div class="mb-2">
        <div class="flex gap-2 mb-2">
            <select
                bind:value={searchType}
                onchange={handleTypeChange}
                class="px-3 py-2 text-sm text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                {disabled}
            >
                <option value="user">User</option>
                <option value="group">Group</option>
            </select>

            <div class="relative flex-1">
                <input
                    type="text"
                    bind:value={searchQuery}
                    oninput={loadSubjects}
                    onfocus={handleFocus}
                    placeholder={searchType === "user"
                        ? "Search users or enter email..."
                        : "Search groups..."}
                    class="w-full px-3 py-2 text-sm text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent pr-10"
                    autocomplete="off"
                    {disabled}
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
                {:else}
                    <svg
                        class="w-5 h-5 absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                        />
                    </svg>
                {/if}

                <!-- Dropdown -->
                {#if showDropdown}
                    <div
                        class="absolute z-10 w-full mt-1 bg-card border border-input rounded-lg shadow-lg max-h-48 overflow-y-auto"
                    >
                        {#if showAddEmailOption}
                            <button
                                type="button"
                                class="w-full px-4 py-2 hover:bg-primary-50 cursor-pointer border-b border-border text-left bg-muted"
                                onclick={addCustomEmail}
                            >
                                <div class="flex items-center">
                                    <div
                                        class="w-8 h-8 rounded-lg flex items-center justify-center mr-3 bg-primary-100"
                                    >
                                        <IconMail
                                            class="w-4 h-4 text-primary-600"
                                        />
                                    </div>
                                    <div>
                                        <div
                                            class="text-sm font-medium text-foreground"
                                        >
                                            Add "{searchQuery}"
                                        </div>
                                        <div class="text-xs text-muted-foreground">
                                            Add external email
                                        </div>
                                    </div>
                                </div>
                            </button>
                        {/if}
                        {#if searchResults.length > 0}
                            {#each searchResults as subject}
                                <button
                                    type="button"
                                    class="w-full px-4 py-2 hover:bg-muted cursor-pointer border-b border-border last:border-b-0 text-left"
                                    onclick={() => selectSubject(subject)}
                                >
                                    <div class="flex items-center">
                                        <div
                                            class="w-8 h-8 rounded-lg flex items-center justify-center mr-3 bg-primary-50"
                                        >
                                            {#if searchType === "user"}
                                                <IconUser
                                                    class="w-4 h-4 text-primary-600"
                                                />
                                            {:else}
                                                <IconUsers
                                                    class="w-4 h-4 text-primary-600"
                                                />
                                            {/if}
                                        </div>
                                        <div>
                                            <div
                                                class="text-sm font-medium text-foreground"
                                            >
                                                {"name" in subject
                                                    ? subject.name
                                                    : subject.username}
                                            </div>
                                            <div class="text-xs text-muted-foreground">
                                                {subject.id}
                                            </div>
                                        </div>
                                    </div>
                                </button>
                            {/each}
                        {:else if !loading && !showAddEmailOption}
                            <div
                                class="px-4 py-3 text-sm text-muted-foreground text-center"
                            >
                                {searchType === "user"
                                    ? "No users found. Enter a valid email to add."
                                    : "No groups found"}
                            </div>
                        {/if}
                    </div>
                {/if}
            </div>
        </div>
    </div>

    <!-- Selected Receivers -->
    {#if selectedItems.length > 0}
        <div class="flex flex-wrap gap-2">
            {#each selectedItems as item, index (item.value)}
                <div
                    class="inline-flex items-center gap-1 px-3 py-1 bg-card border border-input rounded-md text-sm"
                >
                    <div class="w-4 h-4 flex items-center justify-center">
                        {#if item.type === "user"}
                            <IconUser class="w-3 h-3 text-muted-foreground" />
                        {:else}
                            <IconUsers class="w-3 h-3 text-muted-foreground" />
                        {/if}
                    </div>
                    <span class="text-foreground">{item.name}</span>
                    <button
                        type="button"
                        onclick={() => removeReceiver(index)}
                        class="ml-1 text-muted-foreground hover:text-foreground"
                        {disabled}
                    >
                        <IconX class="w-3 h-3" />
                    </button>
                </div>
            {/each}
        </div>
    {:else}
        <div class="text-center text-sm text-muted-foreground">
            No receivers selected
        </div>
    {/if}
</div>
