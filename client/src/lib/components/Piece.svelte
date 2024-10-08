<script lang="ts">
	import type { DraggableOptions } from '$lib/store/dragAndDrop.svelte';
	import { DecodePiece, DecodePieceFull, GetPieceColor } from '$lib/utils/utils';
	import type { Action } from 'svelte/action';

	type DropItem = {
		destinationIndex: number;
		destinationStack: number[];
	};

	let {
		piece,
		tier,
		class: className,
		useAction = () => {},
		useOptions = () => ({}),
		option = { index: -1, stack: [] },
	}: {
		piece: number;
		tier: number;
		class: string;
		useAction?: Action<HTMLElement, DraggableOptions<DropItem | null>>;
		useOptions?: (index: number, stack: number[]) => DraggableOptions<DropItem | null>;
		option?: { index: number; stack: number[] };
	} = $props();

	const encodedPiece = DecodePiece(piece).toLowerCase();
	const color = GetPieceColor(piece);
	const src = `/pieces/${color}${tier}${encodedPiece}.svg`;
</script>

<img
	draggable="false"
	class={className}
	{src}
	alt={DecodePieceFull(piece)}
	use:useAction={useOptions(option.index, option.stack)}
/>

<style lang="scss">
	img {
		touch-action: none;
	}
</style>
