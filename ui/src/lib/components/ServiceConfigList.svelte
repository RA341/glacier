<script lang="ts">
    import ServiceConfigDialog from "$lib/components/ServiceConfigDialog.svelte";
    import {PlusIcon, RefreshCcw} from "@lucide/svelte";
    import {glacierCli} from "$lib/api/api";
    import {ServiceConfigService} from "$lib/gen/service_config/v1/service_config_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {onMount} from "svelte";

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

</script>

<div class="flex items-center gap-5">
    <button
            onclick={() => isOpen = true}
            class="flex items-center gap-2 px-6 py-2.5 bg-frost-500 text-background rounded-xl text-sm font-bold hover:bg-frost-400 transition-all shadow-lg shadow-frost-500/20 active:scale-95"
    >
        <PlusIcon size={18}/>
        Add {ServiceType}
    </button>

    <button
            onclick={refresh}
            class="flex items-center gap-2 px-6 py-2.5 bg-frost-500 text-background rounded-xl text-sm font-bold hover:bg-frost-400 transition-all shadow-lg shadow-frost-500/20 active:scale-95"
    >
        <RefreshCcw size={18}/>
        Refresh
    </button>
</div>

<div>
    {#if ServiceType}
        <!-- Type SELECTION -->
        <div class="grid grid-cols-2 gap-4">
            {#if listRpc.error}
                <div>Error: {listRpc.error}</div>
            {:else if listRpc.loading}
                <div>loading</div>
            {:else}
                {#if !listRpc.value}
                    <div>No supported flavours found</div>
                {:else}
                    {#if !listRpc.value.conf || listRpc.value.conf.length === 0}
                        <div>No configs found</div>
                    {:else}
                        {#each listRpc.value?.conf ?? [] as client}
                            <div>{client.Name}</div>
                            <!--                        <button onclick={() => selectedFlavour = client}-->
                            <!--                                class="flex items-center justify-between p-4 bg-panel border border-border rounded-2xl hover:border-frost-500/50 transition-all group">-->
                            <!--                            <span class="font-bold text-foreground uppercase tracking-wider text-sm">{client}</span>-->
                            <!--                            <ChevronRightIcon size={18} class="text-muted/30 group-hover:text-frost-400"/>-->
                            <!--                        </button>-->
                        {/each}
                    {/if}
                {/if}
            {/if}
        </div>
    {/if}
</div>

<ServiceConfigDialog ServiceType={ServiceType} bind:isOpen={isOpen}/>

