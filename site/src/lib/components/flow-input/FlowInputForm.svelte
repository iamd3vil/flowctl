<script lang="ts">
  import { goto } from '$app/navigation';
  import type { FlowInput } from '$lib/types';
  import { ApiError } from '$lib/apiClient';
  import { handleInlineError, showSuccess } from '$lib/utils/errorHandling';
  import { IconClock, IconPlayerPlay } from '@tabler/icons-svelte';

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

  const mergedInputs = $derived(
    inputs.map(input => {
      if (executionInput && executionInput[input.name] !== undefined) {
        return { ...input, default: String(executionInput[input.name]) };
      }
      return input;
    })
  );

  // Get minimum datetime (now + 1 minute) in local time format for datetime-local input
  const getMinDateTime = () => {
    const now = new Date();
    now.setMinutes(now.getMinutes() + 1);
    return new Date(now.getTime() - now.getTimezoneOffset() * 60000).toISOString().slice(0, 16);
  };

  // Convert local datetime to RFC3339 format
  const toRFC3339 = (localDateTime: string): string => {
    const date = new Date(localDateTime);
    return date.toISOString();
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
      const scheduledAtRFC3339 = toRFC3339(scheduledAt);
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

<div class="bg-white rounded-lg border border-gray-200 overflow-hidden">
  <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
    <h2 class="text-lg font-semibold text-gray-900">Configuration Parameters</h2>
    <p class="text-sm text-gray-600 mt-1">Configure the inputs for this flow execution</p>
  </div>

  <form onsubmit={submit} class="p-6 space-y-6">
    {#if errors.general}
      <div class="p-3 rounded-md bg-danger-50 border border-danger-200">
        <div class="text-sm text-danger-700">{errors.general}</div>
      </div>
    {/if}

    {#each mergedInputs as input (input.name)}
      <div>
        <label for={input.name} class="block text-sm font-medium text-gray-700 mb-2">
          {input.label || input.name}
          {#if input.required}
            <span class="text-red-500">*</span>
          {/if}
        </label>

        {#if errors[input.name]}
          <p class="text-sm text-danger-600 mb-2">{errors[input.name]}</p>
        {/if}

        {#if input.type === 'string' || input.type === 'number'}
          <input
            type={input.type === 'string' ? 'text' : 'number'}
            id={input.name}
            name={input.name}
            value={input.default || ''}
            placeholder={input.description || ''}
            required={input.required}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          />
        {:else if input.type === 'checkbox'}
          <div class="flex items-center">
            <input
              type="checkbox"
              id={input.name}
              name={input.name}
              value="true"
              checked={input.default === 'true'}
              class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
            />
          </div>
        {:else if input.type === 'select' && input.options}
          <select
            id={input.name}
            name={input.name}
            required={input.required}
            value={input.default || ''}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          >
            <option value="">Select an option</option>
            {#each input.options as option}
              <option value={option} selected={option === input.default}>{option}</option>
            {/each}
          </select>
          <!-- {:else if input.type === 'file'}
          <div class="flex items-center">
            <input
              type="file"
              id={input.name}
              name={input.name}
              required={input.required}
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100"
            />
          </div> -->
          {:else if input.type === 'datetime'}
          <div class="flex items-center">
            <input
              type="datetime-local"
              id={input.name}
              name={input.name}
              value={input.default || ''}
              required={input.required}
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100"
            />
          </div>
          {:else if input.type === 'password'}
          <div class="flex items-center">
            <input
              type="password"
              id={input.name}
              name={input.name}
              value={input.default || ''}
              required={input.required}
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100"
            />
          </div>
        {:else}
          <!-- Fallback for other input types -->
          <input
            type="text"
            id={input.name}
            name={input.name}
            value={input.default || ''}
            placeholder={input.description || ''}
            required={input.required}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          />
        {/if}

        {#if input.description}
          <p class="text-sm text-gray-500 mt-1">{input.description}</p>
        {/if}
      </div>
    {/each}

    <!-- Schedule option -->
    <div class="pt-4 border-t border-gray-200">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <IconClock class="w-5 h-5 text-gray-500" />
          <span class="text-sm font-medium text-gray-700">Schedule for later</span>
        </div>
        <button
          type="button"
          onclick={() => scheduleEnabled = !scheduleEnabled}
          class="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 {scheduleEnabled ? 'bg-primary-500' : 'bg-gray-200'}"
          role="switch"
          aria-checked={scheduleEnabled}
        >
          <span
            class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out {scheduleEnabled ? 'translate-x-5' : 'translate-x-0'}"
          ></span>
        </button>
      </div>

      {#if scheduleEnabled}
        <div class="mt-4">
          <label for="scheduled_at" class="block text-sm font-medium text-gray-700 mb-2">
            Run at
            <span class="text-red-500">*</span>
          </label>
          <input
            type="datetime-local"
            id="scheduled_at"
            bind:value={scheduledAt}
            min={getMinDateTime()}
            required={scheduleEnabled}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          />
          <p class="text-sm text-gray-500 mt-1">The flow will be queued and executed at the specified time</p>
        </div>
      {/if}
    </div>

    <div class="flex gap-3 pt-6 border-t border-gray-200">
      <button
        type="button"
        onclick={() => window.history.back()}
        class="flex-1 px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition-colors cursor-pointer"
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
