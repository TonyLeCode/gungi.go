// import { VITE_API_URL } from './../../../../.svelte-kit/ambient.d';
import type { BoardState } from '$lib/store/gameState.js';
import { error } from '@sveltejs/kit';

export async function load({ fetch, params }) {
	const url = `http://${import.meta.env.VITE_API_URL}/getgame/${params.id}`;
	const res = await fetch(url);
	if (!res.ok){
		throw error(500,{
			message: 'Internal Server Error',
		})
	}
	const data: BoardState = await res.json();
	return {
		data,
	};
}
