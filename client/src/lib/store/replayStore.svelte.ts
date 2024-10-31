import { getContext, setContext } from 'svelte';
import { createPaginationStore } from './paginationStore.svelte';
import { BoardStore } from './gameState.svelte';
import type { BoardState } from './gameState.svelte';
import { serializeFen } from '$lib/utils/fenParser';
import { parseMove } from '$lib/utils/historyParser';

//TODO take original fen into account (not current state fen)

class ReplayStore {
	pagination = createPaginationStore(1);
	isActive = $derived(this.pagination.hasNext);
	boardStore;
	// cachedList = $state([]);
	constructor(boardState: BoardState, username: string | null) {
		this.boardStore = new BoardStore(boardState, username);
	}

	setPage(page: number) {
		// TODO
		console.log(page);
		if (page === this.pagination.currentPage) return;
		if (page > this.pagination.currentPage) {
			for (let i = this.pagination.currentPage; i < page; i++) {
				this.next();
			}
		} else if (page < this.pagination.currentPage) {
			for (let i = this.pagination.currentPage; i > page; i--) {
				this.prev();
			}
		}
		// this.pagination.setPage(page);
	}

	setTotalPages(totalPages: number) {
		this.pagination.setTotalPages(totalPages);
		this.pagination.setPage(totalPages);
	}

	prev() {
		console.log(this.pagination.currentPage)
		if (!this.pagination.hasPrev) return;
		const prevMove = this.boardStore.moveHistory[this.pagination.currentPage - 1];
		const parsedPrevMove = parseMove(prevMove);

		const newBoard = $state.snapshot(this.boardStore.boardState);
		const newPlayer1HandList = $state.snapshot(this.boardStore.player1HandList);
		const newPlayer2HandList = $state.snapshot(this.boardStore.player2HandList);
		let newTurnColor = $state.snapshot(this.boardStore.turnColor);
		let newIsPlayer1Ready = $state.snapshot(this.boardStore.isPlayer1Ready);
		let newIsPlayer2Ready = $state.snapshot(this.boardStore.isPlayer2Ready);

		const { moveType } = parsedPrevMove;
		if (moveType === 'move') {
			const { fromCoord, toCoord, piece } = parsedPrevMove;
			newBoard[toCoord].splice(-1, 1);
			newBoard[fromCoord].push(piece);
			newTurnColor = newTurnColor === 'w' ? 'b' : 'w';
		} else if (moveType === 'attack') {
			const { fromCoord, toCoord, fromPiece, toPiece } = parsedPrevMove;
			newBoard[toCoord].splice(-1, 1);
			newBoard[fromCoord].push(fromPiece);
			newBoard[toCoord].push(toPiece);
			newTurnColor = newTurnColor === 'w' ? 'b' : 'w';
		} else if (moveType === 'place') {
			const { toCoord, piece } = parsedPrevMove;

			const handIndex = piece % 13;

			newBoard[toCoord].splice(-1, 1);
			if (this.boardStore.isPlayer1Ready == this.boardStore.isPlayer2Ready) {
				this.boardStore.turnColor === 'w' ? newPlayer2HandList[handIndex]++ : newPlayer1HandList[handIndex]++;
				newTurnColor = newTurnColor === 'w' ? 'b' : 'w';
			} else {
				this.boardStore.turnColor === 'w' ? newPlayer1HandList[handIndex]++ : newPlayer2HandList[handIndex]++;
			}
		} else if (moveType === 'ready') {
			//TODO
			const { playerColor } = parsedPrevMove;
			playerColor === 'w' ? (newIsPlayer1Ready = false) : (newIsPlayer2Ready = false);
			newTurnColor = newTurnColor === 'w' ? 'b' : 'w';
		}

		const newFen = serializeFen(
			newBoard,
			newPlayer1HandList,
			newPlayer2HandList,
			newTurnColor,
			newIsPlayer1Ready,
			newIsPlayer2Ready
		);

		this.boardStore.updateCurrentState(newFen);

		this.pagination.prev();
	}

	next() {
		if (!this.pagination.hasNext) return;
		const nextMove = this.boardStore.moveHistory[this.pagination.currentPage];
		const parsedNextMove = parseMove(nextMove);

		const newBoard = $state.snapshot(this.boardStore.boardState);
		const newPlayer1HandList = $state.snapshot(this.boardStore.player1HandList);
		const newPlayer2HandList = $state.snapshot(this.boardStore.player2HandList);
		let newTurnColor = $state.snapshot(this.boardStore.turnColor);
		let newIsPlayer1Ready = $state.snapshot(this.boardStore.isPlayer1Ready);
		let newIsPlayer2Ready = $state.snapshot(this.boardStore.isPlayer2Ready);

		const { moveType } = parsedNextMove;
		if (moveType === 'move') {
			const { fromCoord, toCoord, piece } = parsedNextMove;
			newBoard[fromCoord].splice(-1, 1);
			newBoard[toCoord].push(piece);
			newTurnColor = newTurnColor === 'w' ? 'b' : 'w';
		} else if (moveType === 'attack') {
			const { fromCoord, toCoord, fromPiece } = parsedNextMove;
			newBoard[fromCoord].splice(-1, 1);
			newBoard[toCoord].splice(-1, 1);
			newBoard[toCoord].push(fromPiece);
			newTurnColor = newTurnColor === 'w' ? 'b' : 'w';
		} else if (moveType === 'place') {
			const { piece, toCoord } = parsedNextMove;

			const handIndex = piece % 13;

			newBoard[toCoord].push(piece);
			if (this.boardStore.isPlayer1Ready == this.boardStore.isPlayer2Ready) {
				this.boardStore.turnColor === 'w' ? newPlayer1HandList[handIndex]-- : newPlayer2HandList[handIndex]--;
				newTurnColor = newTurnColor === 'w' ? 'b' : 'w';
			} else {
				this.boardStore.turnColor === 'w' ? newPlayer1HandList[handIndex]-- : newPlayer2HandList[handIndex]--;

			}
		} else if (moveType === 'ready') {
			//TODO
			const { playerColor } = parsedNextMove;
			playerColor === 'w' ? (newIsPlayer1Ready = true) : (newIsPlayer2Ready = true);
			newTurnColor = newTurnColor === 'w' ? 'b' : 'w';
		}

		const newFen = serializeFen(
			newBoard,
			newPlayer1HandList,
			newPlayer2HandList,
			newTurnColor,
			newIsPlayer1Ready,
			newIsPlayer2Ready
		);

		this.boardStore.updateCurrentState(newFen);

		this.pagination.next();
	}
}

export function setReplayStore(boardState: BoardState, username: string | null) {
	const store = new ReplayStore(boardState, username);
	setContext('replay', store);
	return store;
}

export function getReplayStore() {
	return getContext<ReplayStore>('replay');
}
