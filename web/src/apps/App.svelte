<script lang="ts">
	import { onMount } from 'svelte';
	import ListSurah from '../components/ListSurah.svelte';

	onMount(() => {
		history.pushState(null, document.title, location.href);
	});
</script>

<svelte:window on:popstate={() => history.go(1)} />

<div class="app">
	<ListSurah class="list-surah" />
</div>

<style lang="less">
	.app {
		display: flex;
		flex-flow: row nowrap;
		width: 100%;
		height: 100%;
		max-width: 100%;
		max-height: 100%;
		min-width: 0;
		min-height: 0;
		overflow: hidden;
		color: var(--fg);
		background-color: var(--bg-secondary);
		position: relative;
	}

	.app :global(.list-surah) {
		width: 220px;
		border-right: 1px solid var(--border);
	}
</style>
