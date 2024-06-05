<script lang="ts">
	import TooltipWrapper from '$lib/components/TooltipWrapper.svelte';
	import type { DraggableOptions, Droppable } from '$lib/store/dragAndDrop.svelte';
	import { getGameStore } from '$lib/store/gameState.svelte';
	import { DecodePiece, DecodePieceFull, GetPieceColor } from '$lib/utils/utils';
	import MenuHandPiece from './MenuHandPiece.svelte';

	type DropItem = {
		destinationIndex: number;
		destinationStack: number[];
	};

	let {
		selectedStack,
		ready,
		placeHandMove,
		droppable,
		selectHandPiece
	}: {
		selectedStack: number[];
		ready: () => void;
		placeHandMove: (fromPiece: number, toCoord: number) => void;
		droppable: Droppable<DropItem>;
		selectHandPiece: (piece: number) => void;
	} = $props();

	const boardStore = getGameStore();

	function getImgSrc(piece: number, i: number) {
		const decodedPiece = DecodePiece(piece).toLowerCase();
		const color = GetPieceColor(piece);
		return `/pieces/${color}${i + 1}${decodedPiece}.svg`;
	}

	function getText(piece: number, i: number) {
		const decodedPiece = DecodePieceFull(piece);
		const color = GetPieceColor(piece);
		return `Tier ${i + 1}: ${color === 'w' ? 'White' : 'Black'} ${decodedPiece} `;
	}
</script>

{#snippet hand(label: string, handList: number[], armyCount: number, handCount: number, color: string)}
	<div class="hand-info">
		<div class="label">
			<h3 class:is-user={label === 'Your Hand:'}>{label}</h3>
			<div class="count">{`Board: ${armyCount}  Hand: ${handCount}`}</div>
		</div>
		<div class="stack-container">
			{#each handList as amount, index}
				{#if amount != 0}
					<MenuHandPiece {placeHandMove} {selectHandPiece} {index} {color} {amount} {droppable} />
				{/if}
			{/each}
		</div>
	</div>
{/snippet}

<div class="stack-details">
	<h3>Stack Details:</h3>
	<div class="stack-container details">
		{#if selectedStack.length === 0}
			<span class="not-selected"> Stack not currently selected </span>
		{/if}
		{#each selectedStack as piece, i}
			<TooltipWrapper text={getText(piece, i)}>
				{#snippet children(createRef, interactionProps)}
					<img
						use:createRef
						{...interactionProps.getReferenceProps()}
						class="piece"
						draggable="false"
						src={getImgSrc(piece, i)}
						alt=""
					/>
				{/snippet}
			</TooltipWrapper>
		{/each}
	</div>
</div>
<div class="hand-container" class:reversed-hand-container={boardStore.manualFlip}>
	{@render hand(
		boardStore.userColor === 'b'
			? "Your Hand:"
			: boardStore.userColor === 'spectator'
				? "White's Hand:"
				: 'Your Hand:',
		boardStore.userColor === 'b' ? boardStore.player2HandList : boardStore.player1HandList,
		boardStore.userColor === 'b' ? boardStore.player2ArmyCount : boardStore.player1ArmyCount,
		boardStore.userColor === 'b' ? boardStore.player2HandCount : boardStore.player1HandCount,
		boardStore.userColor === 'b' ? 'b' : 'w'
	)}
	{@render hand(
		boardStore.userColor === 'b'
			? "Opponent's Hand:"
			: boardStore.userColor === 'spectator'
				? "Black's Hand:"
				: "Oponnent's Hand:",
		boardStore.userColor === 'b' ? boardStore.player1HandList : boardStore.player2HandList,
		boardStore.userColor === 'b' ? boardStore.player1ArmyCount : boardStore.player2ArmyCount,
		boardStore.userColor === 'b' ? boardStore.player1HandCount : boardStore.player2HandCount,
		boardStore.userColor === 'b' ? 'w' : 'b'
	)}
</div>
<div class="buttons">
	<!-- TODO resign and request undo -->
	<button class="button-primary">resign</button>
	<button class="button-primary">request undo</button>
	{#if !boardStore.isPlayer1Ready || !boardStore.isPlayer2Ready}
		<button class="button-primary" onclick={ready}>ready</button>
	{/if}
	<button class="button-primary" onclick={() => (boardStore.manualFlip = !boardStore.manualFlip)}>flip board</button>
</div>

<style lang="scss">
	.not-selected {
		font-weight: 300;
	}

	.stack-container {
		position: relative;
		display: flex;
		justify-content: left;
		gap: 0.75rem;
		min-height: 2rem;
		flex-wrap: wrap;
		@media (min-width: 767px) {
			gap: 1rem;
			min-height: 3rem;
		}
	}

	.hand-container {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		@media (min-width: 1200px) {
			flex-direction: column-reverse;
		}
	}

	.reversed-hand-container {
		@media (min-width: 1200px) {
			flex-direction: column;
		}
	}

	h3 {
		margin-bottom: 0.25rem;
	}

	.label {
		display: flex;
	}

	.count {
		margin-left: auto;
	}

	.details {
		justify-content: center;
	}

	.hand-info {
		gap: 1rem;
		// padding: 0.5rem 1rem;
		// padding-bottom: 1rem;
		@media (min-width: 767px) {
			gap: 2rem;
		}
	}

	.buttons {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		justify-content: center;
		padding: 1rem 0;
		@media (min-width: 767px) {
			gap: 1rem;
		}
	}

	.is-user {
		color: rgb(var(--primary));
		font-weight: 600;
	}

	.piece {
		position: relative;
		z-index: 2;
		width: 2rem;
		border-radius: 50%;
		@media (min-width: 767px) {
			width: 3rem;
		}
	}
</style>
