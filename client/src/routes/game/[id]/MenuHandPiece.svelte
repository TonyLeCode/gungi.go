<script lang="ts">
	import { DecodePiece, DecodePieceFull } from '$lib/utils/utils';
	import { draggable } from '$lib/store/dragAndDrop.svelte';
	import TooltipWrapper from '$lib/components/TooltipWrapper.svelte';

	let { index, color, amount }: { index: number; color: string; amount: number } = $props();

	function getPieceImg(piece: number, color: string) {
		const decodedPiece = DecodePiece(piece).toLowerCase();
		return `/pieces/${color}1${decodedPiece}.svg`;
	}

	let text = $derived(`${amount}x ${DecodePieceFull(index)}`);
</script>

<div class="hand">
	<TooltipWrapper {text}>
		{#snippet children(createRef, interactionProps)}
			<img
				use:createRef
				{...interactionProps.getReferenceProps()}
				class="piece"
				draggable="false"
				use:draggable
				src={getPieceImg(index, color)}
				alt=""
			/>
		{/snippet}
	</TooltipWrapper>
	{#if amount > 1}
		<img class="piece-under" draggable="false" src={getPieceImg(index, color)} alt="" />
	{/if}
	<div class="badge">{amount}</div>
</div>

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
