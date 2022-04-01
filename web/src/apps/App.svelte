<script lang="ts">
	import { onMount } from 'svelte';

	onMount(() => {
		history.pushState(null, document.title, location.href);
	});
</script>

<svelte:window
	on:popstate={() => history.go(1)}
	on:beforeunload|preventDefault={(e) => (e.returnValue = '')}
/>

<div class="app" />

<style lang="less">
	.app {
		display: flex;
		flex-flow: column nowrap;
		width: 100%;
		height: 100%;
		max-width: 100%;
		max-height: 100%;
		min-width: 0;
		min-height: 0;
		overflow: hidden;
		color: var(--fg);
		background-color: var(--bg);
		position: relative;
	}
</style>
