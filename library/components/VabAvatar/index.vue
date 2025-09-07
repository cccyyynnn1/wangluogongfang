<script lang="ts" setup>
  import { useUserStore } from '@/store/modules/user'
  import { useSettingsStore } from '@/store/modules/settings'

  import { translate } from '@/i18n'
  import { VabRoute } from '/#/router'

  const route: VabRoute = useRoute()
  const router = useRouter()
  const settingsStore = useSettingsStore()
  const userStore = useUserStore()
  const { avatar, username } = storeToRefs(userStore)
  const { logout } = userStore

  const active = ref(false)

  const handleVisibleChange = (val: boolean) => {
    active.value = val
  }
  const handleCommand = async (command: string) => {
    switch (command) {
      case 'logout':
        await logout()
        await router.replace('/login')
        break
      case 'reset':
        settingsStore.changeToolboxVisible(true, 'reset-password')
        break
      case 'help':
        window.open('#/help')
        break
    }
  }
</script>

<template>
  <!-- <el-dropdown @command="handleCommand" @visible-change="handleVisibleChange"> -->
  <!-- <span class="avatar-dropdown">
    <el-avatar class="user-avatar" :src="avatar" />
    <div class="user-name">
      <span class="hidden-xs-only" :title="username">{{ username }}</span> -->
  <!-- <vab-icon class="vab-dropdown" :class="{ 'vab-dropdown-active': active }" icon="arrow-down-s-line" /> -->
  <!-- </div>
  </span> -->
  <!-- <template #dropdown> -->
  <!-- <el-dropdown-menu>
      <el-dropdown-item command="reset"> -->
  <!-- <vab-icon icon="shield-cross-line" />
    <span>修改密码</span> -->
  <!-- </el-dropdown-item>
      <el-dropdown-item command="logout"> -->
  <!-- <vab-icon icon="logout-circle-r-line" />
    <span>{{ translate('退出登录') }}</span> -->
  <!-- </el-dropdown-item>
    </el-dropdown-menu> -->
  <!-- </template> -->
  <!-- </el-dropdown> -->
  <div class="avatar">
    <!-- <el-avatar class="user-avatar" :src="avatar" style="height: 60px; width: 60px" /> -->
    <div class="user">
      <span style="font-size: 20px; color: #606266; font-weight: 700" :title="username">
        {{ username }}
      </span>
    </div>
    <div class="btn" @click="handleCommand('help')">
      <img alt="" class="img" src="@/assets/tools/word.png" />
      <span>帮助文档</span>
    </div>
    <div class="btn" @click="handleCommand('logout')">
      <img alt="" class="img" src="@/assets/tools/layout.png" />
      <span>{{ translate('退出登录') }}</span>
    </div>
    <div class="btn" @click="handleCommand('reset')">
      <img alt="" class="img" src="@/assets/tools/password.png" />
      <span>修改密码</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .btn {
    margin: 0 auto;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 130px;
    height: 43px;
    background: #ffffff;
    border-radius: 6px;
    margin-top: 10px;
    color: #868da0;
    border: 1px solid #e4eaf0;
    font-size: 13px;
    &:hover {
      cursor: pointer;
    }
    img {
      height: 16px;
      margin-right: 6px;
    }
  }
  .avatar {
    margin: 0px auto;
    text-align: center;
    height: 190px;
    // width: 200px !important;
    // padding: 10px 0;
    .user {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
  .avatar-dropdown {
    display: flex;
    align-content: center;
    align-items: center;
    justify-content: center;
    justify-items: center;

    .user-avatar {
      flex-shrink: 0;
      width: 40px;
      height: 40px;
      margin-left: 15px;
      cursor: pointer;
      border-radius: 50%;
    }

    .user-name {
      position: relative;
      display: flex;
      flex-shrink: 0;
      align-content: center;
      align-items: center;
      height: 40px;
      margin-left: 6px;
      line-height: 40px;
      cursor: pointer;

      span {
        max-width: 100px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      [class*='ri-'] {
        margin-left: 3px !important;
      }
    }
  }
</style>
