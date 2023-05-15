import type { Actions } from './$types';
import { AuthApiError } from '@supabase/supabase-js';
import type { PageServerLoad } from './$types';
import { fail, redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals: { getSession } }) => {
	const session = await getSession();
	if (session) {
		throw redirect(308, '/overview');
	}
};

export const actions = {
	default: async ({ locals, request }) => {
		const body = Object.fromEntries(await request.formData());

		const { data, error } = await locals.supabase.auth.signInWithPassword({
			email: body.email as string,
			password: body.password as string,
		});

		if (error) {
			console.log(error);
			return fail(500, {
				error: 'Server error. Try again later.',
			});
		} else {
			console.log('logged in', data);
		}
		throw redirect(303, '/overview');
	},
} satisfies Actions;
