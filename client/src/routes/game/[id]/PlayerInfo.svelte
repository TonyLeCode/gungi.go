<script lang="ts">
	import { getGameStore } from '$lib/store/gameState.svelte';
	const boardStore = getGameStore();

	let { isOpposite }: { isOpposite: boolean } = $props();

  let displayName = $derived.by(() => {
    const turnPlayerName = boardStore.username;
    const oppositePlayerName = turnPlayerName === boardStore.player1 ? boardStore.player2 : boardStore.player1;
    return boardStore.isViewReversed === isOpposite ? oppositePlayerName : turnPlayerName
  })

  let isPlayerTurn = $derived.by(() => {
    let thisColor = displayName === boardStore.player1 ? 'w' : 'b';
    return boardStore.turnColor === thisColor
  })

  let oppositeClass = isOpposite ? 'opposite' : 'same';
</script>

<div class={`player ${oppositeClass}`}>
	<div class={`indicator ${isPlayerTurn ? 'player-turn' : ''}`}></div>
	{displayName}
  <div class="count">{`Board: ${boardStore.player1ArmyCount}  Hand: ${boardStore.player1HandCount}`}</div>
</div>

<style lang="scss">
	.player {
    display: flex;
		max-width: 45rem;
		margin: 0.5rem;
		@media (min-width: 767px) {
			margin: 0.5rem auto;
		}
		&.opposite {
			margin-top: 0;
		}
		&.same {
			margin-bottom: 0;
		}
	}
  .indicator {
    width: 1rem;
    height: 1rem;
    border-radius: 50%;
    margin-right: 0.5rem;
    border: 4px solid rgb(var(--primary));
    @media (min-width: 767px) {
      width: 1.25rem;
      height: 1.25rem;
    }
  }
  .player-turn {
    background-color: rgb(var(--primary));
  }
  .count {
    margin-left: auto;
  }
</style>
