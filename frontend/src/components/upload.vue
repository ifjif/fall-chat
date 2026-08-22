<template>
  <input id="fallchat-custome-file-input" style="display: none;" :accept="accept" type="file" @change="onchange" />
  <label for="fallchat-custome-file-input">
    <slot />
  </label>
</template>
<script setup>
const props = defineProps({
  accept: {
    type: String, // 接受的文件格式，image/* video/* etc.
  }
})
const emit = defineEmits(['change'])

const onchange = (e) => {
  emit('change', e)
}
</script>
