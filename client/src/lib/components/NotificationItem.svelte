<script lang="ts">
	import { type notificationType } from '$lib/store/notificationStore.svelte';
	import { tweened } from 'svelte/motion';
	import { linear } from 'svelte/easing';

	// export let notification: notificationType;
	let { notification, removeItem }: { notification: notificationType; removeItem: () => void } = $props();

	let duration = tweened(100, {
		duration: 5000,
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
	<div class="container">
		<div class="title">
			{notification.title}
		</div>
		{#if notification.component}
			<svelte:component this={notification.component} />
		{:else}
			<div>
				{@html notification.msg}
			</div>
		{/if}
	</div>
	<button class="close" onclick={removeItem}
		><img draggable="false" src="/closeCircle.svg" alt="exit dialog" width="35px" height="35px" /></button
	>
</div>

<style lang="scss">
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
		min-width: 21rem;
		min-height: 4rem;
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.05);
		display: flex;
		justify-content: center;
		align-items: center;
		background-color: rgb(var(--bg-2));
		border-radius: 8px;
	}
	.close {
		margin-left: auto;
		border-radius: 50%;
		&:focus {
			outline: 2px solid rgba(var(--primary), 0.5);
			outline-offset: -1px;
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
</style>
