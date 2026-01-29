<script lang="ts">
  let { value }: { value: string } = $props();
  
  function getRoleClasses(role: string) {
    switch (role) {
      case 'user':
        return 'bg-success-100 text-success-800 dark:bg-success-900/30 dark:text-success-300';
      case 'reviewer':
        return 'bg-warning-100 text-warning-800 dark:bg-warning-900/30 dark:text-warning-300';
      case 'admin':
        return 'bg-danger-100 text-danger-800 dark:bg-danger-900/30 dark:text-danger-300';
      default:
        return 'bg-subtle text-foreground';
    }
  }
  
  function capitalizeRole(role: string) {
    return role.charAt(0).toUpperCase() + role.slice(1);
  }
</script>

<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getRoleClasses(value)}">
  {capitalizeRole(value)}
</span>