<script lang="ts">
	import { browser } from '$app/environment';
	import { Droppable, draggable } from '$lib/store/dragAndDrop.svelte';
	import type { DraggableOptions } from '$lib/store/dragAndDrop.svelte';
	import { getGameStore } from '$lib/store/gameState.svelte';
	import { getSquareCoords } from '$lib/utils/historyParser';
	import { DecodePiece, GetImage, GetPieceColor, PieceIsPlayerColor, ReverseIndices } from '$lib/utils/utils';

	type DropItem = {
		destinationIndex: number;
		destinationStack: number[];
	};

	let {
		changeSelectedStack,
		promptMoveDialogue,
		movePiece,
		droppable = $bindable(),
	}: {
		changeSelectedStack: (stack: number[]) => void;
		promptMoveDialogue: (fromCoord: number, toCoord: number) => void;
		movePiece: (fromPiece: number, fromCoord: number, toCoord: number) => void;
		droppable: Droppable<DropItem>;
	} = $props();

	const boardStore = getGameStore();

	//TODO, deslect on new board state
	let startSelectIndex = $state(-1);
	let selectedSquareIndex = $state(-1);
	let selectedMoveIndices = $derived.by(() => {
		if (!boardStore.isUserTurn) return [];
		if (selectedSquareIndex === -1) return [];
		if (boardStore.moveListUI[selectedSquareIndex] === undefined) return [];
		return boardStore.moveListUI[selectedSquareIndex];
	});
	let lastMoveHighlightIndex = $derived.by(() => {
		const lastMove = getSquareCoords(boardStore.moveHistory[boardStore.moveHistory.length - 1]);
		return boardStore.isViewReversed ? ReverseIndices(lastMove) : lastMove;
	});

	$effect(() => {
		if (selectedSquareIndex !== -1) {
			changeSelectedStack(boardStore.boardUI[selectedSquareIndex]);
		} else {
			changeSelectedStack([]);
		}
	});

	let fileCoords = $derived(boardStore.isViewReversed ? [1, 2, 3, 4, 5, 6, 7, 8, 9] : [9, 8, 7, 6, 5, 4, 3, 2, 1]);
	let rankCoords = $derived(
		boardStore.isViewReversed
			? ['i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a']
			: ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i']
	);

	let deselectTimeout = $state<number | undefined>(undefined);
	let blockDeselect = $state(false);
	if (browser) {
		// Prevent deselection when clicking on the board
		window.addEventListener('mousedown', (e) => {
			if (deselectTimeout === undefined && selectedSquareIndex !== -1) {
				deselectTimeout = window.setTimeout(() => {
					if (!blockDeselect) {
						selectedSquareIndex = -1;
						deselectTimeout = undefined;
						blockDeselect = false;
					}
					deselectTimeout = undefined;
					blockDeselect = false;
				}, 50);
			}
		});
	}

	function blockDeselection() {
		blockDeselect = true;
	}

	function GetImage2(tier: number, piece: number): string {
		const encodedPiece = DecodePiece(piece).toLowerCase();
		const color = GetPieceColor(piece);
		return `/pieces/${color}${tier}${encodedPiece}.svg`;
	}

	function isActive(stack: number[]): boolean {
		if (!boardStore.isUserTurn) return false;
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

	function moveHandler(fromCoord: number, toCoord: number) {
		const toSquare = boardStore.boardUI[toCoord];
		if (toSquare.length === 0) {
			// is empty
			const fromSquare = boardStore.boardUI[fromCoord];
			const fromPiece = fromSquare[fromSquare.length - 1];
			movePiece(fromPiece, fromCoord, toCoord);
		} else {
			promptMoveDialogue(fromCoord, toCoord);
		}
	}

	function draggableOptions(index: number, stack: number[]): DraggableOptions<DropItem | null> {
		return {
			startEvent: () => {
				startSelectIndex = selectedSquareIndex;
				if (selectedSquareIndex === -1 || (!selectedMoveIndices.includes(index) && selectedSquareIndex !== index)) {
					selectSquareIndex(index);
				}
				blockDeselection();
			},
			dragStartEvent: () => {
				if (selectedSquareIndex !== index) {
					selectSquareIndex(index);
				}
			},
			dragReleaseEvent: (hoverItem) => {
				if (hoverItem === null || hoverItem === undefined) {
					selectSquareIndex(-1);
					return;
				}
				if (selectedMoveIndices.includes(hoverItem.destinationIndex)) {
					// make move
					moveHandler(index, hoverItem.destinationIndex);
				}
				selectSquareIndex(-1);
			},
			shortReleaseEvent: (hoverItem) => {
				if (startSelectIndex === index) {
					selectSquareIndex(-1);
				}
				// make move
				if (selectedSquareIndex !== -1 && selectedMoveIndices.includes(index)) {
					if (hoverItem !== null && hoverItem !== undefined) {
						moveHandler(selectedSquareIndex, hoverItem.destinationIndex);
					}
					selectSquareIndex(-1);
				}
			},
			longReleaseEvent: (hoverItem) => {
				if (startSelectIndex === index) {
					selectSquareIndex(-1);
				}
				// make move
				if (selectedSquareIndex !== -1 && selectedMoveIndices.includes(index)) {
					if (hoverItem !== null && hoverItem !== undefined) {
						moveHandler(selectedSquareIndex, hoverItem.destinationIndex);
					}
					selectSquareIndex(-1);
				}
			},
			releaseEvent: (hoverItem) => {
				startSelectIndex = -1;
			},
			droppable: droppable,
			active: () => {
				return isActive(stack);
			},
		};
	}
</script>

<div class="board">
	{#each boardStore.boardUI as stack, index (String(index) + JSON.stringify(stack))}
		<div
			role="button"
			tabindex="0"
			class="square"
			class:highlight={selectedSquareIndex === index || lastMoveHighlightIndex.includes(index)}
			class:move-highlight={selectedMoveIndices.includes(index) &&
				boardStore.isPlayer1Ready &&
				boardStore.isPlayer2Ready &&
				boardStore.userColor !== 'spectator'}
			onmousedown={() => {
				const stackLength = stack.length;
				const pieceIsPlayerColor = PieceIsPlayerColor(stack[stackLength - 1], boardStore.userColor);
				if (stackLength === 3 && pieceIsPlayerColor) {
					blockDeselection();
					selectedSquareIndex = index;
					return
				}
				if (boardStore.userColor === 'spectator') {
					blockDeselection();
					selectedSquareIndex = index;
					return;
				}
				if (selectedMoveIndices.includes(index)) {
					// make move
					// console.log('make move');
					moveHandler(selectedSquareIndex, index);
					selectSquareIndex(-1);
				} else if (stackLength > 0) {
					// console.log('cannot stack or attack on');
					blockDeselection();
					selectedSquareIndex = index;
				} else {
					// console.log('empty square');
					selectSquareIndex(-1);
				}
			}}
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
				<img draggable="false" class="piece-under" src={GetImage2(stack.length - 1, stack[stack.length - 2])} alt="" />
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
	* {
		touch-action: none;
	}
	.pointer {
		cursor: pointer;
		touch-action: none;
		-ms-touch-action: none;
	}

	.board {
		display: grid;
		grid-template-columns: repeat(9, minmax(20px, 1fr));
		grid-template-rows: repeat(9, minmax(20px, 1fr));
		gap: 2px;
		padding: 2px;
		max-width: 30rem;
		margin-left: auto;
		margin-right: auto;
		aspect-ratio: 1/1;
		background-color: rgb(235, 145, 84);
		position: relative;
		@media (min-width: 767px) {
			margin-bottom: 1.75rem;
			margin-top: 0.75rem;
			box-shadow:
				0px 7px 50px 5px rgba(230, 106, 5, 0.25),
				0px 5px 10px rgba(230, 106, 5, 0.25);
		}
		@media (min-width: 1200px) {
			max-width: 40rem;
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
		color: black;
		@media (min-width: 768px) {
			align-items: center;
			margin-left: -1rem;
			color: inherit;
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
		color: black;
		@media (min-width: 768px) {
			bottom: -1.5rem;
			text-align: center;
			color: inherit;
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
