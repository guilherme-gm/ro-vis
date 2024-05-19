<script setup lang="ts">
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
</script>

<template>
	<BsAccordion>
		<BsAccordionItem title="2012-01-01">
			<BsAccordion>
				<BsAccordionItem title="2012-01-01-ragexe.gpf">
					Lorem, ipsum dolor sit amet consectetur adipisicing elit. Odit laborum suscipit et fuga eligendi amet commodi totam dolore nulla rerum magnam, ipsa sint possimus, minus enim officia minima ad eaque.
				</BsAccordionItem>
			</BsAccordion>
		</BsAccordionItem>
		<BsAccordionItem title="2012-01-03">
			<strong>This is the first item's accordion body.</strong> It is shown by default, until the collapse
			plugin adds the appropriate classes that we use to style each element. These classes control the
			overall appearance, as well as the showing and hiding via CSS transitions. You can modify any of
			this with custom CSS or overriding our default variables. It's also worth noting that just about any
			HTML can go within the <code>.accordion-body</code>, though the transition does limit overflow.
		</BsAccordionItem>
	</BsAccordion>
</template>
