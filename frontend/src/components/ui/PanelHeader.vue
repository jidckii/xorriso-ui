<template>
  <div>
    <!-- Row 1: top toolbar -->
    <div class="flex items-center gap-1 px-2 h-[30px] bg-gray-100 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700/50">
      <slot name="row1" />
    </div>
    <!-- Row 2: path / info bar -->
    <div class="flex items-center gap-1 px-3 h-[34px] bg-gray-50 dark:bg-gray-800/50 border-b border-gray-300 dark:border-gray-700">
      <slot name="row2" />
    </div>
  </div>
</template>
