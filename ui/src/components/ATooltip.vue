<script setup lang="ts">
defineProps<{ text: string; }>();
</script>

<template>
	<a :data-tooltip="text" class="css-tooltip">
		<slot></slot>
	</a>
</template>

<style>
	.css-tooltip {
		position: relative;
	}

	.css-tooltip:hover:after {
		content: attr(data-tooltip);
		background: #000;
		padding: 5px;
		border-radius: 3px;
		display: inline-block;
		position: absolute;
		transform: translate(-50%, -100%);
		margin: 0 auto;
		color: #fff;
		min-width: 100px;
		min-width: 150px;
		top: -5px;
		left: 50%;
		text-align: center;
		font-size: 0.825rem;
	}

	.css-tooltip:hover:before {
		top: -5px;
		left: 50%;
		border: solid transparent;
		content: " ";
		height: 0;
		width: 0;
		position: absolute;
		pointer-events: none;
		border-color: rgba(0, 0, 0, 0);
		border-top-color: #000;
		border-width: 5px;
		margin-left: -5px;
		transform: translate(0, 0px);
	}
</style>
