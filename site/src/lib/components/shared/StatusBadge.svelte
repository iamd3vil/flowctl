<script lang="ts">
  let { value }: { value: string } = $props();

  const capitalizeStatus = (status: string) => {
    if (!status) return '';
    return status.charAt(0).toUpperCase() + status.slice(1);
  };

  const getStatusClasses = (status: string) => {
    switch (status) {
      case 'completed':
      case 'approved':
        return 'bg-badge-success-bg text-badge-success-text';
      case 'pending':
      case 'cancelled':
        return 'bg-badge-warning-bg text-badge-warning-text';
      case 'running':
        return 'bg-badge-info-bg text-badge-info-text';
      case 'errored':
      case 'rejected':
        return 'bg-badge-danger-bg text-badge-danger-text';
      default:
        return 'bg-subtle text-foreground';
    }
  };

  const getDotClasses = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-success-500';
      case 'approved':
        return 'bg-success-500';
      case 'pending':
        return 'bg-warning-500';
      case 'running':
        return 'bg-primary-500';
      case 'cancelled':
        return 'bg-warning-500';
      case 'errored':
        return 'bg-danger-500';
      case 'rejected':
        return 'bg-danger-500';
      default:
        return 'bg-muted-foreground';
    }
  };
</script>

<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getStatusClasses(value)}">
  <span class="w-1.5 h-1.5 mr-1.5 rounded-full {getDotClasses(value)}"></span>
  <span>{capitalizeStatus(value)}</span>
</span>