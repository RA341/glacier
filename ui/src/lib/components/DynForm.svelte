<script lang="ts">
    import {PlusIcon, Trash2Icon} from '@lucide/svelte';
    import {fly} from 'svelte/transition';
    import type {FieldSchema} from "$lib/gen/service_config/v1/service_config_pb";

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
            tempInputs[`${fieldKey}_key`] = "";
            tempInputs[`${fieldKey}_val`] = "";
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
            tempInputs[`${fieldKey}_val`] = "";
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
            <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">
                {field.Name}
            </label>

            {#if field.Type === 'map'}
                <div class="space-y-2 bg-panel/20 p-4 rounded-2xl border border-border">
                    {#each Object.entries(formData[field.InsertKey] || {}) as [k, v] (k)}
                        <div class="flex gap-2">
                            <input disabled value={k} class="flex-1 bg-panel/40 border border-border/50 rounded-lg px-3 py-2 text-xs text-muted font-mono"/>
                            <input bind:value={formData[field.InsertKey][k]} class="flex-1 bg-panel border border-border rounded-lg px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                            <button onclick={() => removeMapEntry(field.InsertKey, k)} class="p-2 text-muted hover:text-red-400">
                                <Trash2Icon size={16}/>
                            </button>
                        </div>
                    {/each}
                    <div class="flex gap-2 pt-2 border-t border-border/50 mt-2">
                        <input bind:value={tempInputs[`${field.InsertKey}_key`]} placeholder="Key" class="flex-1 bg-surface border border-border rounded-lg px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                        <input bind:value={tempInputs[`${field.InsertKey}_val`]} placeholder="Value" class="flex-1 bg-surface border border-border rounded-lg px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                        <button onclick={() => addMapEntry(field.InsertKey)} class="p-2 bg-frost-500 text-background rounded-lg hover:bg-frost-400">
                            <PlusIcon size={16}/>
                        </button>
                    </div>
                </div>

            {:else if field.Type === 'array'}
                <div class="space-y-2 bg-panel/20 p-4 rounded-2xl border border-border">
                    {#each formData[field.InsertKey] || [] as item, i (i)}
                        <div class="flex gap-2">
                            <input bind:value={formData[field.InsertKey][i]} class="flex-1 bg-panel border border-border rounded-lg px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                            <button onclick={() => removeArrayItem(field.InsertKey, i)} class="p-2 text-muted hover:text-red-400">
                                <Trash2Icon size={16}/>
                            </button>
                        </div>
                    {/each}
                    <div class="flex gap-2 pt-2 border-t border-border/50 mt-2">
                        <input bind:value={tempInputs[`${field.InsertKey}_val`]} placeholder={`Add new ${field.Name}...`} class="flex-1 bg-surface border border-border rounded-lg px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                        <button onclick={() => addArrayItem(field.InsertKey)} class="p-2 bg-frost-500 text-background rounded-lg hover:bg-frost-400">
                            <PlusIcon size={16}/>
                        </button>
                    </div>
                </div>

            {:else if field.Type === 'boolean'}
                <div class="flex items-center justify-between p-4 bg-panel/50 border border-border rounded-2xl">
                    <span class="text-sm text-foreground/80">{field.Name}</span>
                    <input type="checkbox" bind:checked={formData[field.InsertKey]} class="w-5 h-5 accent-frost-500 cursor-pointer"/>
                </div>

            {:else}
                <input type={field.Type === 'number' ? 'number' : 'text'} bind:value={formData[field.InsertKey]} placeholder={`Enter ${field.Name}...`} class="w-full bg-panel border border-border rounded-xl py-3 px-4 outline-none focus:border-frost-500 text-sm"/>
            {/if}
        </div>
    {/each}
</div>