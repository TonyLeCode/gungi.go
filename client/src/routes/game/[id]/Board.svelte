<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { reverseList } from '$lib/helpers';
	import type { dragAndDropFunction, dragAndDropItems, dragAndDropOptions, dropFunction } from '$lib/utils/dragAndDrop';
	import { DecodePiece, FenToBoard, GetImage, GetPieceColor, PieceIsPlayerColor } from '$lib/utils/utils';

	// export const boardState = new Array(81).fill(['']);
	export let gameData;
	export let reversed: boolean;
	export let dragAndDrop: dragAndDropFunction;
	export let drop: dropFunction;
	export let playerColor: string;
	export let isPlayerTurn: boolean;

	const dispatch = createEventDispatcher();

	function GetImage2(tier: number, piece: number): string {
		const encodedPiece = DecodePiece(piece).toLowerCase();
		const color = GetPieceColor(piece);
		return `/pieces/${color}${tier}${encodedPiece}.svg`;
	}

	function reverseIfBlack<T>(arr: T[]): T[] {
		if (reversed) {
			return reverseList(arr);
		} else return arr;
	}
	$: boardState = reverseIfBlack(gameData);
	// $: console.log(boardState)
	$: fileCoords = reverseIfBlack([1, 2, 3, 4, 5, 6, 7, 8, 9]);
	$: rankCoords = reverseIfBlack(['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i']);

	function dropOptions(index: number, square: number[]) {
		let correctedIndex = index;
		let piece;
		if (square.length > 0) {
			piece = square[square.length - 1];
		}
		if (reversed) {
			correctedIndex = 80 - index;
		}

		const items = {
			coordIndex: correctedIndex,
			piece: piece,
			// stack: boardState[index],
		};
		return {
			mouseEnterItem: items,
		};
	}

	function dropEvent(items?: dragAndDropItems) {
		if (items?.hoverItem) {
			dispatch('drop', items);
		}
	}

	function dndOptions(index: number, piece: number) {
		let correctedIndex = index;
		function isActive() {
			return PieceIsPlayerColor(piece, playerColor) && isPlayerTurn;
		}
		if (reversed) {
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
</script>

<div class="board">
	{#each boardState as square, index}
		<div
			class="square"
			use:drop={dropOptions(index, square)}
			on:focus={() => {
				console.log('');
			}}
		>
			{#if square.length > 0}
				<img
					draggable="false"
					use:dragAndDrop={dndOptions(index, square[square.length - 1])}
					class={`piece ${PieceIsPlayerColor(square[square.length - 1], playerColor) && isPlayerTurn ? 'pointer' : ''}`}
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
		{#each reverseList(fileCoords) as file}
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
		box-shadow: 0px 7px 50px 5px rgba(230, 106, 5, 0.25);
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

	/* .square:hover::after {
		border-radius: 50%;
		content: '';
		display:block;
		position: absolute;
		left: 0;
		right: 0;
		top: 0;
		bottom: 0;
		background-color:rgba(228, 74, 3, 0.699);
		width: 25px;
		height: 25px;
		margin:auto;
		z-index: 2;
	} */

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
