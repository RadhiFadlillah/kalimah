<script lang="ts">
	import { onMount } from 'svelte';
	import Dialog from '../components/Dialog.svelte';
	import Ayah from '../fragments/Ayah.svelte';
	import Surah from '../fragments/Surah.svelte';
	import ListSurah from '../fragments/ListSurah.svelte';
	import AnswerSheet from '../fragments/AnswerSheet.svelte';
	import type {
		Surah as TSurah,
		Word as TWord,
	} from '../fragments/Surah.svelte';

	// Local variables
	let surahRef: Surah;
	let activeSurah: TSurah | undefined;
	let activeWord: TWord | undefined;

	// Dialog ayah props
	let dlgAyahNumber: number = 0;
	let dlgAyahVisible: boolean = false;

	// Dialog error props
	let dlgErrorVisible: boolean = false;
	let dlgErrorMessage: string = '';

	// Lifecycle function
	onMount(() => {
		history.pushState(null, document.title, location.href);
	});

	// Event handler for list surah
	function handleListSurahLoaded(e: CustomEvent) {
		let surah = e.detail.surah;
		if (activeSurah == null) {
			activeSurah = surah;
		}
	}

	function handleListSurahClick(e: CustomEvent) {
		activeWord = undefined;
		activeSurah = e.detail.surah as TSurah;
	}

	// Event handler for surah
	function handleSurahActived(e: CustomEvent) {
		let newActiveWord = e.detail.word as TWord;
		if (newActiveWord?.id !== activeWord?.id) {
			activeWord = newActiveWord;
		}
	}

	function handleSurahAyahClick(e: CustomEvent) {
		dlgAyahNumber = e.detail.ayah as number;
		dlgAyahVisible = true;
	}

	// Event handler for answer sheet
	function handleAnswerSubmit(e: CustomEvent) {
		let answer = e.detail.answer;
		surahRef?.saveTranslation(activeWord, answer);
	}

	// Event handler for ayah
	function handleAyahError(e: CustomEvent) {
		dlgAyahVisible = false;
		dlgErrorVisible = true;
		dlgErrorMessage = e.detail;
	}

	// Common event handler for fragments
	function handleFragmentError(e: CustomEvent) {
		dlgErrorVisible = true;
		dlgErrorMessage = e.detail;
	}
</script>

<svelte:window on:popstate={() => history.go(1)} />

<div class="app">
	<ListSurah
		class="list-surah"
		active={activeSurah}
		on:loaded={handleListSurahLoaded}
		on:itemclick={handleListSurahClick}
		on:error={handleFragmentError}
	/>
	{#if activeSurah != null}
		<Surah
			bind:this={surahRef}
			class="surah"
			surah={activeSurah}
			on:actived={handleSurahActived}
			on:ayahclick={handleSurahAyahClick}
			on:error={handleFragmentError}
		/>
	{/if}
	{#if activeWord != null}
		<AnswerSheet
			class="answer"
			word={activeWord}
			on:submit={handleAnswerSubmit}
			on:error={handleFragmentError}
		/>
	{/if}

	{#if dlgErrorVisible}
		<Dialog
			title="Error"
			isError={true}
			closable={false}
			on:mainclick={() => (dlgErrorVisible = false)}
		>
			<p slot="content">{dlgErrorMessage}</p>
		</Dialog>
	{/if}

	{#if dlgAyahVisible}
		<Ayah
			ayah={dlgAyahNumber}
			surah={activeSurah?.id || 1}
			on:error={handleAyahError}
			on:close={() => (dlgAyahVisible = false)}
			on:mainclick={() => (dlgAyahVisible = false)}
		/>
	{/if}
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
