<template>
  <div>
    <fl-form ref="formRef" :rules="rules" :model="formData">
      <fl-form-item prop="avatar" label="头像">
        <div>
          <ImageUpload :url="formData.avatar" @change="fileChange"></ImageUpload>
        </div>
      </fl-form-item>
      <fl-form-item label="账号">
        <div>{{ formData.username }}</div>
      </fl-form-item>
      <fl-form-item prop="nickname" label="昵称">
        <fl-input v-model="formData.nickname"></fl-input>
      </fl-form-item>
      <fl-form-item prop="email" label="邮箱">
        <fl-input v-model="formData.email"></fl-input>
      </fl-form-item>
      <fl-form-item label="手机">
        <fl-input v-model="formData.phone"></fl-input>
      </fl-form-item>
      <fl-form-item prop="search_flag" label="搜索标志">
        <fl-radio-group v-model="formData.search_flag" size="small">
          <fl-radio label="所有人可搜" :value="1"></fl-radio>
          <fl-radio label="仅通过id" :value="0"></fl-radio>
        </fl-radio-group>
      </fl-form-item>
      <fl-form-item label="状态">
        <fl-button v-if="formData.status === 1" text disabled type="primary">正常</fl-button>
        <fl-button v-else text disabled type="error">禁用</fl-button>
      </fl-form-item>
      <fl-form-item label="加入时间">
        <fl-button link type="success">{{ formData.created_at }}</fl-button>
      </fl-form-item>
      <fl-form-item>
        <fl-button type="primary" @click="handleProfile">确定</fl-button>
      </fl-form-item>
    </fl-form>
  </div>

</template>
<script setup>
import { updateProfileApi, userProfileApi } from '@/api/user';
import ImageUpload from '@/components/image-upload.vue';
import { h, ref, onMounted } from 'vue'
import { FlMessageBox, FlInput } from 'fall-ui';
import { uploadFile } from '@/utils/file_upload';
import { refreshApi } from '@/api/auth';
import { useUserStore } from '@/stores/user_info';
const userStore = useUserStore()

const formRef = ref(null)
const password = ref('')
let file_ = null

const formData = ref({
  avatar: '',
  username: '',
  nickname: '',
  email: '',
  phone: '',
  search_flag: 1,
  status: 1,
  created_at: ''
})

const rules = {
  avatar: [
    { type: 'string', required: true, message: '请上传头像', trigge: 'change' }
  ],
  nickname: [
    { type: 'string', required: true, message: '请输入昵称', trigge: 'blur' }
  ],
  email: [
    { type: 'string', required: true, message: '请输入邮箱', trigge: 'blur' }
  ],
  search_flag: [
    { type: 'number', required: true, message: '请设置搜索标志', trigge: 'change' }
  ]
}

const handleProfile = () => {
  const fmb = FlMessageBox({
    type: 'warning', title: '请输入密码',
    content: h(FlInput, {
      type: 'password',
      'modelValue': password.value,
      'onUpdate:modelValue': (v) => { password.value = v }
    }),
    onOk: () => { fmb.close(); doProfileUpdate() }
  })
}

const doProfileUpdate = async () => {
  if (!file_ && !formData.value.avatar)
    await formRef.value.validate()

  if (file_) {
    uploadFile(file_, ({ fileUrl }) => {
      formData.value.avatar = fileUrl
      updateProfile()
    })

    return
  }

  updateProfile()
}

const updateProfile = async () => {
  const data = {
    "email": formData.value.email,
    "nickname": formData.value.nickname,
    "password": password.value,
    "phone": formData.value.phone,
    "search_flag": formData.value.search_flag,
    "avatar": formData.value.avatar
  }

  await updateProfileApi(data)
  // 刷新缓存
  const newToken = await refreshApi(null, {
    withCredentials: true
  })
  localStorage.setItem('token', newToken.access_token)
  fetchProfile()
  password.value = ''
}


const fileChange = (file) => {
  file_ = file
}

const fetchProfile = async () => {
  const data = await userProfileApi()
  formData.value.avatar = data.avatar
  formData.value.username = data.username
  formData.value.nickname = data.nickname
  formData.value.email = data.email
  formData.value.phone = data.phone
  formData.value.search_flag = data.search_flag
  formData.value.created_at = new Date(data.created_at)

  userStore.updateInfo(data)
}

onMounted(() => {
  fetchProfile()
})
</script>
