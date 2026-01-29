<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import { apiClient } from "$lib/apiClient";
    import { handleInlineError } from "$lib/utils/errorHandling";
    import { currentUser } from "$lib/stores/auth";
    import { selectedNamespace } from "$lib/stores/namespace";
    import type { Namespace } from "$lib/types";
    import { DEFAULT_PAGE_SIZE } from "$lib/constants";
    import { setContext } from "svelte";
    import {
        permissionChecker,
        type ResourcePermissions,
    } from "$lib/utils/permissions";
    import {
        IconChevronDown,
        IconGridDots,
        IconServer,
        IconKey,
        IconLock,
        IconUsers,
        IconCircleCheck,
        IconClock,
        IconSettings,
        IconChevronsLeft,
        IconChevronsRight,
    } from "@tabler/icons-svelte";
    import UserDropdown from "./UserDropdown.svelte";
    import ThemeToggle from "./ThemeToggle.svelte";
    import Logo from "./Logo.svelte";
    import { APP_VERSION, APP_COMMIT } from "$lib/constants";

    let { namespace }: { namespace: string } = $props();

    let namespaceDropdownOpen = $state(false);
    let namespaces = $state<Namespace[]>([]);
    let searchQuery = $state("");
    let searchResults = $state<Namespace[]>([]);
    let currentNamespace = $state(page.params.namespace || namespace);
    let currentNamespaceId = $state<string>("");
    let searchLoading = $state(false);
    let isCollapsed = $state(false);
    let permissions = $state<{
        flows: ResourcePermissions;
        nodes: ResourcePermissions;
        credentials: ResourcePermissions;
        secrets: ResourcePermissions;
        members: ResourcePermissions;
        approvals: ResourcePermissions;
        history: ResourcePermissions;
    }>({
        flows: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        nodes: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        credentials: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        secrets: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        members: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        approvals: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
        history: {
            canCreate: false,
            canUpdate: false,
            canDelete: false,
            canRead: false,
        },
    });

    const isActiveLink = (section: string): boolean => {
        const currentPath = page.url.pathname;

        if (section === "flows") {
            return currentPath.includes("/flows");
        } else if (section === "nodes") {
            return currentPath.includes("/nodes");
        } else if (section === "credentials") {
            return currentPath.includes("/credentials");
        } else if (section === "secrets") {
            return currentPath.includes("/secrets");
        } else if (section === "members") {
            return currentPath.includes("/members");
        } else if (section === "approvals") {
            return currentPath.includes("/approvals");
        } else if (section === "history") {
            return currentPath.includes("/history");
        } else if (section === "settings") {
            return currentPath.includes("/settings");
        }

        return false;
    };

    const fetchNamespaces = async (filter = "") => {
        try {
            searchLoading = true;
            const data = await apiClient.namespaces.list({
                count_per_page: DEFAULT_PAGE_SIZE,
                filter: filter,
            });
            const results = data.namespaces || [];

            if (filter) {
                searchResults = results;
            } else {
                namespaces = results;
                searchResults = results;

                // Set current namespace ID
                const currentNs = namespaces.find(
                    (ns) => ns.name === namespace,
                );
                if (currentNs) {
                    currentNamespaceId = currentNs.id;
                } else if (namespaces.length > 0) {
                    // If namespace not found, use first available namespace
                    currentNamespaceId = namespaces[0].id;
                }
            }
        } catch (error) {
            handleInlineError(error, "Unable to Load Namespaces");
            if (filter) {
                searchResults = [];
            } else {
                namespaces = [];
                searchResults = [];
            }
        } finally {
            searchLoading = false;
        }
    };

    const handleSearchInput = async () => {
        if (searchQuery.trim()) {
            await fetchNamespaces(searchQuery);
            namespaceDropdownOpen = true;
        } else {
            searchResults = namespaces;
            namespaceDropdownOpen = false;
        }
    };

    const handleSearchFocus = () => {
        if (!searchQuery.trim()) {
            searchResults = namespaces;
        }
        namespaceDropdownOpen = true;
    };

    const checkAllPermissions = async () => {
        if (!$currentUser || !currentNamespaceId) return;

        const resourceMappings = {
            flows: "flow",
            nodes: "node",
            credentials: "credential",
            secrets: "namespace_secret",
            members: "member",
            approvals: "approval",
            history: "execution",
        };

        try {
            const permissionPromises = Object.entries(resourceMappings).map(
                async ([frontendKey, backendResource]) => {
                    const perms = await permissionChecker(
                        $currentUser,
                        backendResource,
                        currentNamespaceId,
                        ["view"],
                    );
                    return { frontendKey, perms };
                },
            );

            const results = await Promise.all(permissionPromises);

            results.forEach(({ frontendKey, perms }) => {
                permissions[frontendKey as keyof typeof permissions] = perms;
            });
        } catch (error) {
            handleInlineError(error, "Unable to Check Sidebar Permissions");
        }
    };

    const selectNamespace = (ns: Namespace) => {
        namespaceDropdownOpen = false;
        searchQuery = "";

        // Don't navigate if already on the same namespace
        if (ns.name === namespace) {
            return;
        }

        // Save selected namespace to store
        selectedNamespace.set(ns.name);

        // Force a full page reload by using window.location
        window.location.href = `/view/${ns.name}/flows`;
    };

    const toggleCollapse = () => {
        isCollapsed = !isCollapsed;
        if (isCollapsed) {
            namespaceDropdownOpen = false;
            searchQuery = "";
        }
    };

    // Handle escape key and outside clicks
    function handleOutsideClick(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest(".namespace-dropdown")) {
            namespaceDropdownOpen = false;
            searchQuery = "";
            searchResults = namespaces;
        }
    }

    // Set initial context
    setContext("namespace", namespace);

    onMount(() => {
        fetchNamespaces();
        checkAllPermissions();
    });

    // Update currentNamespace when namespace prop changes
    $effect(() => {
        currentNamespace = page.params.namespace || namespace;
        // Also save to store whenever namespace changes
        if (currentNamespace) {
            selectedNamespace.set(currentNamespace);
        }
    });

    // Re-check permissions when currentUser or namespace changes
    $effect(() => {
        if ($currentUser && currentNamespaceId) {
            checkAllPermissions();
        }
    });
</script>

<svelte:window on:click={handleOutsideClick} />

<!-- Sidebar Navigation -->
<nav
    class="relative bg-card border-r border-border flex flex-col transition-all duration-300 ease-in-out {isCollapsed
        ? 'w-20'
        : 'w-60'}"
    aria-label="Main navigation"
>
    <!-- Logo -->
    <a href="/">
        <div class="px-6 py-6 flex flex-col items-center">
            {#if isCollapsed}
                <Logo height="h-6" iconOnly={true} />
            {:else}
                <Logo height="h-8" />
                <div class="text-xs text-muted-foreground mt-1">{APP_VERSION}-{APP_COMMIT}</div>
            {/if}
        </div>
    </a>

    <!-- Namespace Dropdown -->
    {#if !isCollapsed}
        <div class="px-4 mb-4 namespace-dropdown">
            <div class="relative">
                <label
                    for="namespace-search"
                    class="block text-xs font-medium text-muted-foreground mb-1 uppercase"
                    >Namespace</label
                >
                <div class="relative">
                    <input
                        type="text"
                        id="namespace-search"
                        bind:value={searchQuery}
                        oninput={handleSearchInput}
                        onfocus={handleSearchFocus}
                        placeholder={currentNamespace || "Search namespaces..."}
                        class="w-full px-3 py-2 text-sm font-medium text-foreground bg-card border border-input rounded-lg hover:border-input focus:border-primary-500 focus:ring-1 focus:ring-primary-500 transition-colors outline-none pr-8"
                        autocomplete="off"
                    />

                    {#if searchLoading}
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
                        <IconChevronDown
                            class="absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground transition-transform {namespaceDropdownOpen
                                ? 'rotate-180'
                                : ''}"
                            size={16}
                        />
                    {/if}
                </div>

                <!-- Dropdown Menu -->
                {#if namespaceDropdownOpen}
                    <div
                        class="absolute z-50 w-full mt-1 bg-card rounded-lg shadow-lg border border-border max-h-48 overflow-y-auto"
                        role="listbox"
                        aria-label="Namespace selection"
                    >
                        <div class="py-1">
                            {#each searchResults as ns (ns.id)}
                                <button
                                    type="button"
                                    role="option"
                                    aria-selected={ns.name === namespace}
                                    onclick={() => selectNamespace(ns)}
                                    class="w-full text-left px-3 py-2 text-sm text-foreground hover:bg-subtle transition-colors cursor-pointer"
                                    class:bg-primary-50={ns.name === namespace}
                                    class:text-primary-600={ns.name ===
                                        namespace}
                                >
                                    {ns.name}
                                </button>
                            {/each}
                            {#if searchResults.length === 0 && !searchLoading}
                                <div
                                    class="px-3 py-2 text-sm text-muted-foreground text-center"
                                >
                                    {searchQuery
                                        ? "No namespaces found"
                                        : "No namespaces available"}
                                </div>
                            {/if}
                            {#if searchLoading}
                                <div
                                    class="px-3 py-2 text-sm text-muted-foreground text-center"
                                >
                                    Searching...
                                </div>
                            {/if}
                        </div>
                    </div>
                {/if}
            </div>
        </div>
    {/if}

    <!-- Navigation -->
    <ul class="flex-1 px-4 space-y-1" role="list">
        {#if permissions.flows.canRead}
            <li>
                <a
                    href="/view/{namespace}/flows"
                    class="flex items-center text-sm font-medium rounded-lg transition-colors {isCollapsed
                        ? 'justify-center px-4 py-3'
                        : 'px-4 py-3'}"
                    class:bg-primary-50={isActiveLink("flows")}
                    class:text-primary-600={isActiveLink("flows")}
                    class:text-foreground={!isActiveLink("flows")}
                    class:hover:bg-subtle={!isActiveLink("flows")}
                    aria-current={isActiveLink("flows") ? "page" : undefined}
                    title={isCollapsed ? "Flows" : ""}
                >
                    <IconGridDots
                        class="text-xl flex-shrink-0 {isCollapsed
                            ? ''
                            : 'mr-3'}"
                        size={20}
                        aria-hidden="true"
                    />
                    {#if !isCollapsed}
                        Flows
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.nodes.canRead}
            <li>
                <a
                    href="/view/{namespace}/nodes"
                    class="flex items-center text-sm font-medium rounded-lg transition-colors {isCollapsed
                        ? 'justify-center px-4 py-3'
                        : 'px-4 py-3'}"
                    class:bg-primary-50={isActiveLink("nodes")}
                    class:text-primary-600={isActiveLink("nodes")}
                    class:text-foreground={!isActiveLink("nodes")}
                    class:hover:bg-subtle={!isActiveLink("nodes")}
                    aria-current={isActiveLink("nodes") ? "page" : undefined}
                    title={isCollapsed ? "Nodes" : ""}
                >
                    <IconServer
                        class="text-xl flex-shrink-0 {isCollapsed
                            ? ''
                            : 'mr-3'}"
                        size={20}
                        aria-hidden="true"
                    />
                    {#if !isCollapsed}
                        Nodes
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.credentials.canRead}
            <li>
                <a
                    href="/view/{namespace}/credentials"
                    class="flex items-center text-sm font-medium rounded-lg transition-colors {isCollapsed
                        ? 'justify-center px-4 py-3'
                        : 'px-4 py-3'}"
                    class:bg-primary-50={isActiveLink("credentials")}
                    class:text-primary-600={isActiveLink("credentials")}
                    class:text-foreground={!isActiveLink("credentials")}
                    class:hover:bg-subtle={!isActiveLink("credentials")}
                    aria-current={isActiveLink("credentials")
                        ? "page"
                        : undefined}
                    title={isCollapsed ? "Credentials" : ""}
                >
                    <IconKey
                        class="text-xl flex-shrink-0 {isCollapsed
                            ? ''
                            : 'mr-3'}"
                        size={20}
                        aria-hidden="true"
                    />
                    {#if !isCollapsed}
                        Credentials
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.secrets.canRead}
            <li>
                <a
                    href="/view/{namespace}/secrets"
                    class="flex items-center text-sm font-medium rounded-lg transition-colors {isCollapsed
                        ? 'justify-center px-4 py-3'
                        : 'px-4 py-3'}"
                    class:bg-primary-50={isActiveLink("secrets")}
                    class:text-primary-600={isActiveLink("secrets")}
                    class:text-foreground={!isActiveLink("secrets")}
                    class:hover:bg-subtle={!isActiveLink("secrets")}
                    aria-current={isActiveLink("secrets")
                        ? "page"
                        : undefined}
                    title={isCollapsed ? "Secrets" : ""}
                >
                    <IconLock
                        class="text-xl flex-shrink-0 {isCollapsed
                            ? ''
                            : 'mr-3'}"
                        size={20}
                        aria-hidden="true"
                    />
                    {#if !isCollapsed}
                        Secrets
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.members.canRead}
            <li>
                <a
                    href="/view/{namespace}/members"
                    class="flex items-center text-sm font-medium rounded-lg transition-colors {isCollapsed
                        ? 'justify-center px-4 py-3'
                        : 'px-4 py-3'}"
                    class:bg-primary-50={isActiveLink("members")}
                    class:text-primary-600={isActiveLink("members")}
                    class:text-foreground={!isActiveLink("members")}
                    class:hover:bg-subtle={!isActiveLink("members")}
                    aria-current={isActiveLink("members") ? "page" : undefined}
                    title={isCollapsed ? "Members" : ""}
                >
                    <IconUsers
                        class="text-xl flex-shrink-0 {isCollapsed
                            ? ''
                            : 'mr-3'}"
                        size={20}
                        aria-hidden="true"
                    />
                    {#if !isCollapsed}
                        Members
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.approvals.canRead}
            <li>
                <a
                    href="/view/{namespace}/approvals"
                    class="flex items-center text-sm font-medium rounded-lg transition-colors {isCollapsed
                        ? 'justify-center px-4 py-3'
                        : 'px-4 py-3'}"
                    class:bg-primary-50={isActiveLink("approvals")}
                    class:text-primary-600={isActiveLink("approvals")}
                    class:text-foreground={!isActiveLink("approvals")}
                    class:hover:bg-subtle={!isActiveLink("approvals")}
                    aria-current={isActiveLink("approvals")
                        ? "page"
                        : undefined}
                    title={isCollapsed ? "Approvals" : ""}
                >
                    <IconCircleCheck
                        class="text-xl flex-shrink-0 {isCollapsed
                            ? ''
                            : 'mr-3'}"
                        size={20}
                        aria-hidden="true"
                    />
                    {#if !isCollapsed}
                        Approvals
                    {/if}
                </a>
            </li>
        {/if}
        {#if permissions.history.canRead}
            <li>
                <a
                    href="/view/{namespace}/history"
                    class="flex items-center text-sm font-medium rounded-lg transition-colors {isCollapsed
                        ? 'justify-center px-4 py-3'
                        : 'px-4 py-3'}"
                    class:bg-primary-50={isActiveLink("history")}
                    class:text-primary-600={isActiveLink("history")}
                    class:text-foreground={!isActiveLink("history")}
                    class:hover:bg-subtle={!isActiveLink("history")}
                    aria-current={isActiveLink("history") ? "page" : undefined}
                    title={isCollapsed ? "History" : ""}
                >
                    <IconClock
                        class="text-xl flex-shrink-0 {isCollapsed
                            ? ''
                            : 'mr-3'}"
                        size={20}
                        aria-hidden="true"
                    />
                    {#if !isCollapsed}
                        History
                    {/if}
                </a>
            </li>
        {/if}
    </ul>

    <div class="px-4 py-2 border-t border-border">
        <div class="flex items-center {isCollapsed ? 'flex-col gap-2' : 'justify-between'}">
            {#if $currentUser && $currentUser.role === "superuser"}
                <a
                    href="/settings"
                    class="flex items-center justify-center p-2 rounded-lg transition-colors"
                    class:bg-primary-50={isActiveLink("settings")}
                    class:text-primary-600={isActiveLink("settings")}
                    class:text-muted-foreground={!isActiveLink("settings")}
                    class:hover:bg-subtle={!isActiveLink("settings")}
                    aria-current={isActiveLink("settings") ? "page" : undefined}
                    title="Settings"
                >
                    <IconSettings size={20} aria-hidden="true" />
                </a>
            {/if}

            <ThemeToggle collapsed={true} />

            <button
                type="button"
                onclick={toggleCollapse}
                class="flex items-center justify-center p-2 text-muted-foreground hover:bg-subtle rounded-lg transition-colors cursor-pointer"
                aria-label={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
                title={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
            >
                {#if isCollapsed}
                    <IconChevronsRight size={20} aria-hidden="true" />
                {:else}
                    <IconChevronsLeft size={20} aria-hidden="true" />
                {/if}
            </button>
        </div>
    </div>

    <!-- User Profile Section -->
    {#if $currentUser}
        <div class="px-4 py-4 border-t border-border">
            <UserDropdown {isCollapsed} />
        </div>
    {/if}
</nav>
