import { CheckCircle2Icon, InfoIcon, AlertTriangleIcon, XCircleIcon } from '@lucide/svelte';
import { getContext, setContext } from 'svelte';
import { setSnackbarCtx } from '$lib/components/snackbar/snackbar-provider.svelte';

export type DialogType = 'success' | 'info' | 'error' | 'warn';

interface DialogOptions {
	title: string;
	body: string;
	type: DialogType;
	variant: 'alert' | 'confirm';
	resolve: (value: boolean) => void;
}

export const dialogIcons: Record<DialogType, any> = {
	success: CheckCircle2Icon,
	info: InfoIcon,
	warn: AlertTriangleIcon,
	error: XCircleIcon
};

class DialogManager {
	active = $state<DialogOptions | null>(null);

	// Variant 1: Simple OK message
	async alert(title: string, body: string, type: DialogType = 'info'): Promise<void> {
		return new Promise((resolve) => {
			this.active = {
				title,
				body,
				type,
				variant: 'alert',
				resolve: () => {
					this.active = null;
					resolve();
				}
			};
		});
	}

	// Variant 2: OK/Cancel returns boolean
	async confirm(title: string, body: string, type: DialogType = 'warn'): Promise<boolean> {
		return new Promise((resolve) => {
			this.active = {
				title,
				body,
				type,
				variant: 'confirm',
				resolve: (val: boolean) => {
					this.active = null;
					resolve(val);
				}
			};
		});
	}
}

export const dialogManager = new DialogManager();

const DIALOG_KEY = Symbol('dialog');

export function setDialogCtx() {
	return setContext(DIALOG_KEY, {
		alert: (title: string, body: string, type?: DialogType) =>
			dialogManager.alert(title, body, type),
		confirm: (title: string, body: string, type?: DialogType) =>
			dialogManager.confirm(title, body, type)
	});
}

export function getDialogCtx() {
	return getContext<ReturnType<typeof setDialogCtx>>(DIALOG_KEY);
}
