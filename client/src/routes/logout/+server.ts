import { error, redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ locals }) => {
	const { error: err } = await locals.supabase.auth.signOut();

	if (err) {
		console.log(err);
		throw error(500, 'something went wrong');
	}

	console.log('success');
	throw redirect(303, '/');
};
