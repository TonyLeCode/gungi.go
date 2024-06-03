<script lang="ts">
	import { getGameStore } from '$lib/store/gameState.svelte';
	import { FloatingArrow, arrow, flip, offset, shift, useDismiss, useFloating, useHover, useInteractions, useRole } from '@skeletonlabs/floating-ui-svelte';
	const boardStore = getGameStore();

	let { isOpposite }: { isOpposite: boolean } = $props();

	let displayName = $derived.by(() => {
		return isOpposite === boardStore.isViewReversed ? boardStore.player1 : boardStore.player2;
	});

	let isPlayerTurn = $derived.by(() => {
		let thisColor = displayName === boardStore.player1 ? 'w' : 'b';
		return boardStore.turnColor === thisColor;
	});

	let oppositeClass = isOpposite ? 'opposite' : 'same';

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
    }
	});
	const role = useRole(floating.context, { role: 'tooltip' });
	const hover = useHover(floating.context, { move: false });
	const dismiss = useDismiss(floating.context);
	const interactions = useInteractions([role, hover, dismiss]);
</script>

<div class={`player ${oppositeClass}`}>
	{#if open}
		<div class="tooltip" bind:this={floating.elements.floating} style={floating.floatingStyles} {...interactions.getFloatingProps}>
			{isPlayerTurn ? "Current Player's turn" : 'Awaiting Turn'}
      <FloatingArrow bind:ref={arrowRef} context={floating.context} fill="rgb(var(--bg-2))" />
		</div>
	{/if}
	<div
		class={`indicator ${isPlayerTurn ? 'player-turn' : ''}`}
		bind:this={floating.elements.reference}
		{...interactions.getReferenceProps()}
	></div>
	<span class:is-user={displayName === boardStore.username}>{displayName}</span>
	<!-- <div class="count">{`Board: ${boardStore.player1ArmyCount}  Hand: ${boardStore.player1HandCount}`}</div> -->
</div>

<style lang="scss">
	.player {
		display: flex;
		max-width: 45rem;
		margin: 0.5rem;
		@media (min-width: 767px) {
			margin: 0.5rem auto;
		}
		&.opposite {
			margin-top: 0;
		}
		&.same {
			margin-bottom: 0;
		}
	}
	.is-user {
		color: rgb(var(--primary));
		font-weight: 600;
	}
	.indicator {
		width: 1rem;
		height: 1rem;
		border-radius: 50%;
		margin-right: 0.5rem;
		border: 4px solid rgb(var(--primary));
		@media (min-width: 767px) {
			width: 1.25rem;
			height: 1.25rem;
		}
	}
	.player-turn {
		background-color: rgb(var(--primary));
	}
	.count {
		margin-left: auto;
	}

  .tooltip{
    background-color: rgb(var(--bg-2));
    padding: 0.25rem 0.5rem;
		z-index: 6;
  }
</style>
