import { writable } from 'svelte/store';

export type notificationType = {
	id: string;
	type: string;
	title: string;
	msg: string | null;
	component?: ConstructorOfATypedSvelteComponent;
};

export const notifications = writable<notificationType[]>([]);

export function AddNotification(notification: notificationType) {
	notifications.update((items) => {
		const newNotifications = [...items, notification];
		return newNotifications;
	});
}

export function RemoveNotification(id: string) {
	notifications.update((items) => {
		const filtered = items.filter((item) => item.id !== id);
		return filtered;
	});
}
