<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { apiClient } from '$lib/apiClient.js';
  import Header from '$lib/components/shared/Header.svelte';
  import FlowMetadata from '$lib/components/flow-create/FlowMetadata.svelte';
  import FlowInputs from '$lib/components/flow-create/FlowInputs.svelte';
  import FlowActions from '$lib/components/flow-create/FlowActions.svelte';
  import ValidationModal from '$lib/components/flow-create/ValidationModal.svelte';
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();
  const namespace = $page.params.namespace as string;

  // Flow state
  let flow = $state({
    metadata: {
      id: '',
      name: '',
      description: '',
      namespace: namespace
    },
    inputs: [] as any[],
    actions: [] as any[]
  });

  // Modal states
  let showValidation = $state(false);
  let validationResult = $state({
    success: false,
    errors: [] as string[]
  });

  // Loading and error states
  let saving = $state(false);
  const availableExecutors = data.availableExecutors;

  onMount(() => {
    // Add first action by default
    addAction();
  });

  function addInput() {
    flow.inputs.push({
      name: '',
      type: 'string',
      label: '',
      description: '',
      required: false,
      default: '',
      validation: '',
      options: [],
      optionsText: ''
    });
  }

  function addAction() {
    const tempId = Date.now() + Math.random();
    flow.actions.push({
      tempId: tempId,
      id: '',
      name: '',
      executor: '',
      with: {},
      on: '',
      variables: [],
      approval: false,
      artifacts: [],
      artifactsText: '',
      condition: '',
      collapsed: false
    });
  }

  function validateFlow() {
    const errors: string[] = [];

    // Validate metadata
    if (!flow.metadata.id) errors.push('Flow ID is required');
    if (!/^[a-zA-Z0-9_]+$/.test(flow.metadata.id)) errors.push('Flow ID must be alphanumeric with underscores only');
    if (!flow.metadata.name) errors.push('Flow Name is required');

    // Validate inputs
    const inputNames = new Set();
    flow.inputs.forEach((input, i) => {
      if (!input.name) errors.push(`Input #${i+1} is missing a name`);
      else if (inputNames.has(input.name)) errors.push(`Duplicate input name: ${input.name}`);
      else inputNames.add(input.name);

      if (!input.type) errors.push(`Input "${input.name || i+1}" is missing a type`);
    });

    // Validate actions
    const actionIds = new Set();
    flow.actions.forEach((action, i) => {
      if (!action.id) errors.push(`Action #${i+1} (${action.name || 'Untitled'}) is missing an ID`);
      else if (actionIds.has(action.id)) errors.push(`Duplicate action ID: ${action.id}`);
      else actionIds.add(action.id);

      if (!action.name) errors.push(`Action #${i+1} is missing a name`);
      if (!action.executor) errors.push(`Action "${action.name || i+1}" is missing an executor`);
    });

    validationResult = { success: errors.length === 0, errors: errors };
    showValidation = true;
    return errors.length === 0;
  }

  async function saveFlow() {
    if (!validateFlow()) return;
    
    saving = true;
    try {
      // Create a POST request to /flows endpoint
      const flowData = {
        metadata: flow.metadata,
        inputs: flow.inputs.filter(i => i.name),
        actions: flow.actions.filter(a => a.id)
      };

      await apiClient.flows.create(namespace, flowData);
      console.log('Flow saved successfully!');
      
    } catch (error) {
      console.error('Error saving flow:', error);
    } finally {
      saving = false;
    }
  }

  // Header actions - create a derived expression
  const headerActions = $derived([
    {
      label: 'Validate',
      onClick: validateFlow,
      variant: 'secondary' as const
    },
    {
      label: saving ? 'Saving...' : 'Save Flow',
      onClick: saveFlow,
      variant: 'primary' as const
    }
  ]);
</script>

<svelte:head>
  <title>Create Flow - {namespace} | Flowctl</title>
</svelte:head>

<div class="flex h-screen bg-gray-50">
  <!-- Main Content -->
  <div class="flex-1 flex flex-col overflow-hidden">
    <Header 
      breadcrumbs={[namespace, 'Flows', 'Create']} 
      actions={headerActions}
    />

    <!-- Page Content -->
    <div class="flex-1 overflow-y-auto">
      <div class="max-w-6xl mx-auto p-6">
        <FlowMetadata bind:metadata={flow.metadata} />
        <FlowInputs bind:inputs={flow.inputs} {addInput} />
        <FlowActions 
          bind:actions={flow.actions} 
          {addAction} 
          {availableExecutors}
        />
      </div>
    </div>
  </div>
</div>

{#if showValidation}
  <ValidationModal 
    bind:show={showValidation}
    {validationResult}
  />
{/if}