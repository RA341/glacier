<script lang="ts">
    import {AlertCircleIcon, LoaderIcon, SaveIcon} from '@lucide/svelte';
    import {callRPC, glacierCli} from "$lib/api/api";
    import {ConfigService} from "$lib/gen/config/v1/config_pb";
    import {createRPCRunner} from "$lib/api/rpc.svelte";
    import {onMount} from "svelte";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";
    import ConfigField from "$lib/components/ConfigField.svelte";
    import {extractValues} from "$lib/components/util.svelte";

    const sm = getSnackbarCtx()
    const libConf = glacierCli(ConfigService)
    let configSchema = $state({});

    let isSaving = $state(false);

    async function handleSave() {
        isSaving = true

        const payload = extractValues(configSchema);

        const {err} = await callRPC(() => libConf.set({
            configSchema: JSON.stringify(payload)
        }))

        if (err) {
            sm.push(`Could not update config: ${err}`, 'error')
        }

        await configRpc.runner()

        isSaving = false
    }

    const configRpc = createRPCRunner(() => libConf.get({}))

    onMount(() => {
        configRpc.runner()
    })

    $effect(() => {
        if (configRpc.value) {
            configSchema = JSON.parse(configRpc.value.configSchema)
        }
    })
</script>

{#if configRpc.loading && !configRpc.value}
    <!-- THEMED LOADING STATE -->
    <div class="flex flex-col items-center justify-center h-[60vh] gap-4 text-muted">
        <LoaderIcon class="animate-spin text-frost-500" size={40}/>
        <p class="text-sm font-bold uppercase tracking-widest animate-pulse">Loading Configuration...</p>
    </div>

{:else if configRpc.error}
    <!-- THEMED ERROR STATE -->
    <div class="max-w-md mx-auto mt-20 p-6 bg-red-500/5 border border-red-500/20 rounded-3xl flex flex-col items-center text-center gap-4">
        <AlertCircleIcon size={48} class="text-red-400"/>
        <div>
            <h3 class="font-bold text-foreground">Failed to load config</h3>
            <p class="text-sm text-muted mt-1">{configRpc.error}</p>
        </div>
        <button onclick={() => configRpc.runner()}
                class="px-6 py-2 bg-panel border border-border rounded-xl text-xs font-bold hover:bg-surface transition-all">
            Retry Connection
        </button>
    </div>

{:else}
    <div class="max-w-4xl mx-auto p-8 min-h-screen">
        <!-- STICKY HEADER -->
        <header class="sticky top-0 z-30 bg-background/80 backdrop-blur-md border-b border-border mb-10 -mx-8 px-8 py-6 flex justify-between items-center">
            <button
                    onclick={handleSave}
                    disabled={isSaving}
                    class="flex items-center gap-2 px-8 py-3 bg-frost-500 text-background rounded-2xl text-sm font-bold hover:bg-frost-400 transition-all shadow-lg shadow-frost-500/20 active:scale-95 disabled:opacity-50"
            >
                {#if isSaving}
                    <LoaderIcon size={18} class="animate-spin"/>
                    Saving...
                {:else}
                    <SaveIcon size={18}/>
                    Save Changes
                {/if}
            </button>
        </header>

        <div class="space-y-12 pb-20">
            <ConfigField schema={configSchema}/>
        </div>
    </div>
{/if}
