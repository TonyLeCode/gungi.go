import { redirect } from '@sveltejs/kit';

export const GET = async (event) => {
	const {
		url,
		locals: { supabase },
	} = event;
	const redirectUrl = '/overview';
	const code = url.searchParams.get('code') as string;

	if (code) {
		const { error } = await supabase.auth.exchangeCodeForSession(code);
		if (!error) {

			const next = url.searchParams.get('next') ?? redirectUrl;
			throw redirect(303, `/${next.slice(1)}`);
		}
		console.log(error);
	}

	// TODO error page
	// return the user to an error page with instructions
	throw redirect(303, '/auth/auth-code-error');
};
