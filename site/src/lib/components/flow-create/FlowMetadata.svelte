<script lang="ts">
  let {
    metadata = $bindable()
  }: {
    metadata: {
      id: string;
      name: string;
      description: string;
      namespace: string;
    }
  } = $props();

  function sanitizeId(value: string) {
    return value.replace(/[^a-zA-Z0-9_]/g, '');
  }

  function updateId(value: string) {
    metadata.id = sanitizeId(value);
  }

  function updateName(value: string) {
    metadata.name = value;
  }

  function updateDescription(value: string) {
    metadata.description = value;  
  }
</script>

<!-- Flow Metadata Section -->
<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
  <h2 class="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
    <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
    Flow Information
  </h2>
  <div class="grid grid-cols-2 gap-4">
    <div>
      <label for="flow-id" class="block text-sm font-medium text-gray-700 mb-2">Flow ID *</label>
      <input 
        type="text" 
        id="flow-id" 
        value={metadata.id}
        oninput={(e) => updateId(e.currentTarget.value)}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        placeholder="unique_flow_id"
      />
      <p class="mt-1 text-xs text-gray-500">Alphanumeric and underscores only</p>
    </div>
    <div>
      <label for="flow-name" class="block text-sm font-medium text-gray-700 mb-2">Flow Name *</label>
      <input 
        type="text" 
        id="flow-name" 
        value={metadata.name}
        oninput={(e) => updateName(e.currentTarget.value)}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        placeholder="My Flow Name"
      />
    </div>
    <div class="col-span-2">
      <label for="flow-description" class="block text-sm font-medium text-gray-700 mb-2">Description</label>
      <textarea 
        id="flow-description" 
        value={metadata.description}
        oninput={(e) => updateDescription(e.currentTarget.value)}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none h-20"
        placeholder="Describe what this flow does..."
      ></textarea>
    </div>
  </div>
</div>