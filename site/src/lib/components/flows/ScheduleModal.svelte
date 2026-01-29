<script lang="ts">
  import { handleInlineError, showSuccess } from '$lib/utils/errorHandling';
  import { autofocus } from '$lib/utils/autofocus';
  import { getTimezones } from '$lib/utils/timezone';
  import { isValidCronExpression } from '$lib/utils';
  import FlowInputFields from '$lib/components/shared/FlowInputFields.svelte';
  import type { FlowInput, UserSchedule, ScheduleCreateReq, ScheduleUpdateReq } from '$lib/types';

  let {
    isEditMode = false,
    schedule = null,
    flowInputs,
    namespace,
    flowId,
    onSave,
    onClose
  }: {
    isEditMode?: boolean;
    schedule?: UserSchedule | null;
    flowInputs: FlowInput[];
    namespace: string;
    flowId: string;
    onSave: (data: ScheduleCreateReq | ScheduleUpdateReq) => Promise<void>;
    onClose: () => void;
  } = $props();

  // Initialize form data based on edit mode
  function getInitialFormData() {
    if (isEditMode && schedule) {
      return {
        cron: schedule.cron,
        timezone: schedule.timezone,
        is_active: schedule.is_active,
        inputs: { ...schedule.inputs }
      };
    }
    // Initialize inputs with defaults for create mode
    const initialInputs: Record<string, any> = {};
    flowInputs.forEach(input => {
      if (input.default) {
        initialInputs[input.name] = input.default;
      }
    });
    return {
      cron: '',
      timezone: 'UTC',
      is_active: true,
      inputs: initialInputs
    };
  }

  let formData = $state(getInitialFormData());
  let loading = $state(false);
  let cronError = $state('');
  const timezones = getTimezones();

  function validateCron() {
    if (!formData.cron.trim()) {
      cronError = 'Cron expression is required';
      return false;
    }
    if (!isValidCronExpression(formData.cron)) {
      cronError = 'Invalid cron expression';
      return false;
    }
    cronError = '';
    return true;
  }

  async function handleSubmit(event: Event) {
    event.preventDefault();
    if (!validateCron()) return;

    loading = true;
    try {
      await onSave(formData);
      showSuccess(
        isEditMode ? 'Schedule Updated' : 'Schedule Created',
        'Operation completed successfully'
      );
      onClose();
    } catch (error) {
      handleInlineError(error, 'Unable to Save Schedule');
    } finally {
      loading = false;
    }
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') onClose();
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="fixed inset-0 z-50 flex items-center justify-center bg-overlay" onclick={onClose}>
  <div class="bg-card rounded-lg shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto m-4" onclick={(e) => e.stopPropagation()}>
    <div class="px-6 py-4 border-b border-border">
      <h2 class="text-xl font-semibold text-foreground">
        {isEditMode ? 'Edit Schedule' : 'Create Schedule'}
      </h2>
    </div>

    <form onsubmit={handleSubmit} class="p-6 space-y-4">
      <div>
        <label class="block mb-1 text-sm font-medium text-foreground">Cron Expression *</label>
        <input
          type="text"
          bind:value={formData.cron}
          onblur={validateCron}
          class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 text-foreground bg-card {cronError ? 'border-danger-300 focus:ring-danger-500' : 'border-input focus:ring-primary-500'}"
          placeholder="0 2 * * *"
          required
          use:autofocus
        />
        {#if cronError}
          <p class="text-sm text-danger-600 mt-1">{cronError}</p>
        {/if}
        <p class="text-xs text-muted-foreground mt-1">
          Examples: <code class="bg-subtle px-1 rounded">0 2 * * *</code> (daily 2AM),
          <code class="bg-subtle px-1 rounded">0 */6 * * *</code> (every 6 hours)
        </p>
      </div>

      <div>
        <label class="block mb-1 text-sm font-medium text-foreground">Timezone *</label>
        <input
          type="text"
          list="tz-list"
          bind:value={formData.timezone}
          class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
          required
        />
        <datalist id="tz-list">
          {#each timezones as tz}
            <option value={tz.tzCode}>{tz.label}</option>
          {/each}
        </datalist>
      </div>

      {#if isEditMode}
        <div class="flex items-center justify-between py-2 border-t border-border">
          <span class="text-sm font-medium text-foreground">Active</span>
          <button
            type="button"
            onclick={() => formData.is_active = !formData.is_active}
            class="relative inline-flex h-6 w-11 rounded-full border-2 border-transparent transition-colors {formData.is_active ? 'bg-primary-500' : 'bg-input'}"
            role="switch"
            aria-checked={formData.is_active}
          >
            <span class="inline-block h-5 w-5 transform rounded-full bg-card shadow transition {formData.is_active ? 'translate-x-5' : 'translate-x-0'}"></span>
          </button>
        </div>
      {/if}

      {#if flowInputs.length > 0}
        <div class="pt-4 border-t border-border">
          <h3 class="text-sm font-semibold text-foreground mb-3">Flow Inputs</h3>
          <FlowInputFields inputs={flowInputs} bind:values={formData.inputs} />
        </div>
      {/if}

      <div class="flex justify-end gap-2 pt-4 border-t border-border">
        <button
          type="button"
          onclick={onClose}
          disabled={loading}
          class="px-4 py-2 text-sm font-medium text-foreground bg-subtle rounded-md hover:bg-subtle-hover disabled:opacity-50 cursor-pointer"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={loading}
          class="px-4 py-2 text-sm font-medium text-white bg-primary-500 rounded-md hover:bg-primary-600 disabled:opacity-50 cursor-pointer"
        >
          {loading ? 'Saving...' : (isEditMode ? 'Update' : 'Create')}
        </button>
      </div>
    </form>
  </div>
</div>
