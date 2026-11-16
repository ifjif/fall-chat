<!-- 外部设置 aspect-ratio时，请让一边auto -->
<template>
  <div class="base-image-container" :style="containerStyle">
    <img v-if="fullUrl" :src="fullUrl" :alt="alt" class="base-image-img" @error="handleError" @load="handleLoad" />
    <!-- 缺省/加载失败时的占位区域 -->
    <div v-else class="base-image-placeholder">
      <slot name="placeholder">
        <span>暂无图片</span>
      </slot>
    </div>
  </div>
</template>

<script setup>
import { computed, defineProps, ref } from 'vue';

// 定义 Props 接口
//interface Props {
//  uri?: string;       // 图片路径或完整URL
//  width?: string | number; // 宽度，支持 '100px', '50%', 100 等
//  height?: string | number; // 高度
//  alt?: string;       // 替代文本
//  baseUrl?: string;   // 基础域名，例如 'https://api.example.com'
//}

const props = defineProps({
  uri: {
    type: String,
    default: ''
  },
  width: {
    type: [String, Number],
    default: '100%'
  },
  height: {
    type: [String, Number],
    default: 'auto'
  },
  alt: {
    type: String,
    default: 'image',
  },
  baseUrl: {
    type: String,
    default: ''
  },
  remoteUrl: {
    type: Function
  }
})
const emit = defineEmits(['load'])

// 状态管理：用于处理图片加载失败的情况
const isLoadError = ref(false);
const remote = ref('')

// 计算完整的 URL
const fullUrl = computed(() => {
  if (!props.uri) return '';

  // 如果已经是完整 URL (http/https) 或者 base64，直接返回
  if (props.uri.startsWith('http')
    || props.uri.startsWith('https:')
    || props.uri.startsWith('data:')
    || props.uri.startsWith('blob:')
  ) {
    return props.uri;
  }

  if (remote.value) return remote.value

  if (props.remoteUrl) {
    props.remoteUrl(props.uri).then(r => {
      remote.value = r
    })
    return remote.value
  }

  // 否则拼接 baseUrl
  // 处理 baseUrl 结尾和 uri 开头的斜杠，防止出现 //
  const base = props.baseUrl.replace(/\/+$/, '');
  const path = props.uri.replace(/^\/+/, '');

  return base ? `${base}/${path}` : props.uri;
});

// 计算容器样式
const containerStyle = computed(() => {
  return {
    width: typeof props.width === 'number' ? `${props.width}px` : props.width,
    height: typeof props.height === 'number' ? `${props.height}px` : props.height,
  };
});

// 处理加载错误
const handleError = () => {
  isLoadError.value = true;
};

// 处理加载成功（如果需要重置错误状态）
const handleLoad = () => {
  isLoadError.value = false;
  emit('load')
};
</script>

<style scoped>
.base-image-container {
  position: relative;
  display: inline-block;
  /* 或者 block，视需求而定 */
  overflow: hidden;
  background-color: #f5f7fa;
  /* 缺省背景色 */
}

.base-image-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  /* 保持比例填充，裁剪多余部分 */
  display: block;
}

.base-image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #909399;
  font-size: 14px;
  background-color: #f5f7fa;
}
</style>
