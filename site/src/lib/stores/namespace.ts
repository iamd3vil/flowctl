import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Get initial value from localStorage if in browser
const storedNamespace = browser ? localStorage.getItem('selectedNamespace') : null;

export const selectedNamespace = writable<string>(storedNamespace || 'default');

// Subscribe to changes and save to localStorage
if (browser) {
	selectedNamespace.subscribe(value => {
		localStorage.setItem('selectedNamespace', value);
	});
}
