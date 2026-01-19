<script lang="ts">

    import {cli} from "$lib/api/api";
    import {LibraryService} from "$lib/gen/library/v1/library_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";

    const libSrv = cli(LibraryService)

    let libRpc = createRPCRunner(() => libSrv.triggerTracker({}))
    function trigger() {
        libRpc.runner()
    }

</script>


Downloads Page

<button
        onclick={trigger}
        class="right-2 p-2 text-muted hover:text-frost-400 transition-colors"
        aria-label="trigger"
>
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
         stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"/>
        <path d="m21 21-4.3-4.3"/>
    </svg>
</button>g