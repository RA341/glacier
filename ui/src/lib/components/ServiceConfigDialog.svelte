<script lang="ts">
    import {
        AlertTriangleIcon,
        ArrowLeftIcon,
        ChevronRightIcon,
        PlusIcon,
        SaveIcon,
        Trash2Icon,
        XIcon
    } from '@lucide/svelte';

    import {fade, fly} from 'svelte/transition';
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {onMount} from "svelte";
    import {callRPC, glacierCli} from "$lib/api/api";
    import {ServiceConfigService} from "$lib/gen/service_config/v1/service_config_pb";

    let {isOpen = $bindable(), ServiceType = ""} = $props();
    let scConfig = glacierCli(ServiceConfigService);

    let formData = $state<Record<string, any>>({});

    let tempInputs = $state<Record<string, string>>({});

    let supportedClientsRpc = createRPCRunner(() => scConfig.getSupportedValues({ServiceType}));
    let getSchemaRpc = createRPCRunner(() => scConfig.getSchema({
        ServiceType: ServiceType,
        Flavour: selectedFlavour!
    }));

    let selectedFlavour = $state<string | null>(null);
    let errorMessage = $state<string | null>(null);

    $effect(() => {
        if (!selectedFlavour) {
            formData = {};
            tempInputs = {};
            return;
        }
        getSchemaRpc.runner().then(() => {
            const initial: Record<string, any> = {};
            getSchemaRpc.value?.fields.forEach(f => {
                if (f.Type === 'boolean') initial[f.InsertKey] = false;
                else if (f.Type === 'number') initial[f.InsertKey] = 0;
                else if (f.Type === 'map') initial[f.InsertKey] = {};
                else if (f.Type === 'array') initial[f.InsertKey] = [];
                else initial[f.InsertKey] = "";
            });
            formData = initial;
        });
    });

    onMount(() => {
        supportedClientsRpc.runner();
    });

    $effect(() => {
        supportedClientsRpc.value
    })

    function addMapEntry(fieldKey: string) {
        const k = tempInputs[`${fieldKey}_key`]?.trim();
        const v = tempInputs[`${fieldKey}_val`]?.trim();
        if (k && v !== undefined) {
            formData[fieldKey][k] = v;
            tempInputs[`${fieldKey}_key`] = "";
            tempInputs[`${fieldKey}_val`] = "";
        }
    }

    function removeMapEntry(fieldKey: string, entryKey: string) {
        delete formData[fieldKey][entryKey];
    }

    function addArrayItem(fieldKey: string) {
        const v = tempInputs[`${fieldKey}_val`]?.trim();
        if (v) {
            formData[fieldKey].push(v);
            tempInputs[`${fieldKey}_val`] = "";
        }
    }

    function removeArrayItem(fieldKey: string, index: number) {
        formData[fieldKey].splice(index, 1);
    }

    function handleClose() {
        isOpen = false;
        selectedFlavour = null;
    }

    let enabled = $state(false)
    let name = $state("");

    async function handleSave() {
        console.log("Submitting:", formData)
        if (!name) {
            console.log("Name is required")
            return
        }

        const {err} = await callRPC(() =>
            scConfig.new({
                    conf: {
                        Enabled: enabled,
                        Name: name,
                        ServiceType: ServiceType,
                        Flavour: selectedFlavour!,
                        Config: new TextEncoder().encode(JSON.stringify(formData)),
                    }
                },
            )
        )
        if (err) {
            errorMessage = err
            return
        }

        handleClose();
    }
</script>

{#if isOpen}
    <div class="fixed inset-0 z-100 flex items-center justify-center p-6" transition:fade={{ duration: 150 }}>
        <div class="absolute inset-0 bg-background/80 backdrop-blur-md" onclick={handleClose}></div>

        <div class="relative w-full max-w-xl bg-surface border border-border rounded-3xl shadow-2xl overflow-hidden flex flex-col min-h-[500px] max-h-[90vh]"
             transition:fly={{ y: 20, duration: 300 }}>

            <!-- HEADER -->
            <div class="flex items-center justify-between p-6 border-b border-border bg-panel/30">
                <div class="flex items-center gap-3">
                    {#if selectedFlavour}
                        <button onclick={() => selectedFlavour = null}
                                class="p-2 hover:bg-panel rounded-lg text-muted transition-colors">
                            <ArrowLeftIcon size={18}/>
                        </button>
                    {/if}
                    <h2 class="text-xl font-bold text-foreground">
                        {selectedFlavour ? `Configure ${selectedFlavour}` : 'Select Client Type'}
                    </h2>
                </div>
                <button onclick={handleClose} class="text-muted hover:text-foreground p-1">
                    <XIcon size={20}/>
                </button>
            </div>

            <!-- BODY -->
            <div class="flex-1 overflow-y-auto p-8">
                {#if !selectedFlavour}
                    <!-- Type SELECTION -->
                    <div class="grid grid-cols-2 gap-4">
                        {#if supportedClientsRpc.error}
                            <div>Error: {supportedClientsRpc.error}</div>
                        {:else if supportedClientsRpc.loading}
                            <div>loading</div>
                        {:else}
                            {#if !supportedClientsRpc.value}
                                <div>No supported flavours found</div>
                            {:else}
                                {#each supportedClientsRpc.value?.values ?? [] as client}
                                    <button onclick={() => selectedFlavour = client}
                                            class="flex items-center justify-between p-4 bg-panel border border-border rounded-2xl hover:border-frost-500/50 transition-all group">
                                        <span class="font-bold text-foreground uppercase tracking-wider text-sm">{client}</span>
                                        <ChevronRightIcon size={18} class="text-muted/30 group-hover:text-frost-400"/>
                                    </button>
                                {/each}
                            {/if}
                        {/if}
                    </div>
                {:else}
                    {#if getSchemaRpc.error}
                        <div>Error: {getSchemaRpc.error}</div>
                    {:else if getSchemaRpc.loading}
                        <div>loading</div>
                    {:else}
                        {#if !getSchemaRpc.value}
                            <div>No supported flavours found</div>
                        {:else}
                            <div class="space-y-8" in:fly={{ x: 20, duration: 300 }}>

                                <span class="text-sm text-foreground/80">Name</span>
                                <input
                                        required
                                        type={'text'}
                                        bind:value={name}
                                        placeholder={`name`}
                                        class="w-full bg-panel border border-border rounded-xl py-3 px-4 outline-none focus:border-frost-500 text-sm"
                                />

                                <span class="text-sm text-foreground/80">Enable</span>
                                <input required type="checkbox" bind:checked={enabled}
                                       class="w-5 h-5 accent-frost-500 cursor-pointer"/>

                                {#each getSchemaRpc.value?.fields ?? [] as field}
                                    <div class="space-y-3">
                                        <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">{field.Name}</label>

                                        {#if field.Type === 'map'}
                                            <div class="space-y-2 bg-panel/20 p-4 rounded-2xl border border-border">
                                                <!-- Existing Entries -->
                                                {#each Object.entries(formData[field.InsertKey] || {}) as [k, v]}
                                                    <div class="flex gap-2">
                                                        <input disabled value={k}
                                                               class="flex-1 bg-panel/40 border border-border/50 rounded-lg px-3 py-2 text-xs text-muted font-mono"/>
                                                        <input bind:value={formData[field.InsertKey][k]}
                                                               class="flex-1 bg-panel border border-border rounded-lg px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                                                        <button onclick={() => removeMapEntry(field.InsertKey, k)}
                                                                class="p-2 text-muted hover:text-red-400">
                                                            <Trash2Icon size={16}/>
                                                        </button>
                                                    </div>
                                                {/each}
                                                <!-- Add New Entry -->
                                                <div class="flex gap-2 pt-2 border-t border-border/50 mt-2">
                                                    <input bind:value={tempInputs[`${field.InsertKey}_key`]}
                                                           placeholder="Key"
                                                           class="flex-1 bg-surface border border-border rounded-lg px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                                                    <input bind:value={tempInputs[`${field.InsertKey}_val`]}
                                                           placeholder="Value"
                                                           class="flex-1 bg-surface border border-border rounded-lg px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                                                    <button onclick={() => addMapEntry(field.InsertKey)}
                                                            class="p-2 bg-frost-500 text-background rounded-lg hover:bg-frost-400">
                                                        <PlusIcon size={16}/>
                                                    </button>
                                                </div>
                                            </div>

                                        {:else if field.Type === 'array'}
                                            <div class="space-y-2 bg-panel/20 p-4 rounded-2xl border border-border">
                                                <!-- Existing Items -->
                                                {#each formData[field.InsertKey] || [] as item, i}
                                                    <div class="flex gap-2">
                                                        <input bind:value={formData[field.InsertKey][i]}
                                                               class="flex-1 bg-panel border border-border rounded-lg px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                                                        <button onclick={() => removeArrayItem(field.InsertKey, i)}
                                                                class="p-2 text-muted hover:text-red-400">
                                                            <Trash2Icon size={16}/>
                                                        </button>
                                                    </div>
                                                {/each}
                                                <!-- Add New Item -->
                                                <div class="flex gap-2 pt-2 border-t border-border/50 mt-2">
                                                    <input bind:value={tempInputs[`${field.InsertKey}_val`]}
                                                           placeholder={`Add new ${field.Name}...`}
                                                           class="flex-1 bg-surface border border-border rounded-lg px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                                                    <button onclick={() => addArrayItem(field.InsertKey)}
                                                            class="p-2 bg-frost-500 text-background rounded-lg hover:bg-frost-400">
                                                        <PlusIcon size={16}/>
                                                    </button>
                                                </div>
                                            </div>

                                        {:else if field.Type === 'boolean'}
                                            <div class="flex items-center justify-between p-4 bg-panel/50 border border-border rounded-2xl">
                                                <span class="text-sm text-foreground/80">{field.Name}</span>
                                                <input type="checkbox" bind:checked={formData[field.InsertKey]}
                                                       class="w-5 h-5 accent-frost-500 cursor-pointer"/>
                                            </div>

                                        {:else}
                                            <div class="relative">
                                                <input
                                                        type={field.Type === 'number' ? 'number' : 'text'}
                                                        bind:value={formData[field.InsertKey]}
                                                        placeholder={`Enter ${field.Name}...`}
                                                        class="w-full bg-panel border border-border rounded-xl py-3 px-4 outline-none focus:border-frost-500 text-sm"
                                                />
                                            </div>
                                        {/if}
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    {/if}
                {/if}
            </div>

            <!-- FOOTER -->
            <div class="p-6 border-t border-border bg-panel/30 flex justify-end gap-3">
                <button onclick={handleClose}
                        class="px-6 py-2.5 rounded-xl border border-border text-sm font-bold text-muted">Cancel
                </button>
                {#if selectedFlavour}
                    <button onclick={handleSave}
                            class="px-8 py-2.5 rounded-xl bg-frost-500 text-background font-bold text-sm hover:bg-frost-400 flex items-center gap-2">
                        <SaveIcon size={16}/>
                        Save Client
                    </button>
                {/if}
            </div>
        </div>

        {#if errorMessage}
            <div class="fixed inset-0 z-[110] flex items-center justify-center p-6" transition:fade={{ duration: 100 }}>
                <!-- Darker backdrop for the error alert -->
                <!-- svelte-ignore a11y_click_events_have_key_events -->
                <!-- svelte-ignore a11y_no_static_element_interactions -->
                <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" onclick={() => errorMessage = null}></div>

                <div
                        class="relative w-full max-w-sm bg-surface border border-red-500/30 rounded-2xl shadow-2xl p-6 flex flex-col gap-4"
                        transition:fly={{ y: 10, duration: 200 }}
                >
                    <div class="flex items-center gap-3 text-red-400">
                        <div class="p-2 bg-red-500/10 rounded-lg">
                            <AlertTriangleIcon size={20} />
                        </div>
                        <h3 class="font-bold">Action Failed</h3>
                    </div>

                    <p class="text-sm text-muted leading-relaxed">
                        {errorMessage}
                    </p>

                    <div class="flex justify-end pt-2">
                        <button
                                onclick={() => errorMessage = null}
                                class="px-5 py-2 bg-panel border border-border rounded-xl text-sm font-bold hover:text-foreground transition-all"
                        >
                            Dismiss
                        </button>
                    </div>
                </div>
            </div>
        {/if}

    </div>
{/if}
