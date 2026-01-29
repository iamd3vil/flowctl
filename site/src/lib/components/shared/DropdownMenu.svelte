<script lang="ts">
  import { IconDotsVertical } from '@tabler/icons-svelte';

  interface MenuItem {
    label: string;
    onClick: () => void;
    variant?: 'default' | 'danger';
  }

  let {
    items,
    position = 'right'
  }: {
    items: MenuItem[];
    position?: 'left' | 'right';
  } = $props();

  let isOpen = $state(false);
  let buttonElement: HTMLButtonElement | undefined = $state();
  let dropdownPosition = $state({ top: 0, left: 0, right: 0 });

  function updateDropdownPosition() {
    if (!buttonElement) return;

    const rect = buttonElement.getBoundingClientRect();
    dropdownPosition = {
      top: rect.bottom + 4,
      left: position === 'left' ? rect.left : 0,
      right: position === 'right' ? window.innerWidth - rect.right : 0
    };
  }

  function handleOutsideClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.dropdown-menu-container')) {
      isOpen = false;
    }
  }

  function handleItemClick(item: MenuItem) {
    item.onClick();
    isOpen = false;
  }

  function toggleDropdown() {
    isOpen = !isOpen;
    if (isOpen) {
      updateDropdownPosition();
    }
  }
</script>

<svelte:window on:click={handleOutsideClick} on:scroll={updateDropdownPosition} on:resize={updateDropdownPosition} />

<div class="relative dropdown-menu-container">
  <button
    type="button"
    bind:this={buttonElement}
    onclick={toggleDropdown}
    class="p-1 hover:bg-subtle rounded cursor-pointer"
    aria-label="Actions menu"
  >
    <IconDotsVertical class="w-5 h-5 text-muted-foreground" />
  </button>

  {#if isOpen}
    <div
      class="fixed w-36 bg-card rounded-md shadow-lg border border-border z-50"
      style="top: {dropdownPosition.top}px; {position === 'right' ? `right: ${dropdownPosition.right}px` : `left: ${dropdownPosition.left}px`}"
      role="menu"
    >
      <div class="py-1 flex flex-col">
        {#each items as item}
          <button
            type="button"
            onclick={() => handleItemClick(item)}
            class="w-full text-left px-4 py-2 text-sm cursor-pointer block {item.variant === 'danger' ? 'text-danger-600 hover:bg-danger-50' : 'text-foreground hover:bg-subtle'}"
            role="menuitem"
          >
            {item.label}
          </button>
        {/each}
      </div>
    </div>
  {/if}
</div>
