<script lang="ts">
	// Components
	import Button from './Button.svelte';
	import LoadingCover from './LoadingCover.svelte';

	// Icons
	import icClose from '@iconify-icons/ic/outline-close';

	// Transition
	import { quintOut } from 'svelte/easing';
	import { fade } from 'svelte/transition';

	// Functions
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	// Properties
	let className: string = '';
	export let title: string;
	export let style: string = '';
	export let closable: boolean = true;
	export let isError: boolean = false;
	export let loading: boolean = false;
	export let mainButtonText: string = 'OK';
	export let secondaryButtonText: string | undefined = undefined;
	export let defaultFocusTargets: string[] = [];
	export { className as class };

	// Actions
	function initFocus(element: HTMLElement) {
		if (Array.isArray(defaultFocusTargets)) {
			for (let i = 0; i < defaultFocusTargets.length; i++) {
				let selector = defaultFocusTargets[i];
				let target = element.querySelector(selector);
				if (target) {
					(target as HTMLElement).focus();
					return;
				}
			}
		}

		let selectors = [
			':scope .content button:enabled',
			':scope .content input:read-write',
			':scope .content textarea:read-write',
			':scope .content label.checkbox-input',
			':scope .footer button:enabled',
			':scope .footer input:read-write',
			':scope .footer textarea:read-write',
		].join(', ');

		let focusableElement: Element | null = null;
		focusableElement = element.querySelector(selectors);

		if (focusableElement != null) {
			(focusableElement as HTMLElement).focus();
		} else {
			element.focus();
		}
	}

	// Event handler
	function handleCloseClick() {
		dispatch('close');
	}

	function handleMainButtonClick() {
		dispatch('mainclick');
	}

	function handleSecondaryButtonClick() {
		dispatch('secondaryclick');
	}

	function handleKeydown(e: KeyboardEvent) {
		if (closable && (e.key === 'Esc' || e.key === 'Escape')) {
			e.preventDefault();
			e.stopPropagation();
			dispatch('close');
		}
	}
</script>

<div
	class="dialog-overlay"
	tabindex="-1"
	on:keydown={handleKeydown}
	on:keydown|capture
	transition:fade={{ easing: quintOut }}
>
	<div {style} class="dialog {className}" class:error={isError} use:initFocus>
		<div class="header">
			<p class="title">{title}</p>
			{#if closable}
				<Button icon={icClose} disabled={loading} on:click={handleCloseClick} />
			{/if}
		</div>

		<div class="content" data-scrollbar tabindex="-1">
			<slot name="content">
				<p>Dialog is empty</p>
			</slot>

			{#if loading}
				<LoadingCover class="loading-cover" />
			{/if}
		</div>

		<div class="footer">
			<slot name="footer">
				<Button
					text={mainButtonText}
					disabled={loading}
					on:click={handleMainButtonClick}
				/>
				{#if secondaryButtonText}
					<Button
						text={secondaryButtonText}
						disabled={loading}
						on:click={handleSecondaryButtonClick}
					/>
				{/if}
			</slot>
		</div>
	</div>
</div>

<style lang="less">
	.dialog-overlay {
		top: 0;
		left: 0;
		position: fixed;
		width: 100%;
		height: 100%;
		padding: 8px;
		display: flex;
		flex-flow: column nowrap;
		align-items: center;
		justify-content: center;
		background: var(--bg-overlay);
		z-index: var(--z-dialog);
	}

	.dialog {
		display: flex;
		position: relative;
		flex-flow: column nowrap;
		color: var(--fg);
		background: var(--bg);
		max-width: 400px;
		max-height: 80vh;
		border-radius: 8px;
		border: 1px solid var(--border);
	}

	.header {
		display: flex;
		flex-flow: row nowrap;
		align-items: center;
		border-bottom: 1px solid var(--border);
		padding: 0 8px;
		gap: 8px;
		flex-shrink: 0;

		.title {
			padding: 16px 0 15px;
			font-size: 1.2rem;
			overflow: hidden;
			text-align: center;
			font-variation-settings: 'wght' 600;
			white-space: nowrap;
			text-overflow: ellipsis;
			align-self: center;
			flex: 1 0;
		}
	}

	.content {
		padding: 16px;
		position: relative;
		flex: 1 0;

		:global(.loading-cover) {
			position: absolute;
			top: 0;
			left: 0;
			right: 0;
			bottom: 0;
			width: auto !important;
			height: auto !important;
			z-index: 1;
		}
	}

	.footer {
		padding: 8px;
		border-top: 1px solid var(--border);
		display: flex;
		flex-flow: row wrap;
		justify-content: center;
		flex-shrink: 0;

		&:empty {
			display: none;
		}
	}

	// Conditional styling
	.dialog.error {
		.header {
			background-color: var(--bg-error);
			border-top-left-radius: 8px;
			border-top-right-radius: 8px;

			:global(button) {
				background-color: transparent;
			}
		}
	}
</style>
