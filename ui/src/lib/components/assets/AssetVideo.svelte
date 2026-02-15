<script lang="ts">
    import type { Snippet } from 'svelte';

    let {
        src,
        class: className = '',
        autoplay = false,
        muted = false,
        loop = false,
        controls = true,
        loadingSlot,
        errorSlot,
        onLoad,
        onError
    }: {
        src: string | undefined | null;
        class?: string;
        autoplay?: boolean;
        muted?: boolean;
        loop?: boolean;
        controls?: boolean;
        loadingSlot?: Snippet;
        errorSlot?: Snippet;
        onLoad?: () => void;
        onError?: (err: Event) => void;
    } = $props();

    let status = $state<'loading' | 'error' | 'success'>('loading');

    $effect(() => {
        if (src) status = 'loading';
        else status = 'error';
    });

    function handleCanPlay() {
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
        <!-- svelte-ignore a11y_media_has_caption -->
        <video
                {src}
                {autoplay}
                {muted}
                {loop}
                {controls}
                playsinline
                class="h-full w-full {status === 'loading' ? 'invisible' : 'visible'}"
                oncanplay={handleCanPlay}
                onerror={handleError}
        ></video>
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