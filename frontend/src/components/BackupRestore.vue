<template>
    <div class="card bg-base-200 shadow-md max-w-md mx-auto">
        <div class="card-body">
            <h2 class="card-title">Backup & Restore Secrets</h2>
            <div class="flex flex-col gap-4">
                <button class="btn btn-primary" @click="backup">Backup</button>

                <div class="flex flex-row gap-2">
                    <input ref="fileInput" type="file" accept="application/json" class="file-input file-input-bordered" @change="onFileChange" />
                    <button class="btn btn-primary" @click="restore" :disabled="!selectedFile">Restore</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { SERVER_URL } from "@/main";
import { toast } from "vue3-toastify";

async function backup() {
    let secrets;
    try {
        const response = await fetch(`${SERVER_URL}/secrets`);
        const data = await response.json();
        secrets = data;
    } catch (error) {
        console.error("Error fetching secrets:", error);
    }

    const blob = new Blob([JSON.stringify(secrets, null, 2)], { type: "application/json" });
    const url = URL.createObjectURL(blob);

    const a = document.createElement("a");
    a.href = url;
    a.download = "secrets.json";
    a.click();
    URL.revokeObjectURL(url);
}

const fileInput = ref<HTMLInputElement | null>(null);
const selectedFile = ref<File | null>(null);
function onFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    selectedFile.value = input.files && input.files.length > 0 ? input.files[0] : null;
}

async function restore() {
    if (!selectedFile.value) {
        toast("No file selected", { type: "error" });
        return;
    }

    let secrets;
    try {
        const json = await selectedFile.value.text();
        secrets = JSON.parse(json);
    } catch (error) {
        console.error("Error parsing JSON:", error);
        toast("Invalid JSON format", { type: "error" });
        return;
    }

    const formData = new FormData();
    formData.append("secrets", JSON.stringify(secrets));

    try {
        const response = await fetch(`${SERVER_URL}/restore`, {
            method: "POST",
            body: formData,
        });
        if (!response.ok) throw new Error("Restore failed");
    } catch (err) {
        toast("Restore failed: " + err.message, { type: "error" });
    }
}
</script>
