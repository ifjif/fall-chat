<template>
  <span :style="styles" class="badge" v-if="show">
    <template v-if="valueType === 'string'">
      {{ value }}
    </template>
    <template v-else-if="valueType === 'number'">
      <template v-if="value > max">{{ max }}+</template>
      <template v-else>
        {{ value }}
      </template>
    </template>
  </span>
</template>
<script setup>
import { computed } from 'vue'
const props = defineProps({
  value: {
    type: [String, Number]
  },
  max: {
    type: [String, Number]
  },
  showZero: {
    type: Boolean,
    default: false
  },
  offset: {
    type: Array,
    validate(v) {
      return v && v.length === 2 && (typeof v[0] === 'number') && (typeof v[1] === 'number')
    },
    default: [0, 0]
  }
})

const valueType = computed(() => {
  if (typeof props.value === 'string') {
    return 'string'
  } else if (typeof props.value === 'number') {
    return 'number'
  }
})

const show = computed(() => {
  return props.showZero || props.value !== 0
})

const styles = computed(() => {
  const styles = {
    top: `${props.offset[0]}px`,
    right: `${props.offset[0]}px`,
  }

  return styles
})
</script>

<style scoped>
.badge {
  position: absolute;
  padding: 0 5px;
  color: white;
  font-size: 12px;
  text-align: center;
  background-color: #f56c6c;
  border-radius: 50px;
  right: 0;
  top: 0;
}
</style>
