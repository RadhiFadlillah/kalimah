<script lang="ts">
	import { onMount } from 'svelte';
	import ListSurah from '../components/ListSurah.svelte';
	import Surah from '../components/Surah.svelte';
	import type { Surah as TSurah } from '../components/Surah.svelte';

	// Local variables
	let activeSurah: number = 0;
	let activeSurahName: string = '';
	let activeWord: number = 0;

	// Lifecycle function
	onMount(() => {
		history.pushState(null, document.title, location.href);
	});

	// Event handler
	function handleListSurahLoaded(e: CustomEvent) {
		let surah = e.detail.surah;
		if (activeSurah <= 0) {
			activeSurah = surah.id;
			activeSurahName = surah.name;
		}
	}

	function handleListSurahClick(e: CustomEvent) {
		let surah = e.detail.surah as TSurah;
		activeSurah = surah.id;
		activeSurahName = surah.name;
	}

	function handleSurahActived(e: CustomEvent) {
		console.log('ACTIVE WORD:', e.detail.word);
	}
</script>

<svelte:window on:popstate={() => history.go(1)} />

<div class="app">
	<ListSurah
		class="list-surah"
		{activeSurah}
		on:loaded={handleListSurahLoaded}
		on:itemclick={handleListSurahClick}
	/>
	{#if activeSurah > 0}
		<Surah
			class="surah"
			surah={activeSurah}
			surahName={activeSurahName}
			on:actived={handleSurahActived}
		/>
	{/if}
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

	.app :global(.surah) {
		flex: 1 0;
		border-right: 1px solid var(--border);
	}
</style>
