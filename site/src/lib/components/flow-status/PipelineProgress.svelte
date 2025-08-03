<script lang="ts">
  type StepStatus = 'pending' | 'running' | 'completed' | 'failed';

  type Step = {
    id: string;
    name: string;
    status: StepStatus;
    description?: string;
  };

  type Props = {
    steps: Step[];
    title?: string;
    orientation?: 'horizontal' | 'vertical';
    size?: 'sm' | 'md' | 'lg';
    showConnectors?: boolean;
  };

  let {
    steps,
    title,
    orientation = 'horizontal',
    size = 'md',
    showConnectors = true
  }: Props = $props();

  const getStepClasses = (status: StepStatus) => {
    const baseClasses = 'rounded-lg border-2 transition-all duration-300';
    const sizeClasses = {
      sm: 'p-3 min-w-[150px]',
      md: 'p-4 min-w-[200px]',
      lg: 'p-6 min-w-[250px]'
    };

    switch (status) {
      case 'failed':
        return `${baseClasses} ${sizeClasses[size]} bg-red-50 border-red-500`;
      case 'completed':
        return `${baseClasses} ${sizeClasses[size]} bg-green-50 border-green-500`;
      case 'running':
        return `${baseClasses} ${sizeClasses[size]} bg-blue-50 border-blue-500`;
      default:
        return `${baseClasses} ${sizeClasses[size]} bg-gray-50 border-gray-300`;
    }
  };

  const getIconClasses = (status: StepStatus) => {
    const baseClasses = 'rounded-full flex items-center justify-center text-white text-xs font-bold';
    const sizeClasses = {
      sm: 'w-5 h-5',
      md: 'w-6 h-6',
      lg: 'w-8 h-8'
    };

    switch (status) {
      case 'failed':
        return `${baseClasses} ${sizeClasses[size]} bg-red-500`;
      case 'completed':
        return `${baseClasses} ${sizeClasses[size]} bg-green-500`;
      case 'running':
        return `${baseClasses} ${sizeClasses[size]} bg-blue-500 animate-pulse`;
      default:
        return `${baseClasses} ${sizeClasses[size]} bg-gray-400`;
    }
  };

  const getIcon = (status: StepStatus): string => {
    switch (status) {
      case 'failed':
        return '✗';
      case 'completed':
        return '✓';
      case 'running':
        return '●';
      default:
        return '-';
    }
  };

  const getStatusMessage = (status: StepStatus): string => {
    switch (status) {
      case 'failed':
        return 'Failed';
      case 'completed':
        return 'Completed';
      case 'running':
        return 'Processing...';
      default:
        return 'Pending';
    }
  };

  const getConnectorClasses = (currentStatus: StepStatus, isLast: boolean) => {
    if (isLast) return 'hidden';

    const baseClasses = orientation === 'horizontal' 
      ? 'absolute right-0 top-1/2 transform translate-x-4 -translate-y-1/2 h-0.5'
      : 'absolute bottom-0 left-1/2 transform -translate-x-1/2 translate-y-4 w-0.5';
    
    const sizeClasses = orientation === 'horizontal' ? 'w-4' : 'h-4';
    
    const colorClasses = currentStatus === 'completed' 
      ? 'bg-green-500' 
      : 'bg-gray-300';

    return `${baseClasses} ${sizeClasses} ${colorClasses}`;
  };

  const getContainerClasses = () => {
    const baseClasses = 'flex';
    const orientationClasses = orientation === 'horizontal' 
      ? 'gap-4 overflow-x-auto pb-2' 
      : 'flex-col gap-4 overflow-y-auto pr-2';

    return `${baseClasses} ${orientationClasses}`;
  };
</script>

<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
  {#if title}
    <h2 class="text-lg font-semibold text-gray-900 mb-4">{title}</h2>
  {/if}
  
  <div class={getContainerClasses()}>
    {#each steps as step, index}
      <div class="{getStepClasses(step.status)} relative">
        <div class="flex justify-between items-center mb-2">
          <span class="font-semibold text-gray-900 {size === 'sm' ? 'text-sm' : size === 'lg' ? 'text-lg' : 'text-base'}">
            {step.name}
          </span>
          <div class={getIconClasses(step.status)}>
            <span class="leading-none flex items-center justify-center">
              {getIcon(step.status)}
            </span>
          </div>
        </div>
        
        <p class="text-sm text-gray-600">
          {step.description || getStatusMessage(step.status)}
        </p>

        <!-- Progress connector -->
        {#if showConnectors}
          <div class={getConnectorClasses(step.status, index === steps.length - 1)}></div>
        {/if}
      </div>
    {/each}
  </div>
</div>