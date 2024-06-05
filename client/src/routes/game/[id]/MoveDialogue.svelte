<script lang="ts">
	import Modal from '$lib/components/Modal.svelte';

	let {
		showModal = $bindable(),
		text,
		attackFn,
		stackFn,
	}: {
		showModal: boolean;
		text: string;
		attackFn: (() => void) | null;
		stackFn: (() => void) | null;
	} = $props();
</script>

<Modal bind:showModal>
	<p>{text}</p>
	<div class="button-container">
		<button class="button-primary" onclick={attackFn} disabled={attackFn === null}>Attack</button>
		<button class="button-primary" onclick={stackFn} disabled={stackFn === null}>Stack</button>
	</div>
</Modal>

<style lang="scss">
	p {
		margin-bottom: 1rem;
		text-align: center;
	}
	.button-container {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}
</style>
