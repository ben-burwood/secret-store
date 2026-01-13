<template>
    <div class="card bg-base-200 shadow-md max-w-md mx-auto">
        <div class="card-body">
            <h2 class="card-title">Import Secrets</h2>
            Import from .env or json
            <input ref="fileInput" type="file" class="file-input file-input-bordered" @change="onFileChange" />
            <button class="btn btn-primary" @click="importFromFile" :disabled="!selectedFile">Import</button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { SERVER_URL } from "@/main";
import { toast } from "vue3-toastify";

const fileInput = ref<HTMLInputElement | null>(null);
const selectedFile = ref<File | null>(null);
function onFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    selectedFile.value = input.files && input.files.length > 0 ? input.files[0] : null;
}

function parseEnv(text: string) {
    const secrets: { key: string; value: string }[] = [];
    for (const line of text.split(/\r?\n/)) {
        const trimmed = line.trim();
        if (!trimmed || trimmed.startsWith("#")) continue;
        const eqIdx = trimmed.indexOf("=");
        if (eqIdx === -1) continue;
        const key = trimmed.slice(0, eqIdx).trim();
        const value = trimmed.slice(eqIdx + 1).trim();
        if (key) secrets.push({ key, value });
    }
    return secrets;
}

function parseJson(json) {
    // If it's an object, convert to array
    if (!Array.isArray(json) && typeof json === "object") {
        return Object.entries(json).map(([key, value]) => ({ key, value }));
    }
    return json;
}

async function importFromFile() {
    if (!selectedFile.value) {
        toast("No file selected", { type: "error" });
        return;
    }

    let secrets;
    try {
        const text = await selectedFile.value.text();
        if (selectedFile.value.name.endsWith(".json")) {
            secrets = parseJson(JSON.parse(text));
        } else if (selectedFile.value.name.endsWith(".env") || selectedFile.value.type === "text/plain") {
            secrets = parseEnv(text);
        }
    } catch (error) {
        console.error("Error parsing File:", error);
        toast("Invalid format", { type: "error" });
        return;
    }

    const formData = new FormData();
    formData.append("secrets", JSON.stringify(secrets));

    try {
        const response = await fetch(`${SERVER_URL}/import`, {
            method: "POST",
            body: formData,
        });
        if (!response.ok) throw new Error();
    } catch (err) {
        toast("Import failed: " + err.message, { type: "error" });
    }
}
</script>
