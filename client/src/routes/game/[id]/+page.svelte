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
	import MoveDialogue from './MoveDialogue.svelte';
	import UndoDialogue from './UndoDialogue.svelte';
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
	let showMoveDialogue = $state(false);
	let moveDialogueText = $state('');
	let attackFn = $state<(() => void) | null>(null);
	let stackFn = $state<(() => void) | null>(null);

	let undoDialogBool = $state(false);
	let selectedStack = $state<number[]>([]);
	let completedBool = $state(false);
	let completedText = $state('');

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
		showMoveDialogue = true;
		const trueFromCoord = boardStore.isViewReversed ? 80 - fromCoord : fromCoord;
		const trueToCoord = boardStore.isViewReversed ? 80 - toCoord : toCoord;
		const [fromFile, fromRank] = IndexToCoords(trueFromCoord);
		const [toFile, toRank] = IndexToCoords(trueToCoord);
		const fromSquare = boardStore.boardState[trueFromCoord];
		const toSquare = boardStore.boardState[trueToCoord];
		const fromPiece = fromSquare[fromSquare.length - 1];
		const toPiece = toSquare[toSquare.length - 1];

		moveDialogueText = `${DecodePieceFull(fromPiece)} ${fromFile.toUpperCase()}${fromRank} to ${DecodePieceFull(toPiece)} ${toFile.toUpperCase()}${toRank}`;

		if (GetPieceColor(fromPiece) !== GetPieceColor(toPiece)) {
			attackFn = () => {
				sendMoveMsg(fromPiece, trueFromCoord, trueToCoord, 2);
				showMoveDialogue = false;
			};
		} else {
			attackFn = null;
		}

		if (toSquare.length !== 3) {
			stackFn = () => {
				sendMoveMsg(fromPiece, trueFromCoord, trueToCoord, 1);
				showMoveDialogue = false;
			};
		} else {
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

	function handleGameMsg(event?: MessageEvent) {
		try {
			const data = JSON.parse(event?.data);
			switch (data.type) {
				case 'game':
					boardStore.updateBoard(data.payload);
					break;
				case 'undoRequest':
					undoDialogBool = true;
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
							type: 'warning',
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
					completedBool = true;
					break;
				case 'gameResign':
					if (data.payload === 'w/r') {
						completedText = 'Black Resigns';
						completedBool = true;
					} else if (data.payload === 'b/r') {
						completedText = 'White Resigns';
						completedBool = true;
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
		let unsub = websocketStore.addMsgListener(handleGameMsg);

		$effect(() => {
			if (websocketStore.state === 'connected') {
				const msg = {
					type: 'joinGame',
					payload: boardStore.id,
				};
				websocketStore.send(msg);
			}
		});

		return () => {
			unsub?.();
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
	<Menu {selectHandPiece} {selectedStack} {ready} {placeHandMove} {droppable} />

	<!-- TODO modal and dialogues -->
	{#if completedBool}
		<Modal bind:showModal={completedBool}>
			<h2 class="completed-text">{completedText}</h2>
		</Modal>
	{/if}
	<MoveDialogue bind:showModal={showMoveDialogue} text={moveDialogueText} {attackFn} {stackFn} />
	<UndoDialogue />
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
</style>
