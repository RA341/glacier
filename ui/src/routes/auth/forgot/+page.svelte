<script lang="ts">
    import {MessageSquareIcon, ArrowLeftIcon, ShieldAlertIcon, ChevronLeftIcon} from '@lucide/svelte';
    import {fade} from 'svelte/transition';
</script>

<div class="space-y-6" in:fade={{ duration: 300 }}>
    <!-- Header Section -->
    <div class="space-y-2">
        <h2 class="text-xl font-bold text-foreground leading-tight">Password Recovery</h2>
        <!--        <p class="text-sm text-muted">Self-service recovery is disabled for security.</p>-->
    </div>

    <!-- Instructions Card -->
    <div class="p-6 bg-panel/40 border border-border rounded-2xl flex flex-col items-center text-center gap-4">
        <div class="p-4 bg-amber-500/10 text-amber-500 rounded-full border border-amber-500/20">
            <MessageSquareIcon size={32} strokeWidth={1.5}/>
        </div>

        <div class="space-y-2">
            <p class="text-sm font-bold text-foreground">
                Administrator Action Required
            </p>
            <p class="text-xs text-muted leading-relaxed px-2">
                Please contact your administrator to request a credential reset.
            </p>
        </div>
    </div>

    <!-- Navigation -->
    <div class="pt-4 border-t border-border">
        <a
                href="/auth/login"
                class="group w-full py-4 bg-panel border border-border text-foreground font-bold rounded-2xl hover:bg-surface hover:border-frost-500/50 transition-all flex items-center justify-center gap-2 active:scale-[0.98]"
        >
            <ChevronLeftIcon size={18} class="group-hover:-translate-x-1 transition-transform"/>
            Back to Login
        </a>
    </div>
</div>
