import { dev } from '$app/environment';
import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from '$env/static/public';
import { createServerClient } from '@supabase/ssr';
import { redirect, type Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';

const supabase: Handle = async ({ event, resolve }) => {
	event.locals.supabase = createServerClient(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY, {
		cookies: {
			getAll: () => event.cookies.getAll(),
			setAll: (cookiesToSet) => {
				cookiesToSet.forEach(({ name, value, options }) => {
					event.cookies.set(name, value, { ...options, path: '/' });
				});
			},
		},
	});

	/**
	 * Unlike `supabase.auth.getSession()`, which returns the session _without_
	 * validating the JWT, this function also calls `getUser()` to validate the
	 * JWT before returning the session.
	 */
	event.locals.safeGetSession = async () => {
		const {
			data: { session },
		} = await event.locals.supabase.auth.getSession();
		if (!session) {
			return { session: null, user: null, username: null };
		}

		const {
			data: { user },
			error,
		} = await event.locals.supabase.auth.getUser();
		if (error) {
			return { session: null, user: null, username: null };
		}

		const fetchUrl = dev
			? `http://${import.meta.env.VITE_API_URL}/username`
			: `https://${import.meta.env.VITE_API_URL}/username`;
		const token = session.access_token;
		const options = {
			headers: {
				Authorization: `Bearer ${token}`,
			},
		};
		const res = await fetch(fetchUrl, options);
		if (!res.ok) throw new Error('Something went wrong in the server. Try again later.');

		const username = await res.json();
		// if (typeof username !== 'string' || !username && event.url.pathname !== '/username') redirect(308, '/username');

		return { session, user, username };
	};

	return resolve(event, {
		filterSerializedResponseHeaders(name) {
			/**
			 * Supabase libraries use the `content-range` and `x-supabase-api-version`
			 * headers, so we need to tell SvelteKit to pass it through.
			 */
			return name === 'content-range' || name === 'x-supabase-api-version';
		},
	});
};

const authGuard: Handle = async ({ event, resolve }) => {
	const { session, user, username } = await event.locals.safeGetSession();
	event.locals.session = session;
	event.locals.user = user;
	event.locals.username = username;

	const protectedPages = ['/overview', '/play'];

	if (!event.locals.session && protectedPages.includes(event.url.pathname)) {
		throw redirect(303, '/');
	}

	if (event.locals.session && event.url.pathname === '/login') {
		throw redirect(303, '/overview');
	}

	if (!event.locals.session) {
		return resolve(event);
	}

	if (
		(typeof username !== 'string' || !username) &&
		event.url.pathname !== '/username' &&
		event.url.pathname !== '/logout'
	) {
		throw redirect(308, '/username');
	}

	return resolve(event);
};

export const handle: Handle = sequence(supabase, authGuard);
