<template>
  <div class="grid-container" :style="{
    gap: typeof gap === 'number' ? `${gap}px` : gap,
    // 核心样式：自动填充，每项最小宽度为 minItemWidth，剩余空间自动均分(变宽)
    // 【核心升级】根据是否传入 maxItemWidth 决定使用哪种拉伸模式
    gridTemplateColumns: maxItemWidth
      ? `repeat(${cols}, minmax(${toPx(minItemWidth)}, ${toPx(maxItemWidth)}))`
      : `repeat(${cols}, minmax(${toPx(minItemWidth)}, 1fr))`,
    justifyItems: 'center',
    justifyContent: 'center',
    alignItems: 'center',
  }">
    <slot />
  </div>
</template>

<script setup>
const props = defineProps({
  cols: {
    type: [String, Number],
    default: 'auto-fit'
  },
  minItemWidth: { // 子项的最小宽度
    type: Number,
    default: 120 // 默认最小宽度 120px
  },
  maxItemWidth: { // 子项的最大拉伸宽度(设置为false则表示不限制,无限拉伸)
    type: [Number, String, Boolean],
    default: false
  },
  gap: {
    type: Number,
    default: 12,
  }
})
// 辅助函数：将传入的数字或字符串统一处理为 CSS 长度值
const toPx = (val) => {
  return typeof val === 'number' ? `${val}px` : val
}
</script>

<style scoped>
.grid-container {
  display: grid;
}

.grid-container>* {
  width: 100%;
  /* 确保子元素填充满 grid 分配的单元格 */
  box-sizing: border-box;
}
</style>
