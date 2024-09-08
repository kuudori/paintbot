<template>
  <div class="vp-toolbar-wrapper">
  <div class="vp-toolbar">

    <slot name="actions">
      <div class="vp-actions">
        <input type="color" class="vp-setting-color" @click="telegram.webapp.HapticFeedback.selectionChanged();"
               v-model="settings.color">
        <slot name="action-undo">
          <button @click="telegram.webapp.HapticFeedback.selectionChanged();emit('undo')" class="vp-action-undo">
            <img src="/src/assets/icons/undo.svg" alt="Undo"/>
          </button>
        </slot>
        <slot name="action-redo">
          <button @click="telegram.webapp.HapticFeedback.selectionChanged();emit('redo')" class="vp-action-redo">
            <img src="/src/assets/icons/redo.svg" alt="Redo"/>
          </button>
        </slot>
      </div>
    </slot>

    <slot name="tools">
      <div class="vp-tools">
        <div
            v-for="tool in visibleTools"
            :key="tool.type"
            class="tool-container"
        >
          <button
              :class="[settings.tool === tool.type ? 'active' : '', `vp-tool-${tool.type}`]"
              @click="selectTool(tool.type)"
              :title="tool.type"
              v-html="tool.icon"
          ></button>
          <transition name="slide-up">
            <div v-if="showThickness === tool.type" class="vp-setting-thickness-wrapper">
              <input
                  type="range"
                  min="1"
                  max="100"
                  class="vp-setting-thickness"
                  v-model="settings.thickness"
                  @input="onThicknessChange"
              >
            </div>
          </transition>
        </div>

        <div class="tools-container">
          <transition name="slide-up">
            <div v-if="showTools" class="tools-dropdown">
              <button
                  v-for="tool in groupedTools"
                  :key="tool.type"
                  :class="[settings.tool === tool.type ? 'active' : '', `vp-tool-${tool.type}`]"
                  @click="selectTool(tool.type)"
                  :title="tool.type"
                  v-html="tool.icon"
              ></button>
            </div>
          </transition>
          <button @click="toggleTools" class="tools-toggle" :class="{ 'active': showTools }">
            <img v-if="!showTools" src="/src/assets/icons/plus.svg" alt="Show more tools"/>
            <img v-else src="/src/assets/icons/minus.svg" alt="Hide tools"/>
          </button>
        </div>
      </div>
    </slot>
  </div>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import type {Settings, Shape, Tool, ToolType} from '@/types'
import telegram from "@/utils/telegram";

const props = defineProps<{ tools: Tool<Shape>[] }>()

const emit = defineEmits<{
  (e: 'set-tool', tool: ToolType): void
  (e: 'save'): void
  (e: 'undo'): void
  (e: 'redo'): void
}>()

const settings = defineModel<Settings>('settings', {
  default: () => ({
    tool: 'freehand',
    thickness: 5,
    color: '#c82d2d'
  })
})

const onThicknessChange = () => {
  telegram.webapp.HapticFeedback.impactOccurred('soft')
}

const showTools = ref(false)
const showThickness = ref<ToolType | null>(null)

const VISIBLE_TOOL_TYPES: ToolType[] = ['freehand', 'eraser', 'textarea', 'move']

const visibleTools = computed(() =>
    props.tools.filter(tool => VISIBLE_TOOL_TYPES.includes(tool.type))
)

const groupedTools = computed(() =>
    props.tools.filter(tool => !VISIBLE_TOOL_TYPES.includes(tool.type))
)

const toggleTools = () => {
  telegram.webapp.HapticFeedback.selectionChanged();
  showThickness.value = null
  showTools.value = !showTools.value
}

const selectTool = (toolType: ToolType) => {
  telegram.webapp.HapticFeedback.selectionChanged();
  if (settings.value.tool === 'freehand') {
    showThickness.value = showThickness.value === toolType ? null : toolType
  } else {
    settings.value.tool = toolType
    showThickness.value = null
  }
  emit('set-tool', toolType)
  showTools.value = false
}

onMounted(() => {
  telegram.webapp.onEvent('mainButtonClicked', save);
  telegram.showSaveButton();
})

async function save() {
  await emit('save')
}
</script>

<style scoped>

.vp-toolbar-wrapper {
  background: antiquewhite;
  height: 110px;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: flex-end;
}

.vp-toolbar {
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  background: antiquewhite;
  z-index: 1000;
  width: 100%;
  margin-bottom: 10px;
}

.vp-actions, .vp-tools {
  display: flex;
  align-items: center;
}

.vp-toolbar button {
  background: none;
  border: none;
  cursor: pointer;

  border-radius: 4px;
}

.vp-toolbar button:hover {
  background-color: rgba(0, 0, 0, 0.1);
}

.vp-toolbar button.active {
  background-color: rgba(0, 0, 0, 0.2);
}

.vp-setting-color {
  border-radius: 50%;
  left: -5%;

  position: relative;
  padding: 0;
  border: none;
  overflow: hidden;
}

.vp-tools {
  display: flex;
  align-items: center;
}

.tool-container {
  position: relative;
}

.tools-container {
  position: relative;
}

.tools-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 45px;
  height: 45px;
  z-index: 1000;
}

.tools-dropdown {
  position: absolute;
  bottom: 105%;
  left: -500%;
  z-index: 1000;
  display: flex;

  margin-bottom: 5px;
}

.vp-setting-thickness-wrapper {
  position: absolute;
  bottom: 100%;
  left: -250%;
  z-index: 1000;
  display: flex;

  margin-bottom: 5px;
}

.vp-setting-thickness {
  -webkit-appearance: none;
  width: 300px;
  height: 10px;
  margin: 20px 0;
  background: #e0e0e0;
  outline: none;
  border-radius: 5px;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.2);
}

.vp-setting-thickness::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 20px;
  height: 20px;
  background: black;
  cursor: pointer;
  border-radius: 50%;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.vp-setting-thickness::-moz-range-thumb {
  width: 20px;
  height: 20px;
  background: black;
  cursor: pointer;
  border-radius: 50%;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>