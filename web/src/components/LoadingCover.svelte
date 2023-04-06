<script lang="ts">
	// Props
	let className: string = '';
	export let style: string = '';
	export { className as class };
</script>

<div {style} class="loading-cover {className}" />

<style lang="less">
	.loading-cover {
		background: var(--bg);
		width: 100%;
		height: 100%;
		opacity: 0.8;
		overflow: hidden;
		z-index: 100;
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
	}
</style>
