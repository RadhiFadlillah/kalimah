<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { quintOut } from 'svelte/easing';
	import { fade } from 'svelte/transition';

	// Props
	let className: string = '';
	export let style: string = '';
	export let padding: number = 0;
	export { className as class };

	// Local variables
	let canvas: HTMLCanvasElement;

	let timer: NodeJS.Timer;
	let mainColor: string = '';
	let clientWidth: number = 0;
	let clientHeight: number = 0;
	let opacityTracker: Record<string, number> = {};
	let fontSize = 14 + padding * 2;

	// Redraw whenever size changed
	$: {
		if (canvas != null) {
			canvas.width = clientWidth;
			canvas.height = clientHeight;
			draw();
		}
	}

	function draw() {
		// Make sure canvas exist
		if (canvas == null || mainColor === '') return;

		// Get canvas context
		let ctx = canvas.getContext('2d');
		if (ctx == null) return;

		// Clear current canvas
		ctx.clearRect(0, 0, canvas.width, canvas.height);

		// Calculate count of rows and columns
		let nRows = Math.ceil(canvas.height / fontSize);
		let nColumns = Math.ceil(canvas.width / fontSize);

		// Draw squares if necessary
		for (let column = 0; column < nColumns; column++) {
			for (let row = 0; row < nRows; row++) {
				// Prepare coordinates
				let key = `${column}x${row}`;
				let x = column * fontSize;
				let y = row * fontSize;

				// Check if we need to regenerate opacity
				let currentOpacity = opacityTracker[key];
				if (currentOpacity == null || Math.random() > 0.8) {
					currentOpacity = Math.random();
					opacityTracker[key] = currentOpacity;
				}

				// Draw square
				ctx.fillStyle = mainColor;
				ctx.globalAlpha = currentOpacity;
				ctx.fillRect(x, y, fontSize, fontSize);
			}
		}
	}

	// Lifecycle
	onMount(() => {
		timer = setInterval(draw, 100);

		let css = getComputedStyle(document.documentElement);
		mainColor = css.getPropertyValue('--border');
	});

	onDestroy(() => {
		if (timer != null) clearInterval(timer);
	});
</script>

<div
	{style}
	class="loading-cover {className}"
	bind:clientWidth
	bind:clientHeight
	transition:fade={{ easing: quintOut }}
>
	<canvas bind:this={canvas} />
</div>

<style lang="less">
	.loading-cover {
		background: var(--bg);
		width: 100%;
		height: 100%;
		overflow: hidden;
		position: relative;
	}
</style>
