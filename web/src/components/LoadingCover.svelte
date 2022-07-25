<script lang="ts">
	// Props
	let className: string = '';
	export let style: string = '';
	export { className as class };
</script>

<div {style} class="loading-cover {className}" />

<style lang="less">
	.loading-cover {
		background: var(--bg-secondary);
		width: 100%;
		height: 100%;
		opacity: 0.8;
		overflow: hidden;
		position: relative;
		pointer-events: none;
	}
</style>
