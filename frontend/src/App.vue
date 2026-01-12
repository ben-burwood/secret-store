<template>
    <Layout v-model:selectedPage="selectedPage">
        <component :is="currentComponent" />
    </Layout>
</template>

<script setup lang="ts">
import Layout from "./components/Layout.vue";
import GenerateSecret from "./components/GenerateSecret.vue";
import SecretList from "./components/SecretList.vue";
import AuthKey from "./components/AuthKey.vue";
import BackupRestore from "./components/BackupRestore.vue";
import { ref, computed } from "vue";
import { Page } from "@/types/page";

const selectedPage = ref(Page.HOME);
const currentComponent = computed(() => {
    switch (selectedPage.value) {
        case Page.HOME:
            return SecretList;
        case Page.GENERATE:
            return GenerateSecret;
        case Page.AUTHENTICATION:
            return AuthKey;
        case Page.SETTINGS:
            return BackupRestore;
        default:
            return SecretList;
    }
});
</script>
