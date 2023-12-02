import { browser } from '$app/environment';
import { writable } from 'svelte/store';

export type notificationType = {
	id: string;
	type: string;
	title: string;
	msg: string | null;
	component?: ConstructorOfATypedSvelteComponent;
};

function createNotificationStore() {
	if (!browser) return;
	const notifications = writable<notificationType[]>([]);

	function AddNotification(notification: notificationType) {
		notifications.update((items) => {
			const newNotifications = [...items, notification];
			return newNotifications;
		});
	}

	function RemoveNotification(id: string) {
		notifications.update((items) => {
			const filtered = items.filter((item) => item.id !== id);
			return filtered;
		});
	}
	return {
		store: notifications,
		add: AddNotification,
		remove: RemoveNotification,
	};
}

export const notifications = createNotificationStore();

function createTopNotificationStore() {
	if (!browser) return;
	const notification = writable<string>('');

	function SetNotification(text: string) {
		notification.set(text);
	}

	return {
		store: notification,
		set: SetNotification,
	};
}

export const topNotification = createTopNotificationStore();
