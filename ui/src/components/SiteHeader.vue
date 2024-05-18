<script setup lang="ts">
import BsNavbar from "@/components/bootstrap/BsNavbar.vue";
import BsNavLink from "./bootstrap/BsNavLink.vue";
</script>

<template>
	<header>
		<BsNavbar brand-text="RO Vis">
			<template #nav-items>
				<BsNavLink to="/" label="Home" />
			</template>
		</BsNavbar>
	</header>
</template>
