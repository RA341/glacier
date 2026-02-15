import { getContext, setContext } from 'svelte';
import { InfoIcon } from '@lucide/svelte';
import CircleCheckBig from '@lucide/svelte/icons/circle-check-big';
import TriangleAlert from '@lucide/svelte/icons/triangle-alert';
import CircleX from '@lucide/svelte/icons/circle-x';

export type SnackbarType = 'success' | 'info' | 'warn' | 'error';

export interface Toast {
	id: string;
	message: string;
	type: SnackbarType;
	duration: number;
	remaining: number;
	startTime: number;
}

export const toastIcons = {
	success: CircleCheckBig,
	info: InfoIcon,
	warn: TriangleAlert,
	error: CircleX
};

const SNACKBAR_KEY = Symbol('snackbar');

export function setSnackbarCtx() {
	let toasts = $state<Toast[]>([]);
	const timers = new Map<string, any>();

	const remove = (id: string) => {
		clearTimeout(timers.get(id));
		timers.delete(id);
		toasts = toasts.filter((t) => t.id !== id);
	};

	const push = (message: string, type: SnackbarType = 'info', duration = 5000) => {
		const id = crypto.randomUUID();
		const newToast: Toast = {
			id,
			message,
			type,
			duration,
			remaining: duration,
			startTime: Date.now()
		};

		toasts.push(newToast);

		const timer = setTimeout(() => remove(id), duration);
		timers.set(id, timer);
	};

	const pause = (id: string) => {
		const toast = toasts.find((t) => t.id === id);
		if (!toast || !timers.has(id)) return;

		clearTimeout(timers.get(id));
		timers.delete(id);

		const elapsed = Date.now() - toast.startTime;
		toast.remaining = Math.max(0, toast.remaining - elapsed);
	};

	const resume = (id: string) => {
		const toast = toasts.find((t) => t.id === id);
		if (!toast || timers.has(id)) return;

		toast.startTime = Date.now();

		const timer = setTimeout(() => remove(id), toast.remaining);
		timers.set(id, timer);
	};

	const context = {
		get toasts() {
			return toasts;
		},
		push,
		remove,
		pause,
		resume
	};

	setContext(SNACKBAR_KEY, context);
	return context;
}

export function getSnackbarCtx() {
	const ctx = getContext<ReturnType<typeof setSnackbarCtx>>(SNACKBAR_KEY);
	if (!ctx) throw new Error('getSnackbarCtx must be used within a provider');
	return ctx;
}
