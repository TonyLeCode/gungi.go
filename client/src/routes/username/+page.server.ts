// import { AuthApiError } from '@supabase/supabase-js';
import type { PageServerLoad } from './$types';
import { redirect, type Actions, fail } from '@sveltejs/kit';
import { z } from 'zod';
import { superValidate } from 'sveltekit-superforms/server';

const schema = z.object({
	username: z.string().min(3).max(24),
});

//TODO custom error message
//TODO should say "Username must contain at least 3 character(s)" instead of "String must contain at least 3 character(s)"
export const load: PageServerLoad = async ({ locals: { getSession }, url }) => {
	const onboard = url.searchParams.get('onboard');

	const session = await getSession();
	if (!session) {
		throw redirect(308, '/');
	}

	if (onboard) {
		const fetchUrl = `http://${import.meta.env.VITE_API_URL}/user/onboarding`;
		const token = session.access_token;
		const options = {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${token}`,
			},
		};
		fetch(fetchUrl, options);
	}

	const form = await superValidate(schema);
	return { form };
};

export const actions: Actions = {
	default: async ({ locals, request }) => {
		const form = await superValidate(request, schema);
		if (!form.valid) {
			return fail(400, { form });
		}

		//TODO unique username validation on backend
		// return setError(form, 'username', 'Username already exists')

		// const { data, error } = await locals.supabase.auth.signUp({
		// 	email: form.data.email,
		// 	password: form.data.password,
		// 	options: {
		// 		data: {
		// 			username: form.data.username,
		// 		},
		// 	},
		// });

		// if (error) {
		// 	console.log(error);
		// 	if (error instanceof AuthApiError && error.status === 400) {
		// 		return fail(400, {
		// 			error: 'Invalid Registration',
		// 		});
		// 	}
		// 	return fail(500, {
		// 		error: 'Server error. Try again later.',
		// 	});
		// } else {
		// 	console.log('registered: ', data);
		// }

		return { form };
	},
};
