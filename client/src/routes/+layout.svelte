<script lang="ts">
	import { invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import type { LayoutData } from './$types';
	import Navbar from '$lib/components/Navbar.svelte';
  import Notifications from '$lib/components/Notifications.svelte';
	import {ws, wsConnState, websocketConnect, df} from '$lib/store/websocket'

	export let data: LayoutData;

	$: ({ supabase, session } = data);
	// $: console.log($wsConnState)
	// $: console.log(supabase)

	onMount(() => {
		const { data } = supabase.auth.onAuthStateChange((event, _session) => {
			if (_session?.expires_at !== session?.expires_at) {
				invalidate('supabase:auth');
			}
		});

		if (session?.access_token) {
			websocketConnect(`ws://${import.meta.env.VITE_API_URL}/ws`, session.access_token)
		}

		return () => {
			data.subscription.unsubscribe();
			$ws?.close()
		}
	});
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>
<Notifications />
<Navbar {session} />
<slot />