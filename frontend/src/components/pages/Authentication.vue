<template>
    <div class="card bg-base-200 shadow-md max-w-lg mx-auto">
        <div class="card-body gap-5">
            <h2 class="card-title">Manage API Key</h2>

            <div v-if="loading">
                <span class="loading loading-dots loading-md"></span>
            </div>
            <div v-else>
                <div v-if="!apiKey" class="text-error font-semibold">No API Key</div>
                <SecretDisplay v-else :secret="apiKey" class="w-full" />
            </div>

            <button class="btn btn-primary" @click="generateApiKey" :disabled="loading">
                {{ apiKey ? "Regenerate" : "Generate" }}
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import SecretDisplay from "@/components/SecretDisplay.vue";
import { toast } from "vue3-toastify";
import { backendFetch } from "@/main";

const loading = ref(false);

const apiKey = ref("");
async function fetchApiKey() {
    loading.value = true;
    try {
        const response = await backendFetch(`/api/key`);
        if (!response.ok) return;
        const data = await response.json();
        apiKey.value = data.key ?? "";
    } catch (error) {
        console.error("Error fetching API key:", error);
        toast("Failed to fetch API key", { type: "error" });
    } finally {
        loading.value = false;
    }
}
onMounted(fetchApiKey);

async function generateApiKey() {
    loading.value = true;
    try {
        const response = await backendFetch(`/api/key/generate`);
        if (!response.ok) return;
        const data = await response.json();
        apiKey.value = data.key;
    } catch (error) {
        console.error("Error generating API key:", error);
        toast("Failed to generate API key", { type: "error" });
    } finally {
        loading.value = false;
    }
}
</script>
