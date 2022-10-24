<script lang="ts">
	// Import functions
	import { onMount, createEventDispatcher } from 'svelte';
	import { getRequest, postRequest } from '../libs/api-request';
	import LoadingCover from '../components/LoadingCover.svelte';
	import type { Word } from './Surah.svelte';
	const dispatch = createEventDispatcher();

	// Props
	let className: string = '';
	export let word: Word | undefined;
	export { className as class };

	// Local variables
	let choices: string[] = [];
	let wrongChoices: string[] = [];
	let dataLoading: boolean = false;

	// API function
	async function loadChoices(word?: Word) {
		if (word == null || dataLoading) return;

		choices = [];
		wrongChoices = [];
		dataLoading = true;

		try {
			choices = await getRequest(`/api/choice/${word?.id}`);
		} catch (err) {
			dispatch('error', String(err));
		}

		dataLoading = false;
	}

	async function submitAnswer(answer: string) {
		if (word == null || dataLoading) return;
		dataLoading = true;

		let resp: any;
		try {
			let body = { word: word, answer: answer };
			resp = await postRequest('/api/answer', body);
		} catch (err) {
			console.error(err);
			resp = -1;
		}

		dataLoading = false;
		if (resp === -1) {
			return;
		}

		if (resp === '0') {
			wrongChoices = [...wrongChoices, answer];
			if (wrongChoices.length >= 5) {
				await loadChoices(word);
				return;
			}
		}

		if (resp === '1') {
			dispatch('submit', { answer: answer });
		}
	}

	// Lifecycle function
	onMount(() => loadChoices(word));

	// Reload data whenever word changed
	$: loadChoices(word);
</script>

<div class="root {className}">
	<p class="arabic">{word?.arabic}</p>
	<div class="container">
		{#each choices as choice}
			<button
				class:wrong={wrongChoices.includes(choice)}
				on:click={() => submitAnswer(choice)}
				>{choice}
			</button>
		{/each}
	</div>
	{#if dataLoading}
		<LoadingCover class="answer-loading" />
	{/if}
</div>

<style lang="less">
	div.root {
		display: flex;
		flex-flow: row nowrap;
		background-color: var(--bg);
		align-items: center;
		position: relative;
	}

	p.arabic {
		padding: 16px;
		font-size: 4rem;
		font-family: 'KFGQPC-HAFS';
		text-align: center;
		color: var(--main);
		direction: rtl;
	}

	div.container {
		display: flex;
		flex: 1 0;
		flex-flow: row wrap;
		align-content: center;
		padding: 8px;

		button {
			font-size: 1.1rem;
			padding: 8px;
			margin: 4px;
			color: var(--fg);
			background-color: transparent;
			border: 1px solid var(--border);
			border-radius: 8px;
			cursor: pointer;

			&:hover,
			&:focus {
				background-color: var(--bg-hover);
			}

			&.wrong {
				color: var(--fg-error);
				background-color: var(--bg-error);
				cursor: pointer;
				pointer-events: none;
			}
		}
	}

	div.root :global(.answer-loading) {
		position: absolute;
	}

	@media screen and (max-width: 600px) {
		div.root {
			flex-flow: column nowrap;
		}

		p.arabic {
			font-size: 3rem;
			padding-bottom: 0;
		}

		div.container {
			justify-content: center;
		}
	}
</style>
