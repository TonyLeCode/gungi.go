<script lang="ts">
	import type { Session } from '@supabase/supabase-js';

	export let session: Session | null;
	$: session = session;
</script>

<nav class="navbar">
	<div class="nav-inner">
		<a class="brand" href="/">Gungi</a>
		<ul class="nav-list">
			{#if session}
				<a href="/overview">overview</a>
				<a href="/play">play</a>
				<!-- <a href="/games">games</a> -->
				<!-- TODO learning, puzzles, resources, library -->
			{/if}
			<!-- <a href="/rules">rules</a> -->
		</ul>
		<ul class="nav-account">
			{#if session}
				<!-- TODO settings, friends, notifications  -->
				<span class="name">{session.user.user_metadata.username}</span>
				<a href="/logout">logout</a>
			{:else}
				<a class="a" href="/login">login</a>
				<a class="a" href="/register">register</a>
			{/if}
		</ul>
	</div>
</nav>

<style>
	nav {
		user-select: none;
	}
	a {
		padding: 0 0.5rem;
	}
	.brand {
		font-weight: 600;
	}
	.name {
		color: rgb(var(--primary));
		margin-right: 1rem;
		position: relative;
		font-weight: 600;
	}

	.navbar {
		padding: 0.75rem 1rem;
		max-width: 120rem;
		margin: auto;
	}
	.nav-inner {
		display: flex;
		max-width: 96rem;
		margin: auto;
	}
	.nav-list {
		display: none;
	}
	.nav-account {
		display: flex;
		gap: 1rem;
		margin-left: auto;
	}

	@media only screen and (min-width: 800px) {
		.nav-list {
			display: flex;
			gap: 0.5rem;
		}
	}
</style>
