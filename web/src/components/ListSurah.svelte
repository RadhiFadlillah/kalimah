<script lang="ts">
	// Import functions
	import { tick, onMount } from 'svelte';
	import { getRequest } from '../libs/api-request';

	// Data types
	interface Surah {
		id: number;
		name: string;
		translation: string;
		translated: boolean;
	}

	// Props
	let className: string = '';
	export { className as class };

	// Local variables
	let listSurah: Surah[] = [];
	let dataLoading: boolean = false;

	// API function
	async function loadData() {
		dataLoading = true;

		try {
			listSurah = await getRequest('/api/surah');
			console.log(listSurah);
		} catch (err) {
			console.error(err);
		}

		dataLoading = true;
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
				aria-disabled={!surah.translated}
			>
				<p class="number">{idx + 1}</p>
				<p class="name">{surah.name}</p>
				<p class="translation">{surah.translation}</p>
			</div>
		{/each}
	</div>
</div>

<style lang="less">
	div.root {
		display: flex;
		flex-flow: column nowrap;
		background-color: var(--bg);
	}

	p.header {
		flex-shrink: 0;
		padding: 8px;
		font-size: 1.2rem;
		font-variation-settings: 'wght' 600;
		border-bottom: 1px solid var(--border);
		text-align: center;
		color: var(--main);
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
	}
</style>
