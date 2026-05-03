import { fileURLToPath, URL } from "node:url";

import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
    plugins: [vue(), tailwindcss()],
    resolve: {
        alias: {
            "@": fileURLToPath(new URL("./src", import.meta.url)),
        },
    },
    server: {
        proxy: {
            "/web": "http://localhost:8080",
        },
    },
    fmt: {
        printWidth: 150,
        tabWidth: 4,
    },
});
