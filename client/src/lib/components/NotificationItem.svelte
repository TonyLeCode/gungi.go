<script lang="ts">
	import { notifications, type notificationType } from '$lib/store/notification';
	import { tweened } from 'svelte/motion';
	import { linear } from 'svelte/easing';

	export let notification: notificationType;
	// console.log(notification);
	let duration = tweened(100, {
		duration: 5000,
		easing: linear,
	});
	duration.set(0);

	$: if ($duration === 0) {
		notifications?.remove(notification.id);
	}
</script>

<div class="notification">
	<div style="width: {$duration}%" class="bar" />
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
	<button
		class="close"
		on:click={() => {
			notifications?.remove(notification.id);
		}}><img draggable="false" src="/closeCircle.svg" alt="exit dialog" width="35px" height="35px" /></button
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
		// border-bottom-right-radius: 8px;
		// border-bottom-left-radius: 8px;
	}
	.notification {
		// border: 1px blue dashed;
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
		// overflow: hidden;
		// border: 2px rgba(var(--primary), 0.5) solid;
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
