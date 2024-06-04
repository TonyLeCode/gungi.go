<script lang="ts">
	import { getNotificationStore, type notificationType } from '$lib/store/notificationStore.svelte';
	import { onMount } from 'svelte';
	import NotificationItem from './NotificationItem.svelte';
	import { nanoid } from 'nanoid';
	import { fly } from 'svelte/transition';

	let notificationsStore = getNotificationStore();
	// $: notificationsStore = notifications?.store ?? writable<notificationType[]>([])

	onMount(() => {
		notificationsStore.add({
			id: nanoid(),
			title: 'Game Accepted',
			type: 'success',
			msg: 'Go to <a class="a-primary" href="/play/lol">game<a>',
		} as notificationType);
		notificationsStore.add({
			id: nanoid(),
			title: 'Undo Request Rejected',
			type: 'warning',
			msg: 'Your undo request has been denied',
		} as notificationType);
		notificationsStore.add({
			id: nanoid(),
			title: 'Something went wrong',
			type: 'error',
			msg: 'Sorry',
		} as notificationType);
	});
</script>

<ul>
	{#each notificationsStore.list as notification (notification.id)}
		<li transition:fly|global={{ x: '16px', duration: 250 }}>
			<NotificationItem {notification} removeItem={() => notificationsStore.remove(notification.id)} />
		</li>
	{/each}
</ul>

<style lang="scss">
	ul {
		position: fixed;
		bottom: 2rem;
		max-width: 24rem;
		margin-left: auto;
		z-index: 10;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		right: 1rem;
		left: 1rem;
		font-size: 0.875rem;
		@media (min-width: 767px) {
			right: 2rem;
			left: 2rem;
			font-size: 1rem;
		}
	}
</style>
