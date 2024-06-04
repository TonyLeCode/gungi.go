<script lang="ts">
	import TooltipWrapper from '$lib/components/TooltipWrapper.svelte';
	import { getGameStore } from '$lib/store/gameState.svelte';
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

	let text = $derived(isPlayerTurn ? "Current Player's turn" : 'Awaiting Turn');
</script>

<div class={`player ${oppositeClass}`}>
	<TooltipWrapper {text}>
		{#snippet children(createRef, interactionProps)}
			<div
				class={`indicator ${isPlayerTurn ? 'player-turn' : ''}`}
				use:createRef
				{...interactionProps.getReferenceProps()}
			></div>
		{/snippet}
	</TooltipWrapper>
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

	// .tooltip{
	//   background-color: rgb(var(--primary));
	//   padding: 0.5rem 1.5rem;
	// 	z-index: 6;
	// 	color: white;
	// }
</style>
