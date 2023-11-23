<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Modal from '$lib/components/Modal.svelte';

	export let showModal: boolean;

	const dispatch = createEventDispatcher();

	function handleAccept() {
		dispatch("undoResponse", {response: "accept"})
		showModal = false;
	}
	function handleReject() {
		dispatch("undoResponse", {response: "reject"})
		showModal = false;
	}
</script>

<Modal bind:showModal backdropExit={false}>
	<p>Your opponent has requested<br> an undo</p>
	<div class="button-container">
		<button on:click={handleAccept} class="button-primary">Accept Undo</button>
		<button on:click={handleReject} class="button-primary">Reject Undo</button>
		<!-- <button on:click={handleAttack} class="button-primary" disabled={disableAttackDialogue}>Attack</button>
		<button on:click={handleStack} class="button-primary" disabled={disableStackDialogue}>Stack</button> -->
	</div>
</Modal>

<style lang="scss">
	p {
    font-size: 1.25rem;
    font-weight: 600;
		margin-bottom: 1rem;
		text-align: center;
	}
	.button-container {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}
</style>
