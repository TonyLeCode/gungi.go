<script>
	import { reverseList } from '$lib/helpers';

	export const boardState = new Array(81).fill(['']);
	const fileCoords = [1, 2, 3, 4, 5, 6, 7, 8, 9];
	const rankCoords = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i'];
</script>

<div class="board">
	{#each boardState as square}
		<div class="square">
			{#if square.length > 0}
				<img class="piece" src="/pieces/b1b.svg" alt="" />
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
	.board {
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
</style>
