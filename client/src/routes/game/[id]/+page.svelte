<script lang="ts">
	import Board from './Board.svelte';
	import { DecodePiece, DecodePieceFull, FenToHand, GetPieceColor, IndexToCoords } from '$lib/utils/utils.js';
	import Hand from './Hand.svelte';
	import Chat from './Chat.svelte';
	import Replay from './Replay.svelte';
	import { dragAndDrop, drop } from '$lib/utils/dragAndDrop';
	import MoveDialogue from './MoveDialogue.svelte';

	export let data;
	// console.log(data?.params.id);
	console.log(data);
	// for(let i in hands[0]){
	// 	const encoded = DecodePiece(i).toLocaleLowerCase()
	// 	console.log(`/pieces/w1${encoded}.svg`)
	// }
	//FenToBoard on board size
	let showMoveDialogue = false;

	function countPiecesOnBoard(fen: string) {
		const pieces = fen.split(' ')[0];
		const matchW = pieces.match(/[A-Z]/g);
		const matchB = pieces.match(/[a-z]/g);
		return [matchW?.length ?? 0, matchB?.length ?? 0];
	}
	function reverseNameIfBlack(isBlack: boolean): string {
		return !isBlack ? data.data.player1 : data.data.player2;
	}

	let moveDialogueText = '';
	let disableAttackDialogue = false;
	let disableStackDialogue = false;
	let menuState = 0;
	// countPiecesOnBoard(data.data.current_state)
	$: hands = FenToHand(data.data.current_state);
	$: [onBoardWhite, onBoardBlack] = countPiecesOnBoard(data.data.current_state);
	$: turnColor = data.data.current_state.split(' ')[2];
	$: turnPlayer = turnColor === 'w' ? data.data.player1 : data.data.player2;
	$: playerColor = data.data.player1 === data.session?.user.user_metadata.username ? 'w' : 'b';
	$: console.log(playerColor);

	function handleDropEvent(event: CustomEvent) {
		// console.log(event.detail);
		disableAttackDialogue = false;
		disableStackDialogue = false;

		const { dragItem, hoverItem } = event.detail;
		let fromCoord = '';
		if (dragItem.coordIndex) {
			const [file, rank] = IndexToCoords(dragItem.coordIndex);
			fromCoord = `From: ${file.toUpperCase()}${rank} \n`;
		} else if (dragItem.from) {
			fromCoord = 'From: Hand \n';
		}
		const [file2, rank2] = IndexToCoords(hoverItem.coordIndex);
		let destinationPieceText = 'No piece at destination';
		if (hoverItem.piece != null) {
			destinationPieceText = `Destination Piece: ${DecodePieceFull(hoverItem.piece)}`;
		}

		if (dragItem.coordIndex) {
			const [file, rank] = IndexToCoords(dragItem.coordIndex);
			const [file2, rank2] = IndexToCoords(hoverItem.coordIndex);
			moveDialogueText = `${DecodePieceFull(
				dragItem.piece
			)} ${file.toUpperCase()}${rank} to ${file2.toUpperCase()}${rank2}`;
		}
		if (GetPieceColor(hoverItem?.piece) == playerColor) {
			disableAttackDialogue = true;
		}
		console.log(hoverItem.stack)
		if (hoverItem.stack?.length == 3){
			disableStackDialogue = true;
		}
		if(hoverItem.stack?.length != 0){
			showMoveDialogue = true;
		} else {
			alert(
				`${fromCoord}From Piece: ${DecodePieceFull(
					dragItem.piece
				)} \nDestination: ${file2.toUpperCase()}${rank2} \n${destinationPieceText}`
			);
		}
	}
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<main>
	<section>
		<Board
			{dragAndDrop}
			{drop}
			{playerColor}
			on:drop={handleDropEvent}
			gameData={data.data}
			reversed={playerColor !== 'w'}
		/>
	</section>
	<aside class="side-menu">
		<div class="game-state">
			{turnColor === 'w' ? 'White' : 'Black'} To Play
		</div>
		<div class="tabs">
			<button
				class="button-ghost"
				on:click={() => {
					menuState = 0;
				}}>hand</button
			>
			<button
				class="button-ghost"
				on:click={() => {
					menuState = 1;
				}}>chat</button
			>
			<button
				class="button-ghost"
				on:click={() => {
					menuState = 2;
				}}>move history</button
			>
		</div>
		{#if menuState === 0}
			<Hand
				on:drop={handleDropEvent}
				{dragAndDrop}
				reversed={playerColor !== 'w'}
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
	<MoveDialogue
		bind:showModal={showMoveDialogue}
		{disableAttackDialogue}
		{disableStackDialogue}
		text={moveDialogueText}
	/>
</main>

<style lang="scss">
	main {
		display: flex;
		max-width: 90rem;
		margin: 0 auto;
	}
	section {
		width: 100%;
		user-select: none;
	}
	aside{
		user-select: none;
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
