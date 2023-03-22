<script lang="ts">
	// Import fragments
	import ListSurah from '../fragments/ListSurah.svelte';
	import Surah from '../fragments/Surah.svelte';
	import AnswerSheet from '../fragments/AnswerSheet.svelte';

	// Import components
	import Header from '../components/Header.svelte';

	// Import dialogs
	import Dialog from '../components/Dialog.svelte';
	import Translation from '../fragments/Translation.svelte';

	// Import type functions
	import type {
		Surah as TSurah,
		Word as TWord,
	} from '../fragments/Surah.svelte';
	import { onMount } from 'svelte';

	// Local variables
	let surahRef: Surah;
	let surahTitle: string;
	let activeSurah: TSurah | undefined;
	let activeWord: TWord | undefined;

	// Dialog translation props
	let dlgTransNumber: number = 0;
	let dlgTransTitle: string | undefined;
	let dlgTransVisible: boolean = false;

	// Dialog error props
	let dlgErrorVisible: boolean = false;
	let dlgErrorMessage: string = '';

	// Reactive variables
	$: headerTitle = activeSurah ? surahTitle : 'Daftar Surah';

	// Lifecycle function
	onMount(() => {
		// This is done to show warning when user trying to close app
		history.pushState(null, document.title, location.href);
	});

	// Event handler for list surah
	function handleListSurahClick(e: CustomEvent) {
		activeWord = undefined;
		activeSurah = e.detail.surah as TSurah;
	}

	// Event handler for surah
	function handleSurahAyahClick(e: CustomEvent) {
		dlgTransNumber = e.detail.ayah as number;
		dlgTransTitle = `${activeSurah?.name} ${dlgTransNumber}`;
		dlgTransVisible = true;
	}

	// Event handler for answer sheet
	function handleAnswerSubmit(e: CustomEvent) {
		let answer = e.detail.answer;
		surahRef?.saveTranslation(activeWord, answer);
	}

	// Event handler for translation
	function handleTranslationError(e: CustomEvent) {
		dlgTransVisible = false;
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
	<Header
		title={headerTitle}
		backVisible={activeSurah != null}
		on:back={() => {
			activeSurah = undefined;
		}}
	/>
	{#if activeSurah == null}
		<ListSurah
			class="list-surah"
			active={activeSurah}
			on:itemclick={handleListSurahClick}
			on:error={handleFragmentError}
		/>
	{:else}
		<Surah
			bind:this={surahRef}
			class="surah"
			surah={activeSurah}
			bind:activeWord
			bind:title={surahTitle}
			on:ayahclick={handleSurahAyahClick}
			on:error={handleFragmentError}
		/>
		{#if activeWord != null}
			<AnswerSheet
				class="answer"
				word={activeWord}
				on:submit={handleAnswerSubmit}
				on:error={handleFragmentError}
			/>
		{/if}
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

	{#if dlgTransVisible}
		<Translation
			ayah={dlgTransNumber}
			title={dlgTransTitle}
			surah={activeSurah?.id || 1}
			on:error={handleTranslationError}
			on:close={() => (dlgTransVisible = false)}
			on:mainclick={() => (dlgTransVisible = false)}
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
		background-color: var(--bg);
		position: relative;
		display: flex;
		flex-flow: column nowrap;
	}

	.app :global(.list-surah),
	.app :global(.surah) {
		flex: 1 0;
	}

	.app :global(.answer) {
		border-top: 1px solid var(--border);
	}
</style>
