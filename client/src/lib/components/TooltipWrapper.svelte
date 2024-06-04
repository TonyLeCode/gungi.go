<script lang="ts">
	import {
		FloatingArrow,
		arrow,
		flip,
		offset,
		shift,
		useDismiss,
		useFloating,
		useHover,
		useInteractions,
		useRole,
		type UseInteractionsReturn,
	} from '@skeletonlabs/floating-ui-svelte';

	import type { Snippet } from 'svelte';

	type childrenType = [createRef: (node: HTMLElement) => void, interactionProps: UseInteractionsReturn];

	let { children, text }: { children: Snippet<childrenType>; text: string } = $props();

	let open = $state(false);

	let arrowRef: HTMLElement | null = $state(null);
	const floating = useFloating({
		get open() {
			return open;
		},
		onOpenChange: (v) => (open = v),
		placement: 'top',
		get middleware() {
			return [offset(10), flip(), shift(), arrowRef && arrow({ element: arrowRef })];
		},
	});
	const role = useRole(floating.context, { role: 'tooltip' });
	const hover = useHover(floating.context, { move: false });
	const dismiss = useDismiss(floating.context);
	const interactions = useInteractions([role, hover, dismiss]);

	function createRef(node: HTMLElement) {
		floating.elements.reference = node;
	}
</script>

{@render children(createRef, interactions)}

{#if open}
	<div
		class="tooltip"
		bind:this={floating.elements.floating}
		style={floating.floatingStyles}
		{...interactions.getFloatingProps}
	>
		{text}
		<FloatingArrow bind:ref={arrowRef} context={floating.context} fill="rgb(var(--primary))" />
	</div>
{/if}

<style lang="scss">
	.tooltip {
		// background-color: rgb(var(--bg-3));
		background-color: rgb(var(--primary));
		color: white;
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.25);
		width: max-content;
		z-index: 4;
		padding: 0.5rem 1.5rem;
	}
</style>
