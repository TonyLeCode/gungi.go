<script lang="ts">
	import { DecodePiece, GetPieceColor } from '$lib/utils/utils';

	export let stack: number[] = [];
	$: imgStack = stack.map((piece, i) => {
		const decodedPiece = DecodePiece(piece).toLowerCase();
		const color = GetPieceColor(piece);
		return `/pieces/${color}${i + 1}${decodedPiece}.svg`;
	});
</script>

<div class="stack-container">
	<div>Stack Details</div>
	<div class="img-container">
		{#each imgStack as imgString}
			<img class="piece" draggable="false" src={imgString} alt="" />
		{/each}
	</div>
</div>

<style lang="scss">
	.stack-container {
		background-color: rgb(var(--bg-2));
		box-shadow: 0px 2px 55px rgba(0, 0, 0, 0.07), 0px 4px 15px rgba(0, 0, 0, 0.05);
		border-radius: 8px;
		padding: 1.5rem 2rem;
		margin-bottom: 1rem;
	}
  .img-container{
    display: flex;
		justify-content: center;
		gap: 1rem;
		min-height: 3rem;
  }
	.piece {
		width: 3rem;
		border-radius: 50%;
	}
</style>
