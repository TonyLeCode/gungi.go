// import { VITE_API_URL } from './../../../.svelte-kit/ambient.d';
// import { getServerSession } from '@supabase/auth-helpers-sveltekit'

import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export interface Game {
	completed: boolean;
	current_state: string;
	date_started: Date;
	fen: {
		String: string;
		Valid: boolean;
	};
	id: string;
	username1: string;
	username2: string;
}

export const load: PageServerLoad = async ({ locals: { getSession }, fetch }) => {
	const session = await getSession();
	const token = session?.access_token;
	// console.log('token', session?.access_token);
	// const url = 'http://localhost:8080/getongoinggamelist';
	const url = `http://${import.meta.env.VITE_API_URL}/getongoinggamelist`;
	const options = {
		headers: {
			Authorization: `Bearer ${token}`,
		},
	};

	const res = await fetch(url, options);
	if (!res.ok) {
		throw error(500, {
			message: 'Internal Server Error',
		});
	}
	const data: Game = await res.json();

	return {
		data: data ?? ([] as Game[]),
	};
};
