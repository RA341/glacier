<script lang="ts">
    import {EyeIcon, EyeOffIcon, InfoIcon, LockIcon} from '@lucide/svelte';
    import Self from './ConfigField.svelte'
    import {splitCamelCase} from '$lib/api/strings'

    let {schema, level = 0} = $props();

    let visibleSecrets = $state<Record<string, boolean>>({});

    function toggleSecret(key: string) {
        visibleSecrets[key] = !visibleSecrets[key];
    }
</script>

<div class="space-y-6" class:ml-6={level > 0} class:border-l={level > 0} class:border-border={level > 0}
     class:pl-6={level > 0}>
    {#each Object.entries(schema) as [key, item]}
        <div class="flex flex-col gap-2">
            {#if item.IsStruct}
                <!-- SECTION HEADER -->
                <div class="pt-4 first:pt-0">
                    <h3 class="text-sm font-black text-frost-400 uppercase tracking-widest flex items-center gap-2">
                        {key}
                        {#if item.Help}
                            <div class="group relative">
                                <InfoIcon size={12} class="text-muted/40 cursor-help"/>
                                <div class="absolute left-full ml-2 top-0 w-48 p-2 bg-panel border border-border rounded text-[10px] font-medium normal-case tracking-normal text-white opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none z-50">
                                    {item.Help}
                                </div>
                            </div>
                        {/if}
                    </h3>
                    <div class="h-px bg-border mt-2 w-full opacity-50"></div>
                </div>

                {#if item.Nested}
                    <Self schema={item.Nested} level={level + 1}/>
                {/if}

            {:else}
                <!-- LEAF INPUTS -->
                <div
                        class="group flex flex-col gap-1.5 p-2 rounded-xl hover:bg-panel/20 transition-colors"
                        class:opacity-60={item.EnvSet}
                >
                    <div class="flex justify-between items-center">
                        <div class="flex items-center gap-2">
                            <label class="text-[11px] font-bold text-muted uppercase tracking-wider">{splitCamelCase(key)}</label>

                            {#if item.EnvSet}
                                <div class="group/env relative">
                                    <LockIcon size={12} class="text-amber-500/60"/>
                                    <!-- ENV TOOLTIP -->
                                    <div class="absolute left-full ml-2 top-1/2 -translate-y-1/2 w-64 p-3 bg-surface border border-amber-500/30 rounded-lg shadow-xl text-[10px] font-medium normal-case tracking-normal text-muted opacity-0 group-hover/env:opacity-100 transition-opacity pointer-events-none z-[60]">
                                        <p class="text-amber-400 font-bold mb-1 uppercase tracking-tighter">Locked by
                                            Environment</p>
                                        This value is set via <code
                                            class="bg-panel px-1 rounded text-foreground">{item.Env}</code>.
                                        Changes must be made via the env var or by removing the env then editing.
                                    </div>
                                </div>
                            {/if}
                        </div>
                        {#if item.Default}
                            <code class="text-[9px] font-mono ">Default: {item.Default}</code>
                        {/if}

                        {#if item.Env}
                            <code class="text-[9px] font-mono ">ENV: {item.Env}</code>
                        {/if}
                    </div>

                    <div class="relative flex items-center gap-2">
                        {#if item.FieldType === 'bool'}
                            <label class="relative inline-flex items-center" class:cursor-not-allowed={item.EnvSet}
                                   class:cursor-pointer={!item.EnvSet}>
                                <input type="checkbox" bind:checked={item.Value} disabled={item.EnvSet}
                                       class="sr-only peer">
                                <div class="w-11 h-6 bg-panel border border-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-muted after:border-border after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-frost-500 peer-checked:after:bg-background opacity-100 peer-disabled:opacity-40"></div>
                            </label>

                        {:else if item.FieldType === 'int'}
                            <input
                                    type="number"
                                    bind:value={item.Value}
                                    disabled={item.EnvSet}
                                    class="w-full bg-panel border border-border rounded-lg py-2 px-3 text-sm outline-none focus:border-frost-500 transition-colors disabled:cursor-not-allowed"
                            />

                        {:else if item.IsSecret}
                            <div class="relative flex-1">
                                <input
                                        type={visibleSecrets[key] ? 'text' : 'password'}
                                        bind:value={item.Value}
                                        disabled={item.EnvSet}
                                        class="w-full bg-panel border border-border rounded-lg py-2 pl-3 pr-10 text-sm outline-none focus:border-frost-500 transition-colors disabled:cursor-not-allowed"
                                />
                                <button
                                        type="button"
                                        onclick={() => toggleSecret(key)}
                                        class="absolute right-2 top-1/2 -translate-y-1/2 text-muted hover:text-frost-400"
                                >
                                    {#if visibleSecrets[key]}
                                        <EyeOffIcon size={16}/>
                                    {:else}
                                        <EyeIcon size={16}/>
                                    {/if}
                                </button>
                            </div>

                        {:else if Array.isArray(item.Value)}
                            <input
                                    type="text"
                                    value={item.Value.join(', ')}
                                    disabled={item.EnvSet}
                                    oninput={(e) => item.Value = e.currentTarget.value.split(',').map(s => s.trim())}
                                    class="w-full bg-panel border border-border rounded-lg py-2 px-3 text-sm outline-none focus:border-frost-500 transition-colors font-mono disabled:cursor-not-allowed"
                            />

                        {:else}
                            <input
                                    type="text"
                                    bind:value={item.Value}
                                    disabled={item.EnvSet}
                                    class="w-full bg-panel border border-border rounded-lg py-2 px-3 text-sm outline-none focus:border-frost-500 transition-colors disabled:cursor-not-allowed"
                            />
                        {/if}
                    </div>

                    {#if item.Help}
                        <p class="text-[10px] italic">{item.Help}</p>
                    {/if}
                </div>
            {/if}
        </div>
    {/each}
</div>
