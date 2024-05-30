<script lang="ts">
	import { Droppable, draggable } from '$lib/store/dragAndDrop.svelte.ts';
	import { getGameStore } from '$lib/store/gameState.svelte';
	import { GetImage } from '$lib/utils/utils';

	const boardStore = getGameStore();

	let fileCoords = [9, 8, 7, 6, 5, 4, 3, 2, 1];
	let rankCoords = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i'];

	$effect(() => {
		if (boardStore.isViewReversed) {
			fileCoords = [1, 2, 3, 4, 5, 6, 7, 8, 9];
			rankCoords = ['i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a'];
		}
	});

	type dropItem = {
		destinationIndex: number;
		destinationStack: number[];
	};
	const droppable = new Droppable<dropItem>();
</script>

<div class="board">
	{#each boardStore.boardUI as square, index (String(index) + JSON.stringify(square))}
		<div
			class="square"
			use:droppable.addDroppable={{ mouseEnterItem: { destinationIndex: index, destinationStack: square } }}
		>
			{#if square.length > 0}
				<img draggable="false" use:draggable class="piece" src={GetImage(square)} alt="" />
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
