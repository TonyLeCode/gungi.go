import { dev } from '$app/environment';
import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ locals: { supabase }, url }) => {
	const {
		data: { session },
	} = await supabase.auth.getSession();
	const token = session?.access_token;
	const reqUrl = dev
		? `http://${import.meta.env.VITE_API_URL}/getongoinggamelist`
		: `https://${import.meta.env.VITE_API_URL}/getongoinggamelist`;
	const options = {
		headers: {
			Authorization: `Bearer ${token}`,
		},
	};

	const offset = url.searchParams.get('offset');
	const res = await fetch(`${reqUrl}?offset=${offset}`, options);

	if (!res.ok) {
		error(500, {
			message: 'Internal Server Error',
		});
	}
	const data = await res.json()

	return new Response(JSON.stringify(data));

	// const res = await fetch(`${url}?offset=${10 * (paginationStore.currentPage - 1)}`, options);
};
