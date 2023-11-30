<script lang="ts">
	export let showModal: boolean;

	let dialog: HTMLDialogElement;
	$: showModal = showModal;
	$: dialog && showModal ? dialog.showModal() : dialog?.close();
</script>

<svelte:window on:keydown={(e) => {
	if(e.key === "Escape"){
		showModal = false
	}
}} />

<dialog
	bind:this={dialog}
	on:close={() => {
		showModal = false;
	}}
>
	<div>
			<button
				class="close"
				on:click={() => {
					dialog.close();
				}}><img draggable="false" src="/closeCircle.svg" alt="exit dialog" width="35px" height="35px" /></button
			>
		<slot />
	</div>
</dialog>

<style lang="scss">
	div {
		padding: 4rem;
		background-color: rgb(var(--bg-2));
	}
	.close {
		position: absolute;
		top: 0.75rem;
		right: 0.75rem;
		border-radius: 50%;
		&:focus {
			outline: 2px solid rgba(var(--primary), 0.5);
			outline-offset: -1px;
		}
	}
	dialog {
		padding: 0;
		border-radius: 4px;
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.07);
		color: inherit;
		&[open] {
			animation: fly-down 250ms ease-out;
		}
		&::backdrop {
			background-color: rgba(146, 146, 146, 0.5);
		}
		&[open]::backdrop {
			animation: fade 150ms ease-out;
		}
	}

	@keyframes fly-down {
		from {
			transform: translateY(-0.5rem);
		}
		to {
			transform: translateY(0);
		}
	}
</style>
