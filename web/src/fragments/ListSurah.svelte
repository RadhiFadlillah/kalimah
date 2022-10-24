<script lang="ts">
	// Import components
	import LoadingCover from '../components/LoadingCover.svelte';

	// Import types
	import type { Surah } from './Surah.svelte';

	// Import functions
	import { getRequest } from '../libs/api-request';
	import { onMount, createEventDispatcher, tick } from 'svelte';
	const dispatch = createEventDispatcher();

	// Props
	let className: string = '';
	export let active: Surah | undefined;
	export let style: string = '';
	export { className as class };

	// Local variables
	let listSurah: Surah[] = [];
	let dataLoading: boolean = false;

	// API function
	async function loadData() {
		dataLoading = true;

		try {
			listSurah = await getRequest('/api/surah');
			await tick();
		} catch (err) {
			dispatch('error', String(err));
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

<div class="root {className}" {style} data-scrollbar>
	<div class="container">
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
	}

	div.container {
		gap: 1px;
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		flex-shrink: 0;
	}

	div.item {
		display: grid;
		gap: 8px;
		padding: 8px;
		grid-template-rows: auto auto;
		grid-template-columns: 28px minmax(0, 1fr);
		outline: 1px solid var(--border);
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
</style>
