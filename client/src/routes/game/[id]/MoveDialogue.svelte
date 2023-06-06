<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Modal from '$lib/components/Modal.svelte';

	export let showModal: boolean;
	export let text: string;
	export let disableAttackDialogue: boolean;
	export let disableStackDialogue: boolean;
	export let moveDialogueInfo: MoveType;

	interface MoveType {
		fromPiece: number;
		fromCoord: number;
		moveType: number;
		toCoord: number;
	}

	const dispatch = createEventDispatcher();

	function handleAttack() {
		dispatch('move', { ...moveDialogueInfo, moveType: 2 });
		showModal = false;
	}
	function handleStack() {
		dispatch('move', { ...moveDialogueInfo, moveType: 1 });
		showModal = false;
	}
</script>

<Modal bind:showModal>
	<p>{text}</p>
	<div class="button-container">
		<button on:click={handleAttack} class="button-primary" disabled={disableAttackDialogue}>Attack</button>
		<button on:click={handleStack} class="button-primary" disabled={disableStackDialogue}>Stack</button>
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
