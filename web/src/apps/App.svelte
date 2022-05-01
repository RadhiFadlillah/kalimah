<script lang="ts">
	import { onMount } from 'svelte';
	import Dialog from '../components/Dialog.svelte';
	import Surah from '../fragments/Surah.svelte';
	import ListSurah from '../fragments/ListSurah.svelte';
	import AnswerSheet from '../fragments/AnswerSheet.svelte';
	import type {
		Surah as TSurah,
		Word as TWord,
	} from '../components/Surah.svelte';

	// Local variables
	let surahRef: Surah;
	let activeSurah: TSurah | undefined;
	let activeWord: TWord | undefined;

	// Lifecycle function
	onMount(() => {
		history.pushState(null, document.title, location.href);
	});

	// Event handler
	function handleListSurahLoaded(e: CustomEvent) {
		let surah = e.detail.surah;
		if (activeSurah == null) {
			activeSurah = surah;
		}
	}

	function handleListSurahClick(e: CustomEvent) {
		activeSurah = e.detail.surah as TSurah;
	}

	function handleWordActived(e: CustomEvent) {
		console.log('ACTIVE WORD:', e.detail.word);
		activeWord = e.detail.word as TWord;
	}

	function handleAnswerSubmit(e: CustomEvent) {
		let answer = e.detail.answer;
		surahRef?.saveTranslation(activeWord, answer);
	}
</script>

<svelte:window on:popstate={() => history.go(1)} />

<div class="app">
	<ListSurah
		class="list-surah"
		active={activeSurah}
		on:loaded={handleListSurahLoaded}
		on:itemclick={handleListSurahClick}
	/>
	{#if activeSurah != null}
		<Surah
			bind:this={surahRef}
			class="surah"
			surah={activeSurah}
			on:actived={handleWordActived}
		/>
	{/if}
	{#if activeWord != null}
		<AnswerSheet
			class="answer"
			word={activeWord}
			on:submit={handleAnswerSubmit}
		/>
	{/if}

	<Dialog title="Kalimah" closable />
</div>

<style lang="less">
	.app {
		width: 100%;
		height: 100%;
		max-width: 100%;
		max-height: 100%;
		min-width: 0;
		min-height: 0;
		overflow: hidden;
		color: var(--fg);
		background-color: var(--border);
		position: relative;
		display: grid;
		gap: 1px;
		grid-template-rows: minmax(0, 1fr) auto;
		grid-template-columns: 220px minmax(0, 1fr);
	}

	.app :global(.list-surah) {
		width: 220px;
		grid-row: 1/-1;
		grid-column: 1;
	}

	.app :global(.surah),
	.app :global(.answer) {
		grid-column: 2;
		max-width: 1080px;
	}
</style>
