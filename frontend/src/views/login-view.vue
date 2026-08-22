<template>
  <div class="login-layout">
    <div class="login-body">
      <fl-form :model="formData">
        <h3 class="login-title">{{ actions[defaultAction] }}</h3>
        <fl-form-item label="用户名">
          <fl-input v-model="formData.username" width="200"></fl-input>
        </fl-form-item>
        <fl-form-item v-if="defaultAction !== 2" label="密码">
          <fl-input v-model="formData.password" width="200" type="password" showPassword></fl-input>
        </fl-form-item>
        <fl-form-item v-if="defaultAction === 2" label="邮箱">
          <fl-input v-model="formData.email" width="200"></fl-input>
        </fl-form-item>
        <fl-form-item v-if="defaultAction === 2" label="验证码">
          <fl-input v-model="formData.code" width="200">
            <template #append>
              <fl-button :disabled="codeBtnDisabled" @click="codeBtnClick">{{ codeBtnText }}</fl-button>
            </template>
          </fl-input>
        </fl-form-item>
        <fl-form-item v-if="defaultAction === 2" label="新密码">
          <fl-input width="200" v-model="formData.new_pass"></fl-input>
        </fl-form-item>
        <div class="login-btn">
          <fl-button type="primary" @click="sendBtn">确认</fl-button>
        </div>
        <div class="login-switch">
          <fl-button v-if="defaultAction !== 0" link type="success" @click="loginSwitch(0)">登录</fl-button>
          <fl-button v-if="defaultAction !== 1" link type="primary" @click="loginSwitch(1)">注册</fl-button>
          <fl-button v-if="defaultAction !== 2" link type="warning" @click="loginSwitch(2)">忘记密码</fl-button>
        </div>
      </fl-form>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { countDown, padZero } from '@/utils/index'
import { forgetPasswordApi, loginApi, registerApi } from '@/api/user'
import { passRestCodeApi } from '@/api/email'
import router from '@/router'
import { useUserStore } from '@/stores/user_info'
const actions = ['登录', '注册', '忘记密码']
const defaultAction = ref(0)
const formData = ref({
  username: "",
  password: "",
  email: "",
  code: "",
  code_id: "",
  new_pass: ""
})

const codeBtnText = ref('获取')
const codeBtnDisabled = ref(false)

const codeBtnClick = async () => {
  codeBtnDisabled.value = true

  const now = new Date()
  const target = new Date(now.getTime() + 1000 * 60)

  countDown(target, (days, hours, minutes, seconds) => {
    codeBtnText.value = padZero(seconds)
  }, () => {
    codeBtnText.value = '获取'
    codeBtnDisabled.value = false
  })

  try {
    const data = await passRestCodeApi(formData.value)
    formData.value.code_id = data.id
  } catch (error) { }

}

const sendBtn = () => {
  switch (defaultAction.value) {
    case 0:
      doLogin()
      break
    case 1:
      doRegister()
      break
    case 2:
      doForgetPassword()
  }
}

const doLogin = async () => {
  try {
    const data = await loginApi(formData.value)
    localStorage.setItem("token", data.access_token)
    useUserStore().loginSuccess()
    router.push('/')
  } catch (err) {
  }
}

const doRegister = async () => {
  try {
    const data = await registerApi(formData.value)
    console.log(data)
  } catch (err) { }
}

const doForgetPassword = async () => {
  try {
    const data = await forgetPasswordApi(formData.value)
    console.log(data)
  } catch (err) { }
}

const loginSwitch = (action) => {
  defaultAction.value = action
}
</script>
<style scoped>
.login-layout {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-body {
  background-color: beige;
  padding: 20px;
}

.login-title {
  margin: 0;
  text-align: center
}

.login-btn {
  margin: 0 auto;
}

.login-switch {
  display: flex;
  justify-content: space-between;
}
</style>
