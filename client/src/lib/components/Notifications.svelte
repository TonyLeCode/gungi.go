<script lang="ts">
	import { getNotificationStore } from '$lib/store/notificationStore.svelte';
	import NotificationItem from './NotificationItem.svelte';
	import { fly } from 'svelte/transition';

	let notificationsStore = getNotificationStore();
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
