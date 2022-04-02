<script lang="ts">
	// Import
	import { onMount, createEventDispatcher, tick } from 'svelte';
	import { getRequest } from '../libs/api-request';
	import LoadingCover from './LoadingCover.svelte';
	import type { Surah } from './Surah.svelte';
	const dispatch = createEventDispatcher();

	// Props
	let className: string = '';
	export let activeSurah: number = 0;
	export { className as class };

	// Local variables
	let listSurah: Surah[] = [];
	let dataLoading: boolean = false;

	// Reactive variable
	$: lastTranslatedSurah = ((): Surah => {
		let idxUntranslated = listSurah.findIndex((s) => s.translated === false);
		return idxUntranslated < 0 ? listSurah[0] : listSurah[idxUntranslated - 1];
	})();

	// API function
	async function loadData() {
		dataLoading = true;

		try {
			listSurah = await getRequest('/api/surah');
			await tick();
			dispatch('loaded', { surah: lastTranslatedSurah });
		} catch (err) {
			console.error(err);
		}

		dataLoading = false;
	}

	function handleItemClick(surah: Surah) {
		dispatch('itemclick', { surah: surah });
	}

	// Lifecycle function
	onMount(() => {
		loadData();
	});
</script>

<div class="root {className}">
	<p class="header">Surah</p>
	<div class="container" data-scrollbar>
		{#each listSurah as surah, idx (surah.id)}
			<div
				class="item"
				role="button"
				tabindex="0"
				on:click={() => handleItemClick(surah)}
				class:active={surah.id === activeSurah}
				aria-disabled={!surah.translated}
			>
				<p class="number">{idx + 1}</p>
				<p class="name">{surah.name}</p>
				<p class="translation">{surah.translation}</p>
			</div>
		{/each}
	</div>
	{#if dataLoading}
		<LoadingCover class="list-loading" />
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

	div.item {
		display: grid;
		gap: 8px;
		padding: 8px;
		grid-template-rows: auto auto;
		grid-template-columns: 28px minmax(0, 1fr);
		border-bottom: 1px solid var(--border);
		cursor: pointer;

		p.number {
			grid-row: 1 / span 2;
			grid-column: 1;
			font-variation-settings: 'wght' 600;

			&::after {
				content: '.';
			}
		}

		p.name,
		p.translation {
			overflow: hidden;
			white-space: nowrap;
			text-overflow: ellipsis;
		}

		p.name {
			font-variation-settings: 'wght' 600;
		}

		p.translation {
			font-size: 0.9rem;
			color: var(--fg-secondary);
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
				font-variation-settings: 'wght' 500;
			}
		}

		&.active {
			background-color: var(--main-bg);

			p {
				color: var(--main);
			}
		}
	}

	div.root :global(.list-loading) {
		z-index: 1;
		position: absolute;
		top: 37px;
		left: 0;
		right: 0;
		bottom: 0;
	}
</style>
