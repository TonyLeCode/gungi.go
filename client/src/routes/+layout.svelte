<script lang="ts">
	import { invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import type { LayoutData } from './$types';
	import Navbar from '$lib/components/Navbar.svelte';
  import Notifications from '$lib/components/Notifications.svelte';

	export let data: LayoutData;

	$: ({ supabase, session } = data);
	// $: console.log(supabase)

	onMount(() => {
		const { data } = supabase.auth.onAuthStateChange((event, _session) => {
			if (_session?.expires_at !== session?.expires_at) {
				invalidate('supabase:auth');
			}
		});

		return () => data.subscription.unsubscribe();
	});
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>
<Notifications />
<Navbar {session} />
<slot />