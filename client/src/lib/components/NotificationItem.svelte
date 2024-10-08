<script lang="ts">
	import { type notificationType } from '$lib/store/notificationStore.svelte';
	import { tweened } from 'svelte/motion';
	import { linear } from 'svelte/easing';
	import { CircleAlert, CircleCheck, CircleSlash, CircleX } from 'lucide-svelte';

	//TODO merge as snippet with notification
	let { notification, removeItem }: { notification: notificationType; removeItem: () => void } = $props();

	let duration = tweened(100, {
		duration: 10000,
		easing: linear,
	});
	duration.set(0);

	$effect(() => {
		if ($duration === 0) {
			removeItem();
		}
	});
</script>

<div class="notification">
	<div style="width: {$duration}%" class="bar"></div>
	{#if notification.type === 'success'}
		<div class={`${notification.type} animation`}>
			<CircleCheck size="30px" strokeWidth="2.5px" />
		</div>
	{:else if notification.type === 'error'}
		<div class={`${notification.type} animation`}>
			<CircleSlash size="30px" strokeWidth="2.5px" />
		</div>
	{:else if notification.type === 'warning'}
		<div class={`${notification.type} animation`}>
			<CircleAlert size="30px" strokeWidth="2.5px" />
		</div>
	{/if}
	<div class="container">
		<div class="title">
			{notification.title}
		</div>
		{#if notification.component}
			{notification.component}
		{:else}
			<div>
				{@html notification.msg}
			</div>
		{/if}
	</div>
	<button class="close" onclick={removeItem}><CircleX size="35px" /></button>
</div>

<style lang="scss">
	.success {
		color: rgb(var(--success));
	}
	.error {
		color: rgb(var(--error));
	}
	.warning {
		color: rgb(var(--warning));
	}

	.bar {
		width: 100%;
		height: 5px;
		background-color: rgb(var(--primary));
		position: absolute;
		bottom: 0;
		left: 0;
	}
	.notification {
		position: relative;
		padding: 1rem 2rem;
		// min-width: 21rem;
		width: 100%;
		min-height: 4rem;
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.25);
		display: flex;
		justify-content: center;
		align-items: center;
		background-color: rgb(var(--bg-2));
		border-radius: 6px;
		gap: 1rem;
	}
	.close {
		margin-left: auto;
		border-radius: 50%;
		color: rgb(var(--primary));
		&:focus {
			outline: 2px solid rgba(var(--primary), 0.5);
			outline-offset: 1px;
		}
		&:hover {
			color: rgb(var(--primary-3));
		}
	}
	.container {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		margin-right: 1rem;
	}
	.title {
		font-weight: bold;
	}

	.animation {
		animation: pulse 1.8s normal;
	}

	@keyframes pulse {
		0% {
			transform: scale(1);
		}
		10% {
			transform: scale(1);
		}
		30% {
			transform: scale(1.4);
		}
		50% {
			transform: scale(1);
		}
		100% {
			transform: scale(1);
		}
	}
</style>
