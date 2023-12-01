<script lang="ts">
	import { invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import type { LayoutData } from './$types';
	import Navbar from '$lib/components/Navbar.svelte';
	import Notifications from '$lib/components/Notifications.svelte';
	import { ws } from '$lib/store/websocket';

	export let data: LayoutData;

	let { supabase, session } = data;
	$: ({ supabase, session } = data);

	onMount(() => {
		const {
			data: { subscription },
		} = supabase.auth.onAuthStateChange((event, _session) => {
			if (_session?.expires_at !== session?.expires_at) {
				invalidate('supabase:auth');
			}
		});
		const token = session?.access_token;
		if (token) {
			//TODO websocket as spectator
			ws?.subscribe((val) => {
				if (val === 'connecting') {
					ws?.authenticate(token);
				}
			});
		}

		return () => {
			subscription.unsubscribe();
			ws?.close();
		};
	});
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>
<Notifications />
<Navbar {session} />
<slot />

<style lang='scss' global>
	@import '../main.scss';
	@import '../normalize.css';
</style>