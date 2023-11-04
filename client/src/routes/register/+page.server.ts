import { AuthApiError } from '@supabase/supabase-js';
import type { PageServerLoad } from './$types';
import { redirect, type Actions, fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals: { getSession } }) => {
	const session = await getSession();
	if (session) {
		throw redirect(308, '/overview');
	}
};

export const actions: Actions = {
	default: async ({ locals, request }) => {
		// const body = Object.fromEntries(await request.formData());
		const bodyData = await request.formData();
		const email = bodyData.get('email');
		const password = bodyData.get('password');
		const username = bodyData.get('username');

		if (!email) return fail(400, { description: 'must have an email' });
		if (!password) return fail(400, { description: 'must have an password' });
		if (!username) return fail(400, { description: 'must have an username' });

		const { data, error } = await locals.supabase.auth.signUp({
			email: email.toString(),
			password: password.toString(),
			options: {
				data: {
					username: username.toString(),
				},
			},
		});

		if (error) {
			console.log(error);
			if(error instanceof AuthApiError && error.status === 400){
				return fail(400, {
					error: "Invalid Registration"
				})
			}
			return fail(500, {
				error: 'Server error. Try again later.',
			});
		} else {
			console.log('registered: ', data);
		}
	},
};
