<script lang="ts">
	import Piece from '$lib/components/Piece.svelte';
	import type { DraggableOptions } from '$lib/store/dragAndDrop.svelte';
	import type { Action } from 'svelte/action';

	type DropItem = {
		destinationIndex: number;
		destinationStack: number[];
	};
	
	type dropOptions<T> = {
		mouseEnterEvent?: () => void;
		mouseLeaveEvent?: () => void;
		mouseEnterItem?: T;
	};

	let {
		boardState,
		isReversed, // Only for coords, does not influence board or events
		showCoord = false,
		highlight = [],
		selectedMoves = [],
		dropAction = () => {},
		dragAction = () => {},
		draggableOptions = () => ({}),
		onMouseDown = () => {},
	}: {
		boardState: number[][];
		isReversed: boolean;
		highlight?: number[];
		selectedMoves?: number[];
		showCoord?: boolean;
		dropAction?: Action<HTMLElement, dropOptions<DropItem>>;
		dragAction?: Action<HTMLElement, DraggableOptions<DropItem | null>>;
		draggableOptions?: (index: number, stack: number[]) => DraggableOptions<DropItem | null>;
		onMouseDown?: (index: number) => void;
	} = $props();

	let fileCoords = $derived(isReversed ? [1, 2, 3, 4, 5, 6, 7, 8, 9] : [9, 8, 7, 6, 5, 4, 3, 2, 1]);
	let rankCoords = $derived(
		isReversed ? ['i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a'] : ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i']
	);
</script>

<div class="board">
	{#each boardState as stack, index (String(index) + JSON.stringify(stack))}
		<div
			role="button"
			tabindex="0"
			class="square"
			class:highlight={highlight.includes(index)}
			class:move-highlight={selectedMoves.includes(index)}
			onmousedown={() => {
				onMouseDown(index);
			}}
			use:dropAction={{
				mouseEnterItem: {
					destinationIndex: index,
					destinationStack: stack,
				},
			}}
		>
			{#if stack.length > 0}
				<Piece
					piece={stack[stack.length - 1]}
					tier={stack.length}
					class={`board-piece`}
					useAction={dragAction}
					useOptions={draggableOptions}
					option={{ index, stack }}
				/>
			{/if}

			{#if stack.length > 1}
				<Piece piece={stack[stack.length - 2]} tier={stack.length - 1} class="board-piece-under" />
			{/if}
		</div>
	{/each}

	{#if showCoord}
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
	{/if}
</div>

<style>
	* {
		touch-action: none;
	}

	:global(.pointer) {
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

	:global(.board-piece) {
		position: relative;
		padding: 4px;
		z-index: 2;
		user-select: none;
		@media (min-width: 767px) {
			padding: 0.375rem;
		}
	}
	:global(.board-piece-under) {
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
