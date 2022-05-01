<script lang="ts">
	import Icon from '@iconify/svelte/dist/offline';
	import type { IconifyIcon } from '@iconify/svelte';

	// Types
	type TooltipPlacement = 'top' | 'left' | 'right' | 'bottom';

	// Props
	let className: string = '';
	export let id: string = '';
	export let text: string = '';
	export let style: string = '';
	export let active: boolean = false;
	export let disabled: boolean = false;
	export let tooltip: string = '';
	export let tooltipPlacement: TooltipPlacement = 'right';
	export let icon: IconifyIcon[] | IconifyIcon | undefined = undefined;
	export let rightIcon: IconifyIcon[] | IconifyIcon | undefined = undefined;
	export { className as class };

	// Local variables
	let buttonElem: HTMLElement;
	let tooltipElem: HTMLElement;

	// Local function
	function getTooltipPosition(placement: TooltipPlacement): [number, number] {
		// Make sure tooltip defined
		if (tooltipElem == null) return [0, 0];

		// Get viewport size
		let viewportWidth = window.innerWidth;
		let viewportHeight = window.innerHeight;

		// Get button location and size
		let buttonRect = buttonElem.getBoundingClientRect();
		let buttonTop = buttonRect.y;
		let buttonLeft = buttonRect.x;
		let buttonWidth = buttonRect.width;
		let buttonHeight = buttonRect.height;
		let buttonVCenter = buttonLeft + buttonWidth / 2;
		let buttonHCenter = buttonTop + buttonHeight / 2;

		// Get tooltip size
		let tooltipWidth = tooltipElem.offsetWidth;
		let tooltipHeight = tooltipElem.offsetHeight;

		// Calculate tooltip position
		const margin = 8;
		let tooltipTop: number = 0;
		let tooltipLeft: number = 0;

		switch (placement) {
			case 'top':
				tooltipTop = buttonTop - margin - tooltipHeight;
				tooltipLeft = buttonVCenter - tooltipWidth / 2;
				if (tooltipTop < 0) {
					tooltipTop = buttonTop + buttonHeight + margin;
				}
				break;

			case 'left':
				tooltipTop = buttonHCenter - tooltipHeight / 2;
				tooltipLeft = buttonLeft - margin - tooltipWidth;
				if (tooltipLeft < 0) {
					tooltipLeft = buttonLeft + buttonWidth + margin;
				}
				break;

			case 'right':
				tooltipTop = buttonHCenter - tooltipHeight / 2;
				tooltipLeft = buttonLeft + buttonWidth + margin;
				if (tooltipLeft + tooltipWidth > viewportWidth) {
					tooltipLeft = buttonLeft - margin - tooltipWidth;
				}
				break;

			case 'bottom':
				tooltipTop = buttonTop + buttonHeight + margin;
				tooltipLeft = buttonVCenter - tooltipWidth / 2;
				if (tooltipTop + tooltipHeight > viewportHeight) {
					tooltipTop = buttonTop - margin - tooltipHeight;
				}
				break;
		}

		// Final adjustment
		if (placement === 'top' || placement === 'bottom') {
			if (tooltipLeft < 0) {
				tooltipLeft = margin;
			} else if (tooltipLeft + tooltipWidth > viewportWidth) {
				tooltipLeft = viewportWidth - tooltipWidth - margin;
			}
		} else if (placement === 'left' || placement === 'right') {
			if (tooltipTop < 0) {
				tooltipTop = margin;
			} else if (tooltipTop + tooltipHeight > viewportHeight) {
				tooltipTop = viewportHeight - tooltipHeight - margin;
			}
		}

		return [tooltipLeft, tooltipTop];
	}

	// Event handler
	function handleMouseEnter() {
		if (tooltipElem == null) return;
		let [left, top] = getTooltipPosition(tooltipPlacement);
		tooltipElem.style.top = `${top}px`;
		tooltipElem.style.left = `${left}px`;
		tooltipElem.style.opacity = '1';
	}

	function handleMouseLeave() {
		if (tooltipElem == null) return;
		tooltipElem.removeAttribute('style');
	}
</script>

<button
	{id}
	{style}
	{disabled}
	class:active
	class={className}
	on:click
	on:dblclick
	on:keydown
	on:mouseenter={handleMouseEnter}
	on:mouseleave={handleMouseLeave}
	bind:this={buttonElem}
>
	<slot>
		{#if icon != null}
			{#if !Array.isArray(icon)}
				<Icon {icon} height={18} />
			{:else}
				{#each icon as ic}
					<Icon icon={ic} height={18} />
				{/each}
			{/if}
		{/if}

		{#if text !== ''}
			<span class="text">{text}</span>
		{/if}

		{#if rightIcon != null}
			{#if !Array.isArray(rightIcon)}
				<Icon icon={rightIcon} height={18} />
			{:else}
				{#each rightIcon as ic}
					<Icon icon={ic} height={18} />
				{/each}
			{/if}
		{/if}
	</slot>

	{#if tooltip !== ''}
		<span class="tooltip" bind:this={tooltipElem}>{tooltip}</span>
	{/if}
</button>

<style lang="less">
	button {
		cursor: pointer;
		color: var(--fg);
		background-color: var(--bg);
		position: relative;
		display: flex;
		flex-flow: row nowrap;
		align-items: center;
		border: 1px solid transparent;
		min-height: 35px;
		padding: 0 8px;
		flex-shrink: 0;
		font-size: 1rem;
		text-align: left;
		gap: 8px;
		font-variation-settings: 'wght' 600;

		&:hover {
			color: var(--fg-hover);
			background-color: var(--bg-hover);
		}

		&:focus {
			border: 1px solid var(--border-focus);
		}

		&[disabled] {
			cursor: default;
			pointer-events: none;
			color: var(--fg-disabled);
		}

		&.active {
			color: var(--fg-active);
			background-color: var(--bg-active);
		}

		&.icon-right :global(svg) {
			order: 2;
		}
	}

	.text {
		flex: 1 0;
		padding: 8px 0;
	}

	.tooltip {
		padding: 8px;
		font-size: 0.9rem;
		color: var(--fg-tooltip);
		background-color: var(--bg-tooltip);
		border-radius: 8px;
		position: fixed;
		left: -100%;
		top: -100%;
		opacity: 0;
		z-index: var(--z-tooltip);
		pointer-events: none;
	}
</style>
