<script lang="ts">
	import { getGameStore } from '$lib/store/gameState.svelte';
	import { DecodePiece, GetPieceColor } from '$lib/utils/utils';

	let { selectedStack }: { selectedStack: number[] } = $props();

	const boardStore = getGameStore();

	let imgStack = $derived.by(() => {
		return selectedStack.map((piece, i) => {
			const decodedPiece = DecodePiece(piece).toLowerCase();
			const color = GetPieceColor(piece);
			return `/pieces/${color}${i + 1}${decodedPiece}.svg`;
		});
	});

	let whiteHandImgs = $derived.by(() => {
		return boardStore.player1HandList.map((amount, piece) => {
			const decodedPiece = DecodePiece(piece).toLowerCase();
			return `/pieces/w1${decodedPiece}.svg`;
		});
	});

  let blackHandImgs = $derived.by(() => {
    return boardStore.player2HandList.map((amount, piece) => {
      const decodedPiece = DecodePiece(piece).toLowerCase();
      return `/pieces/b1${decodedPiece}.svg`;
    });
  });
</script>

{#snippet hand(label: string, imgList: string[], armyCount: number, handCount: number)}
<div class="hand-info">
  <div class="label">
    <h3>{label}</h3>
    <div class="count">{`Board: ${armyCount}  Hand: ${handCount}`}</div>
  </div>
  <div class="stack-container hand">
    {#each imgList as pieceImg}
      <img class="piece" draggable="false" src={pieceImg} alt="" />
    {/each}
  </div>
</div>
{/snippet}

<aside class="side-menu">
	<div class="tabs"></div>
	<div class="stack-details">
		<h3>Stack Details:</h3>
		<div class="stack-container">
			{#each imgStack as stackImgString}
				<img class="piece" draggable="false" src={stackImgString} alt="" />
			{/each}
		</div>
		<div class="hand-container">
			{@render hand("Your Hand:", whiteHandImgs, boardStore.player1ArmyCount, boardStore.player1HandCount)}
      {@render hand("Opponent's Hand:", blackHandImgs, boardStore.player2ArmyCount, boardStore.player2HandCount)}
		</div>
	</div>
	<button onclick={() => (boardStore.manualFlip = !boardStore.manualFlip)}>reverse</button>
</aside>
a

<style lang="scss">
	.side-menu {
		margin: 0 auto;
    margin-top: 0.5rem;
		max-width: 44rem;
		width: 100%;
		padding: 0 2rem;
		@media (min-width: 1200px) {
      margin: auto 0;
			margin-left: auto;
			max-width: 36rem;
			padding: 0;
		}
	}

	.stack-container {
		display: flex;
		justify-content: center;
		gap: 0.75rem;
		min-height: 3rem;
		flex-wrap: wrap;
		@media (min-width: 767px) {
			gap: 1rem;
		}
	}

  .hand-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    @media (min-width: 1200px) {
      flex-direction: column-reverse;
    }
  }

  h3 {
    margin-bottom: 0.5rem;
  }

  .label {
    display: flex;
  }

  .count {
    margin-left: auto;
  }

	.hand {
		justify-content: left;
	}

	.piece {
		width: 2rem;
		border-radius: 50%;
		@media (min-width: 767px) {
			width: 3rem;
		}
	}

	.hand-info {
		gap: 1rem;
		// padding: 0.5rem 1rem;
		// padding-bottom: 1rem;
		@media (min-width: 767px) {
			gap: 2rem;
		}
	}
</style>
