<script lang="ts">
  import { handleInlineError, showSuccess } from '$lib/utils/errorHandling';
  import { apiClient } from '$lib/apiClient';
  import ScheduleModal from './ScheduleModal.svelte';
  import ViewScheduleModal from './ViewScheduleModal.svelte';
  import DeleteModal from '$lib/components/shared/DeleteModal.svelte';
  import DropdownMenu from '$lib/components/shared/DropdownMenu.svelte';
  import { IconClock, IconPlus } from '@tabler/icons-svelte';
  import type { UserSchedule, FlowInput, ScheduleCreateReq, ScheduleUpdateReq } from '$lib/types';

  let {
    namespace,
    flowId,
    flowInputs,
    userSchedulable = false,
    user,
    schedules = [],
    onUpdate,
    canUpdateFlow = false
  }: {
    namespace: string;
    flowId: string;
    flowInputs: FlowInput[];
    userSchedulable: boolean;
    user: any;
    schedules?: UserSchedule[];
    onUpdate?: () => Promise<void>;
    canUpdateFlow?: boolean;
  } = $props();

  let showModal = $state(false);
  let showDelete = $state(false);
  let showViewModal = $state(false);
  let editSchedule = $state<UserSchedule | null>(null);
  let deleteSchedule = $state<UserSchedule | null>(null);
  let viewSchedule = $state<UserSchedule | null>(null);

  // Check if flow has file inputs - user schedules not allowed for flows with file inputs
  let hasFileInputs = $derived(flowInputs.some(input => input.type === 'file'));
  let canCreateSchedule = $derived(userSchedulable && !hasFileInputs);

  function canEdit(schedule: UserSchedule): boolean {
    return canUpdateFlow || (schedule.is_user_created && schedule.created_by === user.id);
  }

  async function handleSave(data: ScheduleCreateReq | ScheduleUpdateReq) {
    if (editSchedule) {
      await apiClient.flows.schedules.update(namespace, flowId, editSchedule.uuid, data);
    } else {
      await apiClient.flows.schedules.create(namespace, flowId, data as ScheduleCreateReq);
    }
    if (onUpdate) await onUpdate();
  }

  async function handleDelete() {
    if (!deleteSchedule) return;
    await apiClient.flows.schedules.delete(namespace, flowId, deleteSchedule.uuid);
    showSuccess('Schedule Deleted', 'Schedule deleted successfully');
    showDelete = false;
    deleteSchedule = null;
    if (onUpdate) await onUpdate();
  }

  async function toggleActive(schedule: UserSchedule) {
    try {
      await apiClient.flows.schedules.update(namespace, flowId, schedule.uuid, {
        cron: schedule.cron,
        timezone: schedule.timezone,
        inputs: schedule.inputs,
        is_active: !schedule.is_active
      });
      showSuccess('Updated', 'Schedule status updated');
      if (onUpdate) await onUpdate();
    } catch (error) {
      handleInlineError(error, 'Failed to update schedule');
    }
  }

  function getMenuItems(schedule: UserSchedule) {
    return [
      {
        label: 'View',
        onClick: () => { viewSchedule = schedule; showViewModal = true; }
      },
      {
        label: 'Edit',
        onClick: () => { editSchedule = schedule; showModal = true; }
      },
      {
        label: schedule.is_active ? 'Deactivate' : 'Activate',
        onClick: () => toggleActive(schedule)
      },
      {
        label: 'Delete',
        onClick: () => { deleteSchedule = schedule; showDelete = true; },
        variant: 'danger' as const
      }
    ];
  }
</script>

<div class="bg-card rounded-lg border border-border">
  <div class="px-4 py-4 border-b border-border flex items-center justify-between">
    <div>
      <h3 class="text-sm font-semibold text-foreground">Schedules</h3>
      <p class="text-xs text-muted-foreground mt-0.5">{schedules.length} {schedules.length === 1 ? 'schedule' : 'schedules'}</p>
    </div>
    {#if canCreateSchedule}
      <button
        type="button"
        onclick={() => { editSchedule = null; showModal = true; }}
        class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-white bg-primary-500 rounded-md hover:bg-primary-600 cursor-pointer"
      >
        <IconPlus class="w-4 h-4" />
        Add
      </button>
    {/if}
  </div>

  {#if schedules.length === 0}
    <div class="flex flex-col items-center py-12">
      <IconClock class="w-12 h-12 text-muted-foreground mb-3" />
      <p class="text-sm text-muted-foreground">No schedules configured</p>
      {#if canCreateSchedule}
        <button
          type="button"
          onclick={() => { editSchedule = null; showModal = true; }}
          class="mt-3 text-sm text-primary-600 hover:text-primary-700 font-medium cursor-pointer"
        >
          Create your first schedule
        </button>
      {/if}
    </div>
  {:else}
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-border">
        <thead class="bg-muted">
          <tr>
            <th scope="col" class="px-4 py-2.5 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Cron</th>
            <th scope="col" class="px-4 py-2.5 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Timezone</th>
            <th scope="col" class="px-4 py-2.5 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Type</th>
            <th scope="col" class="px-4 py-2.5 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Status</th>
            <th scope="col" class="px-4 py-2.5 text-right text-xs font-medium text-muted-foreground uppercase tracking-wider w-20">Actions</th>
          </tr>
        </thead>
        <tbody class="bg-card divide-y divide-border">
          {#each schedules as schedule}
            <tr class="hover:bg-muted transition-colors">
              <td class="px-4 py-3 whitespace-nowrap">
                <code class="bg-subtle px-2 py-0.5 rounded text-xs font-mono text-foreground">{schedule.cron}</code>
              </td>
              <td class="px-4 py-3 whitespace-nowrap text-sm text-foreground">{schedule.timezone}</td>
              <td class="px-4 py-3 whitespace-nowrap">
                <span class="inline-flex px-2 py-0.5 text-xs font-medium rounded {schedule.is_user_created ? 'bg-green-100 text-green-800' : 'bg-blue-100 text-blue-800'}">
                  {schedule.is_user_created ? 'User' : 'System'}
                </span>
              </td>
              <td class="px-4 py-3 whitespace-nowrap">
                {#if schedule.is_user_created}
                  <span class="inline-flex px-2 py-0.5 text-xs font-medium rounded {schedule.is_active ? 'bg-success-100 text-success-800' : 'bg-subtle text-foreground'}">
                    {schedule.is_active ? 'Active' : 'Inactive'}
                  </span>
                {:else}
                  <span class="text-sm text-muted-foreground">-</span>
                {/if}
              </td>
              <td class="px-4 py-3 whitespace-nowrap text-right relative">
                {#if canCreateSchedule && canEdit(schedule)}
                  <div class="inline-flex justify-end">
                    <DropdownMenu items={getMenuItems(schedule)} />
                  </div>
                {:else}
                  <span class="text-muted-foreground">-</span>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

{#if showModal}
  <ScheduleModal
    isEditMode={!!editSchedule}
    schedule={editSchedule}
    {flowInputs}
    {namespace}
    {flowId}
    onSave={handleSave}
    onClose={() => { showModal = false; editSchedule = null; }}
  />
{/if}

{#if showDelete && deleteSchedule}
  <DeleteModal
    title="Delete Schedule"
    itemName={`${deleteSchedule.cron} (${deleteSchedule.timezone})`}
    onConfirm={handleDelete}
    onClose={() => { showDelete = false; deleteSchedule = null; }}
  />
{/if}

{#if showViewModal && viewSchedule}
  <ViewScheduleModal
    schedule={viewSchedule}
    onClose={() => { showViewModal = false; viewSchedule = null; }}
  />
{/if}
