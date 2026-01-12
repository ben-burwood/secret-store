<template>
    <tr :class="{ 'bg-success/10': isEditing }">
        <th>{{ secret.id }}</th>
        <td>
            <template v-if="isEditing">
                <input v-model="editKey" class="input input-bordered w-full" />
            </template>
            <template v-else>
                {{ secret.key }}
            </template>
        </td>
        <td>
            <template v-if="isEditing">
                <input v-model="editValue" class="input input-bordered w-full" @keyup.enter="updateSecret" />
            </template>
            <template v-else>
                <SecretDisplay :secret="secret.value" :showSecret="showSecret" />
            </template>
        </td>
        <td>
            <div v-if="!isEditing" class="flex flex-row gap-2">
                <button class="btn btn-outline btn-primary btn-square" @click="isEditing = true">
                    <Pencil :size="18" />
                </button>
                <button class="btn btn-outline btn-error btn-square" @click="deleteSecret()">
                    <Trash :size="18" />
                </button>
            </div>
            <ConfirmReject v-else @confirm="updateSecret" @reject="isEditing = false" :canConfirm="editKey !== '' && editValue !== ''" />
        </td>
    </tr>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, ref } from "vue";
import { Trash, Pencil } from "lucide-vue-next";
import SecretDisplay from "@/components/SecretDisplay.vue";
import { toast } from "vue3-toastify";
import { SERVER_URL } from "@/main";
import ConfirmReject from "@/components/secrets-table/ConfirmReject.vue";

const props = defineProps<{
    secret: {
        id: string;
        key: string;
        value: string;
    };
    showSecret: boolean;
}>();

const emit = defineEmits(["refresh"]);

const isEditing = ref(false);
const editKey = ref(props.secret.key);
const editValue = ref(props.secret.value);

async function updateSecret() {
    try {
        await fetch(`${SERVER_URL}/secrets/${props.secret.id}`, {
            method: "PATCH",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ key: editKey.value, value: editValue.value }),
        });
        emit("refresh");
    } catch (error) {
        console.error("Error updating secret:", error);
        toast("Error updating secret", { type: "error" });
    } finally {
        isEditing.value = false;
    }
}

async function deleteSecret() {
    try {
        await fetch(`${SERVER_URL}/secrets/${props.secret.id}`, { method: "DELETE" });
        emit("refresh");
    } catch (error) {
        console.error("Error deleting secret:", error);
        toast("Error deleting secret", { type: "error" });
    }
}
</script>
