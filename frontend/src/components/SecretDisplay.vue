<template>
    <label class="input" @click="copySecret">
        <input :type="showSecret ? 'text' : 'password'" :value="secret" readonly />
        <span class="label">
            <ClipboardCheck v-if="copied" color="green" :size="20" />
            <ClipboardPlus v-else :size="20" />
        </span>
    </label>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { ClipboardPlus, ClipboardCheck } from "lucide-vue-next";

const props = defineProps({
    secret: String,
    showSecret: {
        type: Boolean,
        default: true,
    },
});

const copied = ref(false);

function copySecret() {
    navigator.clipboard.writeText(props.secret);
    copied.value = true;
}

watch(copied, () => {
    if (copied.value) {
        setTimeout(() => {
            copied.value = false;
        }, 1000);
    }
});
</script>
