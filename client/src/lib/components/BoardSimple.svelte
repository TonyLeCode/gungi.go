<script lang="ts">
	import { reverseList } from '$lib/helpers';
	import { DecodePiece, FenToBoard, GetImage, GetPieceColor, GetTopStack } from '$lib/utils/utils';

	export let gameData;
	export let isTurn: boolean;
	// console.log(gameData.current_state);
	const boardState = FenToBoard(gameData.current_state);
	$: correctedBoardState = isTurn ? boardState : reverseList(boardState)
	// console.log(boardState);
</script>

<div class="board">
	{#each correctedBoardState as square}
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
		box-shadow: 0px 7px 15px rgba(230, 106, 5, 0.2);
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
