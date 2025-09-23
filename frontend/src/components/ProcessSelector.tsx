import React from 'react';
import { ProcessInfo } from '../../bindings/hptools';

interface ProcessSelectorProps {
  processes: ProcessInfo[];
  selectedProcess: ProcessInfo | null;
  loading: boolean;
  debugMode: boolean;
  onProcessSelect: (process: ProcessInfo | null) => void;
  onRefresh: () => void;
  onDebugModeChange: (debugMode: boolean) => void;
}

export const ProcessSelector: React.FC<ProcessSelectorProps> = ({
  processes,
  selectedProcess,
  loading,
  debugMode,
  onProcessSelect,
  onRefresh,
  onDebugModeChange,
}) => {
  const handleProcessChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const pid = parseInt(e.target.value);
    const process = processes.find(p => p.pid === pid);
    onProcessSelect(process || null);
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-6">
      <h2 className="text-xl font-semibold mb-4">Select Application</h2>
      
      <div className="flex gap-4 mb-4">
        <select
          value={selectedProcess?.pid || ''}
          onChange={handleProcessChange}
          className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Select an application...</option>
          {processes.map((process) => (
            <option key={process.pid} value={process.pid}>
              {process.imageName} - "{process.windowTitle}" (PID: {process.pid}) [{process.windowCount} windows]
            </option>
          ))}
        </select>
        
        <button
          onClick={onRefresh}
          disabled={loading}
          className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 disabled:opacity-50"
        >
          {loading ? '🔄' : '🔄 Refresh'}
        </button>
      </div>

      <div className="flex items-center gap-4 mb-4">
        <label className="flex items-center gap-2">
          <input
            type="checkbox"
            checked={debugMode}
            onChange={(e) => onDebugModeChange(e.target.checked)}
            className="rounded"
          />
          <span className="text-sm text-gray-600">
            Debug Mode (Show all processes with windows)
          </span>
        </label>
        
        {debugMode && (
          <div className="text-xs text-yellow-600 bg-yellow-50 px-2 py-1 rounded">
            ⚠️ Debug mode shows all processes, including system ones
          </div>
        )}
      </div>
    </div>
  );
};