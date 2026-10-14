<script lang="ts">
	import { PlusIcon, Trash2Icon } from '@lucide/svelte';
	import { fly } from 'svelte/transition';
	import type { FieldSchema } from '#lib/gen/service_config/v1/service_config_pb.js';

	interface Props {
		fields: FieldSchema[];
		formData: Record<string, any>;
	}

	let { fields, formData = $bindable() }: Props = $props();
	let tempInputs = $state<Record<string, string>>({});

	function addMapEntry(fieldKey: string) {
		const k = tempInputs[`${fieldKey}_key`]?.trim();
		const v = tempInputs[`${fieldKey}_val`]?.trim();
		if (k && v !== undefined) {
			const currentMap = formData[fieldKey] || {};
			formData[fieldKey] = { ...currentMap, [k]: v };
			tempInputs[`${fieldKey}_key`] = '';
			tempInputs[`${fieldKey}_val`] = '';
		}
	}

	function removeMapEntry(fieldKey: string, entryKey: string) {
		if (!formData[fieldKey]) return;
		const { [entryKey]: _, ...rest } = formData[fieldKey];
		formData[fieldKey] = rest;
	}

	function addArrayItem(fieldKey: string) {
		const v = tempInputs[`${fieldKey}_val`]?.trim();
		if (v) {
			const currentArr = formData[fieldKey] || [];
			formData[fieldKey] = [...currentArr, v];
			tempInputs[`${fieldKey}_val`] = '';
		}
	}

	function removeArrayItem(fieldKey: string, index: number) {
		if (!formData[fieldKey]) return;
		formData[fieldKey] = formData[fieldKey].filter((_: any, i: number) => i !== index);
	}
</script>

<div class="space-y-8" in:fly={{ x: 20, duration: 300 }}>
	{#each fields as field}
		<div class="space-y-3">
			<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase">
				{field.Name}
			</label>

			{#if field.Type === 'map'}
				<div class="space-y-2 rounded-2xl border border-border bg-panel/20 p-4">
					{#each Object.entries(formData[field.InsertKey] || {}) as [k, v] (k)}
						<div class="flex gap-2">
							<input
								disabled
								value={k}
								class="flex-1 rounded-[0.75rem] border border-border/50 bg-panel/40 px-3 py-2 font-mono text-xs text-muted"
							/>
							<input
								bind:value={formData[field.InsertKey][k]}
								class="flex-1 rounded-[0.75rem] border border-border bg-panel px-3 py-2 text-xs outline-none focus:border-frost-500"
							/>
							<button
								onclick={() => removeMapEntry(field.InsertKey, k)}
								class="p-2 text-muted hover:text-red-400"
							>
								<Trash2Icon size={16} />
							</button>
						</div>
					{/each}
					<div class="mt-2 flex gap-2 border-t border-border/50 pt-2">
						<input
							bind:value={tempInputs[`${field.InsertKey}_key`]}
							placeholder="Key"
							class="flex-1 rounded-[0.75rem] border border-border bg-surface px-3 py-2 text-xs outline-none focus:border-frost-500"
						/>
						<input
							bind:value={tempInputs[`${field.InsertKey}_val`]}
							placeholder="Value"
							class="flex-1 rounded-[0.75rem] border border-border bg-surface px-3 py-2 text-xs outline-none focus:border-frost-500"
						/>
						<button
							onclick={() => addMapEntry(field.InsertKey)}
							class="btn-frost rounded-full p-2 text-foreground hover:bg-frost-500 hover:text-background"
						>
							<PlusIcon size={16} />
						</button>
					</div>
				</div>
			{:else if field.Type === 'array'}
				<div class="space-y-2 rounded-2xl border border-border bg-panel/20 p-4">
					{#each formData[field.InsertKey] || [] as item, i (i)}
						<div class="flex gap-2">
							<input
								bind:value={formData[field.InsertKey][i]}
								class="flex-1 rounded-[0.75rem] border border-border bg-panel px-3 py-2 text-xs outline-none focus:border-frost-500"
							/>
							<button
								onclick={() => removeArrayItem(field.InsertKey, i)}
								class="p-2 text-muted hover:text-red-400"
							>
								<Trash2Icon size={16} />
							</button>
						</div>
					{/each}
					<div class="mt-2 flex gap-2 border-t border-border/50 pt-2">
						<input
							bind:value={tempInputs[`${field.InsertKey}_val`]}
							placeholder={`Add new ${field.Name}...`}
							class="flex-1 rounded-[0.75rem] border border-border bg-surface px-3 py-2 text-xs outline-none focus:border-frost-500"
						/>
						<button
							onclick={() => addArrayItem(field.InsertKey)}
							class="btn-frost rounded-full p-2 text-foreground hover:bg-frost-500 hover:text-background"
						>
							<PlusIcon size={16} />
						</button>
					</div>
				</div>
			{:else if field.Type === 'boolean'}
				<div
					class="flex items-center justify-between rounded-2xl border border-border bg-panel/50 p-4"
				>
					<span class="text-sm text-foreground/80">{field.Name}</span>
					<input
						type="checkbox"
						bind:checked={formData[field.InsertKey]}
						class="h-5 w-5 cursor-pointer accent-frost-500"
					/>
				</div>
			{:else}
				<input
					type={field.Type === 'number' ? 'number' : 'text'}
					bind:value={formData[field.InsertKey]}
					placeholder={`Enter ${field.Name}...`}
					class="w-full rounded-xl border border-border bg-panel px-4 py-3 text-sm outline-none focus:border-frost-500"
				/>
			{/if}
		</div>
	{/each}
</div>
