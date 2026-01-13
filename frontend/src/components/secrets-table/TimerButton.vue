<template>
    <button class="btn btn-primary" :style="buttonBg" :disabled="enabled" @click="onClick">
        <slot />
    </button>
</template>

<script setup lang="ts">
import { computed, ref, watch, onUnmounted } from "vue";

const props = defineProps<{
    timeLimit?: number;
}>();

const timeLimit = computed(() => props.timeLimit ?? 10);

const emit = defineEmits(["enabled", "disabled"]);

let interval: ReturnType<typeof setInterval>;
const timeRemaining = ref(0);
const percentRemaining = computed(() => Math.floor((timeRemaining.value / timeLimit.value) * 100));
const enabled = computed(() => percentRemaining.value > 0);
const buttonBg = computed(() => {
    if (!enabled.value) return {};
    const percent = percentRemaining.value;
    const color = "#570df8";
    const faded = hexToRgba(color, 0.3);
    return {
        background: `linear-gradient(90deg, ${faded} ${percent}%, ${color} ${percent}%)`,
    };
});

watch(enabled, (newVal) => {
    if (newVal) {
        emit("enabled");
    } else {
        emit("disabled");
    }
});

function onClick() {
    // Start the timer to begin the timed enable
    if (enabled.value) return;
    if (interval) clearInterval(interval);
    timeRemaining.value = timeLimit.value;
    interval = setInterval(() => {
        timeRemaining.value -= 1;
        if (timeRemaining.value <= 0) {
            timeRemaining.value = 0;
            clearInterval(interval!);
            interval = null;
        }
    }, 1000);
}

onUnmounted(() => {
    if (interval) clearInterval(interval);
});

function hexToRgba(hex: string, alpha: number): string {
    hex = hex.replace("#", "");
    if (hex.length === 3) {
        hex = hex
            .split("")
            .map((x) => x + x)
            .join("");
    }
    const num = parseInt(hex, 16);
    const r = (num >> 16) & 255;
    const g = (num >> 8) & 255;
    const b = num & 255;
    return `rgba(${r},${g},${b},${alpha})`;
}
</script>
