import { error, redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ locals: { supabase } }) => {
	const { error: err } = await supabase.auth.signOut();

	if (err) {
		console.error(err);
		error(500, 'something went wrong');
	}

	console.log('success');
	redirect(303, '/');
};
