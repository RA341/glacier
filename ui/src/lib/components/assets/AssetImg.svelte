<script lang="ts">
    import type {Snippet} from 'svelte';

    let {
        src,
        alt = '',
        class: className = '',
        imgClass: imgClassName = '',
        loadingSlot,
        errorSlot,
        onLoad,
        onError
    }: {
        src: string | undefined | null;
        alt?: string;
        class?: string;
        imgClass?: string;
        loadingSlot?: Snippet;
        errorSlot?: Snippet;
        onLoad?: () => void;
        onError?: (err: Event) => void;
    } = $props();

    let status = $state<'loading' | 'error' | 'success'>('loading');

    // Reset status when the URL changes
    $effect(() => {
        if (src) {
            status = 'loading';
        } else {
            status = 'error';
        }
    });

    function handleLoad() {
        status = 'success';
        onLoad?.();
    }

    function handleError(e: Event) {
        status = 'error';
        onError?.(e);
    }
</script>

<div class="relative {className}">
    {#if src && status !== 'error'}
        <img
                {src}
                {alt}
                class="{imgClassName} {status === 'loading' ? 'invisible' : 'visible'}"
                onload={handleLoad}
                onerror={handleError}
        />
    {/if}

    {#if status === 'loading' && src}
        <div class="absolute inset-0 z-10">
            {@render loadingSlot?.()}
        </div>
    {:else if status === 'error' || !src}
        <div class="absolute inset-0 z-10">
            {@render errorSlot?.()}
        </div>
    {/if}
</div>