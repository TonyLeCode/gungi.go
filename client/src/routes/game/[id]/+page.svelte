<script lang="ts">
	import Board from './Board.svelte';
	import {
		DecodePiece,
		DecodePieceFull,
		FenToBoard,
		FenToHand,
		GetPieceColor,
		IndexToCoords,
	} from '$lib/utils/utils.js';
	import Hand from './Hand.svelte';
	import Chat from './Chat.svelte';
	import Replay from './Replay.svelte';
	import { createDragAndDrop } from '$lib/utils/dragAndDrop';
	import MoveDialogue from './MoveDialogue.svelte';
	import { ws, wsConnState } from '$lib/store/websocket';
	import { onMount } from 'svelte';

	export let data;
	let boardState = data.data;
	$: currentState = FenToBoard(boardState.current_state);
	$: moveList = boardState.moveList;
	// console.log(boardState);
	// console.log(data?.params.id);
	// console.log(data);
	// for(let i in hands[0]){
	// 	const encoded = DecodePiece(i).toLocaleLowerCase()
	// 	console.log(`/pieces/w1${encoded}.svg`)
	// }
	//FenToBoard on board size

	const { dragAndDrop, drop } = createDragAndDrop();

	interface MoveType {
		fromPiece: number;
		fromCoord: number;
		moveType: number;
		toCoord: number;
	}
	let showMoveDialogue = false;
	let moveDialogueInfo: MoveType;

	function countPiecesOnBoard(fen: string) {
		const pieces = fen.split(' ')[0];
		const matchW = pieces.match(/[A-Z]/g);
		const matchB = pieces.match(/[a-z]/g);
		return [matchW?.length ?? 0, matchB?.length ?? 0];
	}
	function reverseNameIfBlack(isBlack: boolean): string {
		return !isBlack ? boardState.player1 : boardState.player2;
	}

	let moveDialogueText = '';
	let disableAttackDialogue = false;
	let disableStackDialogue = false;
	let menuState = 0;
	// countPiecesOnBoard(data.data.current_state)
	$: hands = FenToHand(boardState.current_state);
	$: [onBoardWhite, onBoardBlack] = countPiecesOnBoard(boardState.current_state);
	$: turnColor = boardState.current_state.split(' ')[2];
	$: turnPlayer = turnColor === 'w' ? boardState.player1 : boardState.player2;
	$: playerColor = boardState.player1 === data.session?.user.user_metadata.username ? 'w' : 'b';
	$: isPlayerTurn = turnColor === playerColor;

	function handleDropEvent(event: CustomEvent) {
		// console.log(event.detail);
		disableAttackDialogue = true;
		disableStackDialogue = true;

		const { dragItem, hoverItem } = event.detail;
		if (dragItem?.coordIndex === hoverItem?.coordIndex) {
			return;
		}
		// let fromCoord = '';
		// if (dragItem.coordIndex) {
		// 	const [file, rank] = IndexToCoords(dragItem.coordIndex);
		// 	fromCoord = `From: ${file.toUpperCase()}${rank} \n`;
		// } else if (dragItem.from) {
		// 	fromCoord = 'From: Hand \n';
		// }
		// const [file2, rank2] = IndexToCoords(hoverItem.coordIndex);
		// let destinationPieceText = 'No piece at destination';
		// if (hoverItem.piece != null) {
		// 	destinationPieceText = `Destination Piece: ${DecodePieceFull(hoverItem.piece)}`;
		// }

		if (dragItem.coordIndex) {
			const [file, rank] = IndexToCoords(dragItem.coordIndex);
			const [file2, rank2] = IndexToCoords(hoverItem.coordIndex);
			let piece: number;
			if (dragItem.from === 'hand') {
				piece = dragItem.piece;
			} else {
				const fromSquare = currentState[dragItem?.coordIndex];
				piece = fromSquare[fromSquare.length - 1];
			}
			moveDialogueText = `${DecodePieceFull(piece)} ${file.toUpperCase()}${rank} to ${file2.toUpperCase()}${rank2}`;
		}

		const fromSquare = currentState[hoverItem?.coordIndex];
		if (GetPieceColor(fromSquare[fromSquare.length - 1]) != playerColor) {
			disableAttackDialogue = false;
		}

		const stack = currentState[hoverItem.coordIndex];
		if (stack?.length != 3) {
			disableStackDialogue = false;
		}
		if (stack?.length != 0 && !dragItem.from) {
			const square = currentState[dragItem.coordIndex];

			moveDialogueInfo = {
				fromPiece: square[square.length - 1],
				fromCoord: dragItem.coordIndex,
				moveType: 0,
				toCoord: hoverItem.coordIndex,
			};
			showMoveDialogue = true;
			return;
		} else {
			// alert(
			// 	`${fromCoord}From Piece: ${DecodePieceFull(
			// 		dragItem.piece
			// 	)} \nDestination: ${file2.toUpperCase()}${rank2} \n${destinationPieceText}`
			// );
		}

		if (dragItem.from === 'hand') {
			const move = {
				fromPiece: dragItem.piece,
				fromCoord: -1,
				moveType: 3,
				toCoord: hoverItem.coordIndex,
			};
			const msg = {
				type: 'makeMove',
				payload: move,
			};

			$ws.send(JSON.stringify(msg));
		} else {
			const square = currentState[dragItem.coordIndex];
			const move = {
				fromPiece: square[square.length - 1],
				fromCoord: dragItem.coordIndex,
				moveType: 0,
				toCoord: hoverItem.coordIndex,
			};

			const msg = {
				type: 'makeMove',
				payload: move,
			};

			$ws.send(JSON.stringify(msg));
		}
	}

	function handleMoveEvent(event: CustomEvent) {
		const msg = {
			type: 'makeMove',
			payload: event.detail,
		};
		msg.payload.fromCoord = msg.payload.fromCoord;
		msg.payload.toCoord = msg.payload.toCoord;
		$ws.send(JSON.stringify(msg));
	}

	function handleGameMsg(event: MessageEvent<any>) {
		try {
			const res = JSON.parse(event.data);
			switch (res.type) {
				case 'game':
					boardState = res.payload;
					break;
			}
		} catch (err) {
			console.log(event?.data);
			console.error('Error: ', err);
		}
	}

	onMount(() => {
		ws.subscribe((val) => {
			if (val) {
				$ws.addEventListener('message', handleGameMsg);
			}
		});

		wsConnState.subscribe((val) => {
			if (val === 'connected') {
				const msg = {
					type: 'route',
					payload: 'game',
				};
				$ws.send(JSON.stringify(msg));
				const msg2 = {
					type: 'joinGame',
					payload: boardState.id,
				};
				$ws.send(JSON.stringify(msg2));
			}
		});
		return () => {
			$ws.removeEventListener('message', handleGameMsg);
		};
	});
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<main>
	<section>
		<Board
			{dragAndDrop}
			{drop}
			{moveList}
			{playerColor}
			on:drop={handleDropEvent}
			gameData={currentState}
			reversed={playerColor !== 'w'}
			{isPlayerTurn}
		/>
	</section>
	<aside class="side-menu">
		<div class="game-state">
			{turnColor === 'w' ? 'White' : 'Black'} To Play
		</div>
		<div class="tabs">
			<button
				class={`tab ${menuState === 0 ? 'active' : ''}`}
				on:click={() => {
					menuState = 0;
				}}>hand</button
			>
			<!-- <button
				class={`tab divider ${menuState === 1 ? 'active' : ''}`}
				on:click={() => {
					menuState = 1;
				}}>chat</button
			> -->
			<button
				class={`tab divider ${menuState === 2 ? 'active' : ''}`}
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
				player1={boardState.player1}
				player2={boardState.player2}
				{isPlayerTurn}
				{hands}
				{onBoardBlack}
				{onBoardWhite}
			/>
		{:else if menuState === 1}
			<Chat />
		{:else if menuState === 2}
			<Replay moveHistory={boardState.history.String.split(" ")} />
		{/if}
	</aside>
	<MoveDialogue
		{moveDialogueInfo}
		on:move={handleMoveEvent}
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
	aside {
		user-select: none;
	}

	.tabs {
		/* background-color: red; */
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		/* gap: 1rem; */
		justify-content: center;
		/* margin-left: 10%; */
		margin: 1rem 0;
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.07);
		border-radius: 4px;
		overflow: hidden;
	}

	.tab {
		background-color: rgb(var(--bg-2));
		padding: 0.5rem 1rem;
		&:hover {
			background-color: rgb(var(--primary));
			color: white;
		}
	}
	.divider {
		border-left: 1px rgba(99, 99, 99, 0.2) solid;
	}
	.active {
		background-color: rgb(var(--primary));
		color: white;
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
</style>
