<template>
    <tr>
        <th>+</th>
        <td>
            <input v-model="newKey" class="input input-bordered w-full" placeholder="Key" />
        </td>
        <td>
            <input v-model="newValue" class="input input-bordered w-full" placeholder="Value" @keyup.enter="addSecret" />
        </td>
        <td>
            <ConfirmReject @confirm="addSecret" @reject="cancelAddRow" :canConfirm="newKey !== '' && newValue !== ''" />
        </td>
    </tr>
</template>

<script setup lang="ts">
import { ref, defineEmits } from "vue";
import ConfirmReject from "@/components/secrets-table/ConfirmReject.vue";
import { toast } from "vue3-toastify";
import { SERVER_URL } from "@/main";

const emit = defineEmits(["refresh", "cancel"]);

const newKey = ref("");
const newValue = ref("");
const cancelAddRow = () => {
    newKey.value = "";
    newValue.value = "";
    emit("cancel");
};

async function addSecret() {
    try {
        await fetch(`${SERVER_URL}/secrets/new`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ key, value }),
        });
        emit("refresh");
    } catch (error) {
        console.error("Error adding secret:", error);
        toast("Error adding secret", { type: "error" });
    } finally {
        cancelAddRow();
    }
}
</script>
