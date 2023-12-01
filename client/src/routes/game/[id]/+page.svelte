<script lang="ts" context="module">
	import type { BoardState } from '$lib/store/gameState';
	export const gameStateContext = createService<Writable<BoardState>>('gameState');
	export const completedContext = createService<
		Readable<{
			completed: boolean;
			result: string;
		}>
	>('completed');
	export const player1NameContext = createService<Readable<string>>('player1Name');
	export const player2NameContext = createService<Readable<string>>('player2Name');
	export const userColorContext = createService<Readable<'w' | 'b' | 'spectator'>>('userColor');
	export const moveHistoryContext = createService<Readable<string[]>>('moveHistory');
	export const manualFlipContext = createService<Writable<boolean>>('manualFlip');
	export const isViewReversedContext = createService<Readable<boolean>>('isViewReversed');
	export const turnColorContext = createService<Readable<string>>('turnColor');
	export const player1HandListContext = createService<Readable<number[]>>('player1HandList');
	export const player2HandListContext = createService<Readable<number[]>>('player2HandList');
	export const isPlayer1ReadyContext = createService<Readable<boolean>>('isPlayer1Ready');
	export const isPlayer2ReadyContext = createService<Readable<boolean>>('isPlayer2Ready');
	export const isUserTurnContext = createService<Readable<boolean>>('isUserTurn');
	export const player1ArmyCountContext = createService<Readable<number>>('player1ArmyCount');
	export const player2ArmyCountContext = createService<Readable<number>>('player2ArmyCount');
	export const player1HandCountContext = createService<Readable<number>>('player1HandCount');
	export const player2HandCountContext = createService<Readable<number>>('player2HandCount');
	export const moveListContext = createService<Readable<{ [key: number]: number[] }>>('moveList');
	export const moveListUIContext = createService<Readable<{ [key: number]: number[] }>>('moveListUI');
	export const boardStateContext = createService<Readable<number[][]>>('boardState');
	export const boardUIContext = createService<Readable<number[][]>>('boardUI');
</script>

<script lang="ts">
	import Board from './Board.svelte';
	import { DecodePieceFull, GetPieceColor, IndexToCoords } from '$lib/utils/utils.js';
	import Hand from './Hand.svelte';
	import Chat from './Chat.svelte';
	import Replay from './Replay.svelte';
	import { createDragAndDrop } from '$lib/utils/dragAndDrop';
	import MoveDialogue from './MoveDialogue.svelte';
	import { ws } from '$lib/store/websocket';
	import { onMount } from 'svelte';
	import { createService } from '$lib/store/contextHelper';
	import type { Readable, Writable } from 'svelte/store';
	import { get } from 'svelte/store';
	import { createGameStore } from '$lib/store/gameState';
	import { notifications } from '$lib/store/notification';
	import { nanoid } from 'nanoid';
	import UndoDialogue from './UndoDialogue.svelte';
	import Modal from '$lib/components/Modal.svelte';

	type undoRequests = {
		receiver_username: string;
		sender_username: string;
		status: 'pending' | 'accept' | 'reject';
	};
	export let data;
	const username = data.session?.user.user_metadata.username;
	const undoRequests: undoRequests[] = data.gameData.undo_requests;
	let undoDialogBool = false;
	let completedBool = false;
	let completedText = '';

	for (let i = 0; i < undoRequests.length; i++) {
		if (undoRequests[i].receiver_username === username && undoRequests[i].status === 'pending') {
			undoDialogBool = true;
		} else if (undoRequests[i].sender_username === username && undoRequests[i].status === 'accept') {
			const msg = {
				type: 'completeGameUndo',
				payload: '',
			};
			ws?.send(msg);
			notifications?.add({
				id: nanoid(),
				title: 'Accepted',
				type: 'default',
				msg: 'Your undo request has been accepted',
			});
		} else if (undoRequests[i].sender_username === username && undoRequests[i].status === 'reject') {
			const msg = {
				type: 'completeGameUndo',
				payload: '',
			};
			ws?.send(msg);
			notifications?.add({
				id: nanoid(),
				title: 'Rejected',
				type: 'default',
				msg: 'Your undo request has been rejected',
			});
		}
	}
	const gameStore = createGameStore(data.gameData, data.session?.user.user_metadata.username);
	gameStateContext.set(gameStore.gameState);
	completedContext.set(gameStore.completed);
	player1NameContext.set(gameStore.player1Name);
	player2NameContext.set(gameStore.player2Name);
	userColorContext.set(gameStore.userColor);
	moveHistoryContext.set(gameStore.moveHistory);
	manualFlipContext.set(gameStore.manualFlip);
	isViewReversedContext.set(gameStore.isViewReversed);
	turnColorContext.set(gameStore.turnColor);
	isPlayer1ReadyContext.set(gameStore.isPlayer1Ready);
	isPlayer2ReadyContext.set(gameStore.isPlayer2Ready);
	isUserTurnContext.set(gameStore.isUserTurn);
	player1HandListContext.set(gameStore.player1HandList);
	player2HandListContext.set(gameStore.player2HandList);
	player1ArmyCountContext.set(gameStore.player1ArmyCount);
	player2ArmyCountContext.set(gameStore.player2ArmyCount);
	player1HandCountContext.set(gameStore.player1HandCount);
	player2HandCountContext.set(gameStore.player2HandCount);
	moveListContext.set(gameStore.moveList);
	moveListUIContext.set(gameStore.moveListUI);
	boardStateContext.set(gameStore.boardState);
	boardUIContext.set(gameStore.boardUI);
	const boardState = gameStore.boardState;
	const boardUI = gameStore.boardUI;
	const isPlayer1Ready = gameStore.isPlayer1Ready;
	const isPlayer2Ready = gameStore.isPlayer2Ready;
	const gameState = gameStore.gameState;
	const userColor = gameStore.userColor;
	const turnColor = gameStore.turnColor;
	const completed = gameStore.completed;
	$: console.log($completed)

	const { dragAndDrop, drop } = createDragAndDrop();

	interface MoveType {
		fromPiece: number;
		fromCoord: number;
		moveType: number;
		toCoord: number;
	}
	let showMoveDialogue = false;
	let moveDialogueInfo: MoveType;

	let moveDialogueText = '';
	let disableAttackDialogue = false;
	let disableStackDialogue = false;
	let menuState = 0;
	let stack: number[] = [];

	function handleDropEvent(event: CustomEvent) {
		disableAttackDialogue = true;
		disableStackDialogue = true;

		const { dragItem, hoverItem } = event.detail;
		if (dragItem?.coordIndex === hoverItem?.coordIndex) return;
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
				const fromSquare = $boardState[dragItem?.coordIndex];
				piece = fromSquare[fromSquare.length - 1];
			}
			moveDialogueText = `${DecodePieceFull(piece)} ${file.toUpperCase()}${rank} to ${file2.toUpperCase()}${rank2}`;
		}

		const fromSquare = $boardState[hoverItem?.coordIndex];
		if (GetPieceColor(fromSquare[fromSquare.length - 1]) != $userColor) {
			disableAttackDialogue = false;
		}

		const stack = $boardState[hoverItem.coordIndex];
		if (stack?.length != 3) {
			disableStackDialogue = false;
		}
		if (stack?.length != 0 && !dragItem.from) {
			const square = $boardState[dragItem.coordIndex];

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
				type: 'makeGameMove',
				payload: { ...move, gameID: get(gameState).id },
			};

			ws?.send(msg);
		} else {
			const square = $boardState[dragItem.coordIndex];
			const move = {
				fromPiece: square[square.length - 1],
				fromCoord: dragItem.coordIndex,
				moveType: 0,
				toCoord: hoverItem.coordIndex,
			};

			const msg = {
				type: 'makeGameMove',
				payload: { ...move, gameID: get(gameState).id },
			};

			ws?.send(msg);
		}
	}

	function handleStackClick(event: CustomEvent) {
		stack = get(boardUI)[event.detail];
	}

	function handleMoveEvent(event: CustomEvent) {
		const msg = {
			type: 'makeGameMove',
			payload: { ...event.detail, gameID: get(gameState).id },
		};
		msg.payload.fromCoord = msg.payload.fromCoord;
		msg.payload.toCoord = msg.payload.toCoord;
		ws?.send(msg);
	}

	function handleGameMsg(event?: MessageEvent) {
		try {
			const res = JSON.parse(event?.data);
			switch (res.type) {
				case 'game':
					gameStore.gameState.set(res.payload);
					break;
				case 'undoRequest':
					console.log('request undo');
					undoDialogBool = true;
					break;
				case 'undoResponse':
					if (res.payload === 'accept') {
						notifications?.add({
							id: nanoid(),
							title: 'Accepted',
							type: 'default',
							msg: 'Your undo request has been accepted',
						});
					} else {
						notifications?.add({
							id: nanoid(),
							title: 'Rejected',
							type: 'default',
							msg: 'Your undo request has been rejected',
						});
					}
					const msg = {
						type: 'completeGameUndo',
						payload: '',
					};
					ws?.send(msg);
					break;
				case 'gameEnd':
					completedText = res.payload;
					completedBool = true;
					break;
			}
		} catch (err) {
			console.log(event?.data);
			console.error('Error: ', err);
		}
	}

	function handleResignEvent(event: CustomEvent) {
		console.log('resign');
		const msg = {
			type: 'resign',
		};
		ws?.send(msg);
	}

	function handleUndoEvent(event: CustomEvent) {
		console.log('requestUndo');
		const msg = {
			type: 'requestGameUndo',
		};
		ws?.send(msg);
	}

	function handleUndoResponseEvent(event: CustomEvent) {
		console.log('requestUndo');
		const msg = {
			type: 'responseGameUndo',
			payload: event.detail.response,
		};
		ws?.send(msg);
	}
	function handleReadyEvent(event: CustomEvent) {
		console.log('ready');
		const msg = {
			type: 'makeGameMove',
			payload: {
				gameID: get(gameState).id,
				fromPiece: -1,
				fromCoord: 0,
				moveType: 4,
				toCoord: 0,
			},
		};
		ws?.send(msg);
	}

	onMount(() => {
		const unsubGameMsg = ws?.addMsgListener(handleGameMsg);

		const unsubConnect = ws?.subscribe((val) => {
			if (val === 'connected') {
				const msg2 = {
					type: 'joinGame',
					payload: get(gameStore.gameState).id,
				};
				ws?.send(msg2);
			}
		});
		return () => {
			if (unsubGameMsg) unsubGameMsg();
			if (unsubConnect) unsubConnect();
			const msg = {
				type: 'leaveGame',
				payload: get(gameState).id,
			};
			ws?.send(msg);
		};
	});
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<main>
	<section>
		<Board {dragAndDrop} {drop} on:drop={handleDropEvent} on:stackClick={handleStackClick} />
	</section>
	<aside class="side-menu">
		<div class="game-state">
			{#if $completed.completed}
				{#if $completed.result === 'b'}
					Black Wins By Checkmate
				{:else if $completed.result === 'w'}
					White Wins By Checkmate
				{:else if $completed.result === 'stalement'}
					Stalemate
				{/if}
			{:else}
				{$isPlayer1Ready && $isPlayer2Ready ? '' : 'Drafting -'}
				{$turnColor === 'w' ? 'White' : 'Black'} To Play
			{/if}
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
				on:ready={handleReadyEvent}
				on:resign={handleResignEvent}
				on:undo={handleUndoEvent}
				on:drop={handleDropEvent}
				{stack}
				{dragAndDrop}
			/>
		{:else if menuState === 1}
			<Chat />
		{:else if menuState === 2}
			<Replay />
		{/if}
	</aside>
	{#if completedBool}
		<Modal bind:showModal={completedBool}>
			<h2 class="completed-text">{completedText}</h2>
		</Modal>
	{/if}
	<MoveDialogue
		{moveDialogueInfo}
		on:move={handleMoveEvent}
		bind:showModal={showMoveDialogue}
		{disableAttackDialogue}
		{disableStackDialogue}
		text={moveDialogueText}
	/>
	<UndoDialogue on:undoResponse={handleUndoResponseEvent} bind:showModal={undoDialogBool} />
</main>

<style lang="scss">
	main {
		max-width: 90rem;
		margin: 0 auto;
	}

	section {
		width: 100%;
		user-select: none;
		margin: auto;
		max-width: 50rem;
	}
	aside {
		user-select: none;
	}

	.completed-text {
		font-size: 1.2rem;
		font-weight: 600;
	}
	.tabs {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		justify-content: center;
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
		margin: 0 auto;
		max-width: 44rem;
		width: 100%;
		padding: 0 2rem;
	}

	@media only screen and (min-width: 1200px) {
		main {
			display: flex;
		}
		.side-menu {
			margin-left: auto;
			max-width: 36rem;
			padding: 0;
		}
	}
</style>
