import { useState } from 'react';
import {  ProcessInfo, WindowInfo } from '../../bindings/hptools/internal/models';
import { WailsWindowService } from '../../bindings/hptools/internal/services';
import { UseWindowControlReturn, WindowDimensions } from '../types/window';
import { DEFAULT_DIMENSIONS, STATUS_MESSAGES } from '../constants/window';

export const useWindowControl = (
  setStatus: (status: string) => void
): UseWindowControlReturn => {
  const [dimensions, setDimensionsState] = useState<WindowDimensions>(DEFAULT_DIMENSIONS);
  const [currentWindowInfo, setCurrentWindowInfo] = useState<WindowInfo | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const setDimensions = (newDimensions: Partial<WindowDimensions>) => {
    setDimensionsState(prev => ({ ...prev, ...newDimensions }));
  };

  const clearWindowInfo = () => {
    setCurrentWindowInfo(null);
  };

  const setWindowSize = async (process: ProcessInfo) => {
    if (!process) {
      setStatus(STATUS_MESSAGES.NO_PROCESS_SELECTED);
      return;
    }

    try {
      setLoading(true);
      await WailsWindowService.SetWindowSize(process.pid, dimensions.width, dimensions.height);
      setStatus(STATUS_MESSAGES.WINDOW_RESIZED(dimensions.width, dimensions.height, process.imageName));
      // Refresh window info after resize
      await getWindowInfo(process);
    } catch (error) {
      console.error('Error setting window size:', error);
      setStatus(STATUS_MESSAGES.ERROR(error));
    } finally {
      setLoading(false);
    }
  };

  const setWindowPosition = async (process: ProcessInfo) => {
    if (!process) {
      setStatus(STATUS_MESSAGES.NO_PROCESS_SELECTED);
      return;
    }

    try {
      setLoading(true);
      await WailsWindowService.SetWindowPosition(process.pid, dimensions.x, dimensions.y, dimensions.width, dimensions.height);
      setStatus(STATUS_MESSAGES.WINDOW_MOVED(dimensions.x, dimensions.y, dimensions.width, dimensions.height, process.imageName));
      // Refresh window info after move/resize
      await getWindowInfo(process);
    } catch (error) {
      console.error('Error setting window position:', error);
      setStatus(STATUS_MESSAGES.ERROR(error));
    } finally {
      setLoading(false);
    }
  };

  const getWindowInfo = async (process: ProcessInfo) => {
    if (!process) {
      setStatus(STATUS_MESSAGES.NO_PROCESS_SELECTED);
      return;
    }

    try {
      setLoading(true);
      const info = await WailsWindowService.GetWindowInfo(process.pid);
      if (info) {
        setCurrentWindowInfo(info);
        setStatus(STATUS_MESSAGES.WINDOW_INFO(info.width, info.height, info.x, info.y));
        
        // Update form fields with current values
        setDimensions({
          width: info.width,
          height: info.height,
          x: info.x,
          y: info.y,
        });
      }
    } catch (error) {
      console.error('Error getting window info:', error);
      setStatus(STATUS_MESSAGES.ERROR(error));
    } finally {
      setLoading(false);
    }
  };

  return {
    dimensions,
    currentWindowInfo,
    loading,
    setDimensions,
    setWindowSize,
    setWindowPosition,
    getWindowInfo,
    clearWindowInfo,
  };
};