<template>
    <SecretList v-if="secrets.length > 0" :secrets="secrets" @refresh="fetchSecrets" />
    <div v-else class="flex flex-col items-center justify-center h-full">
        <p class="text-xl font-semibold">Sssh, there's no secrets to be found!</p>

        <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4 mt-32">
            <legend class="fieldset-legend">Add new Secret</legend>

            <label class="label">Key</label>
            <input type="text" class="input" placeholder="Key" v-model="newKey" />

            <label class="label mt-4">Value</label>
            <input type="password" class="input" placeholder="Value" v-model="newValue" @keyup.enter="addSecret" />

            <button class="btn btn-primary mt-4" @click="addSecret">Add</button>
        </fieldset>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import SecretList from "@/components/SecretList.vue";
import { SERVER_URL } from "@/main";
import { toast } from "vue3-toastify";

const secrets = ref([]);
async function fetchSecrets() {
    try {
        const response = await fetch(`${SERVER_URL}/secrets`);
        const data = await response.json();
        secrets.value = data.secrets;
    } catch (error) {
        console.error("Error fetching secrets:", error);
        toast("Failed to fetch secrets", { type: "error" });
    }
}
onMounted(fetchSecrets);

const newKey = ref("");
const newValue = ref("");

async function addSecret() {
    try {
        await fetch(`${SERVER_URL}/secrets/new`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ key: newKey.value, value: newValue.value }),
        });
        fetchSecrets();
    } catch (error) {
        console.error("Error adding secret:", error);
        toast("Error adding secret", { type: "error" });
    }
}
</script>
