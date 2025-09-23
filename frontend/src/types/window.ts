import { ProcessInfo, WindowInfo } from '../../bindings/hptools';

export interface WindowDimensions {
  width: number;
  height: number;
  x: number;
  y: number;
}

export interface SizePreset {
  name: string;
  w: number;
  h: number;
  x?: number;
  y?: number;
}

export interface UseProcessesReturn {
  processes: ProcessInfo[];
  selectedProcess: ProcessInfo | null;
  loading: boolean;
  debugMode: boolean;
  setSelectedProcess: (process: ProcessInfo | null) => void;
  setDebugMode: (debug: boolean) => void;
  fetchProcesses: () => Promise<void>;
}

export interface UseWindowControlReturn {
  dimensions: WindowDimensions;
  currentWindowInfo: WindowInfo | null;
  loading: boolean;
  setDimensions: (dimensions: Partial<WindowDimensions>) => void;
  setWindowSize: (process: ProcessInfo) => Promise<void>;
  setWindowPosition: (process: ProcessInfo) => Promise<void>;
  getWindowInfo: (process: ProcessInfo) => Promise<void>;
  clearWindowInfo: () => void;
}

export interface UseStatusReturn {
  status: string;
  setStatus: (status: string) => void;
  clearStatus: () => void;
}