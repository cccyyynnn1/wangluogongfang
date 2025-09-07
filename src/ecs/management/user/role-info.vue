<script lang="ts">
  export default {
    name: 'RoleInfo',
  }
</script>

<script setup lang="ts">
  defineProps<{
    showName: string
  }>()

  const emit = defineEmits<{
    (e: 'on-saveData', val: string): void
  }>()

  const activeName = ref('info') // tabs选中项

  // 保存
  const saveData = () => {
    emit('on-saveData', 'index')
  }
</script>

<template>
  <div class="role-info-container">
    <el-tabs v-model="activeName">
      <el-tab-pane label="账号绑定" name="info">
        <el-button class="btn" type="primary" @click="saveData">返回</el-button>
        <div class="content">
          <div class="info">
            <div class="text">角色详情</div>
            <el-row class="list">
              <el-col :span="3">
                <div class="grid-content bg-gray">角色名称</div>
              </el-col>
              <el-col :span="3">
                <div class="grid-content">审计人员</div>
              </el-col>
              <el-col :span="3">
                <div class="grid-content bg-gray">描述</div>
              </el-col>
              <el-col :span="9">
                <div class="grid-content">负责操作日志的审计工作</div>
              </el-col>
              <el-col :span="3">
                <div class="grid-content bg-gray">所属部门</div>
              </el-col>
              <el-col :span="3">
                <div class="grid-content">审计部</div>
              </el-col>
            </el-row>
          </div>
          <div class="info">
            <div class="text">权限详情</div>
            <el-row class="list">
              <el-col :span="3">
                <div class="grid-content grid-content-bot bg-gray">
                  模块权限
                </div>
              </el-col>
              <el-col :span="3">
                <div class="grid-content grid-content-bot">
                  <div class="grid-content-item">
                    <div class="item-text">系统管理</div>
                    <div class="item-text">一级分类</div>
                  </div>
                </div>
              </el-col>
              <el-col :span="3">
                <div class="grid-content grid-content-bot bg-gray">
                  <div class="grid-content-item">
                    <div class="grid-content-item-2">
                      <div class="item-text">系统状态</div>
                      <div class="item-text">审计管理</div>
                    </div>
                    <div class="item-text">二级分类</div>
                  </div>
                </div>
              </el-col>
              <el-col :span="9">
                <div class="grid-content grid-content-bot">
                  <div class="grid-content-item">
                    <div class="grid-content-item-2">
                      <div class="item-text-1">集群及节点状态</div>
                      <div class="item-text-1">docker状态</div>
                      <div class="item-text-1">操作日志</div>
                      <div class="item-text-1">系统日志</div>
                    </div>
                    <div class="grid-content-item-2">
                      <div class="item-text">三级分类</div>
                      <div class="item-text">三级分类</div>
                    </div>
                  </div>
                </div>
              </el-col>
              <el-col :span="6">
                <div class="grid-content grid-content-bot bg-gray">
                  <div class="grid-content-item">
                    <div class="grid-content-item-2">
                      <div class="item-text-1">只读</div>
                      <div class="item-text-1">只读</div>
                      <div class="item-text-1">可编辑</div>
                      <div class="item-text-1">只读</div>
                    </div>
                    <div class="grid-content-item-2">
                      <div class="item-text">权限</div>
                      <div class="item-text">权限</div>
                    </div>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style lang="scss" scoped>
  .role-info {
    .btn {
      margin: 10px 0 0 60px;
    }

    .content {
      padding: 0px 120px;

      .info {
        margin-top: 20px;

        .text {
          color: #178cf7;
        }

        .list {
          margin-top: 15px;

          .grid-content {
            display: flex;
            align-items: center;
            justify-content: center;
            height: 40px;
            margin: 0 -1px 0 0;
            line-height: 40px;
            text-align: center;
            border: 1px solid #eee;

            .grid-content-item {
              width: 100%;
              height: 100%;
            }

            .grid-content-item-2 {
              height: 50%;
            }

            .item-text {
              display: flex;
              align-items: center;
              justify-content: center;
              width: 100%;
              height: 50%;
              border-bottom: 1px solid #eee;
            }

            .item-text-1 {
              height: 25%;
              border-bottom: 1px solid #eee;
            }
          }

          .grid-content-bot {
            height: 320px;
          }

          .bg-gray {
            background: #f5f7fa;
          }
        }
      }
    }
  }
</style>
