<script lang="ts">
	import Dialog from '../components/Dialog.svelte';
	import { onMount, createEventDispatcher } from 'svelte';
	import { getRequest } from '../libs/api-request';
	const dispatch = createEventDispatcher();

	// Data type
	interface Ayah {
		id: number;
		arabic: string;
		translation: string;
		tafsir: string;
	}

	// Properties
	export let title: string = 'Tafsir dan Terjemah';
	export let ayah: number;
	export let surah: number;

	// Local variables
	let data: Ayah | undefined;
	let dataLoading: boolean = false;

	// API function
	async function loadData() {
		dataLoading = true;

		try {
			data = await getRequest(`/api/tafsir/surah/${surah}/ayah/${ayah}`);
		} catch (err) {
			dispatch('error', String(err));
		}

		dataLoading = false;
	}

	// Lifecycle function
	onMount(() => loadData());
</script>

<Dialog
	{title}
	loading={dataLoading}
	class="dialog-translation"
	on:close
	on:mainclick
>
	<div slot="content" class="tafsir-content">
		<p class="arabic">{data?.arabic || ''}</p>
		<div class="trans">{@html data?.translation || ''}</div>
		<div class="trans">{@html data?.tafsir || ''}</div>
	</div>
</Dialog>

<style lang="less">
	:global(.dialog-translation) {
		width: 850px;
		max-width: 90vw !important;
		max-height: 90vh !important;
	}

	.tafsir-content {
		margin: -16px 0;

		> * {
			padding: 16px 0;
		}

		> *:not(:last-child) {
			border-bottom: 1px solid var(--border);
		}

		.arabic {
			font-size: 2rem;
			font-family: 'KFGQPC-HAFS';
			text-align: center;
			color: var(--fg);
			direction: rtl;
		}

		.trans {
			font-size: 1rem;
			color: var(--fg);
			line-height: 2;

			:global(em) {
				font-style: italic;
				font-variation-settings: 'wght' 600;
			}

			> :global(*:not(:last-child)) {
				margin-bottom: 16px;
			}
		}
	}
</style>
