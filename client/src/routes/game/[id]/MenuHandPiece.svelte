<script lang="ts">
	import { DecodePiece, DecodePieceFull } from '$lib/utils/utils';
	import { Droppable, draggable, type DraggableOptions } from '$lib/store/dragAndDrop.svelte';
	import TooltipWrapper from '$lib/components/TooltipWrapper.svelte';
	import { getGameStore } from '$lib/store/gameState.svelte';

	type DropItem = {
		destinationIndex: number;
		destinationStack: number[];
	};

	let {
		index,
		color,
		amount,
		droppable,
		selectHandPiece,
    placeHandMove
	}: {
		index: number;
		color: string;
		amount: number;
		droppable: Droppable<DropItem>;
		selectHandPiece: (piece: number) => void;
    placeHandMove: (fromPiece: number, toCoord: number) => void;
	} = $props();

	let boardStore = getGameStore();

	function getPieceImg(piece: number, color: string) {
		const decodedPiece = DecodePiece(piece).toLowerCase();
		return `/pieces/${color}1${decodedPiece}.svg`;
	}

	let text = $derived(`${amount}x ${DecodePieceFull(index)}`);

	function isActive() {
		if (!boardStore.isUserTurn) return false;
		return boardStore.userColor === color;
	}

	function draggableOptions(): DraggableOptions<DropItem | null> {
		return {
      startEvent: (hoverItem) => {
        selectHandPiece(index)
      },
			dragReleaseEvent: (hoverItem) => {
        if (hoverItem !== null && hoverItem !== undefined) {
          //make placement
          let piece = index
          if (boardStore.userColor === 'b') {
            piece += 13
          }
          placeHandMove(piece, hoverItem.destinationIndex)
        }
        selectHandPiece(-1)
      },
			droppable: droppable,
			active: isActive,
		};
	}
	function setattr(node: HTMLElement) {
		node.setAttribute('draggable', 'false');
	}
</script>

<button
	class="hand"
  disabled={!isActive()}
	onmousedown={() => {
    // if (color !== boardStore.userColor) return;
		selectHandPiece(index);
	}}
>
	<TooltipWrapper {text}>
		{#snippet children(createRef, interactionProps)}
			<img
				use:createRef
				class="piece"
				draggable="false"
				use:draggable={draggableOptions()}
				src={getPieceImg(index, color)}
				alt=""
				{...interactionProps.getReferenceProps()}
				use:setattr
			/>
		{/snippet}
	</TooltipWrapper>
	{#if amount > 1}
		<img class="piece-under" draggable="false" src={getPieceImg(index, color)} alt="" />
	{/if}
	<div class="badge">{amount}</div>
</button>

<style lang="scss">
	.hand {
		position: relative;
	}

	.piece {
		position: relative;
		z-index: 2;
		width: 2rem;
		border-radius: 50%;
		@media (min-width: 767px) {
			width: 3rem;
		}
	}

	.piece-under {
		position: absolute;
		left: 0;
		top: 0;
		right: 0;
		z-index: 1;
		-moz-user-select: none;
		-webkit-user-select: none;
		user-select: none;
	}

	.badge {
		user-select: none;
		pointer-events: none;
		text-align: center;
		font-size: 0.8rem;
		background-color: rgb(var(--primary));
		color: white;
		border-radius: 50%;
		width: 20px;
		height: 20px;
		display: flex;
		align-items: center;
		justify-content: center;
		position: absolute;
		right: -6px;
		top: -6px;
		z-index: 2;
		@media (min-width: 767px) {
			right: -4px;
			top: -4px;
			width: 24px;
			height: 24px;
			font-size: 0.875rem;
		}
	}
</style>
