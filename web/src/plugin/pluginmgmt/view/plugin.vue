<template>
  <div>
    <!-- 搜索区域 -->
    <div class="gva-search-box">
      <el-form ref="searchFormRef" :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="插件名称">
          <el-input v-model="searchInfo.name" placeholder="请输入插件名称" />
        </el-form-item>
        <el-form-item label="插件编码">
          <el-input v-model="searchInfo.code" placeholder="请输入插件编码" />
        </el-form-item>
        <el-form-item label="插件类型">
          <el-select v-model="searchInfo.type" placeholder="请选择" clearable>
            <el-option label="内部Go插件" :value="1" />
            <el-option label="外部插件" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="请选择" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="2" />
            <el-option label="异常" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格区域 -->
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" @click="openDialog('add')">新增插件</el-button>
        <el-button
          type="danger"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >批量删除</el-button>
      </div>

      <el-table
        :data="tableData"
        @selection-change="handleSelectionChange"
        border
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="ID" prop="ID" width="80" />
        <el-table-column label="插件名称" prop="name" min-width="150" />
        <el-table-column label="插件编码" prop="code" min-width="120" />
        <el-table-column label="版本" prop="version" width="100" />
        <el-table-column label="类型" prop="type" width="120">
          <template #default="{ row }">
            <el-tag :type="row.type === 1 ? 'success' : 'warning'">
              {{ row.type === 1 ? '内部Go插件' : '外部插件' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" prop="status" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="2"
              @change="changeStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="健康状态" prop="healthStatus" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.healthStatus === 1" type="success">健康</el-tag>
            <el-tag v-else-if="row.healthStatus === 2" type="danger">异常</el-tag>
            <el-tag v-else type="info">未知</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="作者" prop="author" width="120" />
        <el-table-column label="描述" prop="description" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog('edit', row)">编辑</el-button>
            <el-button type="danger" link @click="deletePlugin(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 20, 30, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-drawer
      v-model="dialogVisible"
      :title="dialogTitle"
      size="600px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
        <el-form-item label="插件名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入插件名称" />
        </el-form-item>
        <el-form-item label="插件编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入插件编码" />
        </el-form-item>
        <el-form-item label="版本" prop="version">
          <el-input v-model="formData.version" placeholder="请输入版本号" />
        </el-form-item>
        <el-form-item label="插件类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择" style="width: 100%">
            <el-option label="内部Go插件" :value="1" />
            <el-option label="外部插件" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="作者" prop="author">
          <el-input v-model="formData.author" placeholder="请输入作者" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-input v-model="formData.icon" placeholder="请输入图标地址" />
        </el-form-item>
        <el-form-item label="健康检查URL" prop="healthURL">
          <el-input v-model="formData.healthURL" placeholder="外部插件需要配置" />
        </el-form-item>
        <el-form-item label="Consul服务名" prop="serviceName">
          <el-input v-model="formData.serviceName" placeholder="外部插件需要配置" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入插件描述"
          />
        </el-form-item>
        <el-form-item label="配置" prop="config">
          <el-input
            v-model="formData.config"
            type="textarea"
            :rows="5"
            placeholder="JSON格式配置"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createPlugin,
  deletePlugin,
  updatePlugin,
  getPluginList,
  updatePluginStatus
} from '@/plugin/pluginmgmt/api/plugin'

// 搜索相关
const searchFormRef = ref(null)
const searchInfo = ref({
  name: '',
  code: '',
  type: null,
  status: null
})

// 表格相关
const tableData = ref([])
const multipleSelection = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框相关
const dialogVisible = ref(false)
const dialogTitle = ref('新增插件')
const formRef = ref(null)
const formData = ref({
  name: '',
  code: '',
  version: '1.0.0',
  type: 1,
  author: '',
  icon: '',
  healthURL: '',
  serviceName: '',
  description: '',
  config: ''
})

const rules = {
  name: [{ required: true, message: '请输入插件名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入插件编码', trigger: 'blur' }],
  type: [{ required: true, message: '请选择插件类型', trigger: 'change' }]
}

// 获取列表数据
const getTableData = async () => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  }
  const res = await getPluginList(params)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
  }
}

// 搜索
const onSubmit = () => {
  page.value = 1
  getTableData()
}

// 重置
const onReset = () => {
  searchFormRef.value?.resetFields()
  onSubmit()
}

// 分页变化
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  getTableData()
}

// 选择变化
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 打开对话框
const openDialog = (type, row) => {
  dialogVisible.value = true
  if (type === 'add') {
    dialogTitle.value = '新增插件'
    formData.value = {
      name: '',
      code: '',
      version: '1.0.0',
      type: 1,
      author: '',
      icon: '',
      healthURL: '',
      serviceName: '',
      description: '',
      config: ''
    }
  } else {
    dialogTitle.value = '编辑插件'
    formData.value = { ...row }
  }
}

// 提交表单
const submitForm = async () => {
  await formRef.value.validate()
  
  if (formData.value.ID) {
    await updatePlugin(formData.value)
    ElMessage.success('更新成功')
  } else {
    await createPlugin(formData.value)
    ElMessage.success('创建成功')
  }
  
  dialogVisible.value = false
  getTableData()
}

// 删除单个插件
const deletePlugin = async (row) => {
  await ElMessageBox.confirm('确定要删除该插件吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  
  const res = await deletePlugin({ ids: [row.ID] })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    getTableData()
  }
}

// 批量删除
const onDelete = async () => {
  await ElMessageBox.confirm('确定要删除选中的插件吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  
  const ids = multipleSelection.value.map(item => item.ID)
  const res = await deletePlugin({ ids })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    getTableData()
  }
}

// 修改状态
const changeStatus = async (row) => {
  const res = await updatePluginStatus({
    ID: row.ID,
    status: row.status
  })
  if (res.code === 0) {
    ElMessage.success('状态更新成功')
  } else {
    row.status = row.status === 1 ? 2 : 1
  }
}

onMounted(() => {
  getTableData()
})
</script>

<style scoped>
.gva-search-box {
  margin-bottom: 20px;
}

.gva-table-box {
  padding: 20px;
  background: #fff;
  border-radius: 4px;
}

.gva-btn-list {
  margin-bottom: 20px;
}

.gva-pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
