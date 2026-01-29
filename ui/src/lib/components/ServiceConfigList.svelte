<script lang="ts">
    import ServiceConfigDialog from "$lib/components/ServiceConfigDialog.svelte";
    import {
        AlertCircleIcon,
        Edit3Icon,
        InboxIcon,
        LoaderIcon,
        PlusIcon,
        RefreshCcw,
        Settings2Icon,
        Trash2Icon
    } from "@lucide/svelte";
    import {glacierCli} from "$lib/api/api";
    import {type ServiceConfig, ServiceConfigService} from "$lib/gen/service_config/v1/service_config_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {onMount} from "svelte";
    import {callRPC} from "$lib/api/api.ts";

    let {ServiceType}: { ServiceType: string } = $props();
    let scConfig = glacierCli(ServiceConfigService);
    let isOpen = $state(false)

    let listRpc = createRPCRunner(() => scConfig.list({serviceType: ServiceType}))

    function refresh() {
        listRpc.runner()
    }

    $effect(() => {
        if (!isOpen) {
            refresh()
        }
    })

    onMount(() => {
        refresh();
    })

    let editConf = $state<ServiceConfig | null>(null)

    function openEdit(config: ServiceConfig) {
        editConf = config;
        isOpen = true
    }

    async function deleteClient(config: ServiceConfig) {
        const {err} = await callRPC(() => scConfig.delete({id: config.ID}))
        if (err) {
            console.error(err)
        }

        refresh()
    }
</script>

<div class="space-y-6">
    <!-- TOOLBAR -->
    <div class="flex items-center justify-between px-2">
        <h2 class="text-sm font-bold uppercase tracking-widest text-muted">{ServiceType} Configurations</h2>
        <div class="flex items-center gap-3">
            <button
                    onclick={refresh}
                    class="p-2.5 rounded-xl bg-panel border border-border text-muted hover:text-frost-400 transition-all active:scale-95"
                    title="Refresh List"
            >
                <RefreshCcw size={18} class={listRpc.loading ? 'animate-spin' : ''}/>
            </button>

            <button
                    onclick={() => isOpen = true}
                    class="flex items-center gap-2 px-6 py-2.5 bg-frost-500 text-background rounded-xl text-sm font-bold hover:bg-frost-400 transition-all shadow-lg shadow-frost-500/20 active:scale-95"
            >
                <PlusIcon size={18}/>
                Add {ServiceType}
            </button>
        </div>
    </div>

    <!-- MAIN LIST -->
    <div class="min-h-100">
        {#if listRpc.error}
            <div class="flex flex-col items-center justify-center h-64 text-red-400 gap-3 bg-red-500/5 border border-red-500/10 rounded-3xl">
                <AlertCircleIcon size={32}/>
                <p class="text-sm font-medium">{listRpc.error}</p>
            </div>
        {:else if listRpc.loading && !listRpc.value}
            <div class="flex flex-col items-center justify-center h-64 text-muted gap-4">
                <LoaderIcon class="animate-spin text-frost-500" size={32}/>
                <p class="animate-pulse text-sm">Fetching configurations...</p>
            </div>
        {:else}
            {#if !listRpc.value?.conf || listRpc.value.conf.length === 0}
                <div class="flex flex-col items-center justify-center h-64 border-2 border-dashed border-border rounded-3xl text-muted/30">
                    <InboxIcon size={48} strokeWidth={1} class="mb-2"/>
                    <p class="text-sm font-medium">No {ServiceType} configurations found</p>
                </div>
            {:else}
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {#each listRpc.value.conf as client (client.ID)}
                        <!-- Update the grid to auto-fill based on a fixed minimum width -->
                        <div class="grid grid-cols-[repeat(auto-fill,minmax(320px,1fr))] gap-4">
                            {#each listRpc.value.conf as client (client.ID)}
                                <div class="group relative flex flex-col p-5 bg-panel/30 border border-border rounded-3xl hover:bg-panel/50 hover:border-frost-500/30 transition-all shadow-sm">

                                    <!-- Top Section: Status and Actions -->
                                    <div class="flex justify-between items-start mb-4">
                                        <div class="relative">
                                            <div class="p-3 bg-surface border border-border rounded-2xl text-muted group-hover:text-frost-400 transition-colors">
                                                <Settings2Icon size={24}/>
                                            </div>
                                            <div class="absolute -top-1 -right-1 w-3.5 h-3.5 rounded-full border-2 border-background {client.Enabled ? 'bg-green-500' : 'bg-red-500'}"></div>
                                        </div>

                                        <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                            <button
                                                    onclick={()=>openEdit(client)}
                                                    class="p-2 text-muted hover:text-frost-400 hover:bg-frost-500/10 rounded-lg transition-all">
                                                <Edit3Icon size={16}/>
                                            </button>
                                            <button
                                                    onclick={()=>deleteClient(client)}
                                                    class="p-2 text-muted hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-all">
                                                <Trash2Icon size={16}/>
                                            </button>
                                        </div>
                                    </div>

                                    <!-- Bottom Section: Info -->
                                    <div class="space-y-1">
                                        <h3 class="font-bold text-lg text-foreground truncate">{client.Name}</h3>
                                        <div class="flex items-center gap-2">
                                            <span class="text-[10px] font-bold text-muted uppercase tracking-widest bg-panel border border-border px-2 py-0.5 rounded-md">
                                                {client.Flavour}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    {/each}
                </div>
            {/if}
        {/if}
    </div>
</div>

<ServiceConfigDialog editConf={editConf} serviceType={ServiceType} bind:isOpen={isOpen}/>

<style>
    /* Subtle pulsing for the list while refreshing */
    .animate-spin {
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }
</style>
