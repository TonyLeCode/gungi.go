import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { dev } from '$app/environment';

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
	const url = dev ? `http://${import.meta.env.VITE_API_URL}/getongoinggamelist` : `https://${import.meta.env.VITE_API_URL}/getongoinggamelist`;
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
