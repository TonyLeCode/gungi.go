import { error } from '@sveltejs/kit';
import { dev } from '$app/environment';

type GameData = {
	public_id: string;
	fen: string | null;
	history: string;
	completed: boolean;
	date_started: Date;
	date_finished: Date | null;
	current_state: string;
	ruleset: string;
	result: string | null;
	type: string;
	player1: string;
	player2: string;
	moveList: { [key: string]: number[] };
	undo_requests: [{ sender_username: string; receiver_username: string; status: string }];
};

export async function load({ fetch, params }) {
	const url = dev
		? `http://${import.meta.env.VITE_API_URL}/game/${params.id}`
		: `https://${import.meta.env.VITE_API_URL}/game/${params.id}`;
	const res = await fetch(url);
	if (!res.ok) {
		console.log("errrr")
		error(500, {
			message: 'Internal Server Error',
		});
	}
	const data = await res.json();
	return {
		gameData: data as GameData,
	};
}
