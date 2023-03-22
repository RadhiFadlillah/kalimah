<script lang="ts" context="module">
	export interface Choice {
		text: string;
		isCorrect: boolean;
	}
</script>

<script lang="ts">
	// Import functions
	import { onMount, createEventDispatcher } from 'svelte';
	import { getRequest, postRequest } from '../libs/api-request';
	import LoadingCover from '../components/LoadingCover.svelte';
	import type { Word } from './Surah.svelte';
	import { text } from 'svelte/internal';
	const dispatch = createEventDispatcher();

	// Props
	let className: string = '';
	export let word: Word | undefined;
	export { className as class };

	// Local variables
	let choices: Choice[] = [];
	let wrongChoices: string[] = [];
	let dataLoading: boolean = false;

	// API function
	async function loadChoices(word?: Word) {
		if (word == null || dataLoading) return;

		wrongChoices = [];
		dataLoading = true;

		try {
			choices = await getRequest(`/api/choice/${word?.id}`);
		} catch (err) {
			dispatch('error', String(err));
			choices = [];
		}

		dataLoading = false;
	}

	async function submitAnswer(choice: Choice) {
		if (word == null || dataLoading) return;

		// If choice is incorrect, reload choices if necessary
		if (!choice.isCorrect) {
			wrongChoices = [...wrongChoices, choice.text];
			if (wrongChoices.length >= 5) {
				await loadChoices(word);
			}
			return;
		}

		// If this word is separator, save it in database
		if (word.isSeparator) {
			dataLoading = true;
			let errorOccured = false;

			try {
				await postRequest('/api/answer', word);
			} catch (err) {
				console.error(err);
				errorOccured = true;
			}

			dataLoading = false;
			if (errorOccured) return;
		}

		dispatch('submit', { answer: choice.text });
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
				class:wrong={wrongChoices.includes(choice.text)}
				on:click={() => submitAnswer(choice)}
				>{choice.text}
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
