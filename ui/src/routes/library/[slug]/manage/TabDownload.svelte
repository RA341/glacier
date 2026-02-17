<script lang="ts">
    import {ShieldCheckIcon} from '@lucide/svelte';
    import {type Game, LibraryService} from "$lib/gen/library/v1/library_pb";
    import {callRPC, glacierCli} from "$lib/api/api";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";

    let {game = $bindable()}: { game: Game; refresh: () => Promise<void> } = $props();

    const lib = glacierCli(LibraryService)
    const sm = getSnackbarCtx()

    async function handleRedownload() {
        if (!game?.ID) {
            return
        }

        const {err} = await callRPC(() => lib.redownload({gameId: game.ID}))
        if (err) {
            sm.push(`Could not redownload ${err}`, 'error')
        }
        // await refresh()
    }

</script>

<!-- Download Section -->
<section class="space-y-4">
    <h2 class="text-sm font-bold uppercase tracking-widest text-muted px-2">Download</h2>
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 p-6 bg-surface border border-border rounded-3xl">
        <!-- State Info -->
        <div class="lg:col-span-5 grid grid-cols-2 gap-y-6 gap-x-4 border-r border-border pr-6">
            <div>
                <button onclick={handleRedownload} class="text-frost-400 hover:text-foreground p-1">
                    Retry
                </button>
            </div>
            <div>
                <p class="text-[9px] font-bold text-muted uppercase">Client Type</p>
                <p class="text-sm font-bold">{game?.DownloadState?.Client}</p>
            </div>
            <div>
                <p class="text-[9px] font-bold text-muted uppercase">Download ID</p>
                <p class="text-sm font-mono opacity-60 truncate">{game?.DownloadState?.DownloadId}</p>
            </div>
            <div>
                <p class="text-[9px] font-bold text-muted uppercase">State</p>
                <p class="text-sm font-bold text-frost-400">{game?.DownloadState?.State}</p>
            </div>
            <div>
                <p class="text-[9px] font-bold text-muted uppercase">Progress</p>
                <p class="text-sm font-bold">{game?.DownloadState?.Progress}</p>
            </div>

            <div>
                <p class="text-[9px] font-bold text-muted uppercase">Download Url</p>
                <p class="text-sm font-bold text-frost-400 overflow-hidden">{game?.DownloadState?.DownloadUrl}</p>
            </div>
            <div>
                <p class="text-[9px] font-bold text-muted uppercase">Path</p>
                <p class="text-sm font-bold wrap-break-word">{game?.DownloadState?.DownloadPath}</p>
            </div>
        </div>

        <div class="col-span-2 pt-4 ">
            <p class="text-[9px] font-bold text-muted uppercase mb-2 flex items-center gap-1">
                <ShieldCheckIcon size={12}/>
                Virus Scan
            </p>
            <p class="text-xs text-green-400 font-medium">TODO</p>
        </div>
    </div>
</section>
