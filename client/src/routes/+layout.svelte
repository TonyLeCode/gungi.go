<script lang="ts">
	import { goto, invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import type { LayoutData } from './$types';
	import Navbar from '$lib/components/Navbar.svelte';
	import Notifications from '$lib/components/Notifications.svelte';
	import { ws } from '$lib/store/websocket';
	import TopNotification from '$lib/components/TopNotification.svelte';
	import { topNotification } from '$lib/store/notification';

	// export let data: LayoutData;
	let { data, children } = $props();

	// let { supabase, session } = data;
	// $: ({ supabase, session } = data);
	// $: $ws === 'connected' && session && ws?.authenticate(session.access_token);
	// $: notifStore = topNotification?.store;

	onMount(() => {
		const {
			data: { subscription },
		} = data.supabase.auth.onAuthStateChange((_, newSession) => {
			if (!newSession) {
				/**
				 * Queue this as a task so the navigation won't prevent the
				 * triggering function from completing
				 */
				setTimeout(() => {
					// goto('/', { invalidateAll: true });
				});
			}
			if (newSession?.expires_at !== data.session?.expires_at) {
				invalidate('supabase:auth');
			}
		});

		return () => {
			subscription.unsubscribe;
		};
	});
	// onMount(() => {
	// 	const {
	// 		data: { subscription },
	// 	} = supabase.auth.onAuthStateChange((event, _session) => {
	// 		if (_session?.expires_at !== session?.expires_at) {
	// 			invalidate('supabase:auth');
	// 		}
	// 	});

	// 	return () => {
	// 		subscription.unsubscribe();
	// 		ws?.close();
	// 	};
	// });
</script>

<svelte:head>
	<title>White Monarch Server</title>
</svelte:head>
<!-- {#if $notifStore !== '' && $notifStore !== undefined}
	<TopNotification text={$notifStore} />
{/if} -->
<Notifications />
<Navbar session={data.session} />
<!-- <slot /> -->
{@render children()}

<style lang="scss" global>
	@import '../main.scss';
	@import '../normalize.css';
</style>
