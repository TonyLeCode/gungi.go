<script lang="ts">
	import Board from '$lib/components/Board.svelte';
	import PieceHand from '$lib/components/PieceHand.svelte';
	import { DecodePiece, FenToHand } from '$lib/utils/utils.js';
	import Hand from './Hand.svelte';
	import Chat from './Chat.svelte';
	import Replay from './Replay.svelte';

	export let data;
	// console.log(data?.params.id);
	console.log(data);
	// for(let i in hands[0]){
	// 	const encoded = DecodePiece(i).toLocaleLowerCase()
	// 	console.log(`/pieces/w1${encoded}.svg`)
	// }
	//FenToBoard on board size
	function countPiecesOnBoard(fen: string) {
		const pieces = fen.split(' ')[0];
		const matchW = pieces.match(/[A-Z]/g);
		const matchB = pieces.match(/[a-z]/g);
		return [matchW?.length, matchB?.length];
	}
	function reverseNameIfBlack(isBlack: boolean): string {
		return !isBlack ? data.data.player1 : data.data.player2;
	}

	let menuState = 0;
	// countPiecesOnBoard(data.data.current_state)
	$: hands = FenToHand(data.data.current_state);
	$: [onBoardWhite, onBoardBlack] = countPiecesOnBoard(data.data.current_state);
	$: turnColor = data.data.current_state.split(' ')[2];
	$: turnPlayer = turnColor === 'w' ? data.data.player1 : data.data.player2;
	$: playerColor = data.data.player1 === data.session?.user.user_metadata.username ? 'w' : 'b';
	$: console.log(playerColor);
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<main>
	<section>
		<Board gameData={data.data} reversed={playerColor !== 'w'} />
	</section>
	<aside class="side-menu">
		<div class="game-state">
			{turnColor === 'w' ? 'White' : 'Black'} To Play
		</div>
		<div class="tabs">
			<button class="button-ghost" on:click={() => {menuState = 0}}>hand</button>
			<button class="button-ghost" on:click={() => {menuState = 1}}>chat</button>
			<button class="button-ghost" on:click={() => {menuState = 2}}>move history</button>
		</div>
		{#if menuState === 0}
			<Hand
				{playerColor}
				player1={data.data.player1}
				player2={data.data.player2}
				{hands}
				{onBoardBlack}
				{onBoardWhite}
			/>
		{:else if menuState === 1}
			<Chat />
		{:else if menuState === 2}
			<Replay />
		{/if}
	</aside>
</main>

<style lang="scss">
	main {
		display: flex;
		max-width: 90rem;
		margin: 0 auto;
	}

	.tabs {
		/* background-color: red; */
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		/* gap: 1rem; */
		justify-content: center;
		/* margin-left: 10%; */
		padding: 1rem 0;
	}

	.game-state {
		text-align: center;
		margin-bottom: 1rem;
	}

	.side-menu {
		/* background-color: gray; */
		margin-left: auto;
		max-width: 36rem;
		width: 100%;
		// margin-top: auto;
		// margin-bottom: auto;
	}

	/* section {
		max-width: 70rem;
		margin: 0 auto;
		margin-top: 2rem;
		padding: 0 2rem;
		text-align: center;
	}

	.gameList {
		display: grid;
		gap: 1rem;
		grid-template-columns: repeat(auto-fit, 20rem);
		padding: 1rem;
		justify-content: center;
	} */
</style>
