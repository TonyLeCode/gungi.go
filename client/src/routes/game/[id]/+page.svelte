<script lang="ts">
	import Board from './Board.svelte';
	import { setGameStore } from '$lib/store/gameState.svelte';
	import PlayerInfo from './PlayerInfo.svelte';
	import Menu from './Menu.svelte';
	let { data } = $props();

	const boardStore = setGameStore(data.gameData, data.session?.user.user_metadata.username);

	let selectedStack = $state<number[]>([]);

	function changeSelectedStack(stack: number[]) {
		selectedStack = stack;
	}
</script>

<svelte:head>
	<title>Game: {boardStore.player1} vs {boardStore.player2} | White Monarch Server</title>
</svelte:head>

<main>
	<section>
		<div class="game-state">
			{#if boardStore.completed}
				{#if boardStore.result === 'b'}
					Black Wins By Checkmate
				{:else if boardStore.result === 'w'}
					White Wins By Checkmate
				{:else if boardStore.result === 'b/r'}
					Black Wins By Resignation
				{:else if boardStore.result === 'w/r'}
					White Wins By Resignation
				{:else if boardStore.result === 'draw'}
					Draw
				{/if}
			{:else}
				<span class:turn-indicator={boardStore.isUserTurn}
					>{boardStore.isPlayer1Ready && boardStore.isPlayer2Ready ? '' : 'Drafting -'}
					{boardStore.turnColor === 'w' ? 'White' : 'Black'} To Play</span
				>
			{/if}
		</div>
		<PlayerInfo isOpposite={true} />
		<Board changeSelectedStack={changeSelectedStack} />
		<PlayerInfo isOpposite={false} />
		<!-- <div class="player same">
			<div class={`${boardStore.turnColor === 'w' ? 'w' : 'b'}`}></div>
			{boardStore.isViewReversed && boardStore.userColor === 'w' ? boardStore.player1 : boardStore.player2}
		</div> -->
	</section>
	<Menu selectedStack={selectedStack} />
</main>

<style lang="scss">
	main {
		max-width: 30rem;
		margin: 0 auto;
		font-size: 0.75rem;
		@media (min-width: 608px) {
			font-size: 0.875rem;
		}
		@media (min-width: 767px) {
			font-size: 1rem;
		}
		@media (min-width: 1200px) {
			display: flex;
			max-width: 90rem;
			gap: 2rem;
			padding: 0 2rem;
		}
	}

	section {
		width: 100%;
		user-select: none;
		margin: auto;
		max-width: 50rem;
	}

	.completed-text {
		font-size: 1.2rem;
		font-weight: 600;
	}

	.turn-indicator {
		color: rgb(var(--primary));
		font-weight: 600;
	}

	.divider {
		border-left: 1px rgba(99, 99, 99, 0.2) solid;
	}

	.game-state {
		text-align: center;
		margin-bottom: 0.25rem;
		@media (min-width: 767px) {
			margin: 0.25rem 0;
		}
	}
</style>
