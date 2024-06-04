<script lang="ts">
	import { draggable } from '$lib/store/dragAndDrop.svelte';
	import { getGameStore } from '$lib/store/gameState.svelte';
	import { DecodePiece, GetPieceColor } from '$lib/utils/utils';

	let {
		selectedStack,
		ready,
		placeHandMove,
	}: {
		selectedStack: number[];
		ready: () => void;
		placeHandMove: (fromPiece: number, fromCoord: number, toCoord: number) => void;
	} = $props();

	const boardStore = getGameStore();

	let imgStack = $derived.by(() => {
		return selectedStack.map((piece, i) => {
			const decodedPiece = DecodePiece(piece).toLowerCase();
			const color = GetPieceColor(piece);
			return `/pieces/${color}${i + 1}${decodedPiece}.svg`;
		});
	});

	// let whiteHandImgs = $derived.by(() => {
	// 	return boardStore.player1HandList.map((amount, piece) => {
	// 		const decodedPiece = DecodePiece(piece).toLowerCase();
	// 		return `/pieces/w1${decodedPiece}.svg`;
	// 	});
	// });

	// let blackHandImgs = $derived.by(() => {
	// 	return boardStore.player2HandList.map((amount, piece) => {
	// 		const decodedPiece = DecodePiece(piece).toLowerCase();
	// 		return `/pieces/b1${decodedPiece}.svg`;
	// 	});
	// });
	function getPieceImg(piece: number, color: string) {
		const decodedPiece = DecodePiece(piece).toLowerCase();
		return `/pieces/${color}1${decodedPiece}.svg`;
	}
</script>

{#snippet hand(label: string, handList: number[], armyCount: number, handCount: number, color: string)}
	<!-- TODO add badge and hide if 0 -->
	<div class="hand-info">
		<div class="label">
			<h3 class:is-user={label === 'Your Hand:'}>{label}</h3>
			<div class="count">{`Board: ${armyCount}  Hand: ${handCount}`}</div>
		</div>
		<div class="stack-container">
			{#each handList as amount, index}
				{#if amount != 0}
					<div class="hand">
						<img class="piece" draggable="false" use:draggable src={getPieceImg(index, color)} alt="" />
						{#if amount > 1}
							<img class="piece-under" draggable="false" src={getPieceImg(index, color)} alt="" />
						{/if}
						<div class="badge" title={String(amount)}>{amount}</div>
					</div>
				{/if}
			{/each}
		</div>
	</div>
{/snippet}

<div class="stack-details">
	<h3>Stack Details:</h3>
	<div class="stack-container details">
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
		boardStore.userColor === 'spectator' ? "White's Hand:" : 'Your Hand:',
		boardStore.userColor === 'b' ? boardStore.player2HandList : boardStore.player1HandList,
		boardStore.userColor === 'b' ? boardStore.player2ArmyCount : boardStore.player1ArmyCount,
		boardStore.userColor === 'b' ? boardStore.player2HandCount : boardStore.player1HandCount,
		'w'
	)}
	{@render hand(
		boardStore.userColor === 'spectator' ? "Black's Hand:" : "Opponent's Hand:",
		boardStore.userColor === 'b' ? boardStore.player1HandList : boardStore.player2HandList,
		boardStore.userColor === 'b' ? boardStore.player1ArmyCount : boardStore.player2ArmyCount,
		boardStore.userColor === 'b' ? boardStore.player1HandCount : boardStore.player2HandCount,
		'b'
	)}
</div>
<div class="buttons">
	<!-- TODO resign and request undo -->
	<button class="button-primary">resign</button>
	<button class="button-primary">request undo</button>
	{#if !boardStore.isPlayer1Ready || !boardStore.isPlayer2Ready}
		<button class="button-primary" onclick={ready}>ready</button>
	{/if}
	<button class="button-primary" onclick={() => (boardStore.manualFlip = !boardStore.manualFlip)}>flip board</button>
</div>

<style lang="scss">
	.not-selected {
		font-weight: 300;
	}

	.stack-container {
		position: relative;
		display: flex;
		justify-content: left;
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
		position: relative;
	}

	.details {
		justify-content: center;
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

	.badge {
		user-select: none;
		text-align: center;
		font-size: 0.7rem;
		background-color: rgb(var(--primary));
		color: white;
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
