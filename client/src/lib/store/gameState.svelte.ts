import { reverseList } from '$lib/helpers';
import { FenToBoard } from '$lib/utils/utils';
import { getContext, setContext } from 'svelte';

export interface BoardState {
	completed: boolean;
	current_state: string;
	date_finished: Date | null;
	date_started: Date;
	fen: string | null;
	result: string | null;
	history: string;
	public_id: string;
	moveList: { [key: string]: number[] };
	player1: string;
	player2: string;
	ruleset: string;
	type: string;
	// check: string
}

export class BoardStore {
	completed = $state(false);
	current_state = $state('');
	date_finished = $state<Date | null>(null);
	date_started = $state<Date>();
	fen = $state<string | null>(null);
	result = $state<string | null>('');
	history = $state('');
	public_id = $state('');
	moveList = $state({} as { [key: string]: number[] });
	player1 = $state('');
	player2 = $state('');
	ruleset = $state('');
	type = $state('');

	private currentStateFields = $derived(this.current_state.split(' '));
	private pieces = $derived(this.currentStateFields[0]);
	private hands = $derived(this.currentStateFields[1].split('/'));
	username = $state('');

	player1HandList = $derived.by(() => {
		const hand = this.hands[0];
		const newHand: number[] = [];

		for (let i = 0; i < hand.length; i++) {
			newHand.push(Number(hand[i]));
		}
		return newHand;
	});
	player2HandList = $derived.by(() => {
		const hand = this.hands[1];
		const newHand: number[] = [];

		for (let i = 0; i < hand.length; i++) {
			newHand.push(Number(hand[i]));
		}
		return newHand;
	});
	turnColor = $derived(this.currentStateFields[2]);
	isPlayer1Ready = $derived(this.currentStateFields[3][0] === '1');
	isPlayer2Ready = $derived(this.currentStateFields[3][1] === '1');
	player1ArmyCount = $derived(this.pieces.match(/[A-Z]/g)?.length || 0);
	player2ArmyCount = $derived(this.pieces.match(/[a-z]/g)?.length || 0);
	player1HandCount = $derived(this.player1HandList.reduce((a, b) => a + b));
	player2HandCount = $derived(this.player2HandList.reduce((a, b) => a + b));

	userColor = $derived.by<'w' | 'b' | 'spectator'>(() => {
		if (this.username === this.player1) return 'w';
		if (this.username === this.player2) return 'b';
		return 'spectator';
	});
	moveHistory = $derived(this.history.split(' '));
	manualFlip = $state(false);
	isViewReversed = $derived.by(() => {
		if (this.username !== this.player1 && this.username !== this.player2) return this.manualFlip;
		const isUserWhite = this.username === this.player1;
		return this.manualFlip === isUserWhite;
	});
	isUserTurn = $derived(this.userColor === this.turnColor);
	moveListUI = $derived.by(() => {
		const transformedMoveList: { [key: string]: number[] } = {};
		if (this.isViewReversed) {
			for (const key in this.moveList) {
				const transformedKey = 80 - parseInt(key, 10);
				const transformedValue = this.moveList[key].map((value) => 80 - value);
				transformedMoveList[transformedKey] = transformedValue;
			}
			return transformedMoveList;
		} else {
			return this.moveList;
		}
	});
	boardState = $derived(FenToBoard(this.current_state));
	boardUI = $derived(this.isViewReversed ? reverseList(this.boardState) : this.boardState);

	constructor(initState: BoardState, username: string | null) {
		this.completed = initState.completed;
		this.current_state = initState.current_state;
		this.date_finished = initState.date_finished;
		this.date_started = initState.date_started;
		this.fen = initState.fen;
		this.result = initState.result;
		this.history = initState.history;
		this.public_id = initState.public_id;
		this.moveList = initState.moveList;
		this.player1 = initState.player1;
		this.player2 = initState.player2;
		this.ruleset = initState.ruleset;
		this.type = initState.type;

		this.username = username || '';
	}

	updateBoard(newState: BoardState) {
		this.completed = newState.completed;
		this.current_state = newState.current_state;
		this.date_finished = newState.date_finished;
		this.date_started = newState.date_started;
		this.fen = newState.fen;
		this.result = newState.result;
		this.history = newState.history;
		this.public_id = newState.public_id;
		this.moveList = newState.moveList;
		this.player1 = newState.player1;
		this.player2 = newState.player2;
		this.ruleset = newState.ruleset;
		this.type = newState.type;
	}

	updateCurrentState(fen: string) {
		this.current_state = fen;
	}
}

export function setGameStore(initState: BoardState, username: string | null) {
	const store = new BoardStore(initState, username);
	setContext('gameState', store);
	return store;
}

export function getGameStore() {
	return getContext<BoardStore>('gameState');
}
