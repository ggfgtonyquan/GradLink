<template>
  <div class="chips">
    <button
      v-if="showSort"
      type="button"
      class="chip"
      :class="{ active: sort === 'latest' }"
      @click="$emit('update:sort', 'latest')"
    >
      最新
    </button>
    <button
      v-if="showSort"
      type="button"
      class="chip"
      :class="{ active: sort === 'hot' }"
      @click="$emit('update:sort', 'hot')"
    >
      热门
    </button>
    <button
      v-for="tag in tags"
      :key="tag.id"
      type="button"
      class="chip"
      :class="{ active: selectedIds.includes(tag.id) }"
      @click="$emit('toggle-tag', tag.id)"
    >
      {{ tag.name }}
    </button>
  </div>
</template>

<script setup lang="ts">
import type { Tag } from '../api/types'

withDefaults(
  defineProps<{
    tags: Tag[]
    selectedIds?: number[]
    sort?: 'latest' | 'hot'
    showSort?: boolean
  }>(),
  {
    selectedIds: () => [],
    sort: 'latest',
    showSort: false,
  },
)

defineEmits<{
  'update:sort': ['latest' | 'hot']
  'toggle-tag': [number]
}>()
</script>

<style scoped>
.chips {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 2px 0;
  scrollbar-width: none;
}

.chips::-webkit-scrollbar {
  display: none;
}

.chip {
  white-space: nowrap;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: var(--surface);
  color: var(--muted);
  padding: 7px 11px;
  font-size: 13px;
}

.chip.active {
  border-color: var(--brand);
  background: var(--brand-soft);
  color: var(--brand);
  font-weight: 700;
}
</style>
