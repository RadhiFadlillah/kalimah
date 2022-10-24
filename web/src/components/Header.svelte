<script lang="ts">
	// Import components
	import Button from '../components/Button.svelte';

	// Import icons
	import icBack from '@iconify-icons/ic/outline-arrow-back';
	import icTheme from '@iconify-icons/ic/outline-wb-sunny';
	import icRefresh from '@iconify-icons/ic/outline-refresh';

	// Import functions
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	// Props
	let className: string = '';
	export let title: string = '';
	export let backVisible: boolean = false;
	export { className as class };

	// Local functions
	function handleBack() {
		dispatch('back');
	}

	function reloadPage() {
		window.location.reload();
	}

	function toggleNightMode() {
		let html = document.documentElement;
		let darkMode = localStorage.getItem('dark-mode') === '1';

		if (darkMode) {
			// Already dark, change to light
			localStorage.removeItem('dark-mode');
			html.classList.remove('dark');
		} else {
			// Theme is light, change to dark
			localStorage.setItem('dark-mode', '1');
			html.classList.add('dark');
		}
	}
</script>

<div class="header {className}">
	{#if backVisible}
		<Button icon={icBack} on:click={handleBack} />
	{/if}
	<p>{title}</p>
	<Button icon={icRefresh} on:click={reloadPage} />
	<Button icon={icTheme} on:click={toggleNightMode} />
</div>

<style lang="less">
	div.header {
		display: flex;
		align-items: center;
		flex-flow: row nowrap;
		border-bottom: 1px solid var(--border);
		flex-shrink: 0;

		p {
			flex: 1 0;
			font-size: 1.2rem;
			font-variation-settings: 'wght' 600;
			text-align: center;
			color: var(--main);
			line-height: 36px;
		}
	}
</style>
