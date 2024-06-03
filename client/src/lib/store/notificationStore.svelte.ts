import { browser } from '$app/environment';
import { getContext, setContext } from 'svelte';

export type notificationType = {
	id: string;
	type: string;
	title: string;
	msg: string | null;
	component?: ConstructorOfATypedSvelteComponent;
};

class NotificationStore {
	list = $state<notificationType[]>([]);
	constructor() {
		if (!browser) return;
	}

	add(notification: notificationType) {
		this.list.push(notification);
	}

	remove(id: string) {
		this.list = this.list.filter((item) => item.id !== id);
	}
}

export function setNotificationStore() {
	const store = new NotificationStore();
	setContext('notifications', store);
	return store;
}

export function getNotificationStore() {
	return getContext<NotificationStore>('notifications');
}

class TopNotificationStore {
	notification = $state('');

	SetNotification(text: string) {
		this.notification = text;
	}
}

export function setTopNotificationStore() {
	const store = new TopNotificationStore();
	setContext('top-notification', store);
	return store;
}

export function getTopNotificationStore() {
	return getContext<TopNotificationStore>('top-notification');
}
