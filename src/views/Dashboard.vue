<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-6 py-4">
        <div class="flex items-center gap-3">
          <Shield class="h-8 w-8 text-blue-600" />
          <div>
            <h1 class="text-2xl font-bold text-gray-900">电力无人机安全扫描系统</h1>
            <p class="text-gray-600">专业的无人机设备安全漏洞检测与评估平台</p>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-6 py-8">
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <StatsCard
          title="总资产数"
          :value="systemStore.stats.totalAssets"
          :icon="Server"
          color="blue"
        />
        <StatsCard
          title="在线设备"
          :value="systemStore.stats.onlineAssets"
          :icon="CheckCircle"
          color="green"
        />
        <StatsCard
          title="射频威胁"
          :value="systemStore.stats.rfThreats"
          :icon="Radio"
          color="red"
        />
        <StatsCard
          title="高危漏洞"
          :value="systemStore.stats.criticalVulnerabilities"
          :icon="AlertTriangle"
          color="orange"
        />
      </div>

      <!-- Navigation Tabs -->
      <div class="bg-white rounded-lg shadow-sm border mb-8">
        <nav class="flex space-x-8 px-6" aria-label="Tabs">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm',
              activeTab === tab.id
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            <component :is="tab.icon" class="h-5 w-5 inline mr-2" />
            {{ tab.name }}
          </button>
        </nav>
      </div>

      <!-- Tab Content -->
      <div class="space-y-6">
        <!-- Assets Tab -->
        <div v-if="activeTab === 'assets'">
          <AssetManagement />
        </div>

        <!-- Scanning Tab -->
        <div v-if="activeTab === 'scanning'">
          <ScanningMonitor />
        </div>

        <!-- RF Security Tab -->
        <div v-if="activeTab === 'rf-security'">
          <RFSecurityPanel />
        </div>

        <!-- Reports Tab -->
        <div v-if="activeTab === 'reports'">
          <ReportsPanel />
        </div>

        <!-- Settings Tab -->
        <div v-if="activeTab === 'settings'">
          <SystemSettings />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSystemStore } from '@/stores/system'
import { 
  Shield, 
  Server, 
  CheckCircle, 
  Radio, 
  AlertTriangle,
  Database,
  Activity,
  Radar,
  FileText,
  Settings
} from 'lucide-vue-next'

import StatsCard from '@/components/StatsCard.vue'
import AssetManagement from '@/components/AssetManagement.vue'
import ScanningMonitor from '@/components/ScanningMonitor.vue'
import RFSecurityPanel from '@/components/RFSecurityPanel.vue'
import ReportsPanel from '@/components/ReportsPanel.vue'
import SystemSettings from '@/components/SystemSettings.vue'

const systemStore = useSystemStore()
const activeTab = ref('assets')

const tabs = [
  { id: 'assets', name: '资产管理', icon: Database },
  { id: 'scanning', name: '扫描监控', icon: Activity },
  { id: 'rf-security', name: '射频安全', icon: Radar },
  { id: 'reports', name: '安全报告', icon: FileText },
  { id: 'settings', name: '系统设置', icon: Settings }
]

onMounted(() => {
  systemStore.initialize()
})
</script>
