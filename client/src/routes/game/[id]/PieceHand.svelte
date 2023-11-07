<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { dragAndDropFunction, dragAndDropItems, dragAndDropOptions } from '$lib/utils/dragAndDrop';
	import { DecodePiece, DecodePieceFull, GetPieceColor, IndexToCoords, PieceIsPlayerColor } from '$lib/utils/utils';

	export let color: string;
	export let piece: number;
	export let amount: number;
	export let dragAndDrop: dragAndDropFunction;
	export let playerColor: string;
	export let isPlayerTurn: boolean;

	const dispatch = createEventDispatcher();

	const decodedPiece = DecodePiece(piece).toLowerCase();

	function isActive() {
		return color == playerColor && isPlayerTurn;
	}

	function dropEvent(items?: dragAndDropItems) {
		if (items?.hoverItem) {
			dispatch('drop', items);
		}
	}

	function dndOptions(piece: number) {
		let correctedPiece = piece;
		if (color == 'b') {
			correctedPiece = correctedPiece + 13;
		}
		const items = {
			piece: correctedPiece,
			from: 'hand',
		};
		return {
			releaseEvent: dropEvent,
			setDragItem: items,
			active: isActive,
		};
	}
</script>

<div class={`hand ${color === 'b' ? 'dark-hand' : ''}`}>
	<img
		class={`piece ${color == playerColor && isPlayerTurn ? 'piece-hover pointer' : ''}`}
		draggable="false"
		use:dragAndDrop={dndOptions(piece)}
		src={`/pieces/${color}1${decodedPiece}.svg`}
		alt=""
	/>
	{#if amount > 1}
		<img class="piece-under" draggable="false" src={`/pieces/${color}1${decodedPiece}.svg`} alt="" />
	{/if}
	<div class="badge" title={String(amount)}>
		{amount}
	</div>
</div>

<style>
	img {
		width: 3rem;
		border-radius: 50%;
		box-shadow: 0px 3px 10px rgba(0, 0, 0, 0.1);
	}
	.pointer {
		cursor: pointer;
	}
	.piece-hover:hover {
		outline: rgb(var(--primary)) 5px solid;
	}
	.piece {
		position: relative;
		z-index: 2;
	}
	.piece-under {
		position: absolute;
		left: 0;
		top: 0;
		right: 0;
		z-index: 1;
		user-select: none;
	}
	.hand {
		position: relative;
		--bg-color: rgb(var(--primary));
		--fill-color: rgb(var(--bg));
	}
	.dark-hand {
		--bg-color: rgb(var(--bg));
		--fill-color: rgb(var(--primary));
	}
	.badge {
		user-select: none;
		text-align: center;
		font-size: 0.7rem;
		background-color: var(--bg-color);
		color: var(--fill-color);
		border: 1px solid var(--fill-color);
		border-radius: 50%;
		width: 22px;
		height: 22px;
		display: flex;
		align-items: center;
		justify-content: center;
		position: absolute;
		right: -4px;
		top: -4px;
		z-index: 2;
	}
</style>
