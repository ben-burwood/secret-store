<template>
    <tr>
        <td>
            <input v-model="newKey" class="input input-bordered w-full" placeholder="Key" />
        </td>
        <td>
            <input v-model="newValue" class="input input-bordered w-full" placeholder="Value" @keyup.enter="addSecret" />
        </td>
        <td>
            <input v-model="newTag" list="existing-tags" class="input input-bordered w-full" placeholder="Tag (optional)" @keyup.enter="addSecret" />
        </td>
        <td>
            <ConfirmReject @confirm="addSecret" @reject="cancelAddRow" :canConfirm="newKey !== '' && newValue !== ''" />
        </td>
    </tr>
</template>

<script setup lang="ts">
import { ref } from "vue";
import ConfirmReject from "@/components/ConfirmReject.vue";
import { toast } from "vue3-toastify";
import { backendFetch } from "@/main";

const emit = defineEmits(["refresh", "cancel"]);

const newKey = ref("");
const newValue = ref("");
const newTag = ref("");
const cancelAddRow = () => {
    newKey.value = "";
    newValue.value = "";
    newTag.value = "";
    emit("cancel");
};

async function addSecret() {
    try {
        const res = await backendFetch(`/secrets/new`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ key: newKey.value, value: newValue.value, tag: newTag.value.trim() || null }),
        });
        if (!res.ok) {
            const body = await res.json().catch(() => null);
            const message = body?.error ?? "Error adding secret";
            toast(message, { type: "error" });
            return;
        }
        emit("refresh");
        cancelAddRow();
    } catch (error) {
        console.error("Error adding secret:", error);
        toast("Error adding secret", { type: "error" });
        cancelAddRow();
    }
}
</script>
