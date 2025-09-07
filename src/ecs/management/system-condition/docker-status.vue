<script lang="ts">
  export default {
    name: 'DockerStatus',
  }
</script>

<script setup lang="ts">
  import { Refresh, Search } from '@element-plus/icons-vue'
  import { useTableCopy } from '@/utils'

  const listLoading = ref(false) // 是否加载

  // 表格数据
  const listDate = reactive<object[]>([
    {
      createTime: '2012-12-30 12:34:45',
      dockerNmae: '名称',
      description: '审计管理描述',
      status: 1,
    },
    {
      createTime: '2012-12-30 12:34:45',
      dockerNmae: '名称',
      description: '审计管理描述',
      status: 0,
    },
    {
      createTime: '2012-12-30 12:34:45',
      dockerNmae: '名称',
      description: '审计管理描述',
      status: 0,
    },
    {
      createTime: '2012-12-30 12:34:45',
      dockerNmae: '名称',
      description: '审计管理描述',
      status: 2,
    },
    {
      createTime: '2012-12-30 12:34:45',
      dockerNmae: '名称',
      description: '审计管理描述',
      status: 2,
    },
  ])
  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive({
    total: 0,
    pageNo: 1,
    pageSize: 10,
    title: '',
  })

  const queryForm = reactive({
    role: '',
  })

  const formData = reactive({
    departmentName: '',
  })

  // 刷新时间
  const refresh = () => {}

  // 查询
  const queryData = () => {}

  // 页容量改变
  const handleSizeChange = () => {}

  // 页面改变
  const handleCurrentChange = () => {}

  // 操作
  const handleChange = (val: boolean) => {}
</script>

<template>
  <div class="docker-status">
    <div class="top">
      <div class="left">Docker状态</div>
      <div class="right">
        <span>最近一次更新时间：2022-02-12</span>
        <el-icon class="btn">
          <Refresh @click="refresh" />
        </el-icon>
      </div>
    </div>
    <div class="content">
      <vab-query-form>
        <vab-query-form-left-panel :span="12" />
        <vab-query-form-right-panel :span="12">
          <el-form inline :model="queryForm" @submit.prevent>
            <el-form-item>
              <el-input v-model.trim="queryForm.role" clearable placeholder="输入Docker名称" />
            </el-form-item>
            <el-form-item>
              <el-button :icon="Search" type="primary" @click="queryData">检索</el-button>
            </el-form-item>
          </el-form>
        </vab-query-form-right-panel>
      </vab-query-form>
      <el-table
        v-loading="listLoading"
        :border="true"
        class="my-table"
        :data="listDate"
        default-expand-all
        row-key="id"
        :tree-props="{ children: 'children' }"
        @cell-contextmenu="useTableCopy"
      >
        <el-table-column :align="'center'" label="创建时间" prop="createTime" show-overflow-tooltip width="200" />
        <el-table-column :align="'center'" label="Docker名称" prop="dockerNmae" show-overflow-tooltip width="200" />
        <el-table-column :align="'center'" label="描述" prop="description" show-overflow-tooltip />
        <el-table-column :align="'center'" label="状态" prop="status" show-overflow-tooltip width="200">
          <template #default="{ row }">
            <el-button v-if="row.status === 1" plain type="success">进行中</el-button>
            <el-button v-else-if="row.status === 0" plain type="danger">停止</el-button>
            <el-button v-else-if="row.status === 2" plain type="primary">重启中</el-button>
          </template>
        </el-table-column>
        <el-table-column :align="'center'" label="操作" width="240">
          <template #default="{ row }">
            <el-button size="small" @click="handleChange(false)">重启</el-button>
            <el-button :disabled="!row.parentId" size="small" @click="handleChange(false)">停止</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty class="vab-data-empty" description="暂无数据" />
        </template>
      </el-table>

      <el-pagination
        background
        :current-page="queryPage.pageNo"
        :layout="layout"
        :page-size="queryPage.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="queryPage.total"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .docker-status {
    .top {
      display: flex;
      align-items: center;
      justify-content: space-between;
      height: 60px;
      padding: 0 20px;
      line-height: 60px;
      border-bottom: 1px solid #eee;

      .left {
        font-size: 16px;
        font-weight: 500;
        color: #303133;
      }

      .right {
        span {
          margin-top: -5px;
          font-size: 13px;
          font-weight: 400;
          color: #a9acb3;
        }

        .btn {
          margin: -10px 0 0 10px;
          font-size: 13px;
          color: #0d88fe;
          vertical-align: -webkit-baseline-middle;

          &:hover {
            cursor: pointer;
          }
        }
      }
    }

    .content {
      padding: 20px;
    }
  }
</style>
