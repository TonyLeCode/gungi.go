import { redirect, error, type RequestHandler } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ locals: { supabase }, url }) => {
	const supabaseResponse = await supabase.auth.signInWithOAuth({
		provider: 'google',
		options: {
			redirectTo: `${url.origin}/auth/callback`,
		},
	});

	console.log(decodeURIComponent(supabaseResponse.data.url || ''));

	if (supabaseResponse.error) {
		console.error(supabaseResponse.error);
		error(500, {
			message: 'Server error. Try again later.',
		});
	}

	throw redirect(303, supabaseResponse.data.url);
	// const fetchUrl = dev
	// 	? `http://${import.meta.env.VITE_API_URL}/user/onboarding`
	// 	: `https://${import.meta.env.VITE_API_URL}/user/onboarding`;
	// const token = supabaseResponse.data.session.access_token;
	// const options = {
	// 	headers: {
	// 		Authorization: `Bearer ${token}`,
	// 	},
	// };
	// const res = await fetch(fetchUrl, options);
	// if (res.ok) {
	// 	const hasOnboarded = await res.json();
	// 	if (!hasOnboarded) redirect(308, '/username?onboard=true');
	// }

	// redirect(303, '/overview');
};
