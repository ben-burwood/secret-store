<template>
    <div class="card bg-base-200 shadow-md">
        <div class="card-body">
            <h1 class="card-title">Secret Generator</h1>

            <div>
                <input v-model="length" type="range" min="8" max="100" class="range" />
                <input
                    type="number"
                    v-model="length"
                    class="input input-ghost validator w-min"
                    required
                    placeholder="Length"
                    min="8"
                    max="100"
                    title="Must be between be 8 to 100"
                />
                <p class="validator-hint">Must be between be 8 to 100</p>
            </div>

            <div class="flex flex-col gap-2">
                <label class="label">
                    <input type="checkbox" checked class="toggle toggle-primary" v-model="includeNumbers" />
                    Include Numbers
                </label>
                <label class="label">
                    <input type="checkbox" checked class="toggle toggle-primary" v-model="includeSymbols" />
                    Include Symbols
                </label>
            </div>

            <div class="card-actions mt-5">
                <button class="btn btn-primary" :disabled="loading" @click="generateSecret">
                    <RefreshCcw :size="18" :class="{ 'animate-spin': loading }" />
                    Generate
                </button>
                <SecretDisplay :secret />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import SecretDisplay from "@/components/SecretDisplay.vue";
import { RefreshCcw, ClipboardPlus } from "lucide-vue-next";
import { SERVER_URL } from "@/main";
import { toast } from "vue3-toastify";

const length = ref(32);
const includeNumbers = ref(true);
const includeSymbols = ref(true);

const loading = ref(false);

const secret = ref("test");

async function generateSecret() {
    loading.value = true;
    const params = new URLSearchParams({
        length: length.value.toString(),
        includeNumbers: includeNumbers.value ? "true" : "false",
        includeSymbols: includeSymbols.value ? "true" : "false",
    });
    try {
        const response = await fetch(`${SERVER_URL}/secret/generate?${params.toString()}`);
        const data = await response.json();
        secret.value = data.secret;
    } catch (e) {
        toast(`Error Generating Secret : ${e.message}`, { type: "error" });
        console.error(e);
    } finally {
        loading.value = false;
    }
}
</script>
