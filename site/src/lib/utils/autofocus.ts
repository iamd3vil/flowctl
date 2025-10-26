/**
 * Svelte action to autofocus an element when it's mounted to the DOM
 * Useful for modal inputs that need to be focused when the modal opens
 *
 * Usage:
 * <input use:autofocus />
 */
export function autofocus(node: HTMLElement) {
	setTimeout(() => node.focus(), 0);
}
