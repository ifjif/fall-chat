<template>
  <p @click="$emit('click')" v-if="more" style="position:relative;text-align:center;cursor: pointer;">
    <FlLoading :loading="loading"></FlLoading>
    更多...
  </p>
</template>

<script setup>
import { FlLoading } from 'fall-ui';
const props = defineProps({
  cursor: {
    type: String,
    default: "",
  },
  more: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false,
  }
})

const emit = defineEmits(['click'])
</script>
