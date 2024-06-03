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
	<!-- TODO add badge and hide if 0 -->
	<div class="hand-info">
		<div class="label">
			<h3 class:is-user={label === "Your Hand:"}>{label}</h3>
			<div class="count">{`Board: ${armyCount}  Hand: ${handCount}`}</div>
		</div>
		<div class="stack-container hand">
			{#each imgList as pieceImg}
				<img class="piece" draggable="false" src={pieceImg} alt="" />
			{/each}
		</div>
	</div>
{/snippet}

<div class="stack-details">
	<h3>Stack Details:</h3>
	<div class="stack-container">
		{#if imgStack.length === 0}
			<span class="not-selected"> Stack not currently selected </span>
		{/if}
		{#each imgStack as stackImgString}
			<img class="piece" draggable="false" src={stackImgString} alt="" />
		{/each}
	</div>
</div>
<div class="hand-container" class:reversed-hand-container={boardStore.manualFlip}>
  {@render hand(
    boardStore.userColor === 'spectator' ? "White's Hand:" : "Your Hand:", 
    boardStore.userColor === 'b' ? blackHandImgs : whiteHandImgs,
    boardStore.userColor === 'b' ? boardStore.player2ArmyCount : boardStore.player1ArmyCount, 
    boardStore.userColor === 'b' ? boardStore.player2HandCount : boardStore.player1HandCount)}
  {@render hand(
    boardStore.userColor === 'spectator' ? "Black's Hand:" : "Opponent's Hand:", 
    boardStore.userColor === 'b' ? whiteHandImgs : blackHandImgs, 
    boardStore.userColor === 'b' ? boardStore.player1ArmyCount : boardStore.player2ArmyCount, 
    boardStore.userColor === 'b' ? boardStore.player1HandCount : boardStore.player2HandCount)}
</div>
<div class="buttons">
	<button class="button-primary">resign</button>
	<button class="button-primary">request undo</button>
	{#if !boardStore.isPlayer1Ready || !boardStore.isPlayer2Ready}
		<button class="button-primary">ready</button>
	{/if}
	<button class="button-primary" onclick={() => (boardStore.manualFlip = !boardStore.manualFlip)}>flip board</button>
</div>

<style lang="scss">
	.not-selected {
		font-weight: 300;
	}

	.stack-container {
		display: flex;
		justify-content: center;
		gap: 0.75rem;
		min-height: 2rem;
		flex-wrap: wrap;
		@media (min-width: 767px) {
			gap: 1rem;
			min-height: 3rem;
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

	.reversed-hand-container {
		@media (min-width: 1200px) {
			flex-direction: column;
		}
	}

	h3 {
		margin-bottom: 0.25rem;
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

	.buttons {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		justify-content: center;
		padding: 1rem 0;
		@media (min-width: 767px) {
			gap: 1rem;
		}
	}

	.is-user {
		color: rgb(var(--primary));
		font-weight: 600;
	}
</style>
