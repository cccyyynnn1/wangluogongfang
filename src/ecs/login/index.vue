<template>
  <div
    class="login-container"
    :style="{
      background: `url(${
        loginBg || require('@/assets/login_images/login.png')
      })  center center / cover fixed no-repeat`,
    }"
  >
    <div class="logo-tip">
      <img v-if="loginLogo" alt="" class="bg-logo" :src="loginLogo" />
      <span style="font-size: 22px; font-weight: 600; color: #676284">{{ loginTip }}</span>
    </div>
    <el-row>
      <el-col :lg="14" :md="11" :sm="24" :xl="14" :xs="24">
        <div style="color: transparent"></div>
      </el-col>
      <el-col :lg="9" :md="12" :sm="24" :xl="8" :xs="24">
        <el-form ref="formRef" class="login-form" label-position="left" :model="form" :rules="rules">
          <div class="title-tips">{{ title }}</div>
          <el-form-item prop="username">
            <el-input
              v-model.trim="form.username"
              v-focus
              :placeholder="translateTitle('请输入用户名')"
              tabindex="1"
              type="text"
            >
              <template #prefix>
                <vab-icon icon="user-line" />
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              :key="passwordType"
              ref="passwordRef"
              v-model.trim="form.password"
              :placeholder="translateTitle('请输入密码')"
              tabindex="2"
              :type="passwordType"
              @keyup.enter="handleLogin"
            >
              <template #prefix>
                <vab-icon icon="lock-line" />
              </template>
              <template v-if="passwordType === 'password'" #suffix>
                <vab-icon class="show-password" icon="eye-off-line" @click="handlePassword" />
              </template>
              <template v-else #suffix>
                <vab-icon class="show-password" icon="eye-line" @click="handlePassword" />
              </template>
            </el-input>
          </el-form-item>
          <!-- 验证码验证逻辑需自行开发，如不需要验证码功能建议注释 -->
          <el-form-item v-if="verifyImgStatus" prop="verifyCode">
            <el-input
              v-model.trim="form.verifyCode"
              :placeholder="translateTitle('验证码') + previewText"
              tabindex="3"
              type="text"
              @keyup.enter="handleLogin"
            >
              <template #prefix>
                <vab-icon icon="barcode-box-line" />
              </template>
            </el-input>
            <el-image class="code" :src="codeUrl" @click="changeCode" />
          </el-form-item>
          <el-form-item>
            <el-button class="login-btn" :loading="loading" type="primary" @click="handleLogin">
              {{ translateTitle('登录') }}
            </el-button>
          </el-form-item>
          <el-form-item v-if="!showBtns">
            <div class="btns">
              <template v-for="button in buttons" :key="button.text">
                <el-button
                  v-if="button.disable != 0"
                  link
                  size="large"
                  :type="button.type"
                  @click="thirdpPartyLogin(button.id)"
                >
                  {{ button.text }}
                </el-button>
              </template>
            </div>
          </el-form-item>
        </el-form>
      </el-col>
      <el-col :lg="1" :md="1" :sm="24" :xl="1" :xs="24">
        <div style="color: transparent"></div>
      </el-col>
    </el-row>
    <div class="icp">{{ icp }}</div>
  </div>

  <LogTo
    v-if="visible"
    :mode="remark"
    :visible="visible"
    @on-closeEvent="handleClose"
    @on-loginEvent="loginOtherType"
  />
</template>

<script>
  import LogTo from './components/log_to.vue'
  import { useSettingsStore } from '@/store/modules/settings'
  import { useUserStore } from '@/store/modules/user'
  import { translate } from '@/i18n'
  import { isPassword } from '@/utils/validate'
  import { getCheckActivateApi } from '~/src/api-ecs/public'
  import { getSysInfoApi, getOtherLoginStatusApi, getVerifyImgStatusApi } from '~/src/api-ecs/login'
  import { usePubilcStore } from '@/store/modules/public'
  import { tokenName } from '@/config'
  import sha1 from 'sha1'

  export default defineComponent({
    name: 'Login',
    components: {
      LogTo,
    },
    directives: {
      focus: {
        mounted(el) {
          el.querySelector('input').focus()
        },
      },
    },
    setup() {
      const router = useRouter()
      const userStore = useUserStore()
      const settingsStore = useSettingsStore()
      const publicStore = usePubilcStore()
      const icon = useFavicon()
      let requestId = new Date().getTime()
      const nimLength = 8
      const handleClose = () => {
        state.visible = false
      }

      const thirdpPartyLogin = (id) => {
        state.remark = id
        if (id == 1) {
          state.visible = true
        } else if (id == 3 && state.otherLoginType.ssoUrl) {
          const temp = window.location.origin
          window.location.href = `${state.otherLoginType.ssoUrl}?backUrl=${temp}`
        }
      }
      const loginOtherType = (form) => {
        userStore
          .loginOther({
            ...form,
            // password: sha1(form.password),
          })
          .then((res) => {
            const { data } = res
            if (data?.expire) {
              data?.token && userStore.setToken(data?.token)
              userStore.setExpire(data?.expire)
              return settingsStore.changeToolboxVisible(true, 'reset-password')
            }
            if (data?.token) return userStore.afterLogin(data?.token, tokenName)
          })
          .catch((err) => {
            err?.code === 50 && changeCode()
          })
      }
      const login = (form) => {
        userStore
          .login({
            ...form,
            password: sha1(form.password),
            loginType: 2,
            requestId: requestId,
          })
          .then(async (res) => {
            const { data } = res
            if (data?.expire) {
              userStore.setToken(data?.token)
              userStore.setExpire(data?.expire)
              return settingsStore.changeToolboxVisible(true, 'reset-password')
            }
            if (data?.token) return userStore.afterLogin(data?.token, tokenName)
          })
          .then(() => {
            publicStore.setGlobalLoginMessageList()
          })
          .catch((err) => {
            err?.code === 50 && changeCode()
          })
      }

      const validateUsername = (rule, value, callback) => {
        if ('' === value) callback(new Error(translate('用户名不能为空')))
        else callback()
      }
      const validatePassword = (rule, value, callback) => {
        if (!isPassword(value)) callback(new Error(translate(`密码不能少于${nimLength}位`)))
        else callback()
      }

      const state = reactive({
        loginBg: '',
        loginLogo: '',
        loginTip: '',
        icp: '',
        formRef: null,
        passwordRef: null,
        visible: false,
        remark: '',
        verifyImgStatus: false,
        showBtns: true,
        form: {
          username: '',
          password: '',
          verifyCode: '',
        },
        showstatus: {
          ecs_ldap: 0,
          ecs_radius: 0,
          ecs_sso: 0,
        },
        otherLoginType: {
          ssoUrl: '',
        },
        buttons: [
          { type: 'primary', text: 'LDAP', key: 'ecs_ldap', disable: 0, id: 2 },
          { type: 'primary', text: 'RADIUS', key: 'ecs_radius', disable: 0, id: 1 },
          { type: 'primary', text: 'SSO', key: 'ecs_sso', disable: 0, id: 3 },
        ],
        rules: {
          username: [
            {
              required: true,
              trigger: 'blur',
              validator: validateUsername,
            },
          ],
          password: [
            {
              required: true,
              trigger: 'blur',
              validator: validatePassword,
            },
          ],
          verifyCode: [
            {
              required: true,
              trigger: 'blur',
              message: '验证码不能空',
            },
          ],
        },
        loading: false,
        passwordType: 'password',
        redirect: undefined,
        timer: 0,
        codeUrl: `/v3/authority/auth/getImgVerifyCode?requestId=${requestId}`,
        previewText: '',
        title: '',
      })

      const ssoBackEvent = () => {
        const pramas = window.location.search.split('&')[0].split('?')
        let key = ''
        let token = ''
        pramas?.forEach((item) => {
          key = item.split('=')[0]
          token = item.split('=')[1]
          if (key == 'token' && token) {
            userStore.afterLogin(token, tokenName)
            window.location.href = window.location.origin
          }
        })
      }

      const getVerifyImgStatus = async () => {
        const { data } = await getVerifyImgStatusApi()
        state.verifyImgStatus = data
      }

      onMounted(async () => {
        getSysInfo()
        await getVerifyImgStatus()
        handleCheckActivate()
        await getOtherLoginStatus()
        ssoBackEvent()
        await settingsStore.setSecConfig()
      })

      const getOtherLoginStatus = async () => {
        const { data } = await getOtherLoginStatusApi()
        state.buttons.forEach((item) => {
          for (const key in data) {
            if (item.key == key && key !== 'ecs_sso') {
              item.disable = data[key]
            }
            if (key == 'ecs_sso') {
              item.disable = JSON.parse(data['ecs_sso']?.value).enable
            }
          }
        })
        state.showBtns = state.buttons.every((item) => {
          return item.disable == 0
        })
        state.otherLoginType.ssoUrl = JSON.parse(data.ecs_sso?.value)?.authenticationAddr
      }

      async function getSysInfo() {
        try {
          const { data } = await getSysInfoApi()
          const info = JSON.parse(data[0].value)
          settingsStore.updateState({
            title: info.name,
            logo: info?.logo,
          })
          const title = useTitle()
          title.value = info.name
          icon.value = info?.logo
          state.title = info.name
          state.loginBg = info.loginBg
          state.loginLogo = info.logo
          state.loginTip = info.logginTip
          state.icp = info.icp
        } catch (error) {
          console.error(error)
        }
      }
      const handleCheckActivate = async () => {
        const { data } = await getCheckActivateApi()
        const expUsed = data.expUsed * 1000
        const nowDate = new Date().getTime()
        if (expUsed <= nowDate) {
          router.push('/activate')
        }
      }
      const handlePassword = () => {
        state.passwordType === 'password' ? (state.passwordType = '') : (state.passwordType = 'password')
        nextTick(() => {
          state['passwordRef'].focus()
        })
      }
      const handleLogin = async () => {
        state['formRef'].validate(async (valid) => {
          if (valid)
            try {
              state.loading = true
              login(state.form)
            } finally {
              state.loading = false
            }
        })
      }
      const changeCode = () => {
        requestId = new Date().getTime()
        state.codeUrl = `/v3/authority/auth/getImgVerifyCode?requestId=${requestId}`
        state.form.verifyCode = ''
      }
      onBeforeRouteLeave((to, from, next) => {
        clearInterval(state.timer)
        next()
      })
      return {
        translateTitle: translate,
        ...toRefs(state),
        handlePassword,
        handleLogin,
        changeCode,
        thirdpPartyLogin,
        handleClose,
        loginOtherType,
      }
    },
  })
</script>

<style lang="scss" scoped>
  .login-container {
    height: 100vh;
    background-size: cover;
    position: relative;
  }
  .logo-tip {
    top: 44px;
    left: 36px;
    position: absolute;
    display: flex;
    align-items: center;
  }
  .bg-logo {
    object-fit: contain;
    height: 50px;
    margin-right: 6px;
    background-size: 100% 100% !important;
  }
  .icp {
    bottom: 22px;
    left: 50%;
    transform: translateX(-50%);
    position: absolute;
    font-weight: 400;
    font-size: 14px;
    color: #8883a5;
    line-height: 20px;
  }

  .login-form {
    .btns {
      width: 100%;
      display: flex;
      justify-content: space-around;
      align-items: center;
    }
    position: relative;
    max-width: 100%;
    padding: 4.5vh;
    margin: calc((100vh - 450px) / 2) 5vw 5vw;
    overflow: hidden;
    background: #fff;
    // background: url('~@/assets/login_images/login_form.png');
    // background-size: 100% 100%;
    border-radius: 8px;

    .title {
      font-size: 54px;
      font-weight: 500;
      color: #171e40;
    }

    .title-tips {
      font-weight: 500;
      font-size: 24px;
      color: #676284;
      // line-height: 33px;
      margin-bottom: 25px;
      text-align: center;
    }

    .login-btn {
      display: inherit;
      width: 100%;
      height: 46px;
      margin-top: 5px;
      background: #6954f0;
      font-size: 16px;
      border: 0;

      &:hover {
        opacity: 0.9;
      }
    }

    .tips {
      margin-bottom: 10px;
      font-size: $base-font-size-default;
      color: var(--el-color-white);

      span {
        &:first-of-type {
          margin-right: 16px;
        }
      }
    }

    .title-container {
      position: relative;

      .title {
        margin: 0 auto 40px auto;
        font-size: 34px;
        font-weight: bold;
        color: var(--el-color-primary);
        text-align: center;
      }
    }

    i {
      position: absolute;
      top: 8px;
      left: 15px;
      font-size: 16px;
    }

    .show-password {
      float: right;
      width: 32px;
      height: 32px;
      font-size: 16px;
    }

    :deep() {
      .el-form-item {
        padding-right: 0;
        margin: 18px 0;
        color: #454545;
        background: #f9f9ff;
        border: 1px solid transparent;
        border-radius: 2px;

        &__content {
          min-height: $base-input-height;
          line-height: $base-input-height;
        }

        &__error {
          position: absolute;
          top: 100%;
          left: 18px;
          font-size: $base-font-size-small;
          line-height: 18px;
          color: var(--el-color-error);
        }

        .el-input__wrapper {
          background: #f9f9ff !important;
        }
      }

      .el-input {
        box-sizing: border-box;

        input {
          height: 44px;
          padding-left: 30px;
          line-height: 48px;
          border: 0;
        }

        &__suffix-inner {
          position: absolute;
          right: 65px;
          cursor: pointer;
        }
      }

      .code {
        position: absolute;
        top: 14px;
        right: 4px;
        cursor: pointer;
        border-radius: $base-border-radius;
      }
    }
  }
</style>
