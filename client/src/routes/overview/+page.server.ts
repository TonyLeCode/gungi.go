import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { dev } from '$app/environment';

export interface Game {
	public_id: string;
	fen: string;
	history: string;
	completed: boolean;
	date_started: Date | null;
	date_finished: Date | null;
	current_state: string;
	ruleset: string;
	result: string;
	type: string;
	username1: string;
	username2: string;
	moveList: string;
}

interface GetOverviewResponse {
	ongoingGames: Game[];
	completedGames: Game[];
	gameHistoryCount: number;
}

export const load: PageServerLoad = async ({ locals: { supabase }, fetch }) => {
	const {
		data: { session },
	} = await supabase.auth.getSession();
	const token = session?.access_token;
	const url = dev
		? `http://${import.meta.env.VITE_API_URL}/overview`
		: `https://${import.meta.env.VITE_API_URL}/overview`;
	const options = {
		headers: {
			Authorization: `Bearer ${token}`,
		},
	};

	const res = await fetch(url, options);
	if (!res.ok) {
		error(500, {
			message: 'Internal Server Error',
		});
	}
	const data: GetOverviewResponse = await res.json();

	return {
		data: data ?? ({ ongoingGames: [], completedGames: [], gameHistoryCount: 0 } as GetOverviewResponse),
	};
};
