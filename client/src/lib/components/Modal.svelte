<script lang="ts">
	import { CircleX } from 'lucide-svelte';

	// export let showModal: boolean;
	let { class: classname, children }: { class?: string; children: () => any } = $props();

	let dialog: HTMLDialogElement;

	export function open() {
		dialog.showModal();
	}

	export function close() {
		dialog.close();
	}
</script>

<svelte:window
	on:keydown={(e) => {
		if (e.key === 'Escape') {
			close();
		}
	}}
/>

<dialog bind:this={dialog}>
	<div class={classname}>
		<button
			class="close"
			onclick={() => {
				close();
			}}
		>
			<CircleX />
		</button>
		{@render children()}
	</div>
</dialog>

<style lang="scss">
	div {
		padding: 3rem;
		background-color: rgb(var(--bg-2));
		@media (min-width: 767px) {
			padding: 4rem;
		}
	}
	.close {
		position: absolute;
		top: 0.75rem;
		right: 0.75rem;
		border-radius: 50%;
		color: rgb(var(--primary));
		:global(svg) {
			width: 30px;
			height: 30px;
		}
		&:focus {
			outline: 2px solid rgba(var(--primary), 0.5);
			outline-offset: 1px;
		}
		&:hover {
			color: rgb(var(--primary-3));
		}
		@media (min-width: 767px) {
			:global(svg) {
				width: 35px;
				height: 35px;
			}
		}
	}
	dialog {
		padding: 0;
		border-radius: 2px;
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
