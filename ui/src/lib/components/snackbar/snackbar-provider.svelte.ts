import type { Component } from 'svelte';
import { getContext, setContext } from 'svelte';
import { CircleCheck, CircleX, InfoIcon, TriangleAlert } from '@lucide/svelte';

export type SnackbarType = 'success' | 'info' | 'error' | 'warn';

export interface Toast {
	id: string;
	message: string;
	type: SnackbarType;
	duration: number;
}

export const toastIcons: Record<SnackbarType, Component<any>> = {
	success: CircleCheck,
	info: InfoIcon,
	warn: TriangleAlert,
	error: CircleX
};

class SnackbarManager {
	toasts = $state<Toast[]>([]);

	push = (message: string, type: SnackbarType = 'info', duration = 4000) => {
		const id = crypto.randomUUID();
		this.toasts.push({ id, message, type, duration });

		if (duration > 0) {
			setTimeout(() => this.remove(id), duration);
		}
	};

	remove = (id: string) => {
		this.toasts = this.toasts.filter((t) => t.id !== id);
	};
}

const SNACK_KEY = Symbol('snackbar');

export function setSnackbarCtx() {
	const manager = new SnackbarManager();
	return setContext(SNACK_KEY, manager);
}

export function getSnackbarCtx() {
	const ctx = getContext<SnackbarManager>(SNACK_KEY);
	if (!ctx) throw new Error('getSnackbarCtx must be used within a provider');
	return ctx;
}
