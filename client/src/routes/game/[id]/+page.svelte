<script lang="ts">
	import Board from './Board.svelte';
	import { setGameStore } from '$lib/store/gameState.svelte';
	import PlayerInfo from './PlayerInfo.svelte';
	import Menu from './Menu.svelte';
	import { Crown } from 'lucide-svelte';
	import { getWebsocketStore } from '$lib/store/websocketStore.svelte';
	import { onMount } from 'svelte';
	import { getNotificationStore } from '$lib/store/notificationStore.svelte';
	import { nanoid } from 'nanoid';
	import Modal from '$lib/components/Modal.svelte';
	import { DecodePieceFull, GetPieceColor, IndexToCoords } from '$lib/utils/utils';
	import { Droppable } from '$lib/store/dragAndDrop.svelte';
	let { data } = $props();

	// 0 - move
	// 1 - stack
	// 2 - attack
	// 3 - place
	// 4 - ready

	const boardStore = setGameStore(data.gameData, data.session?.user.user_metadata.username);
	const websocketStore = getWebsocketStore();
	const notificationStore = getNotificationStore();

	type DropItem = {
		destinationIndex: number;
		destinationStack: number[];
	};

	const droppable = new Droppable<DropItem>();
	let selectedHandPiece = $state(-1);

	interface MoveType {
		fromPiece: number;
		fromCoord: number;
		moveType: number;
		toCoord: number;
	}
	let showMoveDialogue: Modal
	let moveDialogueText = $state('');
	let attackFn = $state<(() => void) | null>(null);
	let stackFn = $state<(() => void) | null>(null);

	let showUndoDialogue: Modal
	let selectedStack = $state<number[]>([]);
	let completedDialog: Modal
	let completedText = $state(' '); // can't be empty string because of svelte bug

	function changeSelectedStack(stack: number[]) {
		const x = getWebsocketStore();
		selectedStack = stack;
	}

	function sendMoveMsg(fromPiece: number, fromCoord: number, toCoord: number, moveType: number) {
		const msg = {
			type: 'makeGameMove',
			payload: {
				gameID: boardStore.id,
				fromPiece,
				fromCoord,
				moveType,
				toCoord,
			},
		};

		websocketStore.send(msg);
	}

	// Attacking and Stacking from onboard piece
	function promptMoveDialogue(fromCoord: number, toCoord: number) {
		const trueFromCoord = boardStore.isViewReversed ? 80 - fromCoord : fromCoord;
		const trueToCoord = boardStore.isViewReversed ? 80 - toCoord : toCoord;
		const [fromFile, fromRank] = IndexToCoords(trueFromCoord);
		const [toFile, toRank] = IndexToCoords(trueToCoord);
		const fromSquare = boardStore.boardState[trueFromCoord];
		const toSquare = boardStore.boardState[trueToCoord];
		const fromPiece = fromSquare[fromSquare.length - 1];
		const toPiece = toSquare[toSquare.length - 1];

		if (DecodePieceFull(fromPiece) === 'Fortress' && toSquare.length > 1) return;

		moveDialogueText = `${DecodePieceFull(fromPiece)} ${fromFile.toUpperCase()}${fromRank} to ${DecodePieceFull(toPiece)} ${toFile.toUpperCase()}${toRank}`;

		if (GetPieceColor(fromPiece) !== GetPieceColor(toPiece)) {
			attackFn = () => {
				sendMoveMsg(fromPiece, trueFromCoord, trueToCoord, 2);
				// showMoveDialogue = false;
				showMoveDialogue?.close()
			};
		} else {
			attackFn = null;
		}

		if (toSquare.length !== 3 && DecodePieceFull(fromPiece) !== 'Fortress') {
			stackFn = () => {
				sendMoveMsg(fromPiece, trueFromCoord, trueToCoord, 1);
				// showMoveDialogue = false;
				showMoveDialogue?.close()
			};
		} else {
			stackFn = null;
		}

		if (attackFn || stackFn) {
			// showMoveDialogue = true;
			showMoveDialogue?.open();
		} else {
			attackFn = null;
			stackFn = null;
		}
	}

	// Placing from hand
	function placeHandMove(fromPiece: number, toCoord: number) {
		const trueToCoord = boardStore.isViewReversed ? 80 - toCoord : toCoord;
		sendMoveMsg(fromPiece, -1, trueToCoord, 3);
	}

	function movePiece(fromPiece: number, fromCoord: number, toCoord: number) {
		const trueFromCoord = boardStore.isViewReversed ? 80 - fromCoord : fromCoord;
		const trueToCoord = boardStore.isViewReversed ? 80 - toCoord : toCoord;
		sendMoveMsg(fromPiece, trueFromCoord, trueToCoord, 0);
	}

	function ready() {
		sendMoveMsg(-1, 0, 0, 4);
	}

	function undo() {
		const msg = {
			type: 'requestGameUndo',
		};
		websocketStore.send(msg);
	}

	function resign() {
		const msg = {
			type: 'gameResign',
		};
		websocketStore.send(msg);
	}

	function undoReponse(response: 'accept' | 'reject') {
		const msg = {
			type: 'responseGameUndo',
			payload: response,
		};
		websocketStore.send(msg);
		// showUndoDialogue = false;
		showUndoDialogue?.close();
	}

	function handleGameMsg(event?: MessageEvent) {
		try {
			const data = JSON.parse(event?.data);
			switch (data.type) {
				case 'game':
					boardStore.updateBoard(data.payload);
					break;
				case 'undoRequest':
					// showUndoDialogue = true;
					showUndoDialogue?.open();
					break;
				case 'undoResponse':
					if (data.payload === 'accept') {
						notificationStore.add({
							id: nanoid(),
							title: 'Accepted',
							type: 'success',
							msg: 'Your undo request has been accepted',
						});
					} else {
						notificationStore.add({
							id: nanoid(),
							title: 'Rejected',
							type: 'error',
							msg: 'Your undo request has been rejected',
						});
					}
					const msg = {
						type: 'completeGameUndo',
						payload: '',
					};
					websocketStore.send(msg);
					break;
				case 'gameEnd':
					completedText = data.payload;
					// completedBool = true;
					completedDialog?.open();
					break;
				case 'gameResign':
					if (data.payload === 'w/r') {
						completedText = 'Black Resigns';
						// completedBool = true;
						completedDialog?.open();
					} else if (data.payload === 'b/r') {
						completedText = 'White Resigns';
						// completedBool = true;
						completedDialog?.open();
					}
					break;
			}
		} catch (err) {
			console.log(event?.data);
			console.error('Error: ', err);
		}
	}

	function selectHandPiece(piece: number) {
		selectedHandPiece = piece;
	}

	onMount(() => {
		$effect(() => {
			if (websocketStore.state === 'connected') {
				websocketStore.addMsgListener(handleGameMsg);
				const msg = {
					type: 'joinGame',
					payload: boardStore.id,
				};
				websocketStore.send(msg);

				const undoRequests = data.gameData.undo_requests;
				for (let i = 0; i < undoRequests.length; i++) {
					if (undoRequests[i].receiver_username === boardStore.username && undoRequests[i].status === 'pending') {
						// showUndoDialogue = true;
						showUndoDialogue?.open();
					} else if (undoRequests[i].sender_username === boardStore.username && undoRequests[i].status === 'accept') {
						const msg = {
							type: 'completeGameUndo',
							payload: '',
						};
						websocketStore.send(msg);
						notificationStore.add({
							id: nanoid(),
							title: 'Accepted',
							type: 'success',
							msg: 'Your undo request has been accepted',
						});
					} else if (undoRequests[i].sender_username === boardStore.username && undoRequests[i].status === 'reject') {
						const msg = {
							type: 'completeGameUndo',
							payload: '',
						};
						websocketStore.send(msg);
						notificationStore.add({
							id: nanoid(),
							title: 'Rejected',
							type: 'error',
							msg: 'Your undo request has been rejected',
						});
					}
				}
			}
		});

		return () => {
			websocketStore.removeMsgListener(handleGameMsg);
			const msg = {
				type: 'leaveGame',
				payload: boardStore.id,
			};
			websocketStore.send(msg);
		};
	});
</script>

<svelte:head>
	<title>Game: {boardStore.player1} vs {boardStore.player2} | White Monarch Server</title>
</svelte:head>

<main>
	<section>
		<div class="game-state">
			{#if boardStore.completed}
				<Crown />
				<span class="completed-text">
					{#if boardStore.result === 'b'}
						Black Won By Checkmate
					{:else if boardStore.result === 'w'}
						White Won By Checkmate
					{:else if boardStore.result === 'b/r'}
						Black Won By Resignation
					{:else if boardStore.result === 'w/r'}
						White Won By Resignation
					{:else if boardStore.result === 'draw'}
						Draw
					{/if}
				</span>
			{:else}
				<span class:turn-indicator={boardStore.isUserTurn}
					>{boardStore.isPlayer1Ready && boardStore.isPlayer2Ready ? '' : 'Drafting -'}
					{boardStore.turnColor === 'w' ? 'White' : 'Black'} To Play</span
				>
			{/if}
		</div>
		<PlayerInfo isOpposite={true} />
		<Board {changeSelectedStack} {promptMoveDialogue} {movePiece} {droppable} />
		<PlayerInfo isOpposite={false} />
	</section>
	<Menu {resign} {undo} {selectHandPiece} {selectedStack} {ready} {placeHandMove} {droppable} />


	<Modal bind:this={completedDialog}>
		<h2 class="completed-text completed-dialogue"><Crown />{completedText}</h2>
	</Modal>
	<Modal bind:this={showMoveDialogue}>
		<p class="dialog-p">{moveDialogueText}</p>
		<div class="button-container">
			<button class="button-primary" onclick={attackFn} disabled={attackFn === null}>Attack</button>
			<button class="button-primary" onclick={stackFn} disabled={stackFn === null}>Stack</button>
		</div>
	</Modal>
	<Modal bind:this={showUndoDialogue}>
		<p class="dialog-p">Your opponent has requested an undo</p>
		<div class="button-container">
			<button onclick={() => undoReponse('accept')} class="button-primary">Accept Undo</button>
			<button onclick={() => undoReponse('reject')} class="button-primary">Reject Undo</button>
		</div>
	</Modal>
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

	.game-state :global(svg) {
		width: 1.25rem;
		height: 1.25rem;
		@media (min-width: 608px) {
			width: 1.5rem;
			height: 1.5rem;
		}
		@media (min-width: 767px) {
			width: 2rem;
			height: 2rem;
		}
	}

	.completed-text {
		font-weight: 600;
		margin-top: 6px;
		@media (min-width: 767px) {
			font-size: 1.2rem;
		}
	}

	.turn-indicator {
		color: rgb(var(--primary));
		font-weight: 600;
	}

	.game-state {
		display: flex;
		gap: 0.5rem;
		justify-content: center;
		align-items: center;
		text-align: center;
		margin-bottom: 0.25rem;
		@media (min-width: 767px) {
			margin: 0.25rem 0;
		}
	}

	.completed-dialogue {
		display: flex;
		gap: 0.5rem;
		align-items: center;
		text-align: center;
		:global(svg) {
			width: 20px;
			height: 20px;
		}
		@media (min-width: 767px) {
			:global(svg) {
				width: 30px;
				height: 30px;
			}
		}
	}

	.dialog-p {
		margin-bottom: 1rem;
		text-align: center;
	}

	.button-container {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}
</style>
