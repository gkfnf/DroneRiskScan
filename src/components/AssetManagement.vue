<template>
  <div class="space-y-6">
    <!-- Asset List -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2">
        <div class="bg-white rounded-lg shadow-sm border">
          <div class="p-6 border-b">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <Database class="h-5 w-5" />
                <h3 class="text-lg font-semibold">资产清单</h3>
              </div>
              <button
                @click="showAddAsset = true"
                class="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
              >
                <Plus class="h-4 w-4 mr-2" />
                添加资产
              </button>
            </div>
          </div>
          
          <div class="p-6">
            <div class="space-y-4">
              <div
                v-for="asset in systemStore.assets"
                :key="asset.id"
                class="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50 transition-colors"
              >
                <div class="flex items-center gap-3">
                  <component :is="getAssetIcon(asset.type)" class="h-5 w-5 text-gray-500" />
                  <div>
                    <h4 class="font-medium">{{ asset.name }}</h4>
                    <p class="text-sm text-gray-500">{{ asset.ipAddress }}</p>
                  </div>
                </div>

                <div class="flex items-center gap-3">
                  <div class="flex items-center gap-2">
                    <div :class="[
                      'w-2 h-2 rounded-full',
                      getStatusColor(asset.status)
                    ]" />
                    <span class="text-sm capitalize">{{ getStatusText(asset.status) }}</span>
                  </div>

                  <span :class="[
                    'px-2 py-1 text-xs rounded-full',
                    getVulnerabilityBadgeClass(asset)
                  ]">
                    {{ getVulnerabilityText(asset) }}
                  </span>

                  <button
                    @click="startScan(asset.id)"
                    :disabled="asset.status === 'scanning'"
                    class="inline-flex items-center px-3 py-1 bg-blue-600 text-white text-sm rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                  >
                    <Search class="h-4 w-4 mr-1" />
                    扫描
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Add Asset Form -->
      <div v-if="showAddAsset" class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center gap-2 mb-4">
          <Plus class="h-5 w-5" />
          <h3 class="text-lg font-semibold">添加资产</h3>
        </div>
        
        <form @submit.prevent="handleAddAsset" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">设备名称</label>
            <input
              v-model="newAsset.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="输入设备名称"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">设备类型</label>
            <select
              v-model="newAsset.type"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="drone">无人机</option>
              <option value="gcs">地面控制站</option>
              <option value="server">服务器</option>
              <option value="network">网络设备</option>
              <option value="hangar">机库</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">IP地址</label>
            <input
              v-model="newAsset.ipAddress"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="192.168.1.100"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">位置</label>
            <input
              v-model="newAsset.location"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="设备位置"
            />
          </div>

          <button
            type="submit"
            :disabled="isSubmitting"
            class="w-full inline-flex items-center justify-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Plus class="h-4 w-4 mr-2" />
            {{ isSubmitting ? '添加中...' : '添加资产' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useSystemStore } from '@/stores/system'
import { 
  Database, 
  Plus, 
  Search,
  Zap,
  Activity,
  Server,
  Wifi,
  Shield
} from 'lucide-vue-next'
import type { Asset } from '@/types'

const systemStore = useSystemStore()
const showAddAsset = ref(false)
const isSubmitting = ref(false)

const newAsset = reactive({
  name: '',
  type: 'drone' as Asset['type'],
  ipAddress: '',
  location: '',
  zone: '',
  status: 'offline' as Asset['status']
})

const getAssetIcon = (type: Asset['type']) => {
  const icons = {
    drone: Zap,
    gcs: Activity,
    server: Server,
    network: Wifi,
    hangar: Shield
  }
  return icons[type] || Server
}

const getStatusColor = (status: Asset['status']) => {
  const colors = {
    online: 'bg-green-500',
    offline: 'bg-red-500',
    scanning: 'bg-yellow-500',
    maintenance: 'bg-blue-500'
  }
  return colors[status] || 'bg-gray-500'
}

const getStatusText = (status: Asset['status']) => {
  const texts = {
    online: '在线',
    offline: '离线',
    scanning: '扫描中',
    maintenance: '维护中'
  }
  return texts[status] || status
}

const getVulnerabilityBadgeClass = (asset: Asset) => {
  // 这里应该根据实际的漏洞数量来判断
  const vulnCount = 0 // 临时值
  if (vulnCount === 0) return 'bg-green-100 text-green-800'
  if (vulnCount <= 2) return 'bg-yellow-100 text-yellow-800'
  return 'bg-red-100 text-red-800'
}

const getVulnerabilityText = (asset: Asset) => {
  // 这里应该根据实际的漏洞数量来判断
  const vulnCount = 0 // 临时值
  if (vulnCount === 0) return '安全'
  if (vulnCount <= 2) return '低风险'
  return '高风险'
}

const handleAddAsset = async () => {
  try {
    isSubmitting.value = true
    await systemStore.createAsset({
      ...newAsset,
      createdAt: new Date(),
      updatedAt: new Date()
    })
    
    // 重置表单
    Object.assign(newAsset, {
      name: '',
      type: 'drone' as Asset['type'],
      ipAddress: '',
      location: '',
      zone: '',
      status: 'offline' as Asset['status']
    })
    
    showAddAsset.value = false
  } catch (error) {
    console.error('添加资产失败:', error)
  } finally {
    isSubmitting.value = false
  }
}

const startScan = async (assetId: string) => {
  try {
    await systemStore.startScan(assetId, 'vulnerability')
  } catch (error) {
    console.error('启动扫描失败:', error)
  }
}
</script>
