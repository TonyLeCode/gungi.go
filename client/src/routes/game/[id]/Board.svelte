<script lang="ts">
	import { Droppable, draggable } from '$lib/store/dragAndDrop.svelte.ts';
	import type { DraggableOptions } from '$lib/store/dragAndDrop.svelte.ts';
	import { getGameStore } from '$lib/store/gameState.svelte';
	import { getSquareCoords } from '$lib/utils/historyParser';
	import { DecodePiece, GetImage, GetPieceColor, PieceIsPlayerColor, ReverseIndices } from '$lib/utils/utils';

	const boardStore = getGameStore();

	let selectedSquareIndex = $state(-1);
	let selectedMoveIndices = $derived.by(() => {
		if (selectedSquareIndex === -1) return [];
		return boardStore.moveListUI[selectedSquareIndex]
	});
	let lastMoveHighlightIndex = $derived.by(() => {
		const lastMove = getSquareCoords(boardStore.moveHistory[boardStore.moveHistory.length - 1])
		return boardStore.isViewReversed ? ReverseIndices(lastMove) : lastMove
	})

	let fileCoords = $derived(boardStore.isViewReversed ? [1, 2, 3, 4, 5, 6, 7, 8, 9] : [9, 8, 7, 6, 5, 4, 3, 2, 1]);
	let rankCoords = $derived(
		boardStore.isViewReversed
			? ['i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a']
			: ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i']
	);

	function GetImage2(tier: number, piece: number): string {
		const encodedPiece = DecodePiece(piece).toLowerCase();
		const color = GetPieceColor(piece);
		return `/pieces/${color}${tier}${encodedPiece}.svg`;
	}

	function isActive(stack: number[]): boolean {
		const isPlayerPiece = PieceIsPlayerColor(stack[stack.length - 1], boardStore.userColor);
		const isDraftingPhase = !boardStore.isPlayer1Ready || !boardStore.isPlayer2Ready;

		return isPlayerPiece && !isDraftingPhase;
	}

	function selectSquareIndex(index: number) {
		if (selectedSquareIndex === index) {
			selectedSquareIndex = -1;
		} else {
			selectedSquareIndex = index;
		}
	}

	function draggableOptions(index: number, stack: number[]): DraggableOptions {
		return {
			startEvent: () => {
				selectSquareIndex(index);
			},
			dragReleaseEvent: () => {
				selectSquareIndex(-1);
			},
			active: () => {
				return isActive(stack);
			},
		};
	}

	type dropItem = {
		destinationIndex: number;
		destinationStack: number[];
	};
	const droppable = new Droppable<dropItem>();
</script>

<div class="board">
	{#each boardStore.boardUI as stack, index (String(index) + JSON.stringify(stack))}
		<div
			class="square"
			class:highlight={selectedSquareIndex === index || lastMoveHighlightIndex.includes(index)}
			class:move-highlight={selectedMoveIndices.includes(index) && boardStore.isPlayer1Ready && boardStore.isPlayer2Ready}
			use:droppable.addDroppable={{ mouseEnterItem: { destinationIndex: index, destinationStack: stack } }}
		>
			{#if stack.length > 0}
				<img
					draggable="false"
					use:draggable={draggableOptions(index, stack)}
					class={`piece ${isActive(stack) ? 'pointer' : ''}`}
					src={GetImage(stack)}
					alt=""
				/>
			{/if}

			{#if stack.length > 1}
				<img
					draggable="false"
					class="piece-under"
					src={GetImage2(stack.length - 1, stack[stack.length - 2])}
					alt=""
				/>
			{/if}
		</div>
	{/each}

	<div class="file">
		{#each fileCoords as file}
			<div>{file}</div>
		{/each}
	</div>
	<div class="rank">
		{#each rankCoords as rank}
			<div>{rank}</div>
		{/each}
	</div>
</div>

<style>
	.pointer {
		cursor: pointer;
	}

	.board {
		box-shadow:
			0px 7px 50px 5px rgba(230, 106, 5, 0.25),
			0px 5px 10px rgba(230, 106, 5, 0.25);
		display: grid;
		grid-template-columns: repeat(9, minmax(20px, 1fr));
		grid-template-rows: repeat(9, minmax(20px, 1fr));
		gap: 2px;
		padding: 2px;
		max-width: 40rem;
		margin-left: auto;
		margin-right: auto;
		aspect-ratio: 1/1;
		background-color: rgb(235, 145, 84);
		position: relative;
		@media (min-width: 767px) {
			margin-bottom: 1.75rem;
			margin-top: 0.75rem;
		}
	}

	.square {
		background-color: rgb(254 215 170);
		position: relative;
	}
	.square:hover::before {
		background-color: rgba(255, 131, 82, 0.2);
		border: 4px rgba(255, 131, 82, 0.5) solid;
		content: '';
		display: block;
		position: absolute;
		left: 0;
		right: 0;
		top: 0;
		bottom: 0;
	}

	.board .highlight::before {
		background-color: rgba(255, 131, 82, 0.432);
		border: 4px rgba(255, 131, 82, 0.705) solid;
		content: '';
		display: block;
		position: absolute;
		left: 0;
		right: 0;
		top: 0;
		bottom: 0;
	}

	.board .move-highlight::after {
		border-radius: 50%;
		content: '';
		display: block;
		position: absolute;
		left: 0;
		right: 0;
		top: 0;
		bottom: 0;
		background-color: rgba(228, 74, 3, 0.45);
		width: 25px;
		height: 25px;
		margin: auto;
		z-index: 2;
		user-select: none;
		pointer-events: none;
	}

	.file {
		display: grid;
		position: absolute;
		height: 100%;
		margin-left: 0.25rem;
		@media (min-width: 768px) {
			align-items: center;
			margin-left: -1rem;
		}
	}

	.rank {
		display: grid;
		position: absolute;
		grid-auto-flow: column;
		width: 100%;
		bottom: 0;
		text-align: right;
		right: 0.375rem;
		@media (min-width: 768px) {
			bottom: -1.5rem;
			text-align: center;
		}
	}

	.piece {
		position: relative;
		padding: 4px;
		z-index: 2;
		user-select: none;
		@media (min-width: 767px) {
			padding: 0.375rem;
		}
	}
	.piece-under {
		padding: 0.375rem;
		position: absolute;
		left: 0;
		top: 0;
		right: 0;
		z-index: 1;
		user-select: none;
		@media (min-width: 767px) {
			padding: 0.375rem;
		}
	}
</style>
