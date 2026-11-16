<template>
  <div class="text-ellipsis-pro" :class="{
    'text-ellipsis-pro--expanded': expanded,
    'text-ellipsis-pro--overflow': overflow,
    'text-ellipsis-pro--scroll': mode === 'scroll'
  }">
    <div v-if="mode === 'scroll'" ref="wrapperRef" class="text-ellipsis-pro__scroll-wrapper" @mouseenter="onScrollEnter"
      @mouseleave="onScrollLeave">
      <div ref="innerRef" class="text-ellipsis-pro__scroll-inner" :style="{
        transform: scrolling ? innerRef?.style.transform : 'translateX(0)'
      }">
        <slot>{{ displayText }}</slot>
      </div>
    </div>

    <div v-else ref="contentRef" class="text-ellipsis-pro__content"
      :class="{ 'text-ellipsis-pro__content--expanded': expanded }" :style="clampStyle">
      <slot>{{ displayText }}</slot>
    </div>

    <div v-if="showMask && isClampMode && !expanded && overflow" class="text-ellipsis-pro__mask" :style="maskStyle" />

    <div v-if="expandable && isClampMode && (overflow || expanded)" class="text-ellipsis-pro__actions">
      <slot name="action" :expanded="expanded" :toggle="toggle">
        <button class="text-ellipsis-pro__btn" @click="toggle()">
          {{ expanded ? collapseText : expandText }}
        </button>
      </slot>
    </div>
  </div>
</template>
<script setup>
import {
  ref,
  computed,
  watch,
  onMounted,
  onBeforeUnmount,
  nextTick
} from 'vue'

/*
interface Props {
  text?: string
  lines?: number
  mode?: 'clamp' | 'scroll'
  expandText?: string
  collapseText?: string
  expandable?: boolean
  autoExpand?: boolean
  autoExpandDelay?: number
  collapseDelay?: number
  showMask?: boolean
  scrollSpeed?: number
  scrollPause?: number
}
*/

const props = defineProps({
  text: {
    type: String,
    default: '',
  },
  lines: {
    type: Number,
    default: 2
  },
  mode: {
    type: String,
    default: 'clamp', // scroll 文本溢出时自动横向滚动
  },
  expandText: {
    type: String,
    default: '展开'
  },
  collapseText: {
    type: String,
    default: '收起'
  },
  expandable: {
    type: Boolean,
    default: true,
  },
  autoExpand: {
    type: Boolean,
    default: false, // 检测到溢出后自动展开
  },
  autoExpandDelay: {
    type: Number,
    default: 1500 // 溢出展开延迟
  },
  collapseDelay: {
    type: Number,
    default: 3000 // 展开后延迟自动收起
  },
  showMask: {
    type: Boolean,
    default: true // 收起时显示渐变遮罩
  },
  scrollSpeed: {
    type: Number,
    default: 40 // 滚动速度，越小越快
  },
  scrollPause: {
    type: Number,
    default: 1000 // 滚动到底部后暂停时间
  }
})

const emit = defineEmits(['change', 'overflow', 'scrollEnd'])

const expanded = ref(false)
const overflow = ref(false)
const scrolling = ref(false)
const contentRef = ref(null)
const innerRef = ref(null)
const wrapperRef = ref(null)

let observer = null
let autoExpandTimer = null
let collapseTimer = null
let scrollTimer = null
let scrollPaused = false

const displayText = computed(() => props.text ?? '')

const isClampMode = computed(() => props.mode === 'clamp')

const clampStyle = computed(() => {
  if (isClampMode.value && !expanded.value) {
    return {
      display: '-webkit-box',
      WebkitLineClamp: props.lines,
      WebkitBoxOrient: 'vertical',
      overflow: 'hidden',
      textOverflow: 'ellipsis'
    }
  }
  return {}
})

const maskStyle = computed(() => {
  if (!props.showMask || expanded.value || !overflow.value) return {}

  return {
    position: 'absolute',
    right: 0,
    bottom: 0,
    width: '60px',
    height: '1.6em',
    background: 'linear-gradient(to right, transparent, #fff)'
  }
})

function checkOverflow() {
  const el = contentRef.value
  if (!el) return

  const isOverflow = el.scrollHeight > el.clientHeight + 1
  overflow.value = isOverflow
  emit('overflow', isOverflow)

  if (props.autoExpand && isOverflow && !expanded.value) {
    startAutoExpand()
  }
}

function startAutoExpand() {
  clearAutoExpand()
  autoExpandTimer = window.setTimeout(() => {
    if (overflow.value && !expanded.value) {
      toggle(true)
    }
  }, props.autoExpandDelay)
}

function clearAutoExpand() {
  if (autoExpandTimer) {
    clearTimeout(autoExpandTimer)
    autoExpandTimer = null
  }
}

function startCollapseDelay() {
  clearCollapseDelay()
  collapseTimer = window.setTimeout(() => {
    if (expanded.value) {
      toggle(false)
    }
  }, props.collapseDelay)
}

function clearCollapseDelay() {
  if (collapseTimer) {
    clearTimeout(collapseTimer)
    collapseTimer = null
  }
}

async function update() {
  await nextTick()
  checkOverflow()

  if (props.mode === 'scroll') {
    await nextTick()
    initScroll()
  }
}

watch(() => [props.text, props.lines, props.mode], update)

onMounted(() => {
  update()

  const el = contentRef.value
  if (el && typeof ResizeObserver !== 'undefined') {
    observer = new ResizeObserver(() => checkOverflow())
    observer.observe(el)
  }
})

onBeforeUnmount(() => {
  observer?.disconnect()
  clearAutoExpand()
  clearCollapseDelay()
  stopScroll()
})

function toggle(value) {
  expanded.value = value ?? !expanded.value
  emit('change', expanded.value)

  if (expanded.value) {
    clearCollapseDelay()
    if (props.collapseDelay > 0) {
      startCollapseDelay()
    }
  } else {
    clearCollapseDelay()
  }
}

function initScroll() {
  stopScroll()

  const wrapper = wrapperRef.value
  const inner = innerRef.value
  if (!wrapper || !inner) return

  if (inner.scrollWidth <= wrapper.clientWidth) {
    scrolling.value = false
    return
  }

  scrolling.value = true
  startScroll()
}

function startScroll() {
  stopScroll()

  const wrapper = wrapperRef.value
  const inner = innerRef.value
  if (!wrapper || !inner) return

  let offset = 0
  const maxScroll = inner.scrollWidth - wrapper.clientWidth

  const step = () => {
    if (!scrollPaused) {
      offset += 1

      if (offset >= maxScroll) {
        emit('scrollEnd')

        scrollTimer = window.setTimeout(() => {
          offset = 0
          inner.style.transform = `translateX(0px)`
          startScroll()
        }, props.scrollPause)
        return
      }

      inner.style.transform = `translateX(-${offset}px)`
    }

    scrollTimer = window.setTimeout(step, props.scrollSpeed)
  }

  step()
}

function stopScroll() {
  if (scrollTimer) {
    clearTimeout(scrollTimer)
    scrollTimer = null
  }
}

function onScrollEnter() {
  scrollPaused = true
}

function onScrollLeave() {
  scrollPaused = false
}

/*
 * 
  <div style="background-color:yellow;">
    <!-- 默认多行省略 -->
    <TextEllipsis :text="longText" :collapse-delay="0" />

    <!-- 单行省略，带渐变遮罩 -->
    <TextEllipsis :text="title" :lines="1" />

    <!-- 自动展开，3秒后自动收起 -->
    <TextEllipsis :text="desc" :lines="2" :auto-expand="false" :collapse-delay="3000" />

    <!-- 滚动播放模式 -->
    <TextEllipsis :text="notice" :lines="4" mode="scroll" :scroll-speed="30" :scroll-pause="800" />

    <!-- 自定义操作按钮 -->
    <TextEllipsis :text="longText" :lines="2">
      <template #action="{ expanded, toggle }">
        <button @click="toggle()">
          {{ expanded ? '收起' : '查看更多' }}
        </button>
      </template>
    </TextEllipsis>
  </div>
 *
 *
 */
</script>

<style scoped>
.text-ellipsis-pro {
  position: relative;
  width: 100%;
}

.text-ellipsis-pro__content {
  line-height: 1.6;
  word-break: break-all;
  transition: max-height 0.3s ease;
}

.text-ellipsis-pro__content--expanded {
  display: block;
  overflow: visible;
  text-overflow: clip;
}

.text-ellipsis-pro__mask {
  pointer-events: none;
}

.text-ellipsis-pro__actions {
  margin-top: 4px;
  display: flex;
  align-items: center;
}

.text-ellipsis-pro__btn {
  padding: 0;
  border: none;
  background: none;
  color: #409eff;
  cursor: pointer;
  font-size: inherit;
}

.text-ellipsis-pro__scroll-wrapper {
  overflow: hidden;
  white-space: nowrap;
}

.text-ellipsis-pro__scroll-inner {
  display: inline-block;
  white-space: nowrap;
  will-change: transform;
}
</style>
