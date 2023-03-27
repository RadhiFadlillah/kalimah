<script lang="ts" context="module">
	export interface Surah {
		id: number;
		name: string;
		translation: string;
		translated: boolean;
	}

	export interface Choice {
		text: string;
		isCorrect: boolean;
	}

	export interface Word {
		id: number;
		ayah: number;
		position: number;
		arabic: string;
		translation: string;
		answered: boolean;
		isSeparator: boolean;
		choices: Choice[];
	}

	interface FetchResponse {
		currentPage: number;
		maxPage: number;
		words: Word[];
		disabled: boolean;
	}
</script>

<script lang="ts">
	import LoadingCover from '../components/LoadingCover.svelte';
	import { createEventDispatcher, onMount, tick } from 'svelte';
	import { getRequest } from '../libs/api-request';
	const dispatch = createEventDispatcher();

	// Props
	let className: string = '';
	export { className as class };
	export let surah: Surah | undefined;
	export let activeWord: Word | undefined = undefined;

	// Constants
	const arabicNumerals = '٠١٢٣٤٥٦٧٨٩';

	// Local variables
	let words: Word[] = [];
	let maxPage: number;
	let currentPage: number;
	let pageDisabled: boolean = false;
	let dataLoading: boolean = false;
	let container: HTMLElement;

	// Reactive variables
	$: {
		if (pageDisabled) activeWord = undefined;
		else activeWord = words.find((w) => !w.answered);
	}

	// API function
	async function loadData() {
		words = [];
		dataLoading = true;

		try {
			let url = `/api/words/surah/${surah?.id}/page/${currentPage || 0}`;
			let resp = (await getRequest(url)) as FetchResponse;

			words = resp.words;
			maxPage = resp.maxPage;
			currentPage = resp.currentPage;
			pageDisabled = resp.disabled;

			await tick();
			focusToActive();
		} catch (err) {
			dispatch('error', String(err));
		}

		dataLoading = false;
	}

	// Local function
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
			changePage(currentPage - 1);
		} else if (e.key === '+' && currentPage < maxPage) {
			changePage(currentPage + 1);
		}
	}

	function changePage(page: number) {
		currentPage = page;
		loadData();
	}

	function handleAyahClick(ayah: number) {
		dispatch('ayahclick', { surah: surah?.id, ayah: ayah });
	}

	// Exported function
	export function markAnswered(word: Word | undefined) {
		if (word == null) return;

		let idx = words.findIndex((w) => w.id === word.id);
		if (idx >= 0) words[idx].answered = true;

		let next = word.id + 1;
		let nextIsVisible = words.findIndex((w) => w.id === next) >= 0;
		if (!nextIsVisible && currentPage < maxPage) {
			changePage(currentPage + 1);
		} else {
			focusToActive();
		}
	}

	onMount(() => loadData());
</script>

<div class="root {className}">
	<div
		class="container"
		data-scrollbar
		tabindex="0"
		bind:this={container}
		on:keydown={handleKeydown}
	>
		{#each words as word (word.id)}
			<div
				class="item"
				tabindex="0"
				role="button"
				class:active={word.id === activeWord?.id}
				aria-disabled={!word.answered}
			>
				<p class="arabic">{word.arabic}</p>
				<p class="translation" class:unanswered={!word.answered}>
					{word.translation}
				</p>
			</div>

			{#if word.isSeparator}
				<button
					class="number"
					disabled={!word.answered}
					on:click={() => handleAyahClick(word.ayah)}
				>
					{toArabicNumeral(word.ayah)}
				</button>
			{/if}
		{/each}
	</div>
	{#if maxPage > 1}
		<div class="footer">
			{#each [...Array(maxPage).keys()].map((k) => k + 1) as page}
				<button
					on:click={() => changePage(page)}
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
			color: var(--fg);
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

			p.arabic {
				font-size: 3rem;
				font-family: 'KFGQPC-HAFS';
				text-align: center;
				color: var(--fg);
				direction: rtl;
			}

			p.translation {
				font-size: 0.9rem;
				color: var(--fg-secondary);
				text-align: center;

				&.unanswered {
					visibility: hidden;
				}
			}

			&[aria-disabled='true'] p {
				color: var(--fg-disabled);
			}

			&.active {
				background-color: var(--main-bg);

				p {
					color: var(--main);
				}
			}
		}

		button.number {
			align-self: center;
			font-size: 3.5rem;
			font-family: 'KFGQPC-HAFS';
			padding: 16px;
			color: var(--main);
			background-color: transparent;
			cursor: pointer;

			&:hover,
			&:focus {
				background-color: var(--bg-hover);
			}

			&[disabled] {
				pointer-events: none;
				color: var(--fg-disabled);
			}
		}
	}
</style>
