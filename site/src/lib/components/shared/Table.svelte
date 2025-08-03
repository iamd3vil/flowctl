<script lang="ts" generics="T">
  import type { TableColumn, TableAction } from '$lib/types';

  type Props = {
    columns: TableColumn<T>[];
    data: T[];
    onRowClick?: (row: T) => void;
    actions?: TableAction<T>[];
    loading?: boolean;
    emptyMessage?: string;
    emptyIcon?: string;
    title?: string;
    subtitle?: string;
  };

  let { 
    columns, 
    data, 
    onRowClick, 
    actions = [], 
    loading = false,
    emptyMessage = "No data available",
    emptyIcon,
    title,
    subtitle
  }: Props = $props();

  const getValue = (row: T, column: TableColumn<T>) => {
    const keys = column.key.split('.');
    let value = row as any;
    
    for (const key of keys) {
      if (value && typeof value === 'object') {
        value = value[key];
      } else {
        return undefined;
      }
    }
    
    return value;
  };

  const renderValue = (row: T, column: TableColumn<T>) => {
    const value = getValue(row, column);
    
    if (column.render) {
      return column.render(value, row);
    }
    
    return value ?? '';
  };

  const handleRowClick = (row: T) => {
    if (onRowClick) {
      onRowClick(row);
    }
  };

  const handleActionClick = (action: TableAction<T>, row: T, event: Event) => {
    event.stopPropagation();
    action.onClick(row, event);
  };
</script>

<div class="bg-white rounded-lg border border-gray-200 overflow-hidden shadow-sm">
  {#if title || subtitle}
    <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
      {#if title}
        <h3 class="text-lg font-semibold text-gray-900">{title}</h3>
      {/if}
      {#if subtitle}
        <p class="text-sm text-gray-600 mt-1">{subtitle}</p>
      {/if}
    </div>
  {/if}

  {#if loading}
    <div class="flex items-center justify-center h-64">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <span class="ml-3 text-gray-600">Loading...</span>
    </div>
  {:else if data.length === 0}
    <div class="flex flex-col items-center justify-center h-64 text-center">
      {#if emptyIcon}
        {@html emptyIcon}
      {:else}
        <svg class="w-16 h-16 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/>
        </svg>
      {/if}
      <h3 class="text-lg font-medium text-gray-900 mb-2">{emptyMessage}</h3>
    </div>
  {:else}
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
        <tr>
          {#each columns as column}
            <th 
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider {column.width ? column.width : ''}"
            >
              {column.header}
            </th>
          {/each}
          {#if actions.length > 0}
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-32">
              Actions
            </th>
          {/if}
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        {#each data as row}
          <tr 
            class="hover:bg-gray-50 {onRowClick ? 'cursor-pointer' : ''}" 
            onclick={() => handleRowClick(row)}
          >
            {#each columns as column}
              <td class="px-6 py-4 whitespace-nowrap {column.width ? column.width : ''}">
                {#if column.component}
                  {@const Component = column.component}
                  <Component {row} value={getValue(row, column)} />
                {:else}
                  {@html renderValue(row, column)}
                {/if}
              </td>
            {/each}
            {#if actions.length > 0}
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium w-32">
                {#each actions as action}
                  <button 
                    onclick={(e) => handleActionClick(action, row, e)}
                    class="{action.className || 'text-blue-600 hover:text-blue-800'} mr-3"
                  >
                    {action.label}
                  </button>
                {/each}
              </td>
            {/if}
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>