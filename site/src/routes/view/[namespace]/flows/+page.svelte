<script lang="ts">
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import { apiClient } from "$lib/apiClient";
    import Header from "$lib/components/shared/Header.svelte";
    import Table from "$lib/components/shared/Table.svelte";
    import Pagination from "$lib/components/shared/Pagination.svelte";
    import SearchInput from "$lib/components/shared/SearchInput.svelte";
    import PageHeader from "$lib/components/shared/PageHeader.svelte";
    import { handleInlineError, showSuccess } from "$lib/utils/errorHandling";
    import type { TableColumn, TableAction, FlowListItem, FlowsPaginateResponse, FlowGroupResp } from "$lib/types";
    import { FLOWS_PER_PAGE } from "$lib/constants";
    import {
        permissionChecker,
        type ResourcePermissions,
    } from "$lib/utils/permissions";
    import DeleteModal from "$lib/components/shared/DeleteModal.svelte";
    import GroupEditModal from "$lib/components/shared/GroupEditModal.svelte";

    interface FlowTableRow {
        _kind: 'group' | 'flow';
        name: string;
        description: string;
        slug: string;
        id: string;
        prefix: string;
        step_count: number;
        flow_count: number;
    }

    let { data } = $props();
    let searchValue = $state("");
    let flows = $state<FlowListItem[]>([]);
    let groups = $state<FlowGroupResp[]>([]);
    let pageCount = $state(0);
    let totalCount = $state(0);
    let currentPage = $state(data.currentPage);
    let loading = $state(true);
    let activeGroup = $state<string | null>(data.group || null);
    let permissions = $state<ResourcePermissions>({
        canCreate: false,
        canRead: false,
        canUpdate: false,
        canDelete: false,
    });
    let showDeleteModal = $state(false);
    let flowToDelete = $state<FlowTableRow | null>(null);
    let showEditGroupModal = $state(false);
    let groupToEdit = $state<FlowTableRow | null>(null);

    // Handle the async data from load function
    $effect(() => {
        let cancelled = false;
        activeGroup = data.group || null;

        if (data.group) {
            // Inside a group — load group flows
            data.groupFlowsPromise
                ?.then((result) => {
                    if (!cancelled) {
                        flows = result.flows || [];
                        groups = [];
                        pageCount = 0;
                        totalCount = flows.length;
                        loading = false;
                    }
                })
                .catch((err: Error) => {
                    if (!cancelled) {
                        flows = [];
                        handleInlineError(err, "Unable to Load Group Flows");
                        loading = false;
                    }
                });
        } else {
            // Root view — load groups and paginated flows
            data.groupsPromise
                ?.then((result) => {
                    if (!cancelled) {
                        groups = result.groups || [];
                    }
                })
                .catch(() => {
                    if (!cancelled) groups = [];
                });

            data.flowsPromise
                .then((result: FlowsPaginateResponse) => {
                    if (!cancelled) {
                        flows = result.flows;
                        pageCount = result.page_count;
                        totalCount = result.total_count;
                        loading = false;
                    }
                })
                .catch((err: Error) => {
                    if (!cancelled) {
                        flows = [];
                        pageCount = 0;
                        totalCount = 0;
                        handleInlineError(err, "Unable to Load Flows");
                        loading = false;
                    }
                });
        }

        return () => {
            cancelled = true;
        };
    });

    const goToEditFlow = (flowSlug: string) => {
        goto(`/view/${page.params.namespace}/flows/${flowSlug}/edit`);
    };

    const handleDeleteFlow = (flow: FlowTableRow) => {
        flowToDelete = flow;
        showDeleteModal = true;
    };

    const confirmDeleteFlow = async () => {
        if (!flowToDelete) return;

        try {
            await apiClient.flows.delete(
                page.params.namespace!,
                flowToDelete.slug,
            );
            showSuccess(
                "Flow Deleted",
                `Flow "${flowToDelete.name}" has been deleted successfully`,
            );
            if (activeGroup) {
                await loadGroupFlows(activeGroup);
            } else {
                await loadFlows(searchValue, currentPage);
            }
        } catch (err) {
            handleInlineError(err, "Unable to Delete Flow");
        } finally {
            showDeleteModal = false;
            flowToDelete = null;
        }
    };

    const cancelDelete = () => {
        showDeleteModal = false;
        flowToDelete = null;
    };

    const handleEditGroup = (row: FlowTableRow) => {
        groupToEdit = row;
        showEditGroupModal = true;
    };

    const saveGroupEdit = async (data: { description: string }) => {
        if (!groupToEdit) return;

        await apiClient.flows.groups.update(
            page.params.namespace!,
            groupToEdit.id,
            { name: groupToEdit.prefix, description: data.description },
        );
        showSuccess(
            "Group Updated",
            `Group "${groupToEdit.prefix}" has been updated successfully`,
        );
        // Reload groups
        const result = await apiClient.flows.groups.list(page.params.namespace!);
        groups = result.groups || [];
        showEditGroupModal = false;
        groupToEdit = null;
    };

    const cancelEditGroup = () => {
        showEditGroupModal = false;
        groupToEdit = null;
    };

    // Check permissions on mount
    const checkPermissions = async () => {
        permissions = await permissionChecker(
            data.user!,
            "flow",
            data.namespaceId,
            ["create", "update", "delete"],
            "_",
        );
    };

    const handleAdd = () => {
        goto(`/view/${page.params.namespace}/flows/create`);
    };

    checkPermissions();

    const navigateToGroup = (prefix: string) => {
        goto(`/view/${page.params.namespace}/flows?group=${encodeURIComponent(prefix)}`);
    };

    const navigateToRoot = () => {
        goto(`/view/${page.params.namespace}/flows`);
    };

    const loadGroupFlows = async (group: string) => {
        loading = true;
        try {
            const result = await apiClient.flows.groups.get(page.params.namespace!, group);
            flows = result.flows || [];
            totalCount = flows.length;
            pageCount = 0;
        } catch (err) {
            handleInlineError(err, "Unable to Load Group Flows");
        } finally {
            loading = false;
        }
    };

    const loadFlows = async (filter: string = "", pageNumber: number = 1) => {
        loading = true;

        try {
            const result = await apiClient.flows.list(page.params.namespace!, {
                filter,
                page: pageNumber,
                count_per_page: FLOWS_PER_PAGE,
            });

            flows = result.flows;
            pageCount = result.page_count;
            totalCount = result.total_count;
            currentPage = pageNumber;
        } catch (err) {
            handleInlineError(err, "Unable to Load Flows List");
        } finally {
            loading = false;
        }
    };

    const handleSearch = (query: string) => {
        searchValue = query;
        loadFlows(query, 1);
    };

    const goToPage = (pageNum: number) => {
        loadFlows(searchValue.trim(), pageNum);
    };

    const handlePageChange = (event: CustomEvent<{ page: number }>) => {
        goToPage(event.detail.page);
    };

    const tableData = $derived.by(() => {
        const rows: FlowTableRow[] = [];
        if (!activeGroup) {
            for (const g of groups) {
                rows.push({
                    _kind: 'group',
                    name: g.prefix,
                    description: g.description || '',
                    prefix: g.prefix,
                    flow_count: g.flow_count,
                    slug: '',
                    id: g.id,
                    step_count: 0,
                });
            }
        }
        const visibleFlows = activeGroup ? flows : flows.filter(f => !f.prefix);
        for (const f of visibleFlows) {
            rows.push({
                _kind: 'flow',
                name: f.name,
                description: f.description,
                slug: f.slug,
                id: f.id,
                prefix: f.prefix,
                step_count: f.step_count,
                flow_count: 0,
            });
        }
        return rows;
    });

    const columns: TableColumn<FlowTableRow>[] = [
        {
            key: "name",
            header: "Name",
            sortable: true,
            render: (value: string, row: FlowTableRow) => {
                if (row._kind === 'group') {
                    return `
                    <div class="flex items-center">
                        <div class="flex-shrink-0 h-8 w-8 bg-primary-100 rounded-lg flex items-center justify-center">
                            <svg class="w-4 h-4 text-primary-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>
                        </div>
                        <div class="ml-4">
                            <a href="/view/${page.params.namespace}/flows?group=${encodeURIComponent(row.prefix)}" class="text-sm hover:text-primary-600 hover:underline font-medium text-foreground">
                                ${value}
                            </a>
                            <div class="text-xs text-muted-foreground">${row.flow_count} flow${row.flow_count !== 1 ? 's' : ''}</div>
                        </div>
                    </div>`;
                }
                return `
                <div class="flex items-center">
                    <div class="flex-shrink-0 h-8 w-8 bg-primary-100 rounded-lg flex items-center justify-center">
                        <svg class="w-4 h-4 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                        </svg>
                    </div>
                    <div class="ml-4">
                        <a href="/view/${page.params.namespace}/flows/${row.slug}" class="text-sm hover:text-primary-600 hover:underline font-medium text-foreground">
                            ${value}
                        </a>
                    </div>
                </div>`;
            },
        },
        {
            key: "description",
            header: "Description",
            render: (value: string) => {
                if (!value) return '';
                return `<div class="text-sm text-muted-foreground max-w-xs truncate">${value}</div>`;
            },
        },
    ];

    const actions = $derived.by(() => {
        const actionsList: TableAction<FlowTableRow>[] = [];

        const isFlow = (row: FlowTableRow) => row._kind === 'flow';

        if (permissions.canUpdate) {
            actionsList.push({
                label: "Edit",
                onClick: (row: FlowTableRow) => {
                    if (row._kind === 'group') {
                        handleEditGroup(row);
                    } else {
                        goToEditFlow(row.slug);
                    }
                },
                className: "text-link",
            });
        }

        if (permissions.canDelete) {
            actionsList.push({
                label: "Delete",
                visible: isFlow,
                onClick: (row: FlowTableRow) => handleDeleteFlow(row),
                className: "text-danger-600",
            });
        }

        return actionsList;
    });

    // Breadcrumbs
    const breadcrumbs = $derived.by(() => {
        const crumbs = [
            { label: page.params.namespace! },
            { label: "Flows", url: `/view/${page.params.namespace}/flows` },
        ];
        if (activeGroup) {
            crumbs.push({ label: activeGroup });
        }
        return crumbs;
    });
</script>

<svelte:head>
    <title>{activeGroup ? `${activeGroup} - ` : ''}Flows - {page.params.namespace} - Flowctl</title>
</svelte:head>

<Header breadcrumbs={breadcrumbs}>
    {#snippet children()}
        {#if !activeGroup}
            <SearchInput
                bind:value={searchValue}
                placeholder="Search flows..."
                {loading}
                onSearch={handleSearch}
            />
        {/if}
    {/snippet}
</Header>

<!-- Page Content -->
<div class="p-12">
    <PageHeader
        title={activeGroup ? activeGroup : "Flows"}
        subtitle={activeGroup ? `Flows in the ${activeGroup} group` : "Manage and run your workflows"}
        actions={permissions.canCreate && !activeGroup
            ? [
                  {
                      label: "Add",
                      onClick: handleAdd,
                      variant: "primary",
                      icon: '<svg class="w-4 h-4 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>',
                  },
              ]
            : []}
    />

    <!-- Back to root button (inside group) -->
    {#if activeGroup}
        <button
            onclick={navigateToRoot}
            class="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground mb-4 cursor-pointer"
        >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
            </svg>
            Back to all flows
        </button>
    {/if}

    <!-- Flows Table -->
    <Table
        {columns}
        data={tableData}
        actions={actions}
        {loading}
        emptyMessage={activeGroup
            ? `No flows in the "${activeGroup}" group`
            : searchValue
                ? "Try adjusting your search"
                : "No flows are available in this namespace"}
        emptyIcon={`
        <svg class="mx-auto h-12 w-12 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
        </svg>
      `}
    />

    <!-- Pagination and Count (root view only) -->
    {#if !activeGroup && flows.length > 0}
        <div class="mt-6 flex items-center justify-between">
            <div class="text-sm text-foreground">
                Showing {flows.length} of {totalCount} flows
            </div>
            <Pagination
                {currentPage}
                totalPages={pageCount}
                {loading}
                on:page-change={handlePageChange}
            />
        </div>
    {/if}

    {#if activeGroup && flows.length > 0}
        <div class="mt-6 text-sm text-foreground">
            {flows.length} flow{flows.length !== 1 ? 's' : ''} in this group
        </div>
    {/if}
</div>

<!-- Delete Modal -->
{#if showDeleteModal && flowToDelete}
    <DeleteModal
        title="Delete Flow"
        itemName={flowToDelete.name}
        onConfirm={confirmDeleteFlow}
        onClose={cancelDelete}
    />
{/if}

<!-- Edit Group Modal -->
{#if showEditGroupModal && groupToEdit}
    <GroupEditModal
        groupName={groupToEdit.prefix}
        description={groupToEdit.description}
        onSave={saveGroupEdit}
        onClose={cancelEditGroup}
    />
{/if}
