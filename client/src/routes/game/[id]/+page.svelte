<script lang="ts">
	import Board from '$lib/components/Board.svelte';
	import PieceHand from '$lib/components/PieceHand.svelte';
	import { DecodePiece, FenToHand } from '$lib/utils/utils.js';

	export let data;
	// console.log(data?.params.id);
	console.log(data);
	// for(let i in hands[0]){
	// 	const encoded = DecodePiece(i).toLocaleLowerCase()
	// 	console.log(`/pieces/w1${encoded}.svg`)
	// }
	//FenToBoard on board size
	function countPiecesOnBoard(fen: string){
		const pieces = fen.split(' ')[0]
		const matchW = pieces.match(/[A-Z]/g)
		const matchB = pieces.match(/[a-z]/g)
		return [matchW?.length, matchB?.length]
	}
	function reverseNameIfBlack(isBlack:boolean): string{
		return !isBlack ? data.data.player1 : data.data.player2
	}
	// countPiecesOnBoard(data.data.current_state)
	$: hands = FenToHand(data.data.current_state);
	$: [onBoardWhite, onBoardBlack] = countPiecesOnBoard(data.data.current_state)
	$: turnColor = data.data.current_state.split(' ')[2]
	$: turnPlayer = turnColor === 'w' ? data.data.player1 : data.data.player2
	$: playerColor = data.data.player1 === data.session?.user.user_metadata.username ? 'w' : 'b'
	$: console.log(playerColor)
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
		<!-- <div class="tabs">
			<button>hand</button>
			<button>move history</button>
		</div> -->
		<div class="hands">
			<div class="hand-container">
				<div class="hand-info">
					<h3 class="name">{playerColor === 'w' ? data.data.player2 : data.data.player1}</h3>
					<div>On Board: <span>{playerColor === 'w' ? onBoardBlack : onBoardWhite}</span></div>
					<div>In Hand: <span>{hands[playerColor === 'w' ? 1 : 0].reduce((a,b) => { return a + b})}</span></div>
				</div>
				<div class="hand">
					{#each hands[playerColor === 'w' ? 1 : 0] as amount, i}
						{#if amount != 0}
							<PieceHand color={playerColor === 'w' ? 'b' : 'w'} piece={i} {amount} />
						{/if}
					{/each}
				</div>
			</div>
			<div class="hand-container">
				<div class="hand-info">
					<h3 class="name">{playerColor === 'w' ? data.data.player1 : data.data.player2}</h3>
					<span>On Board: {playerColor === 'w' ? onBoardWhite : onBoardBlack}</span>
					<span>In Hand: {hands[playerColor === 'w' ? 0 : 1].reduce((a,b) => { return a + b})}</span>
				</div>
				<div class="hand">
					{#each hands[playerColor === 'w' ? 0 : 1] as amount, i}
						{#if amount != 0}
							<PieceHand color={playerColor === 'w' ? 'w' : 'b'} piece={i} {amount} />
						{/if}
					{/each}
				</div>
			</div>
			<div class="buttons">
				<button>resign</button>
				<button>request undo</button>
				<button disabled>confirm move</button>
			</div>
		</div>
	</aside>
</main>

<style>
	main {
		display: flex;
		max-width: 90rem;
		margin: 0 auto;
	}

	.buttons {
		/* background-color: red; */
		display: flex;
		gap: 1rem;
		justify-content: center;
		padding: 1rem 0;
	}
	.tabs {
		/* background-color: red; */
		display: flex;
		/* gap: 1rem; */
		justify-content: center;
		/* margin-left: 10%; */
		padding: 1rem 0;
	}

	button{
		background-color: rgb(255, 77, 7);
		color: white;
		padding: .25rem .75rem;
		border-radius: 4px;
		box-shadow: 0px 5px 10px rgba(255, 77, 7, 0.308);
	}

	button:disabled{
		background-color: rgb(184, 184, 184);
		box-shadow: none;
	}

	.game-state {
		text-align: center;
		margin-bottom: 1rem;
	}

	.name {
		/* color: rgb(255, 77, 7); */
		margin-right: auto;
		margin-left: .5rem;
		position: relative;
	}
	.name::before{
		content: '';
		width: 15px;
		height: 15px;
		border-radius: 50%;
		background-color: rgb(255, 77, 7);
		display:block;
		position:absolute;
		left: -1.25rem;
		margin: auto;
		top: 0;
		bottom: 0;
	}

	.side-menu {
		/* background-color: gray; */
		margin-left: auto;
		max-width: 35rem;
		width: 100%;
		margin-top: auto;
		margin-bottom: auto;
	}

	.hand-container{
		display:flex;
		flex-direction: column;
		/* border: 2px solid rgba(255, 77, 7, 0.7); */
		background-color: rgb(245, 245, 245);
		/* background-color: rgb(255, 255, 255); */
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.07);
		border-radius: 16px;
		padding: 1rem 1rem;
	}

	.hands {
		display: grid;
		gap: 1rem;
		
	}
	.hand {
		/* background-color: rgb(199, 199, 199); */
		display: grid;
		grid-template-columns: repeat(auto-fit, 3rem);
		gap: 1rem;
		padding: .5rem .5rem;
	}

	.hand-info {
		display: flex;
		gap: 1rem;
		padding: .5rem 1rem;
		/* justify-content: space-between; */
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
