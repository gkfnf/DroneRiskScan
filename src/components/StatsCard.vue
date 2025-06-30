<template>
  <div class="bg-white rounded-lg shadow-sm border p-6">
    <div class="flex items-center justify-between">
      <div>
        <p class="text-sm font-medium text-gray-600">{{ title }}</p>
        <p :class="[
          'text-2xl font-bold',
          colorClasses[color]
        ]">
          {{ formattedValue }}
        </p>
      </div>
      <component 
        :is="icon" 
        :class="[
          'h-8 w-8',
          iconColorClasses[color]
        ]" 
      />
    </div>
    <div v-if="trend" class="mt-4 flex items-center">
      <component 
        :is="trend.direction === 'up' ? TrendingUp : TrendingDown"
        :class="[
          'h-4 w-4 mr-1',
          trend.direction === 'up' ? 'text-green-500' : 'text-red-500'
        ]"
      />
      <span :class="[
        'text-sm',
        trend.direction === 'up' ? 'text-green-600' : 'text-red-600'
      ]">
        {{ trend.value }}%
      </span>
      <span class="text-sm text-gray-500 ml-1">vs 上周</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { TrendingUp, TrendingDown } from 'lucide-vue-next'

interface Props {
  title: string
  value: number
  icon: any
  color: 'blue' | 'green' | 'red' | 'orange' | 'purple'
  trend?: {
    direction: 'up' | 'down'
    value: number
  }
}

const props = defineProps<Props>()

const colorClasses = {
  blue: 'text-blue-600',
  green: 'text-green-600',
  red: 'text-red-600',
  orange: 'text-orange-600',
  purple: 'text-purple-600'
}

const iconColorClasses = {
  blue: 'text-blue-500',
  green: 'text-green-500',
  red: 'text-red-500',
  orange: 'text-orange-500',
  purple: 'text-purple-500'
}

const formattedValue = computed(() => {
  if (props.value >= 1000) {
    return (props.value / 1000).toFixed(1) + 'K'
  }
  return props.value.toString()
})
</script>
