<script lang="ts">
  import type { ApprovalResp } from '$lib/types';

  let { row, onClick }: { row: ApprovalResp; onClick?: (approval: ApprovalResp) => void } = $props();
</script>

<div class="flex items-center">
  <div class="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center mr-3">
    <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
    </svg>
  </div>
  <div>
    <a href="#" class="text-sm hover:underline font-medium text-gray-900 cursor-pointer hover:text-primary-600 transition-colors" onclick={(e) => { e.preventDefault(); onClick?.(row); }}>{row.action_id}</a>
    <div class="text-sm text-gray-500">{row.id}</div>
  </div>
</div>
