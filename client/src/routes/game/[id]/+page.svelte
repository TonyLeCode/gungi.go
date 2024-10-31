<script lang="ts">
	import Board from '../../../lib/components/Board.svelte';
	import { setGameStore } from '$lib/store/gameState.svelte';
	import PlayerInfo from './PlayerInfo.svelte';
	import Menu from './Menu.svelte';
	import { Crown } from 'lucide-svelte';
	import { getWebsocketStore } from '$lib/store/websocketStore.svelte';
	import { onMount } from 'svelte';
	import { getNotificationStore } from '$lib/store/notificationStore.svelte';
	import { nanoid } from 'nanoid';
	import Modal from '$lib/components/Modal.svelte';
	import { DecodePieceFull, GetPieceColor, IndexToCoords, PieceIsPlayerColor, ReverseIndices } from '$lib/utils/utils';
	import { Droppable, draggable, type DraggableOptions } from '$lib/store/dragAndDrop.svelte';
	import { getSquareCoords } from '$lib/utils/historyParser';
	import { browser } from '$app/environment';
	import { setReplayStore } from '$lib/store/replayStore.svelte';
	let { data } = $props();

	// 0 - move
	// 1 - stack
	// 2 - attack
	// 3 - place
	// 4 - ready

	type Selection =
		| {
				state: 'none';
		  }
		| {
				// a board piece has been selected
				state: 'selectedBoardPiece';
				index: number;
		  }
		| {
				// a board piece had previously been selected
				// now seeing if user drags or releases
				state: 'awaitingBoardPiece';
				currIndex: number;
				prevIndex: number;
		  }
		| {
				// a hand piece has been selected for dropping
				state: 'selectedHandPiece';
				piece: number;
		  }
		| {
				// show stack details when other selections
				// do not fit
				state: 'stackDetails';
				index: number;
		  };

	let username = data.username as string
	const boardStore = setGameStore(data.gameData, username);
	const websocketStore = getWebsocketStore();
	const notificationStore = getNotificationStore();
	const replayStore = setReplayStore(data.gameData, username);
	replayStore.setTotalPages(boardStore.moveHistory.length);
	// $effect(() => {
	// 	replayStore.setTotalPages(boardStore.moveHistory.length);
	// 	console.log("setting total pages")
	// });

	type DropItem = {
		destinationIndex: number;
		destinationStack: number[];
	};

	const droppable = new Droppable<DropItem>();
	let selection = $state<Selection>({ state: 'none' });
	let selectedStack = $derived.by<number[]>(() => {
		if (selection.state === 'selectedBoardPiece') {
			return boardStore.boardUI[selection.index];
		} else if (selection.state === 'awaitingBoardPiece') {
			return boardStore.boardUI[selection.prevIndex];
		} else if (selection.state === 'stackDetails') {
			return boardStore.boardUI[selection.index];
		}
		return [];
	});
	let highlight = $derived.by(() => {
		const arr = [];
		if (selection.state === 'selectedBoardPiece') {
			arr.push(selection.index);
		} else if (selection.state === 'awaitingBoardPiece') {
			arr.push(selection.currIndex);
		}
		let lastMove = getSquareCoords(boardStore.moveHistory[boardStore.moveHistory.length - 1]);
		lastMove = boardStore.isViewReversed ? ReverseIndices(lastMove) : lastMove;
		lastMove.forEach((index) => arr.push(index));
		return arr;
	});
	let selectedMoves = $derived.by(() => {
		if (selection.state === 'selectedBoardPiece') {
			return boardStore.moveListUI[selection.index];
		} else if (selection.state === 'awaitingBoardPiece') {
			return boardStore.moveListUI[selection.prevIndex];
		}
		return [];
	});

	let deselectTimeout = $state<number | undefined>(undefined);
	let blockDeselect = $state(false);
	if (browser) {
		window.addEventListener('mousedown', (e) => {
			if (selection.state !== 'none' && deselectTimeout === undefined) {
				deselectTimeout = window.setTimeout(() => {
					if (!blockDeselect) {
						selection = { state: 'none' };
					}
					deselectTimeout = undefined;
					blockDeselect = false;
				}, 50);
			}
		});
	}

	let showMoveDialogue: ReturnType<typeof Modal>;
	let moveDialogueText = $state('');
	let attackFn = $state<(() => void) | null>(null);
	let stackFn = $state<(() => void) | null>(null);

	let showUndoDialogue: ReturnType<typeof Modal>;
	let completedDialog: ReturnType<typeof Modal>;
	let completedText = $state(' '); // can't be empty string because of svelte bug

	function sendMoveMsg(fromPiece: number, fromCoord: number, toCoord: number, moveType: number) {
		const msg = {
			type: 'makeGameMove',
			payload: {
				// gamePublicID: boardStore.public_id,
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
		if (fromCoord === toCoord) return;

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
				showMoveDialogue?.close();
			};
		} else {
			attackFn = null;
		}

		if (toSquare.length !== 3 && DecodePieceFull(fromPiece) !== 'Fortress') {
			stackFn = () => {
				sendMoveMsg(fromPiece, trueFromCoord, trueToCoord, 1);
				showMoveDialogue?.close();
			};
		} else {
			stackFn = null;
		}

		if (attackFn || stackFn) {
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
		showUndoDialogue?.close();
	}
	function handleGameMsg(event?: MessageEvent) {
		try {
			const data = JSON.parse(event?.data);
			switch (data.type) {
				case 'game':
					boardStore.updateBoard(data.payload);
					replayStore.boardStore.updateBoard(data.payload);
					replayStore.setTotalPages(boardStore.moveHistory.length);
					break;
				case 'undoRequest':
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
					completedDialog?.open();
					break;
				case 'gameResign':
					if (data.payload === 'w/r') {
						completedText = 'Black Resigns';
						completedDialog?.open();
					} else if (data.payload === 'b/r') {
						completedText = 'White Resigns';
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
		blockDeselect = true;
		if (piece === -1) {
			selection = { state: 'none' };
			return;
		}
		selection = {
			state: 'selectedHandPiece',
			piece: piece,
		};
	}

	function isActive(index: number): boolean {
		if (!boardStore.isUserTurn) return false;
		const stack = boardStore.boardUI[index];
		const isPlayerPiece = PieceIsPlayerColor(stack[stack.length - 1], boardStore.userColor);
		const isDraftingPhase = !boardStore.isPlayer1Ready || !boardStore.isPlayer2Ready;

		if (stack.length > 0 && !boardStore.moveListUI[index]) return false;

		return isPlayerPiece && !isDraftingPhase;
	}

	function moveHandler(fromCoord: number, toCoord: number) {
		const toSquare = boardStore.boardUI[toCoord];
		if (!boardStore.moveListUI[fromCoord].includes(toCoord)) return;
		if (toSquare.length === 0) {
			const fromSquare = boardStore.boardUI[fromCoord];
			const fromPiece = fromSquare[fromSquare.length - 1];
			movePiece(fromPiece, fromCoord, toCoord);
		} else {
			promptMoveDialogue(fromCoord, toCoord);
		}
	}

	function draggableOptions(index: number, stack: number[]): DraggableOptions<DropItem | null> {
		// Should select on start
		// If already selected, should store previous index and wait
		// Unless new selection is not a possible move, then make new selection
		// If drag is released, should deselect
		// Unless dropped as a move, then prompt move
		// If short or long pressed, should stay selected
		// If previously selected, should reselect
		// If selected same piece, should unselect
		// Unless selection is a move, then prompt move
		return {
			startEvent: () => {
				blockDeselect = true;
				if (selection.state !== 'selectedBoardPiece') {
					selection = {
						state: 'selectedBoardPiece',
						index: index,
					};
					return;
				}
				if (selectedMoves.includes(index) || index === selection.index) {
					selection = {
						state: 'awaitingBoardPiece',
						currIndex: index,
						prevIndex: selection.index,
					};
				} else {
					selection = {
						state: 'selectedBoardPiece',
						index: index,
					};
				}
			},
			dragStartEvent: () => {
				if (selection.state === 'awaitingBoardPiece') {
					selection = {
						state: 'selectedBoardPiece',
						index: index,
					};
				}
			},
			shortReleaseEvent: () => {
				if (selection.state === 'awaitingBoardPiece') {
					if (selection.currIndex === selection.prevIndex) {
						selection = {
							state: 'none',
						};
					} else {
						moveHandler(selection.prevIndex, selection.currIndex);
						selection = {
							state: 'none',
						};
					}
				}
			},
			longReleaseEvent: () => {
				if (selection.state === 'awaitingBoardPiece') {
					if (selection.currIndex === selection.prevIndex) {
						selection = {
							state: 'none',
						};
					} else {
						moveHandler(selection.prevIndex, selection.currIndex);
						selection = {
							state: 'none',
						};
					}
				}
			},
			dragReleaseEvent: (hoverItem) => {
				if (hoverItem === null || hoverItem === undefined) {
					selection = { state: 'none' };
					return;
				}
				if (selection.state === 'selectedBoardPiece') moveHandler(index, hoverItem.destinationIndex);

				selection = { state: 'none' };
			},
			releaseEvent: (hoverItem) => {},
			droppable: droppable,
			active: () => {
				return isActive(index);
			},
		};
	}

	onMount(() => {
		$effect(() => {
			if (websocketStore.state === 'connected') {
				websocketStore.addMsgListener(handleGameMsg);
				const msg = {
					type: 'joinGame',
					payload: boardStore.public_id,
				};
				websocketStore.send(msg);

				const undoRequests = data.gameData.undo_requests;
				for (let i = 0; i < undoRequests.length; i++) {
					if (undoRequests[i].receiver_username === boardStore.username && undoRequests[i].status === 'pending') {
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
				payload: boardStore.public_id,
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
		<p>{selection.state}</p>
		<div class="game-state">
			{#if replayStore.isActive}
				<div>Viewing Game History</div>
			{:else if boardStore.completed && !replayStore.isActive}
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
		{#if replayStore.isActive}
			<!-- <p>Replaying</p> -->
			<Board
				boardState={replayStore.boardStore.boardUI}
				isReversed={replayStore.boardStore.isViewReversed}
				showCoord={true}
			/>
		{:else}
			<Board
				boardState={boardStore.boardUI}
				isReversed={boardStore.isViewReversed}
				showCoord={true}
				dragAction={draggable}
				dropAction={droppable.addDroppable.bind(droppable)}
				{highlight}
				{selectedMoves}
				{draggableOptions}
				onMouseDown={(index) => {
					const targetSquare = boardStore.boardUI[index];
					const isEnemyStack = !PieceIsPlayerColor(targetSquare[targetSquare.length - 1], boardStore.userColor);
					if (selectedMoves.includes(index) && selection.state === 'selectedBoardPiece') {
						blockDeselect = true;
						moveHandler(selection.index, index);
						selection = {
							state: 'none',
						};
					} else if (isEnemyStack && targetSquare.length > 0) {
						blockDeselect = true;
						selection = {
							state: 'stackDetails',
							index: index,
						};
					} else if (!isEnemyStack && targetSquare.length > 0 && !boardStore.isUserTurn) {
						blockDeselect = true;
						selection = {
							state: 'stackDetails',
							index: index,
						};
					} else if (targetSquare.length > 0 && !boardStore.moveListUI[index]) {
						blockDeselect = true;
						selection = {
							state: 'stackDetails',
							index: index,
						}
					}
				}}
			/>
		{/if}
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
