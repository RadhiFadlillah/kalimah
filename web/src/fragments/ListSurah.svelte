<script lang="ts">
	// Import components
	import Button from '../components/Button.svelte';
	import LoadingCover from '../components/LoadingCover.svelte';

	// Import icons and types
	import type { Surah } from './Surah.svelte';
	import icTheme from '@iconify-icons/ic/outline-wb-sunny';
	import icRefresh from '@iconify-icons/ic/outline-refresh';

	// Import functions
	import { getRequest } from '../libs/api-request';
	import { onMount, createEventDispatcher, tick } from 'svelte';
	const dispatch = createEventDispatcher();

	// Props
	let className: string = '';
	export let active: Surah | undefined;
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
			dispatch('error', String(err));
		}

		dataLoading = false;
	}

	function handleItemClick(surah: Surah) {
		dispatch('itemclick', { surah: surah });
	}

	function toggleNightMode() {
		let html = document.documentElement;
		let darkMode = localStorage.getItem('dark-mode') === '1';

		if (darkMode) {
			// Already dark, change to light
			localStorage.removeItem('dark-mode');
			html.classList.remove('dark');
		} else {
			// Theme is light, change to dark
			localStorage.setItem('dark-mode', '1');
			html.classList.add('dark');
		}
	}

	// Lifecycle function
	onMount(() => {
		loadData();
	});
</script>

<div class="root {className}">
	<div class="header">
		<p>Surah</p>
		<Button
			icon={icRefresh}
			disabled={dataLoading}
			on:click={() => {
				window.location.reload();
			}}
		/>
		<Button
			icon={icTheme}
			disabled={dataLoading}
			on:click={() => toggleNightMode()}
		/>
	</div>
	<div class="container" data-scrollbar>
		{#each listSurah as surah, idx (surah.id)}
			<div
				class="item"
				role="button"
				tabindex="0"
				on:click={() => handleItemClick(surah)}
				class:active={surah.id === active?.id}
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

	div.header {
		display: flex;
		align-items: center;
		flex-flow: row nowrap;
		border-bottom: 1px solid var(--border);
		flex-shrink: 0;

		p {
			flex: 1 0;
			font-size: 1.2rem;
			font-variation-settings: 'wght' 600;
			text-align: center;
			color: var(--main);
			line-height: 36px;
		}
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
