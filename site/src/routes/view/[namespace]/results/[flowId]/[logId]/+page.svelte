<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import Header from '$lib/components/shared/Header.svelte';
  import StatusBadge from '$lib/components/shared/StatusBadge.svelte';
  import Alert from '$lib/components/flow-status/Alert.svelte';
  import Tabs from '$lib/components/shared/Tabs.svelte';
  import PipelineProgress from '$lib/components/flow-status/PipelineProgress.svelte';
  import LogsView from '$lib/components/flow-status/LogsView.svelte';
  import type { PageData } from './$types';
  import type { FlowMetaResp, ExecutionSummary } from '$lib/types';

  let { data }: { 
    data: {
      namespace: string;
      flowId: string; 
      logId: string;
      flowMeta?: FlowMetaResp;
      executionSummary?: ExecutionSummary;
      error?: string;
    }
  } = $props();

  // Flow execution state
  let status = $state<'running' | 'completed' | 'awaiting_approval' | 'errored'>('running');
  let currentActionIndex = $state(-1);
  let completedActions = $state<number[]>([]);
  let failedActionIndex = $state(-1);
  let logOutput = $state('');
  let results = $state<Record<string, any>>({});
  let showApproval = $state(false);
  let approvalID = $state<string | null>(null);
  let errorMessage = $state<string | null>(null);
  let activeTab = $state('logs');
  let startTime = $state('');
  let flowName = $state('');

  // WebSocket connection
  let ws: WebSocket | null = null;
  let hasReceivedMessages = $state(false);

  // Derived values
  let namespace = $derived(data.namespace);
  let flowId = $derived(data.flowId);
  let logId = $derived(data.logId);
  let actions = $derived(data.flowMeta?.actions || []);

  // Transform actions into pipeline steps
  let pipelineSteps = $derived(
    actions.map((action, index) => ({
      id: action.id,
      name: action.name || `Step ${index + 1}`,
      status: getActionStatus(index)
    }))
  );

  // Tab configuration
  let tabs = $derived([
    { id: 'logs', label: 'Real-time Logs' },
    { 
      id: 'output', 
      label: 'Execution Output',
      badge: Object.keys(results).length > 0 ? Object.keys(results).length : undefined
    }
  ]);

  const connectWebSocket = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const wsUrl = `${protocol}://${window.location.host}/api/v1/${namespace}/logs/${logId}`;
    
    ws = new WebSocket(wsUrl);

    ws.onmessage = (event) => {
      hasReceivedMessages = true;
      let msg = {};
      try { 
        msg = JSON.parse(event.data); 
      } catch (e) {
        console.error('Failed to parse WebSocket message:', e);
      }
      processMessage(msg);
    };

    ws.onclose = (event) => {
      console.log('WebSocket closed:', event);
      if (event.code === 1000) {
        status = 'completed';
        for (let i = 0; i < actions.length; i++) {
          if (!completedActions.includes(i)) {
            completedActions.push(i);
          }
        }
      } else if (event.reason) {
        errorMessage = event.reason;
        status = 'errored';
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  };

  const reconstructProgress = (currentActionId: string, executionStatus: string) => {
    let currentActionIndex = actions.findIndex(action => action.id === currentActionId);
    if (currentActionIndex === -1) return;

    for (let i = 0; i < currentActionIndex; i++) {
      if (!completedActions.includes(i)) {
        completedActions.push(i);
      }
    }

    if (executionStatus === 'completed') {
      for (let i = 0; i < actions.length; i++) {
        if (!completedActions.includes(i)) {
          completedActions.push(i);
        }
      }
      currentActionIndex = -1;
    } else if (executionStatus === 'errored') {
      failedActionIndex = currentActionIndex;
      currentActionIndex = -1;
    } else if (executionStatus === 'running' || executionStatus === 'pending') {
      currentActionIndex = currentActionIndex;
    } else if (executionStatus === 'pending_approval') {
      currentActionIndex = currentActionIndex;
      status = 'awaiting_approval';
    }
  };

  const processMessage = (msg: any) => {
    if (msg.action_id) {
      const actionIndex = actions.findIndex(a => a.id === msg.action_id);
      if (actionIndex !== -1) {
        if (currentActionIndex !== -1 && currentActionIndex !== actionIndex) {
          if (!completedActions.includes(currentActionIndex)) {
            completedActions.push(currentActionIndex);
          }
        }
        currentActionIndex = actionIndex;
      }
    }

    switch (msg.message_type) {
      case 'log':
        logOutput += (msg.value || '') + '\n';
        break;
      case 'result':
        results = { ...results, ...(msg.results || {}) };
        if (currentActionIndex !== -1 && !completedActions.includes(currentActionIndex)) {
          completedActions.push(currentActionIndex);
        }
        break;
      case 'error':
        errorMessage = msg.value || "An error occurred.";
        status = 'errored';
        if (currentActionIndex !== -1) {
          failedActionIndex = currentActionIndex;
        }
        break;
      case 'approval':
        approvalID = msg.value;
        showApproval = true;
        status = 'awaiting_approval';
        break;
      default:
        logOutput += (msg.value || '') + '\n';
    }
  };

  const goBack = () => {
    goto(`/view/${namespace}/flows`);
  };

  const getActionStatus = (index: number): 'pending' | 'running' | 'completed' | 'failed' => {
    if (index === failedActionIndex) return 'failed';
    if (completedActions.includes(index)) return 'completed';
    if (index === currentActionIndex && status === 'running') return 'running';
    return 'pending';
  };

  const dismissError = () => {
    errorMessage = null;
  };

  const dismissApproval = () => {
    showApproval = false;
  };

  // Initialize component
  onMount(() => {
    if (data.executionSummary) {
      const execStatus = data.executionSummary.status;
      status = execStatus === 'pending' ? 'running' : 
               execStatus === 'pending_approval' ? 'awaiting_approval' : 
               execStatus;
      startTime = new Date(data.executionSummary.started_at).toLocaleString();
      flowName = data.executionSummary.flow_name || data.flowMeta?.meta?.name || '';
      
      if (data.executionSummary.current_action_id) {
        reconstructProgress(data.executionSummary.current_action_id, data.executionSummary.status);
      }
    } else {
      flowName = data.flowMeta?.meta?.name || ''
      startTime = new Date().toLocaleString();
    }

    connectWebSocket();
  });

  onDestroy(() => {
    if (ws) {
      ws.close();
    }
  });
</script>

<svelte:head>
  <title>Flow Execution - {flowName || 'Loading...'}</title>
</svelte:head>

<div class="flex h-screen bg-gray-50">
  <main class="flex-1 flex flex-col overflow-hidden">
    <!-- Header -->
    <header class="bg-white border-b border-gray-200 px-6 py-4 flex justify-between items-center shadow-sm">
      <div class="flex items-center gap-5">
        <div class="flex items-center text-sm text-gray-500">
          <span>Flows</span>
          <span class="mx-2">/</span>
          <span class="text-gray-900">{flowName || 'Loading...'}</span>
          <span class="mx-2">/</span>
          <span class="text-gray-900">Execution Status</span>
        </div>
      </div>
      <div class="flex items-center gap-4">
        <div class="flex items-center gap-2">
          <span class="text-sm text-gray-500">Status:</span>
          <StatusBadge value={status} />
        </div>
        <button onclick={goBack} class="inline-flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
          </svg>
          Back to Flows
        </button>
      </div>
    </header>

    <!-- Page Content -->
    <div class="flex-1 overflow-y-auto p-6 bg-gray-50">
      <div class="max-w-7xl mx-auto">
        <!-- Flow Info Card -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
          <div class="flex justify-between items-start">
            <div>
              <h1 class="text-2xl font-semibold text-gray-900">{flowName || 'Loading...'}</h1>
              <p class="text-gray-600 mt-1">Started at {startTime}</p>
            </div>
            <div class="text-right">
              <p class="text-sm text-gray-500">Execution ID</p>
              <p class="font-mono text-sm text-gray-900">{logId}</p>
            </div>
          </div>
        </div>

        <!-- Pipeline Progress -->
        <div class="mb-6">
          <PipelineProgress 
            steps={pipelineSteps} 
            title="Pipeline Progress"
            orientation="horizontal"
            size="md"
          />
        </div>

        <!-- Tabs Content -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <Tabs {tabs} bind:activeTab />

          <!-- Tab Content -->
          <div class="p-6">
            {#if activeTab === 'logs'}
              <LogsView 
                bind:logs={logOutput}
                isRunning={status === 'running'}
                height="h-96"
                theme="dark"
                autoScroll={true}
                showCursor={true}
              />
            {:else if activeTab === 'output'}
              {#if Object.keys(results).length > 0}
                <div class="overflow-hidden rounded-lg border border-gray-200">
                  <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                      <tr>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Variable</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Value</th>
                      </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                      {#each Object.entries(results) as [key, value]}
                        <tr class="hover:bg-gray-50 transition-colors">
                          <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 font-mono">{key}</td>
                          <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            <code class="px-2 py-1 bg-green-50 text-green-700 rounded">{value}</code>
                          </td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                </div>
              {:else}
                <div class="text-center py-12 text-gray-500">
                  <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                  </svg>
                  <p class="mt-2">No output variables yet</p>
                </div>
              {/if}
            {/if}
          </div>
        </div>

        <!-- Alerts -->
        {#if showApproval}
          <div class="mt-6">
            <Alert
              type="warning"
              title="Approval Required"
              message="This flow requires manual approval to continue execution."
              dismissible={true}
              onDismiss={dismissApproval}
              actions={[
                {
                  label: "Review Request",
                  href: `/view/${namespace}/approvals/${approvalID}`,
                  primary: true
                }
              ]}
            />
          </div>
        {/if}

        {#if errorMessage}
          <div class="mt-6">
            <Alert
              type="error"
              title="Error"
              message={errorMessage}
              dismissible={true}
              onDismiss={dismissError}
            />
          </div>
        {/if}
      </div>
    </div>
  </main>
</div>