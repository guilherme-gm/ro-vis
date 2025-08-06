<script setup lang="ts">
import type { LocalizableString } from "@/models/LocalizableString";
import { LocalizableStringComponent } from "./LocalizableString";
import { ref } from "vue";

const props = defineProps<{ string: LocalizableString }>();

let localizableString: LocalizableStringComponent | null = null;
if (props.string.Kind === 'localized') {
	localizableString = LocalizableStringComponent.use(ref(props.string));
}
</script>

<template>
	<span v-if="string.Kind === 'static'">{{ string.Value }}</span>
	<span v-else-if="localizableString && string.Kind === 'localized'">{{ localizableString.displayString.value }}</span>
	<span v-else>[Failed to load string]</span>
</template>
