import { error } from '@sveltejs/kit';

export async function load({ fetch, params }) {
	const url = `http://${import.meta.env.VITE_API_URL}/game/${params.id}`;
	const res = await fetch(url);
	if (!res.ok) {
		throw error(500, {
			message: 'Internal Server Error',
		});
	}
	const data = await res.json();
	return {
		gameData: data,
	};
}
