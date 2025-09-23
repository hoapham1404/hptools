import { SizePreset } from '../types/window';

export const DEFAULT_DIMENSIONS = {
  width: 800,
  height: 600,
  x: 100,
  y: 100,
} as const;

export const INPUT_LIMITS = {
  width: { min: 100, max: 3840 },
  height: { min: 100, max: 2160 },
  position: { min: 0, max: 5000 },
} as const;

export const SIZE_PRESETS: SizePreset[] = [
  { name: "HD", w: 1280, h: 720 },
  { name: "FHD", w: 1920, h: 1080 },
  { name: "Square", w: 800, h: 800 },
  { name: "Small", w: 640, h: 480 },
  { name: "Wide", w: 1200, h: 600 },
  { name: "Custom 1", w: 1860, h: 1000, x: 30, y: 30 },
] as const;

export const STATUS_MESSAGES = {
  NO_PROCESS_SELECTED: 'Please select a process first',
  PROCESS_SELECTED: (imageName: string, windowTitle: string) => 
    `Selected: ${imageName} - "${windowTitle}"`,
  PROCESSES_FOUND: (count: number, isDebug: boolean) => 
    `Found ${count} ${isDebug ? 'processes with windows' : 'application processes'}`,
  WINDOW_RESIZED: (width: number, height: number, imageName: string) => 
    `✅ Set window size to ${width}x${height} for ${imageName}`,
  WINDOW_MOVED: (x: number, y: number, width: number, height: number, imageName: string) => 
    `✅ Set window position to (${x}, ${y}) and size to ${width}x${height} for ${imageName}`,
  WINDOW_INFO: (width: number, height: number, x: number, y: number) => 
    `📏 Current window: ${width}x${height} at position (${x}, ${y})`,
  ERROR: (error: unknown) => `❌ Error: ${error}`,
} as const;