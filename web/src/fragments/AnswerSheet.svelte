<script lang="ts" context="module">
</script>

<script lang="ts">
	// Import functions
	import { createEventDispatcher } from 'svelte';
	import { postRequest } from '../libs/api-request';
	import LoadingCover from '../components/LoadingCover.svelte';
	import type { Word, Choice } from './Surah.svelte';
	const dispatch = createEventDispatcher();

	// Props
	let className: string = '';
	export let word: Word | undefined;
	export { className as class };

	// Local variables
	let wrongChoices: string[] = [];
	let dataLoading: boolean = false;

	// API function
	async function submitAnswer(choice: Choice, word?: Word) {
		if (word == null || dataLoading) return;

		// If choice is incorrect, stop
		if (!choice.isCorrect) {
			wrongChoices = [...wrongChoices, choice.text];
			if (wrongChoices.length >= 5) wrongChoices = [];
			return;
		}

		// If this word is separator, track it in database
		if (word.isSeparator) {
			dataLoading = true;
			let errorOccured = false;

			try {
				await postRequest('/api/track', word);
			} catch (err) {
				console.error(err);
				errorOccured = true;
			}

			dataLoading = false;
			if (errorOccured) return;
		}

		dispatch('answered');
	}

	$: {
		word;
		(document.activeElement as HTMLElement).blur();
	}
</script>

<div class="root {className}">
	<p class="arabic">{word?.arabic}</p>
	<div class="container">
		{#each word?.choices || [] as choice}
			<button
				class:wrong={wrongChoices.includes(choice.text)}
				on:click={() => submitAnswer(choice, word)}
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
		flex-flow: column nowrap;
		background-color: var(--bg);
		position: relative;
	}

	p.arabic {
		padding: 8px;
		font-size: 3rem;
		font-family: 'KFGQPC-HAFS';
		text-align: center;
		color: var(--main);
		direction: rtl;
	}

	div.container {
		flex: 1 0;
		display: grid;
		gap: 1px;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));

		button {
			font-size: 1.1rem;
			padding: 16px;
			color: var(--fg);
			background-color: var(--bg);
			outline: 1px solid var(--border);
			font-variation-settings: 'wght' 600;
			cursor: pointer;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;

			&:active,
			&:focus {
				color: var(--fg);
				background-color: var(--bg);
			}

			&.wrong {
				color: var(--bg);
				cursor: pointer;
				pointer-events: none;
			}
		}
	}

	div.root :global(.answer-loading) {
		position: absolute;
	}
</style>
