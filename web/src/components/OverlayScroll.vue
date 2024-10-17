<script setup lang="ts">
import { OverlayScrollbarsComponent } from 'overlayscrollbars-vue'
</script>
<template>
    <OverlayScrollbarsComponent
        :options="{
            scrollbars: {
                theme: 'os-theme-light',
                autoHide: 'leave'
            }
        }"
    >
        <slot></slot>
    </OverlayScrollbarsComponent>
</template>
