<script lang="ts">
	import { createEventDispatcher, onDestroy } from 'svelte';
	import { reverseList } from '$lib/helpers';
	import type { dragAndDropFunction, dragAndDropItems, dragAndDropOptions, dropFunction } from '$lib/utils/dragAndDrop';
	import { DecodePiece, GetImage, GetPieceColor, PieceIsPlayerColor, ReverseIndex, ReverseIndices } from '$lib/utils/utils';
	import { get } from 'svelte/store';

	import {
		boardUIContext,
		completedContext,
		isViewReversedContext,
		isPlayer1ReadyContext,
		isPlayer2ReadyContext,
		isUserTurnContext,
		turnColorContext,
		userColorContext,
		moveListUIContext,
		moveListContext,
		moveHistoryContext,
	} from './+page.svelte';
	import { getSquareCoords } from '$lib/utils/historyParser';
	const boardUI = boardUIContext.get();
	const completed = completedContext.get();
	const userColor = userColorContext.get();
	const isViewReversed = isViewReversedContext.get();
	const isPlayer1Ready = isPlayer1ReadyContext.get();
	const isPlayer2Ready = isPlayer2ReadyContext.get();
	const isUserTurn = isUserTurnContext.get();
	const moveList = moveListContext.get();
	const turnColor = turnColorContext.get();
	const moveListUI = moveListUIContext.get();
	const moveHistory = moveHistoryContext.get();

	$: lastMoveHighlight = getSquareCoords($moveHistory[$moveHistory.length - 1])
	$: lastMoveHighlightUI = $isViewReversed ? ReverseIndices(lastMoveHighlight) : lastMoveHighlight

	export let dragAndDrop: dragAndDropFunction;
	export let drop: dropFunction;

	const dispatch = createEventDispatcher();

	function GetImage2(tier: number, piece: number): string {
		const encodedPiece = DecodePiece(piece).toLowerCase();
		const color = GetPieceColor(piece);
		return `/pieces/${color}${tier}${encodedPiece}.svg`;
	}

	function reverseArrayView<T>(arr: T[]): T[] {
		if ($isViewReversed) {
			return reverseList(arr);
		} else return arr;
	}
	let fileCoords = [9, 8, 7, 6, 5, 4, 3, 2, 1];
	let rankCoords = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i'];
	const unsubscribe = isViewReversed.subscribe((val) => {
		if (val) {
			fileCoords = [1, 2, 3, 4, 5, 6, 7, 8, 9];
			rankCoords = ['i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a'];
		} else {
			fileCoords = [9, 8, 7, 6, 5, 4, 3, 2, 1];
			rankCoords = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i'];
		}
	});

	onDestroy(unsubscribe);

	function dropOptions(index: number, square: number[]) {
		let correctedIndex = index;
		let piece;
		if (square.length > 0) {
			piece = square[square.length - 1];
		}
		if (get(isViewReversed)) {
			correctedIndex = 80 - index;
		}

		const items = {
			coordIndex: correctedIndex,
			piece: piece,
		};
		return {
			mouseEnterItem: items,
		};
	}

	function dropEvent(items?: dragAndDropItems) {
		if (items?.hoverItem) {
			dispatch('drop', items);
		}
		console.log(items);
		highlightIndex = -1;
		moveIndices = [];
	}

	function dndOptions(index: number, piece: number) {
		let correctedIndex = index;
		function isActive() {
			if (!get(isPlayer1Ready) || !get(isPlayer2Ready)) {
				return false;
			}
			return PieceIsPlayerColor(piece, get(userColor)) && get(isUserTurn) && !get(completed).completed;
		}
		if (get(isViewReversed)) {
			correctedIndex = 80 - index;
		}
		const square = {
			coordIndex: correctedIndex,
			piece: piece,
		};
		return {
			releaseEvent: dropEvent,
			setDragItem: square,
			active: isActive,
		};
	}

	let moveIndices: number[] = [];
	let highlightIndex: number;

	function onClick(index: number) {
		const board = get(boardUI)
		const square = board[index];

		if (get(userColor) != get(turnColor) || get(completed).completed) {
			highlightIndex = -1;
			moveIndices = [];
			return;
		}

		if (square.length === 0) {
			highlightIndex = -1;
			moveIndices = [];
			return;
		}

		if (GetPieceColor(square[square.length - 1]) != get(turnColor)) {
			highlightIndex = -1;
			moveIndices = [];
			return;
		}

		highlightIndex = index;
		moveIndices = get(moveListUI)[highlightIndex];

		// Fortress can't stack
		if (square[square.length - 1] % 13 === 4) {
			moveIndices = moveIndices.filter((moveIndex) => {
				const attackedSquare = board[moveIndex]
				return attackedSquare.length <= 1;
			})
		}
	}

	function moveHighlight(index: number): boolean {
		return $moveList[highlightIndex]?.includes(index);
	}

	function handleStackClick(index: number) {
		dispatch('stackClick', index);
	}
</script>

<div class="board">
	{#each $boardUI as square, index (String(index) + JSON.stringify(square))}
		<div
			role="button"
			tabindex="0"
			on:mousedown={() => {
				onClick(index);
			}}
			on:mousedown={() => {
				handleStackClick(index);
			}}
			class="square"
			class:highlight={highlightIndex == index || lastMoveHighlightUI.includes(index)}
			class:move-highlight={moveIndices?.includes(index) && $isPlayer1Ready && $isPlayer2Ready}
			use:drop={dropOptions(index, square)}
			on:focus={() => {
				console.log('');
			}}
		>
			{#if square.length > 0}
				<img
					draggable="false"
					use:dragAndDrop={dndOptions(index, square[square.length - 1])}
					class={`piece ${
						PieceIsPlayerColor(square[square.length - 1], $userColor) && $isUserTurn && !completed ? 'pointer' : ''
					}`}
					src={GetImage(square)}
					alt=""
				/>
				{#if square.length > 1}
					<img
						draggable="false"
						class="piece-under"
						src={GetImage2(square.length - 1, square[square.length - 2])}
						alt=""
					/>
				{/if}
			{/if}
		</div>
	{/each}

	<div class="file">
		{#each fileCoords as file}
			<div class="">{file}</div>
		{/each}
	</div>

	<div class="rank">
		{#each rankCoords as rank}
			<div class="">{rank}</div>
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
		max-width: 45rem;
		aspect-ratio: 1/1;
		background-color: rgb(235, 145, 84);
		padding: 2px;
		margin: 2rem;
		position: relative;
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
		margin-left: -1rem;
		align-items: center;
	}

	.rank {
		display: grid;
		position: absolute;
		grid-auto-flow: column;
		width: 100%;
		bottom: -1.5rem;
		text-align: center;
	}

	.piece {
		padding: 0.375rem;
		position: relative;
		z-index: 2;
		user-select: none;
	}
	.piece-under {
		padding: 0.375rem;
		position: absolute;
		left: 0;
		top: 0;
		right: 0;
		z-index: 1;
		user-select: none;
	}
</style>
