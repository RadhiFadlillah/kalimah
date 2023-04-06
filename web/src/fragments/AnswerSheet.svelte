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
			font-variation-settings: 'wght' 600;
			cursor: pointer;

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
