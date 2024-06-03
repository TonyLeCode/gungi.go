<script lang="ts">
	import { CircleX } from 'lucide-svelte';

	// export let showModal: boolean;
	let { showModal = $bindable(), children }: { showModal: boolean, children: () => any } = $props();

	let dialog: HTMLDialogElement;
	$effect(() => {
		dialog && showModal ? dialog.showModal() : dialog?.close();
	})
</script>

<svelte:window
	on:keydown={(e) => {
		if (e.key === 'Escape') {
			showModal = false;
		}
	}}
/>

<dialog
	bind:this={dialog}
	onclose={() => {
		showModal = false;
	}}
>
	<div>
		<button
			class="close"
			onclick={() => {
				dialog.close();
			}}
		>
			<CircleX size="35px" />
		</button>
		{@render children()}
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
		color: rgb(var(--primary));
		&:focus {
			outline: 2px solid rgba(var(--primary), 0.5);
			outline-offset: 1px;
		}
		&:hover {
			color: rgb(var(--primary-3));
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
			background-color: rgba(var(--bg-4), 0.85);
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
