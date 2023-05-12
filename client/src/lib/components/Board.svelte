<script lang="ts">
	import { reverseList } from '$lib/helpers';
  import { FenToBoard, FenToHand, GetImage } from '$lib/utils/utils';

	// export const boardState = new Array(81).fill(['']);
	export let gameData;
	export let reversed: boolean;

	function reverseIfBlack<T>(arr:T[]):T[]{
		if(reversed){
			return reverseList(arr)
		} else return arr
	}
	let boardState = reverseIfBlack(FenToBoard(gameData.current_state));
	let fileCoords = reverseIfBlack([1, 2, 3, 4, 5, 6, 7, 8, 9]);
	let rankCoords = reverseIfBlack(['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i']);
</script>

<div class="board">
	{#each boardState as square}
		<div class="square">
			{#if square.length > 0}
			<img class="piece" src={GetImage(square)} alt="" />
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
	img{
		/* box-shadow: 0px 7px 15px rgba(230, 106, 5, 0.527); */
		/* border-radius: 50%; */
	}

	.board {
		box-shadow: 0px 7px 50px 5px rgba(230, 106, 5, 0.25);
		display: grid;
		grid-template-columns: repeat(9, minmax(20px, 1fr));
		grid-template-rows: repeat(9, minmax(20px, 1fr));
		gap: 2px;
		max-width: 45rem;
		aspect-ratio: 1/1;
		background-color: rgb(226, 147, 49);
		padding: 2px;
		margin: 2rem;
		position: relative;
	}

	.square {
		background-color: rgb(254 215 170);
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
	}
	.piece:hover {
		background-color: rgba(255, 77, 7, 0.479);
	}
</style>
