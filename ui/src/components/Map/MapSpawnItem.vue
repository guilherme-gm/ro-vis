<script setup lang="ts">
import type { MapSpawn } from "@/models/Map";
import LocalizableString from "../LocalizableString/LocalizableString.vue";
import { computed } from "vue";

const props = defineProps({
	spawn: {
		type: Object as () => MapSpawn,
		required: true,
	},
	kind: {
		type: String,
		validator: (value: string) => ['added', 'removed', 'unchanged', 'none'].includes(value),
		default: 'none',
	},
});

const hiddenPart = computed(() => props.spawn.Name2 ? ` (#${props.spawn.Name2} / #${props.spawn.SpriteId})` : `(#${props.spawn.SpriteId})`);

const elementTable: Record<number, string> = {
	0: 'Neutral',
	1: 'Water',
	2: 'Earth',
	3: 'Fire',
	4: 'Wind',
	5: 'Poison',
	6: 'Holy',
	7: 'Dark',
	8: 'Ghost',
	9: 'Undead',
};

const raceTable: Record<number, string> = {
	0: 'Formless',
	1: 'Undead',
	2: 'Brute',
	3: 'Plant',
	4: 'Insect',
	5: 'Fish',
	6: 'Demon',
	7: 'DemiHuman',
	8: 'Angel',
	9: 'Dragon',
	10: 'Player',
	11: 'Boss',
	12: 'NonBoss',
	14: 'NonDemiHuman',
	15: 'NonPlayer',
	16: 'DemiPlayer',
	17: 'NonDemiPlayer',
};

const sizeTable: Record<number, string> = {
	0: 'Small',
	1: 'Medium',
	2: 'Large',
};

const element = computed(() => {
	const level = Math.floor(props.spawn.Element / 20);
	const ele = props.spawn.Element % 20;
	return `${elementTable[ele]} ${level}`;
});

const race = computed(() => raceTable[props.spawn.Race]);
const size = computed(() => sizeTable[props.spawn.Size]);

</script>

<template>
	<div class="d-flex flex-column" :class="{ [`diff-${kind}`]: true }">
		<span>
			({{ spawn.Type }}) {{ spawn.Amount }}x <LocalizableString :string="spawn.Name1" /> {{ hiddenPart }}
		</span>
		<span>
			(Ele: {{ element }} / Race: {{ race }} / Size: {{ size }})
		</span>
	</div>
</template>
