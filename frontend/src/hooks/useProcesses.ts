import { useState, useEffect } from 'react';
import { WindowService, ProcessInfo } from '../../bindings/hptools';
import { UseProcessesReturn } from '../types/window';
import { STATUS_MESSAGES } from '../constants/window';

export const useProcesses = (
  setStatus: (status: string) => void
): UseProcessesReturn => {
  const [processes, setProcesses] = useState<ProcessInfo[]>([]);
  const [selectedProcess, setSelectedProcess] = useState<ProcessInfo | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [debugMode, setDebugMode] = useState<boolean>(false);

  const fetchProcesses = async () => {
    try {
      setLoading(true);
      const apps = debugMode 
        ? await WindowService.GetAllProcessesWithWindows()
        : await WindowService.GetApplicationProcesses();
      setProcesses(apps);
      setStatus(STATUS_MESSAGES.PROCESSES_FOUND(apps.length, debugMode));
    } catch (error) {
      console.error('Error fetching processes:', error);
      setStatus(STATUS_MESSAGES.ERROR(error));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProcesses();
  }, [debugMode]);

  const handleSetSelectedProcess = (process: ProcessInfo | null) => {
    setSelectedProcess(process);
    if (process) {
      setStatus(STATUS_MESSAGES.PROCESS_SELECTED(process.imageName, process.windowTitle));
    } else {
      setStatus('No process selected');
    }
  };

  return {
    processes,
    selectedProcess,
    loading,
    debugMode,
    setSelectedProcess: handleSetSelectedProcess,
    setDebugMode,
    fetchProcesses,
  };
};