<script lang="ts">
	import { reverseList } from '$lib/helpers';
	import { EncodePiece, FenToBoard, GetPieceColor, GetTopStack } from '$lib/utils/utils';

	export let gameData;
	// console.log(gameData.current_state);
	const boardState = FenToBoard(gameData.current_state);
	// console.log(boardState);
	function GetImage(stack: number[]): string {
		const topPiece = GetTopStack(stack);
		const encodedPiece = EncodePiece(topPiece).toLowerCase();
		const color = GetPieceColor(topPiece);
		return `/pieces/${color}${stack.length}${encodedPiece}.svg`;
	}
</script>

<div class="board">
	{#each boardState as square}
		<div class="square">
			{#if square.length > 0}
				<img class="piece" draggable='false' src={GetImage(square)} alt="" />
			{/if}
		</div>
	{/each}
</div>

<style>
	.board {
		display: grid;
		grid-template-columns: repeat(9, minmax(20px, 1fr));
		grid-template-rows: repeat(9, minmax(20px, 1fr));
		gap: 2px;
		/* max-width: 45rem; */
		aspect-ratio: 1/1;
		background-color: rgb(226, 147, 49);
		padding: 2px;
		position: relative;
	}

	.square {
		background-color: rgb(254 215 170);
	}

	.piece {
		padding: 1px;
		user-select: none;
	}
</style>
