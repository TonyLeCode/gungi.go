<script lang="ts">
	import type { Session } from '@supabase/supabase-js';

	export let session: Session | null;
	$: session = session;

	function switchTheme() {
		const theme = localStorage.getItem('theme');
		if (theme === 'dark') {
			localStorage.setItem('theme', 'light');
			document.documentElement.classList.remove('dark');
		} else {
			localStorage.setItem('theme', 'dark');
			document.documentElement.classList.add('dark');
		}
	}
</script>

<nav class="navbar">
	<div class="nav-inner">
		<a class="brand" href="/"><img src="/gungi-logo.svg" alt="logo"></a>
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
			// TODO theme switch button
			<button class="button-primary" on:click={switchTheme}>Theme</button>
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
		font-size: 0.85rem;
		@media (min-width: 767px) {
      font-size: 1rem;
    }
	}
	a {
		padding: 0 0.5rem;
	}
	.brand {
		/* font-weight: 600; */
		width: 45px;
	}
	.name {
		color: rgb(var(--primary));
		margin-right: 1rem;
		position: relative;
		font-weight: 600;
	}

	.navbar {
		padding: 0.375rem 0.5rem;
		max-width: 120rem;
		margin: auto;
		@media (min-width: 767px) {
			padding: 0.75rem 1rem;
		}
	}
	.nav-inner {
		display: flex;
		max-width: 96rem;
		margin: auto;
		align-items: center;
	}
	.nav-list {
		display: none;
	}
	.nav-account {
		display: none;
		gap: 1rem;
		margin-left: auto;
		@media (min-width: 767px) {
			display: flex;
		}
	}

	@media only screen and (min-width: 800px) {
		.nav-list {
			display: flex;
			gap: 0.5rem;
		}
	}
</style>
