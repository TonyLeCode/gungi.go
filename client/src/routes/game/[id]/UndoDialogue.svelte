<script lang="ts">
	import Modal from '$lib/components/Modal.svelte';

	let {
		showModal = $bindable(),
		undoReponse,
	}: {
		showModal: boolean;
		undoReponse: (response: 'accept' | 'reject') => void;
	} = $props();
</script>

<Modal bind:showModal>
	<p>Your opponent has requested an undo</p>
	<div class="button-container">
		<button onclick={() => undoReponse('accept')} class="button-primary">Accept Undo</button>
		<button onclick={() => undoReponse('reject')} class="button-primary">Reject Undo</button>
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
