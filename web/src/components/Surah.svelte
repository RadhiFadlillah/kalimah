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
	export let surah: Surah | undefined;
	export { className as class };

	// Constants
	const ayahPerPage = 30;
	const arabicNumerals = '٠١٢٣٤٥٦٧٨٩';

	// Local variables
	let words: Word[] = [];
	let pageCount: number;
	let currentPage: number;

	let container: HTMLElement;
	let dataLoading: boolean = false;

	// Reactive variables
	$: headerTitle = ((): string => {
		if (surah == null) return '';
		if (dataLoading) return surah!.name;

		let nWords = words.length;
		let nTranslatedWords = words.filter((w) => w.translation !== '').length;
		let progress = Math.round((nTranslatedWords / nWords) * 100) || 0;
		return `(${progress}%) • ${surah!.name} • ${nTranslatedWords} / ${nWords}`;
	})();
	$: visibleWords = words.filter((w) => {
		let startingAyah = (currentPage - 1) * ayahPerPage + 1;
		let endingAyah = currentPage * ayahPerPage;
		return w.ayah >= startingAyah && w.ayah <= endingAyah;
	});
	$: activeWord = ((): Word | undefined => {
		let realActiveWord = words.find((w) => w.translation === '');
		return visibleWords.find((w) => w.id === realActiveWord?.id);
	})();

	// API function
	async function loadData(surah?: number) {
		words = [];
		dataLoading = true;

		try {
			words = await getRequest(`/api/surah/${surah}/word`);
			let maxAyah: number = words
				.map((w: Word) => w.ayah)
				.reduce((p, c) => (p > c ? p : c), 1);
			let realActiveWord = words.find((w) => w.translation === '');

			pageCount = Math.ceil(maxAyah / ayahPerPage);
			if (realActiveWord == null) currentPage = 1;
			else currentPage = Math.ceil(realActiveWord.ayah / ayahPerPage);

			await tick();
			focusToActive();
		} catch (err) {
			console.error(err);
		}

		dataLoading = false;
	}

	// Local function
	function isAyahSeparator(idx: number): boolean {
		let currentAyah = visibleWords[idx]?.ayah;
		let nextAyah = visibleWords[idx + 1]?.ayah;
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

	function focusToActive() {
		let active = container.querySelector('.item.active') as HTMLElement;
		let scrollTop = active == null ? 0 : active.offsetTop;
		container.scrollTop = scrollTop - 36;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Home' && !e.ctrlKey) {
			e.preventDefault();
			e.stopPropagation();
			focusToActive();
		} else if (e.key === '-' && currentPage > 1) {
			handlePagination(currentPage - 1);
		} else if (e.key === '+' && currentPage < pageCount) {
			handlePagination(currentPage + 1);
		}
	}

	function handlePagination(page: number) {
		currentPage = page;
		tick().then(() => focusToActive());
	}

	// Lifecycle function
	onMount(() => loadData(surah?.id));

	// Reload data whenever surah changed
	$: loadData(surah?.id);

	// Dispatch active word whenever it changed
	$: dispatch('actived', { word: activeWord });
</script>

<div class="root {className}">
	<p class="header">{headerTitle}</p>
	<div
		class="container"
		data-scrollbar
		tabindex="0"
		bind:this={container}
		on:keydown={handleKeydown}
	>
		{#each visibleWords as word, idx (word.id)}
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
	{#if pageCount > 1}
		<div class="footer">
			{#each [...Array(pageCount).keys()].map((k) => k + 1) as page}
				<button
					on:click={() => handlePagination(page)}
					class:active={page === currentPage}
					>{page}
				</button>
			{/each}
		</div>
	{/if}
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

	div.footer {
		padding: 0 8px;
		flex-shrink: 0;
		border-top: 1px solid var(--border);
		display: flex;
		flex-flow: row wrap;
		align-items: center;
		justify-content: center;

		button {
			font-size: 1rem;
			padding: 8px;
			margin: 8px;
			background-color: transparent;
			cursor: pointer;

			&:hover,
			&:focus {
				background-color: var(--bg-hover);
			}

			&.active {
				color: var(--main);
				background-color: var(--main-bg);
				pointer-events: none;
				cursor: default;
			}
		}
	}

	div.container {
		display: flex;
		flex: 1 0;
		flex-flow: row-reverse wrap;
		justify-content: center;
		align-content: flex-start;

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
