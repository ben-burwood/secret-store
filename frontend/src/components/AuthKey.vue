<template>
    <div class="card bg-base-200 shadow-md max-w-lg mx-auto">
        <div class="card-body gap-5">
            <h2 class="card-title">Manage Auth Key</h2>

            <div v-if="!authKey" class="text-error font-semibold">No Auth Key</div>
            <SecretDisplay v-else :secret="authKey" class="w-full" />

            <button class="btn btn-primary" @click="generateAuthKey" :disabled="loading">
                {{ authKey ? "Regenerate" : "Generate" }}
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import SecretDisplay from "./SecretDisplay.vue";
import { toast } from "vue3-toastify";
import { SERVER_URL } from "@/main";

const authKey = ref("");
async function fetchAuthKey() {
    try {
        const response = await fetch(`${SERVER_URL}/auth/key`);
        const data = await response.json();
        authKey.value = data.key;
    } catch (error) {
        console.error("Error fetching Auth key:", error);
        toast("Failed to fetch Auth key", { type: "error" });
    }
}
onMounted(fetchAuthKey);

const loading = ref(false);
async function generateAuthKey() {
    loading.value = true;
    try {
        const response = await fetch(`${SERVER_URL}/auth/key/generate`);
        const data = await response.json();
        authKey.value = data.key;
    } catch (error) {
        console.error("Error generating Auth key:", error);
        toast("Failed to generate Auth key", { type: "error" });
    } finally {
        loading.value = false;
    }
}
</script>
