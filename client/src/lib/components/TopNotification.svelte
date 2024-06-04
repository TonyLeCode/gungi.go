<script lang="ts">
	import { getTopNotificationStore } from '$lib/store/notificationStore.svelte';
	import { getWebsocketStore } from '$lib/store/websocketStore.svelte';
	import { CircleAlert } from 'lucide-svelte';
	import NotificationItem from './NotificationItem.svelte';
	import { onDestroy } from 'svelte';

	let topNotificationStore = getTopNotificationStore();
	let websocketStore = getWebsocketStore();

	let timer: number | null = $state(null);
	let isExpanded = $state(true);

	$effect(() => {
		if (websocketStore.state === 'closed') {
			topNotificationStore.SetNotification('You are disconnected! Please refresh or try again later.');
			$effect.root(() => {
				if (timer) {
					window.clearTimeout(timer);
					timer = null;
				}
			});
			timer = window.setTimeout(() => {
				isExpanded = false;
				timer = null;
			}, 5000);
		}
	});

	onDestroy(() => {
		if (timer) window.clearTimeout(timer);
	});
</script>

{#if topNotificationStore.notification !== ''}
	<div class="fly-up-fade" class:contract={!isExpanded}>
		<CircleAlert style="flex-shrink: 0" size="30px" />
		<span class="text">{topNotificationStore.notification}</span>
	</div>
{/if}

<style lang="scss">
	div {
		position: fixed;
		top: 1rem;
		text-align: center;
		left: 5px;
		right: 5px;
		display: flex;
		justify-content: center;
		align-items: center;
		background-color: rgb(var(--primary));
		color: rgb(var(--bg));
		padding: 1rem;
		margin: 0 auto;
		max-width: 60rem;
		border-radius: 4px;
		z-index: 5;
		box-shadow:
			0px 5px 15px rgba(0, 0, 0, 0.07),
			0px 2px 5px rgba(0, 0, 0, 0.05);
		animation-duration: 400ms;
		transition:
			// max-width 50ms,
			padding 200ms,
			margin 250ms;
		overflow: hidden;
		color: white;
		max-height: 5rem;
	}

	.contract {
		max-width: 3rem;
		margin-top: 0.25rem;
		padding: 0.5rem;
		.text {
			display: none;
		}
		&:hover {
			margin-top: 0;
			max-width: 60rem;
			padding: 1rem;
			.text {
				display: block;
			}
		}
	}

	.text {
		margin-left: 0.5rem;
		display: block;
		// text-wrap: nowrap;
		@media (min-width: 767px) {
			text-wrap: nowrap;
		}
	}
</style>
