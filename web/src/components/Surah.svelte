<script lang="ts" context="module">
	export interface Surah {
		id: number;
		name: string;
		translation: string;
		translated: boolean;
	}

	export interface Word {
		id: number;
		ayah: number;
		position: number;
		arabic: string;
		translation: string;
	}
</script>

<script lang="ts">
	// Import functions
	import { onMount, createEventDispatcher, tick } from 'svelte';
	import { getRequest } from '../libs/api-request';
	import LoadingCover from './LoadingCover.svelte';
	const dispatch = createEventDispatcher();

	// Props
	let className: string = '';
	export let surah: number = 2;
	export let surahName: string = 'Al-Baqarah';
	export { className as class };

	// Constants
	const arabicNumerals = '٠١٢٣٤٥٦٧٨٩';

	// Local variables
	let words: Word[] = [];
	let container: HTMLElement;
	let dataLoading: boolean = false;

	// Reactive variables
	$: nWords = words.length;
	$: nTranslatedWords = words.filter((w) => w.translation !== '').length;
	$: progress = Math.round((nTranslatedWords / nWords) * 100) || 0;
	$: activeWord = words.find((w) => w.translation === '');

	// API function
	async function loadData(surah: number) {
		dataLoading = true;
		words = [];

		try {
			words = await getRequest(`/api/surah/${surah}/word`);
			await tick();
			focusToUntranslated();
		} catch (err) {
			console.error(err);
		}

		dataLoading = false;
	}

	// Local function
	function isAyahSeparator(idx: number): boolean {
		let currentAyah = words[idx]?.ayah;
		let nextAyah = words[idx + 1]?.ayah;
		return currentAyah !== nextAyah;
	}

	function toArabicNumeral(number: number): string {
		let result: string = '';
		for (const c of number.toString()) {
			let idx = parseInt(c, 10);
			if (!isNaN(idx)) result += arabicNumerals[idx];
		}
		return result;
	}

	function focusToUntranslated() {
		let unfocused = container.querySelector('.item[aria-disabled="true"]');
		if (unfocused != null) unfocused.scrollIntoView();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Home' && !e.ctrlKey) {
			e.preventDefault();
			e.stopPropagation();
			focusToUntranslated();
		}
	}

	// Lifecycle function
	onMount(() => loadData(surah));

	$: {
		loadData(surah);
	}

	$: {
		dispatch('actived', { word: activeWord });
	}
</script>

<div class="root {className}">
	<p class="header">
		({progress}%) • {surahName} • {nTranslatedWords} / {nWords}
	</p>
	<div
		class="container"
		data-scrollbar
		tabindex="0"
		bind:this={container}
		on:keydown={handleKeydown}
	>
		{#each words as word, idx (word.id)}
			<div
				class="item"
				tabindex="0"
				role="button"
				class:active={word.id === activeWord?.id}
				aria-disabled={word.translation === ''}
			>
				<p class="arabic">{word.arabic}</p>
				<p class="translation">{word.translation}</p>
			</div>

			{#if isAyahSeparator(idx)}
				<p class="number" class:disabled={word.translation === ''}>
					{toArabicNumeral(word.ayah)}
				</p>
			{/if}
		{/each}
	</div>
	{#if dataLoading}
		<LoadingCover class="surah-loading" />
	{/if}
</div>

<style lang="less">
	div.root {
		display: flex;
		flex-flow: column nowrap;
		background-color: var(--bg);
		position: relative;
	}

	p.header {
		flex-shrink: 0;
		padding: 0 8px;
		font-size: 1.2rem;
		font-variation-settings: 'wght' 600;
		border-bottom: 1px solid var(--border);
		text-align: center;
		color: var(--main);
		line-height: 36px;
	}

	div.container {
		display: flex;
		flex-flow: row-reverse wrap;
		justify-content: center;
	}

	div.item {
		display: flex;
		flex-flow: column nowrap;
		margin: 8px;
		padding: 8px;
		cursor: pointer;

		p.arabic {
			font-size: 3rem;
			font-family: 'KFGQPC-HAFS';
			text-align: center;
			color: var(--fg);
		}

		p.translation {
			font-size: 0.9rem;
			color: var(--fg-secondary);
			text-align: center;
		}

		&:hover,
		&:focus {
			background-color: var(--bg-hover);
		}

		&[aria-disabled='true'] {
			pointer-events: none;
			cursor: default;

			p {
				color: var(--fg-disabled);
			}
		}

		&.active {
			background-color: var(--main-bg);

			p {
				color: var(--main);
			}
		}
	}

	p.number {
		align-self: center;
		font-size: 3.5rem;
		font-family: 'KFGQPC-HAFS';
		padding: 16px;
		color: var(--main);
		border-style: solid;
		pointer-events: none;

		&.disabled {
			color: var(--fg-disabled);
		}
	}

	div.root :global(.surah-loading) {
		z-index: 1;
		position: absolute;
		top: 37px;
		left: 0;
		right: 0;
		bottom: 0;
	}
</style>
