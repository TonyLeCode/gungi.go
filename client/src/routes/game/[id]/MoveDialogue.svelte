<script lang="ts">
	import Modal from '$lib/components/Modal.svelte';

	export let showModal: boolean;
	export let text: string;
	export let disableAttackDialogue: boolean;
	export let disableStackDialogue: boolean;
</script>

<Modal bind:showModal>
  <p>{text}</p>
  <div class="button-container">
    <button class="button-primary" disabled={disableAttackDialogue}>Attack</button>
    <button class="button-primary" disabled={disableStackDialogue}>Stack</button>
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
