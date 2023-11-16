<script lang="ts">
	import Modal from '$lib/components/Modal.svelte';
	type Info = {
		id: string;
		host: string;
		description: string;
		type: string;
		color: string;
		rules: string;
	};

	export let showModal: boolean;
	export let info: Info;
	export let accept: (roomid: string) => void;
</script>

<Modal bind:showModal>
	<div class="container">
		<h3>Challenge</h3>
		<div class="grid">
			<div class="label">Host:</div>
			<div class="name">{info?.host}</div>
			<div class="label">Description:</div>
			<div>
				{info?.description}
			</div>
			<div class="label">Type:</div>
			<div>
				{info?.type}
			</div>
			<div class="label">Host Color:</div>
			<div>
				{info?.color}
			</div>
			<div class="label">Rules:</div>
			<div>
				{info?.rules}
			</div>
		</div>
		<button
			on:click={() => {
				accept(info.id);
				showModal = false;
			}}
			class="button-primary">Accept</button
		>
	</div>
</Modal>

<style lang="scss">
	h3 {
		font-size: 1.25rem;
		font-weight: 600;
		margin-bottom: 1rem;
		text-align: center;
	}
	.grid {
		display: grid;
		grid-template-columns: 7rem 1fr;
		max-width: 30rem;
		word-break: break-all;
		gap: 0.75rem;
		margin-bottom: 1rem;
	}
	.container {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		min-width: 25rem;
		text-transform: capitalize;
	}
	.name {
		font-weight: 600;
	}
	.label {
		text-align: right;
	}
</style>
