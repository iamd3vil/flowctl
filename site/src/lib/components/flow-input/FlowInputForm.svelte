<script lang="ts">
  import { goto } from '$app/navigation';
  import type { FlowInput } from '$lib/types';
  import { ApiError } from '$lib/apiClient';
  import { handleInlineError, showSuccess } from '$lib/utils/errorHandling';
  import { getTimezones } from '$lib/utils/timezone';
  import { DateTime } from 'luxon';
  import { IconClock, IconPlayerPlay } from '@tabler/icons-svelte';
  import FlowInputFields from '$lib/components/shared/FlowInputFields.svelte';

  let {
    inputs,
    namespace,
    flowId,
    executionInput = null,
    onScheduled
  }: {
    inputs: FlowInput[],
    namespace: string,
    flowId: string,
    executionInput?: Record<string, any> | null,
    onScheduled?: () => void
  } = $props();

  let loading = $state(false);
  let errors = $state<Record<string, string>>({});
  let scheduleEnabled = $state(false);
  let scheduledAt = $state('');
  let scheduledTimezone = $state(Intl.DateTimeFormat().resolvedOptions().timeZone);

  const timezones = getTimezones();

  const mergedInputs = $derived(
    inputs.map(input => {
      if (executionInput && executionInput[input.name] !== undefined) {
        return { ...input, default: String(executionInput[input.name]) };
      }
      return input;
    })
  );

  // Convert datetime-local string + IANA timezone to RFC3339
  const toRFC3339 = (localDateTime: string, timezone: string): string => {
    return DateTime.fromISO(localDateTime, { zone: timezone }).toISO() ?? localDateTime;
  };

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    loading = true;
    errors = {};

    const form = event.target as HTMLFormElement;
    const formData = new FormData(form);

    // Build URL with scheduled_at query param if scheduling is enabled
    let url = `/api/v1/${namespace}/trigger/${flowId}`;
    if (scheduleEnabled && scheduledAt) {
      const scheduledAtRFC3339 = toRFC3339(scheduledAt, scheduledTimezone);
      if (new Date(scheduledAtRFC3339) <= new Date()) {
        errors = { general: 'Scheduled time must be in the future' };
        loading = false;
        return;
      }
      url += `?scheduled_at=${encodeURIComponent(scheduledAtRFC3339)}`;
    }

    try {
      const response = await fetch(url, {
        method: 'POST',
        body: formData,
        credentials: 'include',
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));

        // Handle validation errors with field details - show inline
        if (errorData.details && errorData.details.field && errorData.details.error) {
          errors[errorData.details.field] = errorData.details.error;
        }
        // For non-validation errors, use common error handling
        else {
          const apiError = new ApiError(response.status, response.statusText, errorData);
          handleInlineError(apiError, 'Unable to Start Flow');
        }
      } else {
        const data = await response.json();
        // If scheduled, show success message and stay on page
        if (data.scheduled_at) {
          const scheduledDate = new Date(data.scheduled_at);
          showSuccess('Flow Scheduled', `Flow will run at ${scheduledDate.toLocaleString()}`);
          scheduleEnabled = false;
          scheduledAt = '';
          onScheduled?.();
        } else {
          // Immediate execution - redirect to results page
          goto(`/view/${namespace}/results/${flowId}/${data.exec_id}`);
        }
      }
    } catch (error) {
      handleInlineError(error, 'Unable to Start Flow');
    } finally {
      loading = false;
    }
  };

</script>

<div class="bg-card rounded-lg border border-border overflow-hidden">
  <div class="px-6 py-4 border-b border-border bg-muted">
    <h2 class="text-lg font-semibold text-foreground">Configuration Parameters</h2>
    <p class="text-sm text-muted-foreground mt-1">Configure the inputs for this flow execution</p>
  </div>

  <form onsubmit={submit} class="p-6 space-y-6">
    {#if errors.general}
      <div class="p-3 rounded-md bg-danger-50 border border-danger-200">
        <div class="text-sm text-danger-700">{errors.general}</div>
      </div>
    {/if}

    <FlowInputFields inputs={mergedInputs} {errors} useFormData={true} />

    <!-- Schedule option -->
    <div class="pt-4 border-t border-border">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <IconClock class="w-5 h-5 text-muted-foreground" />
          <span class="text-sm font-medium text-foreground">Schedule for later</span>
        </div>
        <button
          type="button"
          onclick={() => scheduleEnabled = !scheduleEnabled}
          class="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 {scheduleEnabled ? 'bg-primary-500' : 'bg-input'}"
          role="switch"
          aria-checked={scheduleEnabled}
        >
          <span
            class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-card shadow ring-0 transition duration-200 ease-in-out {scheduleEnabled ? 'translate-x-5' : 'translate-x-0'}"
          ></span>
        </button>
      </div>

      {#if scheduleEnabled}
        <div class="mt-4 space-y-3">
          <div>
            <label for="scheduled_at" class="block text-sm font-medium text-foreground mb-2">
              Run at
              <span class="text-red-500">*</span>
            </label>
            <input
              type="datetime-local"
              id="scheduled_at"
              bind:value={scheduledAt}
              required={scheduleEnabled}
              class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>
          <div>
            <label for="scheduled_timezone" class="block text-sm font-medium text-foreground mb-2">
              Timezone
            </label>
            <input
              type="text"
              id="scheduled_timezone"
              list="timezone-list"
              bind:value={scheduledTimezone}
              class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              placeholder="Search or select timezone..."
            />
            <datalist id="timezone-list">
              {#each timezones as tz}
                <option value={tz.tzCode}>{tz.label}</option>
              {/each}
            </datalist>
          </div>
          <p class="text-sm text-muted-foreground">The flow will be queued and executed at the specified time</p>
        </div>
      {/if}
    </div>

    <div class="flex gap-3 pt-6 border-t border-border">
      <button
        type="button"
        onclick={() => window.history.back()}
        class="flex-1 px-4 py-2 bg-card border border-input text-foreground rounded-md hover:bg-muted transition-colors cursor-pointer"
      >
        Cancel
      </button>
      <button
        type="submit"
        disabled={loading || (scheduleEnabled && !scheduledAt)}
        class="flex-1 inline-flex items-center justify-center gap-2 px-4 py-2 bg-primary-500 text-white rounded-md hover:bg-primary-600 transition-colors disabled:opacity-50 cursor-pointer"
      >
        {#if loading}
          <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {scheduleEnabled ? 'Scheduling...' : 'Starting Flow...'}
        {:else if scheduleEnabled}
          <IconClock class="w-5 h-5" />
          Schedule
        {:else}
          <IconPlayerPlay class="w-5 h-5" />
          Run Now
        {/if}
      </button>
    </div>
  </form>
</div>
