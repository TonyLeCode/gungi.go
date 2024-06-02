import { error } from '@sveltejs/kit';
import { dev } from '$app/environment';

export async function load({ fetch, params }) {
	const url = dev ? `http://${import.meta.env.VITE_API_URL}/game/${params.id}` : `https://${import.meta.env.VITE_API_URL}/game/${params.id}`;
	const res = await fetch(url);
	if (!res.ok) {
		error(500, {
        			message: 'Internal Server Error',
        		});
	}
	const data = await res.json();
	return {
		gameData: data,
	};
}
